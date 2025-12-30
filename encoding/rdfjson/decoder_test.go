package rdfjson

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

var testingBnode = blanknodeutil.NewStringMapper()

func assertEquals(t *testing.T, expected, actual rdfio.StatementList) {
	var lazyCompare = [2]*bytes.Buffer{
		bytes.NewBuffer(nil),
		bytes.NewBuffer(nil),
	}

	for i, entities := range [2]rdfio.StatementList{expected, actual} {
		ctx := context.Background()
		encoder, err := ntriples.NewEncoder(lazyCompare[i])
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for eIdx, e := range entities {
			triple := e.GetTriple()

			if triple.Subject == nil {
				triple.Subject = testingBnode.MapBlankNodeIdentifier(fmt.Sprintf("b%d", eIdx))
			}

			if err := encoder.PutTriple(ctx, triple); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		if err := encoder.Close(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if lazyCompare[0].String() != lazyCompare[1].String() {
		t.Fatalf("expected does not match actual\n\n# EXPECTED\n%s\n# ACTUAL\n%s", lazyCompare[0].String(), lazyCompare[1].String())
	}
}

func TestExamples(t *testing.T) {
	// https://www.w3.org/TR/rdf-json/
	for _, testcase := range []struct {
		Name     string
		Snippet  string
		Expected rdfio.StatementList
	}{
		{
			Name: "5/1",
			Snippet: `{
  "http://example.org/about" : {
      "http://purl.org/dc/terms/title" : [ { "value" : "Anna's Homepage", 
                                             "type" : "literal", 
                                             "lang" : "en" } ] 
  }
}`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://example.org/about"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object: rdf.Literal{
							LexicalForm: "Anna's Homepage",
							Datatype:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#langString"),
							Tags: map[rdf.LiteralTag]string{
								rdf.LanguageLiteralTag: "en",
							},
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 30, LineColumn: cursorio.TextLineColumn{1, 28}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 41, LineColumn: cursorio.TextLineColumn{2, 6}},
							Until: cursorio.TextOffset{Byte: 73, LineColumn: cursorio.TextLineColumn{2, 38}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{2, 43}},
							Until: cursorio.TextOffset{Byte: 236, LineColumn: cursorio.TextLineColumn{4, 60}},
						},
					},
				},
			},
		},
		{
			Name: "5/3",
			Snippet: `{
  "http://example.org/about" : {
      "http://purl.org/dc/terms/title" : [ { "value" : "Anna's Homepage", 
                                             "type" : "literal", 
                                             "lang" : "en" },
                                           { "value" : "Annas hjemmeside", 
                                             "type" : "literal", 
                                             "lang" : "da" } ] 
  }
}`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://example.org/about"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object: rdf.Literal{
							LexicalForm: "Anna's Homepage",
							Datatype:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#langString"),
							Tags: map[rdf.LiteralTag]string{
								rdf.LanguageLiteralTag: "en",
							},
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 30, LineColumn: cursorio.TextLineColumn{1, 28}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 41, LineColumn: cursorio.TextLineColumn{2, 6}},
							Until: cursorio.TextOffset{Byte: 73, LineColumn: cursorio.TextLineColumn{2, 38}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{2, 43}},
							Until: cursorio.TextOffset{Byte: 236, LineColumn: cursorio.TextLineColumn{4, 60}},
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://example.org/about"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object: rdf.Literal{
							LexicalForm: "Annas hjemmeside",
							Datatype:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#langString"),
							Tags: map[rdf.LiteralTag]string{
								rdf.LanguageLiteralTag: "da",
							},
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 30, LineColumn: cursorio.TextLineColumn{1, 28}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 41, LineColumn: cursorio.TextLineColumn{2, 6}},
							Until: cursorio.TextOffset{Byte: 73, LineColumn: cursorio.TextLineColumn{2, 38}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 281, LineColumn: cursorio.TextLineColumn{5, 43}},
							Until: cursorio.TextOffset{Byte: 440, LineColumn: cursorio.TextLineColumn{7, 60}},
						},
					},
				},
			},
		},
		{
			Name: "5/5",
			Snippet: `{
  "http://example.org/about" : {
      "http://purl.org/dc/terms/title" : [ { "value" : "<p xmlns=\"http://www.w3.org/1999/xhtml\"><b>Anna's</b> Homepage>/p>", 
                                             "type" : "literal", 
                                             "datatype" : "http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral" } ] 
  }
}`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://example.org/about"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object: rdf.Literal{
							LexicalForm: "<p xmlns=\"http://www.w3.org/1999/xhtml\"><b>Anna's</b> Homepage>/p>",
							Datatype:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 30, LineColumn: cursorio.TextLineColumn{1, 28}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 41, LineColumn: cursorio.TextLineColumn{2, 6}},
							Until: cursorio.TextOffset{Byte: 73, LineColumn: cursorio.TextLineColumn{2, 38}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{2, 43}},
							Until: cursorio.TextOffset{Byte: 344, LineColumn: cursorio.TextLineColumn{4, 115}},
						},
					},
				},
			},
		},
		{
			Name: "5/7",
			// spec is missing a closing curly
			Snippet: `{
  "http://example.org/about" : {
      "http://purl.org/dc/terms/creator" : [ { "value" : "_:anna", 
                                               "type" : "bnode" } ] },
  "_:anna" : {
      "http://xmlns.com/foaf/0.1/name" : [ { "value" : "Anna", 
                                             "type" : "literal" } ] 
  }
}`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://example.org/about"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 30, LineColumn: cursorio.TextLineColumn{1, 28}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 41, LineColumn: cursorio.TextLineColumn{2, 6}},
							Until: cursorio.TextOffset{Byte: 75, LineColumn: cursorio.TextLineColumn{2, 40}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 80, LineColumn: cursorio.TextLineColumn{2, 45}},
							Until: cursorio.TextOffset{Byte: 168, LineColumn: cursorio.TextLineColumn{3, 65}},
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdobject.String("Anna"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 176, LineColumn: cursorio.TextLineColumn{4, 2}},
							Until: cursorio.TextOffset{Byte: 184, LineColumn: cursorio.TextLineColumn{4, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 195, LineColumn: cursorio.TextLineColumn{5, 6}},
							Until: cursorio.TextOffset{Byte: 227, LineColumn: cursorio.TextLineColumn{5, 38}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 232, LineColumn: cursorio.TextLineColumn{5, 43}},
							Until: cursorio.TextOffset{Byte: 318, LineColumn: cursorio.TextLineColumn{6, 65}},
						},
					},
				},
			},
		},
		{
			Name: "5/9",
			Snippet: `{
  "_:anna" : {
      "http://xmlns.com/foaf/0.1/homepage" : [ { "value" : "http://example.org/anna", 
                                                 "type" : "uri" } ] 
  }
}`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/homepage"),
						Object:    rdf.IRI("http://example.org/anna"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 12, LineColumn: cursorio.TextLineColumn{1, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{2, 6}},
							Until: cursorio.TextOffset{Byte: 59, LineColumn: cursorio.TextLineColumn{2, 42}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 64, LineColumn: cursorio.TextLineColumn{2, 47}},
							Until: cursorio.TextOffset{Byte: 169, LineColumn: cursorio.TextLineColumn{3, 65}},
						},
					},
				},
			},
		},
		{
			Name: "5/11",
			Snippet: `{
  "_:anna" : {
      "http://xmlns.com/foaf/0.1/name" : [ { "value" : "Anna", 
                                             "type" : "literal" } ],
      "http://xmlns.com/foaf/0.1/homepage" : [ { "value" : "http://example.org/anna", 
                                                 "type" : "uri" } ] 
  }
}`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdobject.String("Anna"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 12, LineColumn: cursorio.TextLineColumn{1, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{2, 6}},
							Until: cursorio.TextOffset{Byte: 55, LineColumn: cursorio.TextLineColumn{2, 38}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 60, LineColumn: cursorio.TextLineColumn{2, 43}},
							Until: cursorio.TextOffset{Byte: 146, LineColumn: cursorio.TextLineColumn{3, 65}},
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/homepage"),
						Object:    rdf.IRI("http://example.org/anna"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 4, LineColumn: cursorio.TextLineColumn{1, 2}},
							Until: cursorio.TextOffset{Byte: 12, LineColumn: cursorio.TextLineColumn{1, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 156, LineColumn: cursorio.TextLineColumn{4, 6}},
							Until: cursorio.TextOffset{Byte: 192, LineColumn: cursorio.TextLineColumn{4, 42}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 197, LineColumn: cursorio.TextLineColumn{4, 47}},
							Until: cursorio.TextOffset{Byte: 302, LineColumn: cursorio.TextLineColumn{5, 65}},
						},
					},
				},
			},
		},
		{
			Name:     "5/13",
			Snippet:  `{ }`,
			Expected: rdfio.StatementList{},
		},
	} {
		t.Run(testcase.Name, func(t *testing.T) {
			out, err := rdfio.CollectStatementsErr(NewDecoder(
				bytes.NewBufferString(testcase.Snippet),
				DecoderConfig{}.SetCaptureTextOffsets(true),
			))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertEquals(t, testcase.Expected, out)
		})
	}
}
