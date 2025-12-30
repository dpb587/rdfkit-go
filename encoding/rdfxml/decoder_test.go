package rdfxml

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfobject"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

var testBnode = blanknodeutil.NewStringMapper()

func assertEqual(t *testing.T, expectedList, actualList rdfio.StatementList) {
	bnEtoA := map[rdf.BlankNode]rdf.BlankNode{}
	bnAtoE := map[rdf.BlankNode]rdf.BlankNode{}

	for idx, expected := range expectedList {
		if len(actualList) <= idx {
			t.Fatalf("expected %dth triple, but got %d", idx, len(actualList))
		}

		actual := actualList[idx]

		expectedTriple := expected.GetTriple()
		actualTriple := actual.GetTriple()

		switch expectedTripleSubject := expectedTriple.Subject.(type) {
		case rdf.BlankNode:
			actualTripleSubject, ok := actualTriple.Subject.(rdf.BlankNode)
			if !ok {
				t.Fatalf("expected subject to be blank node, but got %T", actualTriple.Subject)
			}

			if bn, ok := bnEtoA[expectedTripleSubject]; ok {
				if bn != actualTripleSubject {
					t.Fatalf("expected subject to be %v, but got %v", bn, actualTripleSubject)
				}
			} else if _, ok := bnAtoE[actualTripleSubject]; ok {
				t.Fatalf("expected subject to be %v, but got %v", bn, actualTripleSubject)
			} else {
				bnEtoA[expectedTripleSubject] = actualTripleSubject
				bnAtoE[actualTripleSubject] = expectedTripleSubject
			}
		case rdf.IRI:
			if !expectedTripleSubject.TermEquals(actualTriple.Subject) {
				t.Fatalf("expected subject to be %v, but got %v", expectedTripleSubject, actualTriple.Subject)
			}
		default:
			t.Fatalf("unexpected type: %T", expectedTripleSubject)
		}

		switch expectedTriplePredicate := expectedTriple.Predicate.(type) {
		case rdf.IRI:
			if !expectedTriplePredicate.TermEquals(actualTriple.Predicate) {
				t.Fatalf("expected predicate to be %v, but got %v", expectedTriplePredicate, actualTriple.Predicate)
			}
		default:
			t.Fatalf("unexpected type: %T", expectedTriplePredicate)
		}

		switch expectedTripleObject := expectedTriple.Object.(type) {
		case rdf.BlankNode:
			actualTripleObject, ok := actualTriple.Object.(rdf.BlankNode)
			if !ok {
				t.Fatalf("expected subject to be blank node, but got %T", actualTriple.Object)
			}

			if bn, ok := bnEtoA[expectedTripleObject]; ok {
				if bn != actualTripleObject {
					t.Fatalf("expected subject to be %v, but got %v", bn, actualTripleObject)
				}
			} else if _, ok := bnAtoE[actualTripleObject]; ok {
				t.Fatalf("expected subject to be %v, but got %v", bn, actualTripleObject)
			} else {
				bnEtoA[expectedTripleObject] = actualTripleObject
				bnAtoE[actualTripleObject] = expectedTripleObject
			}
		case rdf.IRI:
			if !expectedTripleObject.TermEquals(actualTriple.Object) {
				t.Fatalf("expected object to be %v, but got %v", expectedTripleObject, actualTriple.Object)
			}
		case rdf.Literal:
			if !expectedTripleObject.TermEquals(actualTriple.Object) {
				t.Fatalf("expected object to be %v, but got %v", expectedTripleObject, actualTriple.Object)
			}
		default:
			t.Fatalf("unexpected type: %T", expectedTripleObject)
		}

		// expectedLocation := expected.(encoding.TripleLocationBinding).GetEncodingLocation()
		// actualLocation := actual.(encoding.TripleLocationBinding).GetEncodingLocation()

		// if expectedLocation.Object.LookupString() != actualLocation.Object.LookupString() {
		// 	t.Fatalf("expected location to be %s, but got %#+v", expectedLocation.Object.LookupString(), actualLocation.Object)
		// }
	}
}

func patchExampleMissingXmlns(s string) string {
	return strings.Replace(s, `<rdf:Description`, `<rdf:Description xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:ex="http://example.org/stuff/1.0/" `, 1)
}

func TestSpecNonNormative(t *testing.T) {
	for _, testcase := range []struct {
		Name     string
		Snippet  string
		Expected rdfdescription.ResourceList
	}{
		{
			Name: "2.2/3",
			Snippet: strings.Join([]string{
				// parser does not support multiple elements; wrap them for testing
				`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:ex="http://example.org/stuff/1.0/">`,
				// literal example
				`<rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar">
  <ex:editor>
    <rdf:Description>
      <ex:homePage>
        <rdf:Description rdf:about="http://purl.org/net/dajobe/">
        </rdf:Description>
      </ex:homePage>
    </rdf:Description>
  </ex:editor>
</rdf:Description>

<rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar">
  <ex:editor>
    <rdf:Description>
      <ex:fullName>Dave Beckett</ex:fullName>
    </rdf:Description>
  </ex:editor>
</rdf:Description>

<rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar">
  <dc:title>RDF 1.1 XML Syntax</dc:title>
</rdf:Description>`,
				`</rdf:RDF>`,
			}, ""),
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.AnonResourceStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							AnonResource: rdfdescription.AnonResource{
								Statements: rdfdescription.StatementList{
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/homePage"),
										Object:    rdf.IRI("http://purl.org/net/dajobe/"),
									},
								},
							},
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.AnonResourceStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							AnonResource: rdfdescription.AnonResource{
								Statements: rdfdescription.StatementList{
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
										Object:    xsdobject.String("Dave Beckett"),
									},
								},
							},
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    xsdobject.String("RDF 1.1 XML Syntax"),
						},
					},
				},
			},
		},
		{
			Name: "2.3/4",
			Snippet: patchExampleMissingXmlns(`<rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar">
  <ex:editor>
    <rdf:Description>
      <ex:homePage>
        <rdf:Description rdf:about="http://purl.org/net/dajobe/">
        </rdf:Description>
      </ex:homePage>
      <ex:fullName>Dave Beckett</ex:fullName>
    </rdf:Description>
  </ex:editor>
  <dc:title>RDF 1.1 XML Syntax</dc:title>
</rdf:Description>`),
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.AnonResourceStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							AnonResource: rdfdescription.AnonResource{
								Statements: rdfdescription.StatementList{
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/homePage"),
										Object:    rdf.IRI("http://purl.org/net/dajobe/"),
									},
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
										Object:    xsdobject.String("Dave Beckett"),
									},
								},
							},
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    xsdobject.String("RDF 1.1 XML Syntax"),
						},
					},
				},
			},
		},
		{
			Name: "2.4/5",
			Snippet: patchExampleMissingXmlns(`<rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar">
  <ex:editor>
    <rdf:Description>
      <ex:homePage rdf:resource="http://purl.org/net/dajobe/"/>
      <ex:fullName>Dave Beckett</ex:fullName>
    </rdf:Description>
  </ex:editor>
  <dc:title>RDF 1.1 XML Syntax</dc:title>
</rdf:Description>`),
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.AnonResourceStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							AnonResource: rdfdescription.AnonResource{
								Statements: rdfdescription.StatementList{
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/homePage"),
										Object:    rdf.IRI("http://purl.org/net/dajobe/"),
									},
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
										Object:    xsdobject.String("Dave Beckett"),
									},
								},
							},
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    xsdobject.String("RDF 1.1 XML Syntax"),
						},
					},
				},
			},
		},
		{
			Name: "2.5/6",
			Snippet: patchExampleMissingXmlns(`<rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar"
           dc:title="RDF 1.1 XML Syntax">
  <ex:editor>
    <rdf:Description ex:fullName="Dave Beckett">
      <ex:homePage rdf:resource="http://purl.org/net/dajobe/"/>
    </rdf:Description>
  </ex:editor>
</rdf:Description>`),
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    xsdobject.String("RDF 1.1 XML Syntax"),
						},
						rdfdescription.AnonResourceStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							AnonResource: rdfdescription.AnonResource{
								Statements: rdfdescription.StatementList{
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
										Object:    xsdobject.String("Dave Beckett"),
									},
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/homePage"),
										Object:    rdf.IRI("http://purl.org/net/dajobe/"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "2.6/7",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:ex="http://example.org/stuff/1.0/">

  <rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar"
             dc:title="RDF1.1 XML Syntax">
    <ex:editor>
      <rdf:Description ex:fullName="Dave Beckett">
        <ex:homePage rdf:resource="http://purl.org/net/dajobe/" />
      </rdf:Description>
    </ex:editor>
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							// https://www.w3.org/TR/2004/REC-rdf-syntax-grammar-20040210/example07.nt has incorrect value
							Object: xsdobject.String("RDF1.1 XML Syntax"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b0"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
							Object:    xsdobject.String("Dave Beckett"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/homePage"),
							Object:    rdf.IRI("http://purl.org/net/dajobe/"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							Object:    testBnode.MapBlankNodeIdentifier("b0"),
						},
					},
				},
			},
		},
		{
			Name: "2.7/8",
			Snippet: `<?xml version="1.0" encoding="utf-8"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:dc="http://purl.org/dc/elements/1.1/">

  <rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar">
    <dc:title>RDF 1.1 XML Syntax</dc:title>
    <dc:title xml:lang="en">RDF 1.1 XML Syntax</dc:title>
    <dc:title xml:lang="en-US">RDF 1.1 XML Syntax</dc:title>
  </rdf:Description>

  <rdf:Description rdf:about="http://example.org/buecher/baum" xml:lang="de">
    <dc:title>Der Baum</dc:title>
    <dc:description>Das Buch ist außergewöhnlich</dc:description>
    <dc:title xml:lang="en">The Tree</dc:title>
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    xsdobject.String("RDF 1.1 XML Syntax"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    rdfobject.LangString("en", "RDF 1.1 XML Syntax"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    rdfobject.LangString("en-US", "RDF 1.1 XML Syntax"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/buecher/baum"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    rdfobject.LangString("de", "Der Baum"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/description"),
							Object:    rdfobject.LangString("de", "Das Buch ist außergewöhnlich"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    rdfobject.LangString("en", "The Tree"),
						},
					},
				},
			},
		},
		{
			Name: "2.8/9",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:ex="http://example.org/stuff/1.0/">

  <rdf:Description rdf:about="http://example.org/item01"> 
    <ex:prop rdf:parseType="Literal" xmlns:a="http://example.org/a#">
      <a:Box required="true">
        <a:widget size="10" />
        <a:grommit id="23" />
      </a:Box>
    </ex:prop>
  </rdf:Description>

</rdf:RDF>`,
			// nt file uses different xml marshal conventions; semantically equivalent?
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/item01"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/prop"),
							// TODO prefer original xmlns map
							// TODO prefer original self-closing?
							Object: rdf.Literal{
								Datatype: "http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral",
								LexicalForm: `
      <Box xmlns="http://example.org/a#" required="true">
        <widget xmlns="http://example.org/a#" size="10"></widget>
        <grommit xmlns="http://example.org/a#" id="23"></grommit>
      </Box>
    `,
							},
						},
					},
				},
			},
		},
		{
			Name: "2.9/10",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:ex="http://example.org/stuff/1.0/">

  <rdf:Description rdf:about="http://example.org/item01">
    <ex:size rdf:datatype="http://www.w3.org/2001/XMLSchema#int">123</ex:size>
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/item01"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/size"),
							Object: rdf.Literal{
								Datatype:    "http://www.w3.org/2001/XMLSchema#int",
								LexicalForm: "123",
							},
						},
					},
				},
			},
		},
		{
			Name: "2.10/11",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:ex="http://example.org/stuff/1.0/">

  <rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar"
             dc:title="RDF 1.1 XML Syntax">
    <ex:editor rdf:nodeID="abc"/>
  </rdf:Description>

  <rdf:Description rdf:nodeID="abc" ex:fullName="Dave Beckett">
    <ex:homePage rdf:resource="http://purl.org/net/dajobe/"/>
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							// example11.nt has incorrect value
							Object: xsdobject.String("RDF 1.1 XML Syntax"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							Object:    testBnode.MapBlankNodeIdentifier("b0"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b0"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
							Object:    xsdobject.String("Dave Beckett"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/homePage"),
							Object:    rdf.IRI("http://purl.org/net/dajobe/"),
						},
					},
				},
			},
		},
		{
			Name: "2.11/12",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:ex="http://example.org/stuff/1.0/">
  <rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar"
                   dc:title="RDF 1.1 XML Syntax">
    <ex:editor rdf:parseType="Resource">
      <ex:fullName>Dave Beckett</ex:fullName>
      <ex:homePage rdf:resource="http://purl.org/net/dajobe/"/>
    </ex:editor>
  </rdf:Description>
</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							// example12.nt has incorrect value
							Object: xsdobject.String("RDF 1.1 XML Syntax"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							Object:    testBnode.MapBlankNodeIdentifier("b0"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b0"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
							Object:    xsdobject.String("Dave Beckett"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/homePage"),
							Object:    rdf.IRI("http://purl.org/net/dajobe/"),
						},
					},
				},
			},
		},
		{
			Name: "2.12/13",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:ex="http://example.org/stuff/1.0/">

  <rdf:Description rdf:about="http://www.w3.org/TR/rdf-syntax-grammar"
            dc:title="RDF 1.1 XML Syntax">
    <ex:editor ex:fullName="Dave Beckett" />
            <!-- Note the ex:homePage property has been ignored for this example -->
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/TR/rdf-syntax-grammar"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							// example13.nt has incorrect value
							Object: xsdobject.String("RDF 1.1 XML Syntax"),
						},
						rdfdescription.AnonResourceStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/editor"),
							AnonResource: rdfdescription.AnonResource{
								Statements: rdfdescription.StatementList{
									rdfdescription.ObjectStatement{
										Predicate: rdf.IRI("http://example.org/stuff/1.0/fullName"),
										Object:    xsdobject.String("Dave Beckett"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "2.13/14",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:ex="http://example.org/stuff/1.0/">

  <rdf:Description rdf:about="http://example.org/thing">
    <rdf:type rdf:resource="http://example.org/stuff/1.0/Document"/>
    <dc:title>A marvelous thing</dc:title>
  </rdf:Description>
</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/thing"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
							Object:    rdf.IRI("http://example.org/stuff/1.0/Document"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    xsdobject.String("A marvelous thing"),
						},
					},
				},
			},
		},
		{
			Name: "2.13/15",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:ex="http://example.org/stuff/1.0/">

  <ex:Document rdf:about="http://example.org/thing">
    <dc:title>A marvelous thing</dc:title>
  </ex:Document>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/thing"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
							Object:    rdf.IRI("http://example.org/stuff/1.0/Document"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
							Object:    xsdobject.String("A marvelous thing"),
						},
					},
				},
			},
		},
		{
			Name: "2.14/16",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:ex="http://example.org/stuff/1.0/"
            xml:base="http://example.org/here/">

  <rdf:Description rdf:ID="snack">
    <ex:prop rdf:resource="fruit/apple"/>
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/here/#snack"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/prop"),
							Object:    rdf.IRI("http://example.org/here/fruit/apple"),
						},
					},
				},
			},
		},
		{
			Name: "2.15/17",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">

  <rdf:Seq rdf:about="http://example.org/favourite-fruit">
    <rdf:_1 rdf:resource="http://example.org/banana"/>
    <rdf:_2 rdf:resource="http://example.org/apple"/>
    <rdf:_3 rdf:resource="http://example.org/pear"/>
  </rdf:Seq>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/favourite-fruit"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
							Object:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#Seq"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#_1"),
							Object:    rdf.IRI("http://example.org/banana"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#_2"),
							Object:    rdf.IRI("http://example.org/apple"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#_3"),
							Object:    rdf.IRI("http://example.org/pear"),
						},
					},
				},
			},
		},
		{
			Name: "2.15/18",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">

  <rdf:Seq rdf:about="http://example.org/favourite-fruit">
    <rdf:li rdf:resource="http://example.org/banana"/>
    <rdf:li rdf:resource="http://example.org/apple"/>
    <rdf:li rdf:resource="http://example.org/pear"/>
  </rdf:Seq>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/favourite-fruit"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
							Object:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#Seq"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#_1"),
							Object:    rdf.IRI("http://example.org/banana"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#_2"),
							Object:    rdf.IRI("http://example.org/apple"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#_3"),
							Object:    rdf.IRI("http://example.org/pear"),
						},
					},
				},
			},
		},
		{
			Name: "2.16/19",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:ex="http://example.org/stuff/1.0/">

  <rdf:Description rdf:about="http://example.org/basket">
    <ex:hasFruit rdf:parseType="Collection">
      <rdf:Description rdf:about="http://example.org/banana"/>
      <rdf:Description rdf:about="http://example.org/apple"/>
      <rdf:Description rdf:about="http://example.org/pear"/>
    </ex:hasFruit>
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/basket"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/hasFruit"),
							Object:    testBnode.MapBlankNodeIdentifier("b0"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b0"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
							Object:    rdf.IRI("http://example.org/banana"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
							Object:    testBnode.MapBlankNodeIdentifier("b1"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b1"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
							Object:    rdf.IRI("http://example.org/apple"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
							Object:    testBnode.MapBlankNodeIdentifier("b2"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b2"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
							Object:    rdf.IRI("http://example.org/pear"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
							Object:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
						},
					},
				},
			},
		},
		{
			Name: "2.17/20",
			Snippet: `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
            xmlns:ex="http://example.org/stuff/1.0/"
            xml:base="http://example.org/triples/">
  <rdf:Description rdf:about="http://example.org/">
    <ex:prop rdf:ID="triple1">blah</ex:prop>
  </rdf:Description>

</rdf:RDF>`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/stuff/1.0/prop"),
							Object:    xsdobject.String("blah"),
							// BindingLocation: encoding.TripleLocation{
							// 	Object: &cursorutil.CursorRange{
							// 		From: cursorutil.CursorPosition{
							// 			Byte: 269,
							// 			Text: cursorutil.TextLineColumn{5, 30},
							// 		},
							// 		Until: cursorutil.CursorPosition{
							// 			Byte: 273,
							// 			Text: cursorutil.TextLineColumn{5, 34},
							// 		},
							// 	},
							// },
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/triples/#triple1"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
							Object:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#Statement"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#subject"),
							Object:    rdf.IRI("http://example.org/"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#predicate"),
							Object:    rdf.IRI("http://example.org/stuff/1.0/prop"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#object"),
							Object:    xsdobject.String("blah"),
						},
					},
				},
			},
		},
	} {
		t.Run(testcase.Name, func(t *testing.T) {
			// if testcase.Name != "2.7/8" {
			// 	return
			// }

			out, err := rdfio.CollectStatementsErr(NewDecoder(
				bytes.NewBufferString(testcase.Snippet),
				DecoderConfig{}.
					SetCaptureTextOffsets(true),
			))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := rdfio.StatementList{}

			for _, e := range testcase.Expected.AsTriples() {
				expected = append(expected, rdfioutil.Statement{
					Triple: e,
				})
			}

			var success bool

			t.Cleanup(func() {
				if success {
					return
				}

				fmt.Fprintf(os.Stderr, "### ACTUAL\n")

				for _, stmt := range out {
					fmt.Fprintf(os.Stderr, "%s\n", stmt.GetTriple())
				}

				fmt.Fprintf(os.Stderr, "### EXPECTED\n")

				for _, stmt := range expected {
					fmt.Fprintf(os.Stderr, "%s\n", stmt.GetTriple())
				}
			})

			assertEqual(t, expected, out)

			success = true
		})
	}
}

func TestSpecTestcase(t *testing.T) {
	for _, testcase := range []struct {
		Name          string
		OptionBaseURL string
		Snippet       string
		Expected      rdfdescription.ResourceList
	}{
		{
			Name: "test001.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: xml:base applies to an rdf:ID on an 
               rdf:Description element.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test001.rdf,v 1.1 2014/02/20 20:36:30 sandro Exp $

-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         xml:base="http://example.org/dir/file">

 <rdf:Description rdf:ID="frag" eg:value="v" />

</rdf:RDF>
`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/dir/file#frag"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/value"),
							Object:    xsdobject.String("v"),
						},
					},
				},
			},
		},
		{
			Name: "test002.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: xml:base applies to an rdf:resource attribute.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test002.rdf,v 1.1 2014/02/20 20:36:30 sandro Exp $

-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         xml:base="http://example.org/dir/file">

 <rdf:Description>
   <eg:value rdf:resource="relFile" />
 </rdf:Description>

</rdf:RDF>
`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b0"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/value"),
							Object:    rdf.IRI("http://example.org/dir/relFile"),
						},
					},
				},
			},
		},
		{
			Name: "test004.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: xml:base applies to an rdf:ID on a property element.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test004.rdf,v 1.1 2014-02-20 20:36:31 sandro Exp $
-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         xml:base="http://example.org/dir/file">

 <rdf:Description>
  <eg:value rdf:ID="frag">v</eg:value>
 </rdf:Description>

</rdf:RDF>
`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: testBnode.MapBlankNodeIdentifier("b0"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/value"),
							Object:    xsdobject.String("v"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/dir/file#frag"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Type_Property,
							Object:    rdfiri.Statement_Class,
						},
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Subject_Property,
							Object:    testBnode.MapBlankNodeIdentifier("b0"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Predicate_Property,
							Object:    rdf.IRI("http://example.org/value"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Object_Property,
							Object:    xsdobject.String("v"),
						},
					},
				},
			},
		},
		{
			Name: "test008.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: example of empty same document ref resolution.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test008.rdf,v 1.1 2014-02-20 20:36:31 sandro Exp $
-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         xml:base="http://example.org/dir/file">

 <eg:type rdf:about="" />

</rdf:RDF>
`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/dir/file"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Type_Property,
							Object:    rdf.IRI("http://example.org/type"),
						},
					},
				},
			},
		},
		{
			Name: "test009.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: Example of relative uri with absolute path resolution.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test009.rdf,v 1.1 2014-02-20 20:36:31 sandro Exp $
-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         xml:base="http://example.org/dir/file">

 <eg:type rdf:about="/absfile" />

</rdf:RDF>

`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/absfile"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Type_Property,
							Object:    rdf.IRI("http://example.org/type"),
						},
					},
				},
			},
		},
		{
			Name: "test011.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: Example of xml:base with no path component.
  Note: The algorithm in RFC 2396 does not handle this case correctly.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test011.rdf,v 1.1 2014-02-20 20:36:32 sandro Exp $
-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         xml:base="http://example.org">

 <eg:type rdf:about="relfile" />

</rdf:RDF>

`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/relfile"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Type_Property,
							Object:    rdf.IRI("http://example.org/type"),
						},
					},
				},
			},
		},
		{
			Name: "test013.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: With an xml:base with fragment the fragment is ignored.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test013.rdf,v 1.1 2014-02-20 20:36:32 sandro Exp $
-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         xml:base="http://example.org/dir/file#frag">

 <eg:type rdf:about="" />
 <rdf:Description rdf:ID="foo" >
   <eg:value rdf:resource="relpath" />
 </rdf:Description>

</rdf:RDF>
`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/dir/file"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdfiri.Type_Property,
							Object:    rdf.IRI("http://example.org/type"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/dir/file#foo"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/value"),
							Object:    rdf.IRI("http://example.org/dir/relpath"),
						},
					},
				},
			},
		},
		{
			Name:          "test014.rdf",
			OptionBaseURL: "http://www.w3.org/2013/RDFXMLTests/xmlbase/test014.rdf",
			Snippet: `<?xml version="1.0"?>

<!--
  Copyright World Wide Web Consortium, (Massachusetts Institute of
  Technology, Institut National de Recherche en Informatique et en
  Automatique, Keio University).
 
  All Rights Reserved.
 
  Please see the full Copyright clause at
  <http://www.w3.org/Consortium/Legal/copyright-software.html>

  Description: two identical rdf:ID's are allowed, as long as they
               refer to different resources.
  Author: Jeremy Carroll (jjc@hpl.hp.com)

  $Id: test014.rdf,v 1.1 2014-02-20 20:36:32 sandro Exp $
-->

<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:eg="http://example.org/"
         >

 <rdf:Description xml:base="http://example.org/dir/file"
                rdf:ID="frag" eg:value="v" />
 <rdf:Description rdf:ID="frag" eg:value="v" />

</rdf:RDF>
`,
			Expected: rdfdescription.ResourceList{
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://example.org/dir/file#frag"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/value"),
							Object:    xsdobject.String("v"),
						},
					},
				},
				rdfdescription.SubjectResource{
					Subject: rdf.IRI("http://www.w3.org/2013/RDFXMLTests/xmlbase/test014.rdf#frag"),
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://example.org/value"),
							Object:    xsdobject.String("v"),
						},
					},
				},
			},
		},
	} {
		t.Run(testcase.Name, func(t *testing.T) {
			dopt := DecoderConfig{}

			if len(testcase.OptionBaseURL) > 0 {
				dopt = dopt.SetBaseURL(testcase.OptionBaseURL)
			}

			out, err := rdfio.CollectStatementsErr(NewDecoder(
				bytes.NewBufferString(testcase.Snippet),
				dopt,
			))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := rdfio.StatementList{}

			for _, e := range testcase.Expected.AsTriples() {
				expected = append(expected, rdfioutil.Statement{
					Triple: e,
				})
			}

			assertEqual(t, expected, out)
		})
	}
}
