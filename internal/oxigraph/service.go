package oxigraph

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/dpb587/rdfkit-go/internal/oxigraph/internal"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/x/storage/sparql"
)

var localrand = rand.New(rand.NewSource(time.Now().UnixMicro()))

type Format string

const (
	NQuadsFormat   Format = "nq"
	NTriplesFormat Format = "nt"
	TurtleFormat   Format = "ttl"
	TrigFormat     Format = "trig"
)

type ServiceOptions struct {
	Exec     string
	BindPort int
	Datadir  string
}

type Service struct {
	opts     ServiceOptions
	datadir  string
	p        *os.Process
	bindPort string
}

func NewService(opts ServiceOptions) *Service {
	if len(opts.Exec) == 0 {
		opts.Exec = "oxigraph"
	}

	s := &Service{
		opts: opts,
	}

	return s
}

func (s *Service) requireDatadir() (string, error) {
	if len(s.datadir) > 0 {
		return s.datadir, nil
	} else if len(s.opts.Datadir) > 0 {
		s.datadir = s.opts.Datadir

		err := os.MkdirAll(s.datadir, 0700)
		if err != nil {
			return "", err
		}

		return s.datadir, nil
	}

	datadir, err := os.MkdirTemp("", "oxigraph-")
	if err != nil {
		return "", err
	}

	s.datadir = datadir

	return s.datadir, nil
}

func (s *Service) NewClient() (*sparql.Client, error) {
	err := s.requireProcess()
	if err != nil {
		return nil, err
	}

	return sparql.NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%s/query", s.bindPort)), nil
}

func (s *Service) Close() error {
	if s.p != nil {
		err := s.p.Signal(os.Interrupt)
		if err != nil {
			return fmt.Errorf("process: signal: %v", err)
		}

		_, err = s.p.Wait()
		if err != nil {
			return fmt.Errorf("process: wait: %v", err)
		}
	}

	if len(s.opts.Datadir) == 0 {
		err := os.RemoveAll(s.datadir)
		if err != nil {
			return fmt.Errorf("datadir: rm: %v", err)
		}
	}

	return nil
}

func (s *Service) requireProcess() error {
	if len(s.bindPort) > 0 {
		return nil
	}

	datadir, err := s.requireDatadir()
	if err != nil {
		return err
	}

	for attempt := 0; attempt < 3; attempt++ {
		bindPort := strconv.Itoa(localrand.Intn(65535-1024) + 1024) // TODO conflicts

		serveCommand := exec.Command(
			s.opts.Exec,
			"serve",
			"--location", datadir,
			"--bind", fmt.Sprintf("127.0.0.1:%s", bindPort),
		)

		serveCommand.Stderr = os.Stderr
		serveCommand.Stdout = os.Stderr

		err = serveCommand.Start()
		if err != nil {
			return fmt.Errorf("start server: %v", err)
		}

		s.p = serveCommand.Process

		err = internal.WaitForPortOpen("127.0.0.1", bindPort, 10*time.Second)
		if err != nil {
			s.p.Kill()
			s.p.Wait()

			if rerr := s.Close(); rerr != nil {
				// TODO log
			}

			continue // retry
		}

		res, err := http.DefaultClient.Get(fmt.Sprintf("http://127.0.0.1:%s/", bindPort))
		if err != nil || res.StatusCode != http.StatusOK {
			s.p.Kill()
			s.p.Wait()

			if rerr := s.Close(); rerr != nil {
				// TODO log
			}

			continue // retry
		}

		s.bindPort = bindPort

		return nil
	}

	return errors.New("failed to start server")
}

type ImportOptions struct {
	Base  rdf.IRI
	Graph rdf.GraphNameValue
}

func (s *Service) ImportReader(r io.Reader, format Format, opts ImportOptions) error {
	datadir, err := s.requireDatadir()
	if err != nil {
		return err
	}

	loadCommand := exec.Command(
		s.opts.Exec,
		"load",
		"--location", datadir,
		"--file", "/dev/stdin",
		"--format", string(format),
	)
	loadCommand.Stdin = r
	loadCommand.Stdout = os.Stderr
	loadCommand.Stderr = os.Stderr

	if len(opts.Base) > 0 {
		loadCommand.Args = append(loadCommand.Args, "--base", string(opts.Base))
	}

	if opts.Graph != nil {
		switch g := opts.Graph.(type) {
		case rdf.IRI:
			loadCommand.Args = append(loadCommand.Args, "--graph", string(g))
		default:
			return fmt.Errorf("unsupported graph: %T", opts.Graph)
		}
	}

	err = loadCommand.Run()
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	return nil
}
