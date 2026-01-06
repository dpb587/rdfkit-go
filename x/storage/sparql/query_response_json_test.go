package sparql

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/internal/ptr"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

// https://www.w3.org/TR/2013/REC-sparql11-results-json-20130321/
func TestDecodeQueryResponseJSON_SpecNonNormative(t *testing.T) {
	for _, tc := range []struct {
		Name   string
		Input  string
		Output *QueryResponse
		Error  string
	}{
		{
			Name: "2/1",
			Input: `{
  "head": { "vars": [ "book" , "title" ]
  } ,
  "results": { 
    "bindings": [
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book6" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Half-Blood Prince" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book7" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Deathly Hallows" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book5" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Order of the Phoenix" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book4" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Goblet of Fire" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book2" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Chamber of Secrets" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book3" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Prisoner Of Azkaban" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book1" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Philosopher's Stone" }
      }
    ]
  }
}`,
			Output: &QueryResponse{
				Head: QueryResponseHead{
					Variables: QueryResponseHeadVariableList{
						{
							Name: "book",
						},
						{
							Name: "title",
						},
					},
				},
				Results: &QueryResponseResultList{
					{
						Bindings: QueryResponseResultBindingMap{
							"book": QueryResponseResultBinding{
								Name: "book",
								Term: rdf.IRI("http://example.org/book/book6"),
							},
							"title": QueryResponseResultBinding{
								Name: "title",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Harry Potter and the Half-Blood Prince",
								},
							},
						},
					},
					{
						Bindings: QueryResponseResultBindingMap{
							"book": QueryResponseResultBinding{
								Name: "book",
								Term: rdf.IRI("http://example.org/book/book7"),
							},
							"title": QueryResponseResultBinding{
								Name: "title",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Harry Potter and the Deathly Hallows",
								},
							},
						},
					},
					{
						Bindings: QueryResponseResultBindingMap{
							"book": QueryResponseResultBinding{
								Name: "book",
								Term: rdf.IRI("http://example.org/book/book5"),
							},
							"title": QueryResponseResultBinding{
								Name: "title",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Harry Potter and the Order of the Phoenix",
								},
							},
						},
					},
					{
						Bindings: QueryResponseResultBindingMap{
							"book": QueryResponseResultBinding{
								Name: "book",
								Term: rdf.IRI("http://example.org/book/book4"),
							},
							"title": QueryResponseResultBinding{
								Name: "title",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Harry Potter and the Goblet of Fire",
								},
							},
						},
					},
					{
						Bindings: QueryResponseResultBindingMap{
							"book": QueryResponseResultBinding{
								Name: "book",
								Term: rdf.IRI("http://example.org/book/book2"),
							},
							"title": QueryResponseResultBinding{
								Name: "title",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Harry Potter and the Chamber of Secrets",
								},
							},
						},
					},
					{
						Bindings: QueryResponseResultBindingMap{
							"book": QueryResponseResultBinding{
								Name: "book",
								Term: rdf.IRI("http://example.org/book/book3"),
							},
							"title": QueryResponseResultBinding{
								Name: "title",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Harry Potter and the Prisoner Of Azkaban",
								},
							},
						},
					},
					{
						Bindings: QueryResponseResultBindingMap{
							"book": QueryResponseResultBinding{
								Name: "book",
								Term: rdf.IRI("http://example.org/book/book1"),
							},
							"title": QueryResponseResultBinding{
								Name: "title",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Harry Potter and the Philosopher's Stone",
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "2/2",
			Input: `{ 
  "head" : { } ,
  "boolean" : true
}`,
			Output: &QueryResponse{
				Head:    QueryResponseHead{},
				Boolean: ptr.Value(true),
			},
		},
		{
			Name: "5/1",
			Input: `{
   "head": {
       "link": [
           "http://www.w3.org/TR/rdf-sparql-XMLres/example.rq"
           ],
       "vars": [
           "x",
           "hpage",
           "name",
           "mbox",
           "age",
           "blurb",
           "friend"
           ]
       },
   "results": {
       "bindings": [
               {
                   "x" : { "type": "bnode", "value": "r1" },

                   "hpage" : { "type": "uri", "value": "http://work.example.org/alice/" },

                   "name" : {  "type": "literal", "value": "Alice" } ,
                   
		   "mbox" : {  "type": "literal", "value": "" } ,

                   "blurb" : {
                     "datatype": "http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral",
                     "type": "literal",
                     "value": "<p xmlns=\"http://www.w3.org/1999/xhtml\">My name is <b>alice</b></p>"
                   },

                   "friend" : { "type": "bnode", "value": "r2" }
               },
               {
                   "x" : { "type": "bnode", "value": "r2" },
                   
                   "hpage" : { "type": "uri", "value": "http://work.example.org/bob/" },
                   
                   "name" : { "type": "literal", "value": "Bob", "xml:lang": "en" },

                   "mbox" : { "type": "uri", "value": "mailto:bob@work.example.org" },

                   "friend" : { "type": "bnode", "value": "r1" }
               }
           ]
       }
}`,
			Output: &QueryResponse{
				Head: QueryResponseHead{
					Links: QueryResponseHeadLinkList{
						{
							Href: "http://www.w3.org/TR/rdf-sparql-XMLres/example.rq",
						},
					},
					Variables: QueryResponseHeadVariableList{
						{
							Name: "x",
						},
						{
							Name: "hpage",
						},
						{
							Name: "name",
						},
						{
							Name: "mbox",
						},
						{
							Name: "age",
						},
						{
							Name: "blurb",
						},
						{
							Name: "friend",
						},
					},
				},
				Results: &QueryResponseResultList{
					{
						Bindings: QueryResponseResultBindingMap{
							"x": QueryResponseResultBinding{
								Name: "x",
								Term: testingBnode.NewStringBlankNode("r1"),
							},
							"hpage": QueryResponseResultBinding{
								Name: "hpage",
								Term: rdf.IRI("http://work.example.org/alice/"),
							},
							"name": QueryResponseResultBinding{
								Name: "name",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "Alice",
								},
							},
							"mbox": QueryResponseResultBinding{
								Name: "mbox",
								Term: rdf.Literal{
									Datatype:    xsdiri.String_Datatype,
									LexicalForm: "",
								},
							},
							"blurb": QueryResponseResultBinding{
								Name: "blurb",
								Term: rdf.Literal{
									Datatype:    "http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral",
									LexicalForm: "<p xmlns=\"http://www.w3.org/1999/xhtml\">My name is <b>alice</b></p>",
								},
							},
							"friend": QueryResponseResultBinding{
								Name: "friend",
								Term: testingBnode.NewStringBlankNode("r2"),
							},
						},
					},
					{
						Bindings: QueryResponseResultBindingMap{
							"x": QueryResponseResultBinding{
								Name: "x",
								Term: testingBnode.NewStringBlankNode("r2"),
							},
							"hpage": QueryResponseResultBinding{
								Name: "hpage",
								Term: rdf.IRI("http://work.example.org/bob/"),
							},
							"name": QueryResponseResultBinding{
								Name: "name",
								Term: rdf.Literal{
									Datatype:    rdfiri.LangString_Datatype,
									LexicalForm: "Bob",
									Tag: rdf.LanguageLiteralTag{
										Language: "en",
									},
								},
							},
							"mbox": QueryResponseResultBinding{
								Name: "mbox",
								Term: rdf.IRI("mailto:bob@work.example.org"),
							},
							"friend": QueryResponseResultBinding{
								Name: "friend",
								Term: testingBnode.NewStringBlankNode("r1"),
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			d := NewQueryResponseDecoderJSON(strings.NewReader(tc.Input))
			d.BlankNodeTable = testingBnode

			res, err := d.Decode()
			if err == nil && len(tc.Error) > 0 {
				t.Errorf("expected error, but got nil")
			} else if err != nil {
				if err.Error() != tc.Error {
					t.Errorf("unexpected error: %s", err)
				}
			} else if _e, _a := tc.Output.Head, res.Head; !reflect.DeepEqual(_e, _a) {
				t.Errorf("expected [%v], but got: %v", _e, _a)
			}

			if _e, _a := tc.Output.Results, res.Results; _e != nil || _a != nil {
				if _e != nil && _a != nil {
					if _el, _al := len(*_e), len(*_a); _el != _al {
						t.Errorf("expected length [%v], but got: %v", _el, _al)
					} else {
						for idx, re := range *_e {
							ra := (*_a)[idx]
							if _ral, _ael := len(re.Bindings), len(ra.Bindings); _ral != _ael {
								t.Errorf("expected length [%v], but got: %v", _ral, _ael)
							}

							for k, v := range re.Bindings {
								assertQueryResponseResultBindingEqual(t, v, ra.Bindings[k])
							}
						}
					}
				} else if _e != nil {
					t.Errorf("expected [%v], but got: nil", _e)
				} else {
					t.Errorf("expected nil, but got: %v", _a)
				}
			}

			if _e, _a := tc.Output.Boolean, res.Boolean; _e != nil || _a != nil {
				if _e != nil && _a != nil {
					if !reflect.DeepEqual(_e, _a) {
						t.Errorf("expected [%v], but got: %v", _e, _a)
					}
				} else if _e != nil {
					t.Errorf("expected [%v], but got: nil", _e)
				} else {
					t.Errorf("expected nil, but got: %v", _a)
				}
			}
		})
	}
}
