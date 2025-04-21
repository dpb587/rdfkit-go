package jsonld

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdfio"
)

var documentLoader = jsonldtype.NewCachingDocumentLoader(jsonldtype.NewDefaultDocumentLoader(http.DefaultClient))

func TestOne(t *testing.T) {
	statements, err := rdfio.CollectStatementsErr(
		NewDecoder(
			strings.NewReader(`
{
  "@context": {
    "@base": "http://example.org/"
  },
  "http://example.org/vocab/at": {"@id": "@"},
  "http://example.org/vocab/foo.bar": {"@id": "@foo.bar"},
  "http://example.org/vocab/ignoreme": {"@id": "@ignoreMe"}
}
`),
			DecoderConfig{}.SetDocumentLoader(documentLoader),
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Fprintf(os.Stderr, "%#+v\n", statements)
}

func TestTwo(t *testing.T) {
	statements, err := rdfio.CollectStatementsErr(
		NewDecoder(
			strings.NewReader(`
{
  "@context": {
    "brick": "https://brickschema.org/schema/Brick#",
    "csvw": "http://www.w3.org/ns/csvw#",
    "dc": "http://purl.org/dc/elements/1.1/",
    "dcam": "http://purl.org/dc/dcam/",
    "dcat": "http://www.w3.org/ns/dcat#",
    "dcmitype": "http://purl.org/dc/dcmitype/",
    "dcterms": "http://purl.org/dc/terms/",
    "doap": "http://usefulinc.com/ns/doap#",
    "foaf": "http://xmlns.com/foaf/0.1/",
    "odrl": "http://www.w3.org/ns/odrl/2/",
    "org": "http://www.w3.org/ns/org#",
    "owl": "http://www.w3.org/2002/07/owl#",
    "prof": "http://www.w3.org/ns/dx/prof/",
    "prov": "http://www.w3.org/ns/prov#",
    "qb": "http://purl.org/linked-data/cube#",
    "rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
    "rdfs": "http://www.w3.org/2000/01/rdf-schema#",
    "schema": "https://schema.org/",
    "sh": "http://www.w3.org/ns/shacl#",
    "skos": "http://www.w3.org/2004/02/skos/core#",
    "sosa": "http://www.w3.org/ns/sosa/",
    "ssn": "http://www.w3.org/ns/ssn/",
    "time": "http://www.w3.org/2006/time#",
    "vann": "http://purl.org/vocab/vann/",
    "void": "http://rdfs.org/ns/void#",
    "xsd": "http://www.w3.org/2001/XMLSchema#"
  },
  "@graph": [
    {
      "@id": "schema:name",
      "@type": "rdf:Property",
      "owl:equivalentProperty": {
        "@id": "dcterms:title"
      },
      "rdfs:comment": "The name of the item.",
      "rdfs:label": "name",
      "rdfs:subPropertyOf": {
        "@id": "rdfs:label"
      },
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": {
        "@id": "schema:Text"
      }
    },
    {
      "@id": "schema:identifier",
      "@type": "rdf:Property",
      "owl:equivalentProperty": {
        "@id": "dcterms:identifier"
      },
      "rdfs:comment": "The identifier property represents any kind of identifier for any kind of [[Thing]], such as ISBNs, GTIN codes, UUIDs etc. Schema.org provides dedicated properties for representing many of these, either as textual strings or as URL (URI) links. See [background notes](/docs/datamodel.html#identifierBg) for more details.\n        ",
      "rdfs:label": "identifier",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": [
        {
          "@id": "schema:Text"
        },
        {
          "@id": "schema:PropertyValue"
        },
        {
          "@id": "schema:URL"
        }
      ]
    },
    {
      "@id": "schema:Thing",
      "@type": "rdfs:Class",
      "rdfs:comment": "The most generic type of item.",
      "rdfs:label": "Thing"
    },
    {
      "@id": "schema:description",
      "@type": "rdf:Property",
      "owl:equivalentProperty": {
        "@id": "dcterms:description"
      },
      "rdfs:comment": "A description of the item.",
      "rdfs:label": "description",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": [
        {
          "@id": "schema:Text"
        },
        {
          "@id": "schema:TextObject"
        }
      ]
    },
    {
      "@id": "schema:subjectOf",
      "@type": "rdf:Property",
      "rdfs:comment": "A CreativeWork or Event about this Thing.",
      "rdfs:label": "subjectOf",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:inverseOf": {
        "@id": "schema:about"
      },
      "schema:rangeIncludes": [
        {
          "@id": "schema:Event"
        },
        {
          "@id": "schema:CreativeWork"
        }
      ],
      "schema:source": {
        "@id": "https://github.com/schemaorg/schemaorg/issues/1670"
      }
    },
    {
      "@id": "schema:disambiguatingDescription",
      "@type": "rdf:Property",
      "rdfs:comment": "A sub property of description. A short description of the item used to disambiguate from other, similar items. Information from other properties (in particular, name) may be necessary for the description to be useful for disambiguation.",
      "rdfs:label": "disambiguatingDescription",
      "rdfs:subPropertyOf": {
        "@id": "schema:description"
      },
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": {
        "@id": "schema:Text"
      }
    },
    {
      "@id": "schema:additionalType",
      "@type": "rdf:Property",
      "rdfs:comment": "An additional type for the item, typically used for adding more specific types from external vocabularies in microdata syntax. This is a relationship between something and a class that the thing is in. Typically the value is a URI-identified RDF class, and in this case corresponds to the\n    use of rdf:type in RDF. Text values can be used sparingly, for cases where useful information can be added without their being an appropriate schema to reference. In the case of text values, the class label should follow the schema.org <a href=\"https://schema.org/docs/styleguide.html\">style guide</a>.",
      "rdfs:label": "additionalType",
      "rdfs:subPropertyOf": {
        "@id": "rdf:type"
      },
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": [
        {
          "@id": "schema:Text"
        },
        {
          "@id": "schema:URL"
        }
      ]
    },
    {
      "@id": "schema:url",
      "@type": "rdf:Property",
      "rdfs:comment": "URL of the item.",
      "rdfs:label": "url",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": {
        "@id": "schema:URL"
      }
    },
    {
      "@id": "schema:mainEntityOfPage",
      "@type": "rdf:Property",
      "rdfs:comment": "Indicates a page (or other CreativeWork) for which this thing is the main entity being described. See [background notes](/docs/datamodel.html#mainEntityBackground) for details.",
      "rdfs:label": "mainEntityOfPage",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:inverseOf": {
        "@id": "schema:mainEntity"
      },
      "schema:rangeIncludes": [
        {
          "@id": "schema:CreativeWork"
        },
        {
          "@id": "schema:URL"
        }
      ]
    },
    {
      "@id": "schema:sameAs",
      "@type": "rdf:Property",
      "rdfs:comment": "URL of a reference Web page that unambiguously indicates the item's identity. E.g. the URL of the item's Wikipedia page, Wikidata entry, or official website.",
      "rdfs:label": "sameAs",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": {
        "@id": "schema:URL"
      }
    },
    {
      "@id": "schema:image",
      "@type": "rdf:Property",
      "rdfs:comment": "An image of the item. This can be a [[URL]] or a fully described [[ImageObject]].",
      "rdfs:label": "image",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": [
        {
          "@id": "schema:ImageObject"
        },
        {
          "@id": "schema:URL"
        }
      ]
    },
    {
      "@id": "schema:alternateName",
      "@type": "rdf:Property",
      "rdfs:comment": "An alias for the item.",
      "rdfs:label": "alternateName",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": {
        "@id": "schema:Text"
      }
    },
    {
      "@id": "schema:potentialAction",
      "@type": "rdf:Property",
      "rdfs:comment": "Indicates a potential Action, which describes an idealized action in which this thing would play an 'object' role.",
      "rdfs:label": "potentialAction",
      "schema:domainIncludes": {
        "@id": "schema:Thing"
      },
      "schema:rangeIncludes": {
        "@id": "schema:Action"
      }
    }
  ]
}`),
			DecoderConfig{}.SetDocumentLoader(documentLoader),
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Fprintf(os.Stderr, "%#+v\n", statements)
}
