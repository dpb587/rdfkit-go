package sparql

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

// TODO stream/lex
type modelQueryResponseJSON struct {
	Head struct {
		Vars []string `json:"vars"`
		Link []string `json:"link"`
	} `json:"head"`
	Boolean *bool `json:"boolean"`
	Results *struct {
		Bindings []map[string]struct {
			Type     string `json:"type"`
			Value    string `json:"value"`
			Lang     string `json:"xml:lang"`
			Datatype string `json:"datatype"`
		} `json:"bindings"`
	} `json:"results"`
}

type QueryResponseDecoderJSON struct {
	BlankNodeTable blanknodeutil.StringMapper

	r io.Reader
}

func NewQueryResponseDecoderJSON(r io.Reader) *QueryResponseDecoderJSON {
	return &QueryResponseDecoderJSON{
		BlankNodeTable: blanknodeutil.NewStringMapper(),
		r:              r,
	}
}

func DecodeQueryResponseJSON(r io.Reader) (*QueryResponse, error) {
	return NewQueryResponseDecoderJSON(r).Decode()
}

func (d *QueryResponseDecoderJSON) Decode() (*QueryResponse, error) {
	var raw modelQueryResponseJSON

	err := json.NewDecoder(d.r).Decode(&raw)
	if err != nil {
		return nil, err
	}

	qr := &QueryResponse{}

	for _, v := range raw.Head.Vars {
		qr.Head.Variables = append(qr.Head.Variables, QueryResponseHeadVariable{
			Name: v,
		})
	}

	for _, l := range raw.Head.Link {
		qr.Head.Links = append(qr.Head.Links, QueryResponseHeadLink{
			Href: l,
		})
	}

	if raw.Results != nil {
		rawResults := QueryResponseResultList{}

		for _, b := range raw.Results.Bindings {
			qrr := QueryResponseResult{
				Bindings: QueryResponseResultBindingMap{},
			}

			for k, v := range b {
				b := QueryResponseResultBinding{
					Name: k,
				}

				switch v.Type {
				case "uri":
					b.Term = rdf.IRI(v.Value)
				case "literal":
					t := rdf.Literal{
						Datatype:    xsdiri.String_Datatype,
						LexicalForm: v.Value,
					}

					if len(v.Datatype) > 0 {
						t.Datatype = rdf.IRI(v.Datatype)
					} else if len(v.Lang) > 0 {
						t.Datatype = rdfiri.LangString_Datatype
						t.Tag = rdf.LanguageLiteralTag{
							Language: v.Lang,
						}
					}

					b.Term = t
				case "bnode":
					b.Term = d.BlankNodeTable.MapBlankNodeIdentifier(v.Value)
				default:
					return nil, errors.New("unknown binding type: " + v.Type)
				}

				qrr.Bindings[k] = b
			}

			rawResults = append(rawResults, qrr)
		}

		qr.Results = &rawResults
	} else if raw.Boolean != nil {
		qr.Boolean = raw.Boolean
	} else {
		return nil, errors.New("missing property: results or boolean")
	}

	return qr, nil
}
