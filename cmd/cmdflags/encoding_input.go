package cmdflags

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/htmljsonld"
	"github.com/dpb587/rdfkit-go/encoding/htmlmicrodata"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/encoding/rdfa"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml"
	"github.com/dpb587/rdfkit-go/encoding/trig"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/internal/ioutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

type EncodingInput struct {
	Type string
	Path string

	FallbackOpener RemoteOpenerFunc

	DefaultBase string

	DocumentLoaderJSONLD jsonldtype.DocumentLoader
}

func (f EncodingInput) openReader() (string, io.ReadCloser, http.Header, error) {
	if f.Path == "-" {
		return "file:///dev/stdin", io.NopCloser(os.Stdin), nil, nil
	}

	inFile, err := os.OpenFile(f.Path, os.O_RDONLY, 0)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && f.FallbackOpener != nil {
			f1, f2, f3 := f.FallbackOpener(f.Path, err)

			return f.Path, f1, f2, f3
		}

		return "", nil, nil, fmt.Errorf("open: %v", err)
	}

	return f.Path, inFile, nil, nil
}

func (f EncodingInput) Open() (*EncodingInputHandle, error) {
	return f.openTee(nil)
}

func (f EncodingInput) OpenTee(w io.Writer) (*EncodingInputHandle, error) {
	return f.openTee(w)
}

func (f EncodingInput) openTee(w io.Writer) (*EncodingInputHandle, error) {
	rcPath, rc, remoteHeader, err := f.openReader()
	if err != nil {
		return nil, err
	}

	readHasher := sha256.New()
	readCloser := rc.Close
	readReader := io.TeeReader(rc, readHasher)

	if w != nil {
		readReader = io.TeeReader(readReader, w)
	}

	fType := f.Type

	if len(fType) == 0 && remoteHeader != nil {
		switch strings.SplitN(remoteHeader.Get("Content-Type"), ";", 2)[0] {
		case "application/json":
			fType = "jsonld"
		case "application/ld+json":
			fType = "jsonld"
		case "application/n-quads":
			fType = "nquads"
		case "application/n-triples":
			fType = "ntriples"
		case "application/rdf+json":
			fType = "rdfjson"
		case "application/rdf+xml":
			fType = "rdfxml"
		case "application/trig":
			fType = "trig"
		case "text/html", "text/xhtml+xml", "application/xhtml+xml":
			fType = "html"
		case "text/turtle":
			fType = "turtle"
		}
	}

	if len(fType) == 0 {
		fType = "trig" // most generic default

		var magicRead = make([]byte, 4096)

		magicReadN, err := io.ReadFull(readReader, magicRead)
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
				return nil, err
			}
		}

		magicBytes := magicRead[:magicReadN]

		if regexp.MustCompile(`^\s*\{`).Match(magicBytes) {
			if regexp.MustCompile(`[^\\]"@[a-z]+"\s*:`).Match(magicBytes) {
				fType = "jsonld"
			} else if regexp.MustCompile(`^\s*\{\s*"[^"]+"\s*:\s*\{\s*"[^"]+"\s*:\s*\[\s*\{\s*"(datatype|lang|type|value)"`).Match(magicBytes) {
				fType = "rdfjson"
			} else {
				fType = "jsonld"
			}
		} else if regexp.MustCompile(`^(<[^>]+>|\s)*<html[\s>]`).Match(magicBytes) {
			fType = "html"
		} else if regexp.MustCompile(`^(<[^>]+>|\s)*<\?xml`).Match(magicBytes) {
			fType = "rdfxml"
		} else if regexp.MustCompile(`<rdf:RDF `).Match(magicBytes) {
			fType = "rdfxml"
		}

		readReader = io.MultiReader(bytes.NewReader(magicBytes), readReader)
	}

	handle := &EncodingInputHandle{
		ReadPath:   rcPath,
		readHasher: readHasher,
		reader:     ioutil.NewReaderCloser(readReader, readCloser),
	}

	if len(f.DefaultBase) == 0 {
		if len(rcPath) > 0 {
			f.DefaultBase = fmt.Sprintf("file://%s", rcPath)
		}
	}

	parseHtmlDocument := func() (*html.Document, error) {
		htmlOptions := html.DocumentConfig{}.SetCaptureTextOffsets(true)

		if len(f.DefaultBase) > 0 {
			htmlOptions = htmlOptions.SetLocation(f.DefaultBase)
		}

		htmlDocument, err := html.ParseDocument(handle.reader, htmlOptions)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("html: %v", err)
		}

		return htmlDocument, nil
	}

	switch fType {
	case "html":
		htmlDocument, err := parseHtmlDocument()
		if err != nil {
			return nil, err
		}

		htmlJsonld, err := htmljsonld.NewDecoder(
			htmlDocument,
			htmljsonld.DecoderConfig{}.
				SetDocumentLoader(f.DocumentLoaderJSONLD),
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("htmljsonld: %v", err)
		}

		htmlMicrodata, err := htmlmicrodata.NewDecoder(
			htmlDocument,
			htmlmicrodata.DecoderConfig{}.
				SetVocabularyResolver(htmlmicrodata.ItemtypeVocabularyResolver),
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("htmlmicrodata: %v", err)
		}

		htmlRdfa, err := rdfa.NewDecoder(htmlDocument)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("rdfa: %v", err)
		}

		handle.Format = "html"
		handle.Decoder = rdfioutil.NewStatementIteratorIterator(htmlJsonld, htmlMicrodata, htmlRdfa)
	case "htmljsonld":
		htmlDocument, err := parseHtmlDocument()
		if err != nil {
			return nil, err
		}

		handle.Format = "htmljsonld"
		handle.Decoder, err = htmljsonld.NewDecoder(
			htmlDocument,
			htmljsonld.DecoderConfig{}.SetDocumentLoader(f.DocumentLoaderJSONLD),
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("htmljsonld: %v", err)
		}
	case "htmlmicrodata":
		htmlDocument, err := parseHtmlDocument()
		if err != nil {
			return nil, err
		}

		handle.Format = "htmlmicrodata"
		handle.Decoder, err = htmlmicrodata.NewDecoder(
			htmlDocument,
			htmlmicrodata.DecoderConfig{}.
				SetVocabularyResolver(htmlmicrodata.ItemtypeVocabularyResolver),
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("htmlmicrodata: %v", err)
		}
	case "jsonld", "json-ld":
		jsonldOptions := jsonld.DecoderConfig{}.
			SetCaptureTextOffsets(true).
			SetParserOptions(inspectjson.TokenizerOptions{}.Lax(true)).
			SetDocumentLoader(jsonldtype.NewDefaultDocumentLoader(http.DefaultClient))

		if len(f.DefaultBase) > 0 {
			jsonldOptions = jsonldOptions.SetDefaultBase(f.DefaultBase)
		}

		handle.Format = "jsonld"
		handle.Decoder, err = jsonld.NewDecoder(
			handle.reader,
			jsonldOptions,
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("jsonld: %v", err)
		}
	case "nquads", "nq":
		handle.Format = "nquads"
		handle.Decoder, err = nquads.NewDecoder(
			handle.reader,
			nquads.DecoderConfig{}.
				SetCaptureTextOffsets(true),
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("nquads: %v", err)
		}
	case "ntriples", "nt":
		handle.Format = "ntriples"
		handle.Decoder, err = ntriples.NewDecoder(
			handle.reader,
			ntriples.DecoderConfig{}.
				SetCaptureTextOffsets(true),
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("ntriples: %v", err)
		}
	case "rdfa":
		htmlDocument, err := parseHtmlDocument()
		if err != nil {
			return nil, err
		}

		handle.Format = "rdfa"
		handle.Decoder, err = rdfa.NewDecoder(htmlDocument)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("rdfa: %v", err)
		}
	case "rdfjson", "rj":
		handle.Format = "rdfjson"
		handle.Decoder, err = rdfjson.NewDecoder(
			handle.reader,
			rdfjson.DecoderConfig{}.
				SetCaptureTextOffsets(true),
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("rdfjson: %v", err)
		}
	case "rdfxml":
		rdfxmlOptions := rdfxml.DecoderConfig{}.
			SetCaptureTextOffsets(true)

		if len(f.DefaultBase) > 0 {
			rdfxmlOptions = rdfxmlOptions.SetBaseURL(f.DefaultBase)
		}

		handle.Format = "rdfxml"
		handle.Decoder, err = rdfxml.NewDecoder(
			handle.reader,
			rdfxmlOptions,
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("rdfxml: %v", err)
		}
	case "trig":
		trigOptions := trig.DecoderConfig{}.
			SetCaptureTextOffsets(true).
			SetBaseDirectiveListener(func(data trig.DecoderEvent_BaseDirective_Data) {
				handle.DecodedBase = append(handle.DecodedBase, data.Value)
			}).
			SetPrefixDirectiveListener(func(data trig.DecoderEvent_PrefixDirective_Data) {
				handle.DecodedPrefixMappings = append(handle.DecodedPrefixMappings, iriutil.PrefixMapping{
					Prefix:   data.Prefix,
					Expanded: data.Expanded,
				})
			})

		if len(f.DefaultBase) > 0 {
			trigOptions = trigOptions.SetDefaultBase(f.DefaultBase)
		}

		handle.Format = "trig"
		handle.Decoder, err = trig.NewDecoder(
			handle.reader,
			trigOptions,
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("trig: %v", err)
		}
	case "turtle", "ttl":
		turtleOptions := turtle.DecoderConfig{}.
			SetCaptureTextOffsets(true).
			SetBaseDirectiveListener(func(data turtle.DecoderEvent_BaseDirective_Data) {
				handle.DecodedBase = append(handle.DecodedBase, data.Value)
			}).
			SetPrefixDirectiveListener(func(data turtle.DecoderEvent_PrefixDirective_Data) {
				handle.DecodedPrefixMappings = append(handle.DecodedPrefixMappings, iriutil.PrefixMapping{
					Prefix:   data.Prefix,
					Expanded: data.Expanded,
				})
			})

		if len(f.DefaultBase) > 0 {
			turtleOptions = turtleOptions.SetDefaultBase(f.DefaultBase)
		}

		handle.Format = "turtle"
		handle.Decoder, err = turtle.NewDecoder(
			handle.reader,
			turtleOptions,
		)
		if err != nil {
			readCloser()

			return nil, fmt.Errorf("turtle: %v", err)
		}
	default:
		rc.Close()

		return nil, fmt.Errorf("unknown format: %s", f.Type)
	}

	return handle, nil
}
