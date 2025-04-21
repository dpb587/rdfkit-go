package cmdflags

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	devencodingdiscard "github.com/dpb587/rdfkit-go/internal/devencoding/discard"
	devencodingrdfioutil "github.com/dpb587/rdfkit-go/internal/devencoding/rdfioutil"
	"github.com/dpb587/rdfkit-go/internal/ioutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
)

type EncodingOutput struct {
	Path string
	Type string

	DefaultBase string
}

func (f EncodingOutput) NewStatementWriter() (*EncodingOutputHandle, error) {
	b := &EncodingOutputHandle{}

	var writeCloser func() error
	var writeWriter io.Writer

	if f.Path == "-" {
		writeCloser = func() error { return nil }
		writeWriter = os.Stdout

		b.WritePath = "file:///dev/stdout"
	} else {
		outFile, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, fmt.Errorf("open: %v", err)
		}

		writeCloser = outFile.Close
		writeWriter = outFile

		b.WritePath = f.Path
	}

	writeHasher := sha256.New()
	writeWriter = io.MultiWriter(writeHasher, writeWriter)

	b.writeHasher = writeHasher
	b.writer = ioutil.NewWriterCloser(writeWriter, writeCloser)

	var err error

	switch f.Type {
	case "devdiscard":
		b.Encoder = devencodingdiscard.Encoder
	case "devrdfioutil":
		b.Encoder = devencodingrdfioutil.NewEncoder(
			writeWriter,
			devencodingrdfioutil.EncoderOptions{},
		)
	case "nquads", "nq":
		b.Encoder, err = nquads.NewEncoder(
			writeWriter,
		)
		if err != nil {
			writeCloser()

			return nil, fmt.Errorf("nquads: %v", err)
		}
	case "nquads/ascii":
		b.Encoder, err = nquads.NewEncoder(
			writeWriter,
			nquads.EncoderConfig{}.
				SetASCII(true),
		)
		if err != nil {
			writeCloser()

			return nil, fmt.Errorf("nquads: %v", err)
		}
	case "ntriples", "nt":
		b.Encoder, err = ntriples.NewEncoder(
			writeWriter,
			ntriples.EncoderConfig{},
		)
		if err != nil {
			writeCloser()

			return nil, fmt.Errorf("ntriples: %v", err)
		}
	case "ntriples/ascii":
		b.Encoder, err = ntriples.NewEncoder(
			writeWriter,
			ntriples.EncoderConfig{}.
				SetASCII(true),
		)
		if err != nil {
			writeCloser()

			return nil, fmt.Errorf("ntriples: %v", err)
		}
	case "rdfjson", "rj":
		b.Encoder, err = rdfjson.NewEncoder(
			writeWriter,
		)
		if err != nil {
			writeCloser()

			return nil, fmt.Errorf("rdfjson: %v", err)
		}
	case "turtle", "ttl":
		b.Encoder, err = turtle.NewEncoder(
			writeWriter,
			turtle.EncoderConfig{}.
				SetBase(f.DefaultBase).
				SetPrefixes(iriutil.NewPrefixMap(rdfacontext.InitialContext()...)),
		)
		if err != nil {
			writeCloser()

			return nil, fmt.Errorf("turtle: %v", err)
		}
	default:
		return nil, fmt.Errorf("unknown format: %s", f.Type)
	}

	return b, nil
}
