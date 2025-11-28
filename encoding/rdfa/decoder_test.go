package rdfa

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfliteral"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdliteral"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/rdfio/rdfioutil"
)

var testingBnode = blanknodeutil.NewStringMapper()

func lazyAssertEquals(t *testing.T, expected, actual rdfio.StatementList) {
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

	lines0 := strings.Split(lazyCompare[0].String(), "\n")
	slices.SortFunc(lines0, strings.Compare)

	lines1 := strings.Split(lazyCompare[1].String(), "\n")
	slices.SortFunc(lines1, strings.Compare)

	if str0, str1 := strings.Join(lines0, "\n"), strings.Join(lines1, "\n"); str0 != str1 {
		t.Fatalf("expected does not match actual\n\n# EXPECTED\n%s\n# ACTUAL\n%s", str0, str1)
	}
}

// https://www.w3.org/TR/html-rdfa/
func TestW3trHtmlRdfaNonNormative(t *testing.T) {
	for _, testcase := range []struct {
		Name     string
		Snippet  string
		Expected rdfio.StatementList
	}{
		{
			Name: "3.4/Example 3",
			Snippet: `<p xmlns:ex="http://example.org/vocab#"
   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
 Two rectangles (the example markup for them are stored in a triple):
 <svg xmlns="http://www.w3.org/2000/svg" property="ex:markup" datatype="rdf:XMLLiteral"><rect width="300" height="100" style="fill:rgb(0,0,255);stroke-width:1; stroke:rgb(0,0,0)"/><rect width="50" height="50" style="fill:rgb(255,0,0);stroke-width:2;stroke:rgb(0,0,0)"/></svg>
</p>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://example.org/vocab#markup"),
						Object: rdf.Literal{
							LexicalForm: `<rect xmlns="http://www.w3.org/2000/svg" width="300" height="100" style="fill:rgb(0,0,255);stroke-width:1; stroke:rgb(0,0,0)"/><rect xmlns="http://www.w3.org/2000/svg" width="50" height="50" style="fill:rgb(255,0,0);stroke-width:2;stroke:rgb(0,0,0)"/>`,
							Datatype:    rdfiri.XMLLiteral_Datatype,
						},
					},
				},
			},
		},
		{
			// spec is missing "The User" from the expected triple object
			Name: "3.4/Example 5",
			Snippet: `<p xmlns:ex="http://example.org/vocab#"
   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
   xmlns:fb="http://www.facebook.com/2008/fbml">
 This is how you markup a user in FBML:
 <span property="ex:markup" datatype="rdf:XMLLiteral"><span><fb:user uid="12345">The User</fb:user></span></span>
</p>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://example.org/vocab#markup"),
						Object: rdf.Literal{
							LexicalForm: `<span xmlns:fb="http://www.facebook.com/2008/fbml"><fb:user uid="12345">The User</fb:user></span>`,
							Datatype:    rdfiri.XMLLiteral_Datatype,
						},
					},
				},
			},
		},
		{
			Name: "3.5/Example 8",
			Snippet: `<div vocab="http://schema.org/">
  <div resource="#muse" typeof="rdfa:Pattern">
    <link property="image" href="Muse1.jpg"/>
    <link property="image" href="Muse2.jpg"/>
    <link property="image" href="Muse3.jpg"/>
    <span property="name">Muse</span>
  </div>

  <p typeof="MusicEvent">
    <link property="rdfa:copy" href="#muse"/>
    Muse at the United Center.
    <time property="startDate" datetime="2013-03-03">March 3rd 2013</time>, 
    <a property="location" href="#united">United Center, Chicago, Illinois</a>
    ...
  </p>

  <p typeof="MusicEvent">
    <link property="rdfa:copy" href="#muse"/>
    Muse at the Target Center.
    <time property="startDate" datetime="2013-03-07">March 7th 2013</time>, 
    <a property="location" href="#target">Target Center, Minneapolis, Minnesota</a>
    ...
  </p>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/MusicEvent"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/startDate"),
						Object: rdf.Literal{
							LexicalForm: "2013-03-03",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/location"),
						Object:    rdf.IRI("#united"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/MusicEvent"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/startDate"),
						Object: rdf.Literal{
							LexicalForm: "2013-03-07",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/location"),
						Object:    rdf.IRI("#target"),
					},
				},
				// pattern
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/image"),
						Object:    rdf.IRI("Muse1.jpg"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/image"),
						Object:    rdf.IRI("Muse2.jpg"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/image"),
						Object:    rdf.IRI("Muse3.jpg"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/name"),
						Object:    xsdliteral.NewString("Muse"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/image"),
						Object:    rdf.IRI("Muse1.jpg"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/image"),
						Object:    rdf.IRI("Muse2.jpg"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/image"),
						Object:    rdf.IRI("Muse3.jpg"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/name"),
						Object:    xsdliteral.NewString("Muse"),
					},
				},
			},
		},
		{
			Name: "3.5/Example 10",
			Snippet: `<div vocab="http://schema.org/">
  <div typeof="Person">
    <link property="rdfa:copy" href="#lennon"/>
    <link property="rdfa:copy" href="#band"/>
  </div>
  <p resource="#lennon" typeof="rdfa:Pattern"> 
    Name: <span property="name">John Lennon</span>
  <p>
  <div resource="#band" typeof="rdfa:Pattern">
    <div property="band" typeof="MusicGroup">
      <link property="rdfa:copy" href="#beatles"/>
    </div>
  </div>
  <div resource="#beatles" typeof="rdfa:Pattern">
    <p>Band: <span property="name">The Beatles</span></p>
    <p>Size: <span property="size">4</span> players</p>
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/Person"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/name"),
						Object:    xsdliteral.NewString("John Lennon"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://schema.org/band"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/MusicGroup"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/name"),
						Object:    xsdliteral.NewString("The Beatles"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://schema.org/size"),
						Object:    xsdliteral.NewString("4"),
					},
				},
			},
		},
	} {
		t.Run(testcase.Name, func(t *testing.T) {
			htmlDocument, err := html.ParseDocument(
				bytes.NewBufferString(testcase.Snippet),
				html.DocumentConfig{}.SetCaptureTextOffsets(true),
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out, err := rdfio.CollectStatementsErr(NewDecoder(htmlDocument))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lazyAssertEquals(t, testcase.Expected, out)
		})
	}
}

// https://www.w3.org/TR/rdfa-core/
func TestW3trRdfaCoreNonNormative(t *testing.T) {
	for _, testcase := range []struct {
		Name     string
		Snippet  string
		Expected rdfio.StatementList
	}{
		{
			Name: "8.1.1.1/Example 47",
			Snippet: `<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <title>Jo's Friends and Family Blog</title>
    <link rel="foaf:primaryTopic" href="#bbq" />
    <meta property="dc:creator" content="Jo" />
  </head>
  <body>
    ...
  </body>
</html>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/primaryTopic"),
						Object:    rdf.IRI("#bbq"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Jo"),
					},
				},
			},
		},
		{
			Name: "8.1.1.1/Example 49",
			Snippet: `<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <title>Jo's Blog</title>
  </head>
  <body>
    <h1><span property="dc:creator">Jo</span>'s blog</h1>
    <p>
      Welcome to my blog.
    </p>
  </body>
</html>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Jo"),
					},
				},
			},
		},
		{
			// TODO: spec incorrectly expects //jo/blog#bbq with double slash
			Name: "8.1.1.1/Example 51",
			Snippet: `<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <base href="http://www.example.org/jo/blog" />
    <title>Jo's Friends and Family Blog</title>
    <link rel="foaf:primaryTopic" href="#bbq" />
    <meta property="dc:creator" content="Jo" />
  </head>
  <body>
    ...
  </body>
</html>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://www.example.org/jo/blog"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/primaryTopic"),
						Object:    rdf.IRI("http://www.example.org/jo/blog#bbq"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://www.example.org/jo/blog"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Jo"),
					},
				},
			},
		},
		{
			Name: "8.1.1.2/Example 53",
			Snippet: `<html xmlns="http://www.w3.org/1999/xhtml"
      prefix="cal: http://www.w3.org/2002/12/cal/ical#">
  <head>
    <title>Jo's Friends and Family Blog</title>
    <link rel="foaf:primaryTopic" href="#bbq" />
    <meta property="dc:creator" content="Jo" />
  </head>
  <body>
    <p about="#bbq" typeof="cal:Vevent">
      I'm holding
      <span property="cal:summary">
        one last summer barbecue
      </span>,
      on
      <span property="cal:dtstart" content="2015-09-16T16:00:00-05:00" 
            datatype="xsd:dateTime">
        September 16th at 4pm
      </span>.
    </p>
  </body>
</html>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/primaryTopic"),
						Object:    rdf.IRI("#bbq"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Jo"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("#bbq"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://www.w3.org/2002/12/cal/ical#Vevent"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("#bbq"),
						Predicate: rdf.IRI("http://www.w3.org/2002/12/cal/ical#summary"),
						Object:    xsdliteral.NewString("\n        one last summer barbecue\n      "),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("#bbq"),
						Predicate: rdf.IRI("http://www.w3.org/2002/12/cal/ical#dtstart"),
						Object: rdf.Literal{
							LexicalForm: "2015-09-16T16:00:00-05:00",
							Datatype:    xsdiri.DateTime_Datatype,
						},
					},
				},
			},
		},
		{
			Name: "8.1.1.2/Example 55",
			Snippet: `John knows
<a about="mailto:john@example.org"
  rel="foaf:knows" href="mailto:sue@example.org">Sue</a>.

Sue knows
<a about="mailto:sue@example.org"
  rel="foaf:knows" href="mailto:jim@example.org">Jim</a>.`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("mailto:john@example.org"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/knows"),
						Object:    rdf.IRI("mailto:sue@example.org"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("mailto:sue@example.org"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/knows"),
						Object:    rdf.IRI("mailto:jim@example.org"),
					},
				},
			},
		},
		{
			Name: "8.1.1.2/Example 57",
			Snippet: `<div about="photo1.jpg">
  this photo was taken by
  <span property="dc:creator">Mark Birbeck</span>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("photo1.jpg"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Mark Birbeck"),
					},
				},
			},
		},
		{
			Name: "8.1.1.3/Example 59",
			Snippet: `<div about="http://dbpedia.org/resource/Albert_Einstein" typeof="foaf:Person">
  <span property="foaf:name">Albert Einstein</span>
  <span property="foaf:givenName">Albert</span>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://xmlns.com/foaf/0.1/Person"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/givenName"),
						Object:    xsdliteral.NewString("Albert"),
					},
				},
			},
		},
		{
			Name: "8.1.1.3/Example 61",
			Snippet: `<div about="http://dbpedia.org/resource/Albert_Einstein">
  <div rel="dbp:birthPlace" 
      resource="http://dbpedia.org/resource/German_Empire"
      typeof="http://schema.org/Country">
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/property/birthPlace"),
						Object:    rdf.IRI("http://dbpedia.org/resource/German_Empire"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/German_Empire"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/Country"),
					},
				},
			},
		},
		{
			Name: "8.1.1.3/Example 63",
			Snippet: `<div typeof="foaf:Person">
  <span property="foaf:name">Albert Einstein</span>
  <span property="foaf:givenName">Albert</span>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://xmlns.com/foaf/0.1/Person"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/givenName"),
						Object:    xsdliteral.NewString("Albert"),
					},
				},
			},
		},
		{
			Name: "8.1.1.3/Example 65",
			Snippet: `<div resource="_:a" typeof="foaf:Person">
  <span property="foaf:name">Albert Einstein</span>
  <span property="foaf:givenName">Albert</span>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://xmlns.com/foaf/0.1/Person"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/givenName"),
						Object:    xsdliteral.NewString("Albert"),
					},
				},
			},
		},
		{
			Name: "8.1.1.3/Example 66",
			Snippet: `<div about="http://dbpedia.org/resource/Albert_Einstein">
  <div rel="dbp:birthPlace" typeof="http://schema.org/Country">
    <span property="dbp:conventionalLongName">the German Empire</span>
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/property/birthPlace"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/Country"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://dbpedia.org/property/conventionalLongName"),
						Object:    xsdliteral.NewString("the German Empire"),
					},
				},
			},
		},
		{
			Name: "8.1.1.4.1/Example 71",
			Snippet: `<div about="http://dbpedia.org/resource/Albert_Einstein">
  <span property="foaf:name">Albert Einstein</span>
  <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
  <div rel="dbp:birthPlace" resource="http://dbpedia.org/resource/German_Empire">
    <span property="dbp:conventionalLongName">the German Empire</span>
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/property/birthPlace"),
						Object:    rdf.IRI("http://dbpedia.org/resource/German_Empire"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/German_Empire"),
						Predicate: rdf.IRI("http://dbpedia.org/property/conventionalLongName"),
						Object:    xsdliteral.NewString("the German Empire"),
					},
				},
			},
		},
		{
			Name: "8.1.1.4.1/Example 72",
			Snippet: `<div about="http://dbpedia.org/resource/Albert_Einstein">
  <span property="foaf:name">Albert Einstein</span>
  <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
  <div rel="dbp:birthPlace" resource="http://dbpedia.org/resource/German_Empire">
    <span property="dbp:conventionalLongName">the German Empire</span>
    <span rel="dbp-owl:capital" resource="http://dbpedia.org/resource/Berlin" />
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/property/birthPlace"),
						Object:    rdf.IRI("http://dbpedia.org/resource/German_Empire"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/German_Empire"),
						Predicate: rdf.IRI("http://dbpedia.org/property/conventionalLongName"),
						Object:    xsdliteral.NewString("the German Empire"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/German_Empire"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/capital"),
						Object:    rdf.IRI("http://dbpedia.org/resource/Berlin"),
					},
				},
			},
		},
		{
			Name: "8.1.1.4.2/Example 75",
			Snippet: `<div about="http://dbpedia.org/resource/Baruch_Spinoza" rel="dbp-owl:influenced">
  <div about="http://dbpedia.org/resource/Albert_Einstein">
    <span property="foaf:name">Albert Einstein</span>
    <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
			},
		},
		{
			Name: "8.1.1.4.2/Example 77",
			Snippet: `<div about="http://dbpedia.org/resource/Baruch_Spinoza" rel="dbp-owl:influenced">
  <div>
    <span property="foaf:name">Albert Einstein</span>
    <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				// TODO: decoder currently generates duplicate influenced triple
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
			},
		},
		{
			Name: "8.1.1.4.2/Example 79",
			Snippet: `<div about="http://dbpedia.org/resource/Baruch_Spinoza" rel="dbp-owl:influenced">
  <span property="foaf:name">Albert Einstein</span>
  <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				// TODO: decoder currently generates duplicate influenced triple
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
			},
		},
		{
			Name: "8.1.1.4.2/Example 80",
			Snippet: `<div about="http://dbpedia.org/resource/Baruch_Spinoza">
  <div rel="dbp-owl:influenced">
    <span property="foaf:name">Albert Einstein</span>
    <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				// TODO: decoder currently generates duplicate influenced triple
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
			},
		},
		{
			Name: "8.2/Example 88",
			Snippet: `<div about="http://dbpedia.org/resource/Albert_Einstein" rel="dbp-owl:residence">
  <span about="http://dbpedia.org/resource/German_Empire" />
  <span about="http://dbpedia.org/resource/Switzerland" />
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/residence"),
						Object:    rdf.IRI("http://dbpedia.org/resource/German_Empire"),
					},
				},
				// TODO: decoder should also complete incomplete triple for second @about
				// rdfioutil.Statement{
				// 	Triple: rdf.Triple{
				// 		Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
				// 		Predicate: rdf.IRI("http://dbpedia.org/ontology/residence"),
				// 		Object:    rdf.IRI("http://dbpedia.org/resource/Switzerland"),
				// 	},
				// },
			},
		},
		{
			Name: "8.2/Example 91",
			Snippet: `<div about="http://dbpedia.org/resource/Baruch_Spinoza">
  <div rel="dbp-owl:influenced">
    <div typeof="foaf:Person">
      <span property="foaf:name">Albert Einstein</span>
      <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
    </div>
    <div typeof="foaf:Person">
      <span property="foaf:name">Arthur Schopenhauer</span>
      <span property="dbp:dateOfBirth" datatype="xsd:date">1788-02-22</span>
    </div>          
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://xmlns.com/foaf/0.1/Person"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://xmlns.com/foaf/0.1/Person"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Arthur Schopenhauer"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1788-02-22",
							Datatype:    xsdiri.Date_Datatype,
						},
					},
				},
			},
		},
		{
			// TODO bug-ish decoder emits duplicate influenced triples
			Name: "8.2/Example 94",
			Snippet: `<div about="http://dbpedia.org/resource/Baruch_Spinoza" rel="dbp-owl:influenced">
  <span property="foaf:name">Albert Einstein</span>
  <span property="dbp:dateOfBirth" datatype="xsd:date">1879-03-14</span>
  <div rel="dbp-owl:residence">
    <span about="http://dbpedia.org/resource/German_Empire" />
    <span about="http://dbpedia.org/resource/Switzerland" />
  </div>
</div>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://dbpedia.org/property/dateOfBirth"),
						Object: rdf.Literal{
							LexicalForm: "1879-03-14",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#date"),
						},
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Baruch_Spinoza"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/influenced"),
						Object:    testingBnode.MapBlankNodeIdentifier("b0"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://dbpedia.org/ontology/residence"),
						Object:    rdf.IRI("http://dbpedia.org/resource/German_Empire"),
					},
				},
			},
		},
		{
			Name: "8.3.1.1/Example 101",
			Snippet: `<meta about="http://internet-apps.blogspot.com/"
      property="dc:creator" content="Mark Birbeck" />`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://internet-apps.blogspot.com/"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Mark Birbeck"),
					},
				},
			},
		},
		{
			Name: "8.3.1.1/Example 102",
			Snippet: `<span about="http://internet-apps.blogspot.com/"
      property="dc:creator">Mark Birbeck</span>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://internet-apps.blogspot.com/"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Mark Birbeck"),
					},
				},
			},
		},
		{
			Name: "8.3.1.1/Example 104",
			Snippet: `<span about="http://internet-apps.blogspot.com/"
      property="dc:creator" content="Mark Birbeck">John Doe</span>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://internet-apps.blogspot.com/"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    xsdliteral.NewString("Mark Birbeck"),
					},
				},
			},
		},
		{
			Name: "8.3.1.1.1/Example 106",
			Snippet: `<meta about="http://example.org/node"
  property="ex:property" xml:lang="fr" content="chat" />`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://example.org/node"),
						Predicate: rdf.IRI("http://example.org/property"),
						Object:    rdfliteral.NewLangString("fr", "chat"),
					},
				},
			},
		},
		{
			// spec is incorrect on expectation which ignored the @prefix override
			Name: "8.3.1.1.1/Example 107",
			Snippet: `<html xmlns="http://www.w3.org/1999/xhtml" 
      prefix="ex: http://www.example.com/ns/" xml:lang="fr">
  <head>
    <title xml:lang="en">Example</title>
    <meta about="http://example.org/node"
      property="ex:property" content="chat" />
  </head>
  ...
</html>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://example.org/node"),
						Predicate: rdf.IRI("http://www.example.com/ns/property"),
						Object:    rdfliteral.NewLangString("fr", "chat"),
					},
				},
			},
		},
		{
			Name: "8.3.1.2/Example 108",
			Snippet: `<span property="cal:dtstart" content="2015-09-16T16:00:00-05:00" 
      datatype="xsd:dateTime">
  September 16th at 4pm
</span>.`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://www.w3.org/2002/12/cal/ical#dtstart"),
						Object: rdf.Literal{
							LexicalForm: "2015-09-16T16:00:00-05:00",
							Datatype:    xsdiri.DateTime_Datatype,
						},
					},
				},
			},
		},
		{
			// spec does not include whitespace
			Name: "8.3.1.3/Example 111",
			Snippet: `<h2 property="dc:title" datatype="rdf:XMLLiteral">
  E = mc<sup>2</sup>: The Most Urgent Problem of Our Time
</h2>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI(""),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object: rdf.Literal{
							LexicalForm: "\n  E = mc<sup>2</sup>: The Most Urgent Problem of Our Time\n",
							Datatype:    rdfiri.XMLLiteral_Datatype,
						},
					},
				},
			},
		},
		{
			// not sure what this is actually testing? there is no <sup> tag as mentioned in the description
			Name: "8.3.1.3/Example 113",
			Snippet: `<p>You searched for <strong>Einstein</strong>:</p>
<p about="http://dbpedia.org/resource/Albert_Einstein">
  <span property="foaf:name" datatype="">Albert <strong>Einstein</strong></span>
  (b. March 14, 1879, d. April 18, 1955) was a German-born theoretical physicist.
</p>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://dbpedia.org/resource/Albert_Einstein"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
						Object:    xsdliteral.NewString("Albert Einstein"),
					},
				},
			},
		},
		{
			Name: "8.3.2.1/Example 115",
			Snippet: `<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <title>On Crime and Punishment</title>
    <base href="http://www.example.com/candp.xhtml" />
  </head>
  <body>
    <blockquote about="#q1" rel="dc:source" resource="urn:ISBN:0140449132" >
      <p id="q1">
        Rodion Romanovitch! My dear friend! If you go on in this way
        you will go mad, I am positive! Drink, pray, if only a few drops!
      </p>
    </blockquote>
  </body>
</html>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://www.example.com/candp.xhtml#q1"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/source"),
						Object:    rdf.IRI("urn:ISBN:0140449132"),
					},
				},
			},
		},
		{
			Name: "8.3.2.2/Example 117",
			Snippet: `<link about="mailto:john@example.org"
      rel="foaf:knows" href="mailto:sue@example.org" />`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("mailto:john@example.org"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/knows"),
						Object:    rdf.IRI("mailto:sue@example.org"),
					},
				},
			},
		},
		{
			Name: "8.3.2.2/Example 118",
			Snippet: `<img about="http://www.blogger.com/profile/1109404"
    src="photo1.jpg" rev="dc:creator" rel="foaf:img"/>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("photo1.jpg"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    rdf.IRI("http://www.blogger.com/profile/1109404"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("http://www.blogger.com/profile/1109404"),
						Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/img"),
						Object:    rdf.IRI("photo1.jpg"),
					},
				},
			},
		},
		{ // spec expectation example 121 has invalid syntax
			Name: "8.4/Example 123",
			Snippet: `<p prefix="bibo: http://purl.org/ontology/bibo/ dc: http://purl.org/dc/terms/" typeof="bibo:Chapter">
  "<span property="dc:title">Semantic Annotation and Retrieval</span>" by
   <a inlist="" property="dc:creator" 
                href="http://ben.adida.net/#me">Ben Adida</a>,
   <a inlist="" property="dc:creator" 
                href="http://twitter.com/markbirbeck">Mark Birbeck</a>, and
   <a inlist="" property="dc:creator" 
                href="http://www.ivan-herman.net/foaf#me">Ivan Herman</a>. 
</p>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://purl.org/ontology/bibo/Chapter"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object:    xsdliteral.NewString("Semantic Annotation and Retrieval"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    rdf.IRI("http://ben.adida.net/#me"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    testingBnode.MapBlankNodeIdentifier("b2"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    rdf.IRI("http://twitter.com/markbirbeck"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    testingBnode.MapBlankNodeIdentifier("b3"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    rdf.IRI("http://www.ivan-herman.net/foaf#me"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
					},
				},
			},
		},
		{
			Name: "8.4/Example 124",
			Snippet: `<p prefix="bibo: http://purl.org/ontology/bibo/ dc: http://purl.org/dc/terms/" typeof="bibo:Chapter">
  "<span property="dc:title">Semantic Annotation and Retrieval</span>", by
  <span inlist="" property="dc:creator" resource="http://ben.adida.net/#me">Ben Adida</span>,
  <span inlist="" property="dc:creator">Mark Birbeck</span>, and
  <span inlist="" property="dc:creator" resource="http://www.ivan-herman.net/foaf#me">Ivan Herman</span>.
</p>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://purl.org/ontology/bibo/Chapter"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object:    xsdliteral.NewString("Semantic Annotation and Retrieval"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    rdf.IRI("http://ben.adida.net/#me"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    testingBnode.MapBlankNodeIdentifier("b2"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    xsdliteral.NewString("Mark Birbeck"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    testingBnode.MapBlankNodeIdentifier("b3"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    rdf.IRI("http://www.ivan-herman.net/foaf#me"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
					},
				},
			},
		},
		{
			Name: "8.4/Example 126",
			Snippet: `<p prefix="bibo: http://purl.org/ontology/bibo/ dc: http://purl.org/dc/terms/" typeof="bibo:Chapter">
  "<span property="dc:title">Semantic Annotation and Retrieval</span>", by
  <span inlist="" rel="dc:creator" resource="http://ben.adida.net/#me">Ben Adida</span>,
  <span inlist="" property="dc:creator">Mark Birbeck</span>, and
  <span inlist="" rel="dc:creator" resource="http://www.ivan-herman.net/foaf#me">Ivan Herman</span>.
</p>`,
			Expected: rdfio.StatementList{
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://purl.org/ontology/bibo/Chapter"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/title"),
						Object:    xsdliteral.NewString("Semantic Annotation and Retrieval"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://purl.org/dc/terms/creator"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    rdf.IRI("http://ben.adida.net/#me"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    testingBnode.MapBlankNodeIdentifier("b2"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    xsdliteral.NewString("Mark Birbeck"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    testingBnode.MapBlankNodeIdentifier("b3"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#first"),
						Object:    rdf.IRI("http://www.ivan-herman.net/foaf#me"),
					},
				},
				rdfioutil.Statement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"),
						Object:    rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"),
					},
				},
			},
		},
		// Example 127 is invalid; inlist node doesn't containing any properties as intended
	} {
		t.Run(testcase.Name, func(t *testing.T) {
			htmlDocument, err := html.ParseDocument(
				bytes.NewBufferString(testcase.Snippet),
				html.DocumentConfig{}.SetCaptureTextOffsets(true),
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out, err := rdfio.CollectStatementsErr(NewDecoder(htmlDocument, DecoderConfig{}.
				SetDefaultPrefixes(iriutil.NewPrefixMap(
					iriutil.PrefixMapping{
						Prefix:   "bibo",
						Expanded: "http://purl.org/ontology/bibo/",
					},
					iriutil.PrefixMapping{
						Prefix:   "cal",
						Expanded: "http://www.w3.org/2002/12/cal/ical#",
					},
					iriutil.PrefixMapping{
						Prefix:   "cc",
						Expanded: "http://creativecommons.org/ns#",
					},
					iriutil.PrefixMapping{
						Prefix:   "dbp",
						Expanded: "http://dbpedia.org/property/",
					},
					iriutil.PrefixMapping{
						Prefix:   "dbp-owl",
						Expanded: "http://dbpedia.org/ontology/",
					},
					iriutil.PrefixMapping{
						Prefix:   "dbr",
						Expanded: "http://dbpedia.org/resource/",
					},
					iriutil.PrefixMapping{
						Prefix:   "dc",
						Expanded: "http://purl.org/dc/terms/",
					},
					iriutil.PrefixMapping{
						Prefix:   "ex",
						Expanded: "http://example.org/",
					},
					iriutil.PrefixMapping{
						Prefix:   "foaf",
						Expanded: "http://xmlns.com/foaf/0.1/",
					},
					iriutil.PrefixMapping{
						Prefix:   "owl",
						Expanded: "http://www.w3.org/2002/07/owl#",
					},
					iriutil.PrefixMapping{
						Prefix:   "rdf",
						Expanded: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
					},
					iriutil.PrefixMapping{
						Prefix:   "rdfa",
						Expanded: "http://www.w3.org/ns/rdfa#",
					},
					iriutil.PrefixMapping{
						Prefix:   "rdfs",
						Expanded: "http://www.w3.org/2000/01/rdf-schema#",
					},
					iriutil.PrefixMapping{
						Prefix:   "xhv",
						Expanded: "http://www.w3.org/1999/xhtml/vocab#",
					},
					iriutil.PrefixMapping{
						Prefix:   "xsd",
						Expanded: "http://www.w3.org/2001/XMLSchema#",
					},
				))))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lazyAssertEquals(t, testcase.Expected, out)
		})
	}
}
