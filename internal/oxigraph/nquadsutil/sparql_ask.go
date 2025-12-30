package nquadsutil

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

func NewSpaqlAsk(r io.Reader) ([]byte, error) {
	out := &bytes.Buffer{}

	bn := blanknodeutil.NewStringerInt64()
	rr, err := nquads.NewDecoder(r)
	if err != nil {
		return nil, err
	}

	bygraph := map[rdf.GraphNameValue][]rdf.Triple{}

	for rr.Next() {
		quad := rr.Quad()

		bygraph[quad.GraphName] = append(bygraph[quad.GraphName], quad.Triple)
	}

	if err := rr.Err(); err != nil {
		return nil, err
	}

	bnodes := map[string]struct{}{}

	for gn, triples := range bygraph {
		if gn != rdf.DefaultGraph {
			fmt.Fprintf(out, "\nGRAPH ")

			switch gn := gn.(type) {
			case rdf.IRI:
				_, err := nquads.WriteIRI(out, gn, false)
				if err != nil {
					return nil, err
				}
			case rdf.BlankNode:
				bnode := bn.GetBlankNodeIdentifier(gn)

				_, err := fmt.Fprintf(out, "?%s", bnode)
				if err != nil {
					return nil, err
				}

				bnodes[bnode] = struct{}{}
			default:
				return nil, fmt.Errorf("unsupported graph name type: %T", gn)
			}

			fmt.Fprintf(out, " {\n")
		}

		for _, t := range triples {
			fmt.Fprintf(out, "  ")

			switch s := t.Subject.(type) {
			case rdf.IRI:
				_, err := nquads.WriteIRI(out, s, false)
				if err != nil {
					return nil, err
				}
			case rdf.BlankNode:
				bnode := bn.GetBlankNodeIdentifier(s)

				_, err := fmt.Fprintf(out, "?%s", bnode)
				if err != nil {
					return nil, err
				}

				bnodes[bnode] = struct{}{}
			default:
				return nil, fmt.Errorf("unsupported subject type: %T", s)
			}

			fmt.Fprintf(out, " ")

			switch p := t.Predicate.(type) {
			case rdf.IRI:
				_, err := nquads.WriteIRI(out, p, false)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unsupported predicate type: %T", p)
			}

			fmt.Fprintf(out, " ")

			switch o := t.Object.(type) {
			case rdf.IRI:
				_, err := nquads.WriteIRI(out, o, false)
				if err != nil {
					return nil, err
				}
			case rdf.BlankNode:
				bnode := bn.GetBlankNodeIdentifier(o)

				_, err := fmt.Fprintf(out, "?%s", bnode)
				if err != nil {
					return nil, err
				}

				bnodes[bnode] = struct{}{}
			case rdf.Literal:
				_, err := nquads.WriteLiteral(out, o, false)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unsupported object type: %T", o)
			}

			fmt.Fprintf(out, " .\n")
		}

		if gn != rdf.DefaultGraph {
			fmt.Fprintf(out, "}\n")
		}

		for bnode := range bnodes {
			fmt.Fprintf(out, "FILTER isBlank(?%s)\n", bnode)
		}
	}

	// fmt.Fprintf(os.Stderr, "\n===\n%v\n===\n", out.String())

	return out.Bytes(), nil
}
