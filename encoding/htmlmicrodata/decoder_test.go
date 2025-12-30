package htmlmicrodata

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/triples"
)

var testingBnode = blanknodeutil.NewStringMapper()

func microdataLivingAssertEquals(t *testing.T, expected, actual rdf.TripleList) {
	var lazyCompare = [2]*bytes.Buffer{
		bytes.NewBuffer(nil),
		bytes.NewBuffer(nil),
	}

	for i, entities := range [2]rdf.TripleList{expected, actual} {
		ctx := context.Background()
		encoder, err := ntriples.NewEncoder(lazyCompare[i])
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for eIdx, triple := range entities {
			if triple.Subject == nil {
				triple.Subject = testingBnode.MapBlankNodeIdentifier(fmt.Sprintf("b%d", eIdx))
			}

			if err := encoder.AddTriple(ctx, triple); err != nil {
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

func TestMicrodataLivingNonNormative(t *testing.T) {
	// https://html.spec.whatwg.org/multipage/microdata.html
	// "Last Updated 25 July 2024"
	for _, testcase := range []struct {
		Name     string
		Snippet  string
		Expected encodingtest.TripleStatementList
	}{
		{
			Name: "5.1.2/1",
			Snippet: `<div itemscope>
 <p>My name is <span itemprop="name">Elizabeth</span>.</p>
</div>

<div itemscope>
 <p>My name is <span itemprop="name">Daniel</span>.</p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Elizabeth"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 81, LineColumn: cursorio.TextLineColumn{2, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 46, LineColumn: cursorio.TextLineColumn{1, 30}},
							Until: cursorio.TextOffset{Byte: 52, LineColumn: cursorio.TextLineColumn{1, 36}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 53, LineColumn: cursorio.TextLineColumn{1, 37}},
							Until: cursorio.TextOffset{Byte: 62, LineColumn: cursorio.TextLineColumn{1, 46}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Daniel"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 83, LineColumn: cursorio.TextLineColumn{4, 0}},
							Until: cursorio.TextOffset{Byte: 161, LineColumn: cursorio.TextLineColumn{6, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 129, LineColumn: cursorio.TextLineColumn{5, 30}},
							Until: cursorio.TextOffset{Byte: 135, LineColumn: cursorio.TextLineColumn{5, 36}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 136, LineColumn: cursorio.TextLineColumn{5, 37}},
							Until: cursorio.TextOffset{Byte: 142, LineColumn: cursorio.TextLineColumn{5, 43}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/2",
			Snippet: `<div itemscope>
 <p>My <em>name</em> is <span itemprop="name">E<strong>liz</strong>abeth</span>.</p>
</div>

<section>
 <div itemscope>
  <aside>
   <p>My name is <span itemprop="name"><a href="/?user=daniel">Daniel</a></span>.</p>
  </aside>
 </div>
</section>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Elizabeth"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 107, LineColumn: cursorio.TextLineColumn{2, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 55, LineColumn: cursorio.TextLineColumn{1, 39}},
							Until: cursorio.TextOffset{Byte: 61, LineColumn: cursorio.TextLineColumn{1, 45}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 62, LineColumn: cursorio.TextLineColumn{1, 46}},
							Until: cursorio.TextOffset{Byte: 88, LineColumn: cursorio.TextLineColumn{1, 72}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Daniel"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 120, LineColumn: cursorio.TextLineColumn{5, 1}},
							Until: cursorio.TextOffset{Byte: 250, LineColumn: cursorio.TextLineColumn{9, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 178, LineColumn: cursorio.TextLineColumn{7, 32}},
							Until: cursorio.TextOffset{Byte: 184, LineColumn: cursorio.TextLineColumn{7, 38}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 185, LineColumn: cursorio.TextLineColumn{7, 39}},
							Until: cursorio.TextOffset{Byte: 219, LineColumn: cursorio.TextLineColumn{7, 73}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/3",
			Snippet: `<div itemscope>
 <p>My name is <span itemprop="name">Neil</span>.</p>
 <p>My band is called <span itemprop="band">Four Parts Water</span>.</p>
 <p>I am <span itemprop="nationality">British</span>.</p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Neil"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 207, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 46, LineColumn: cursorio.TextLineColumn{1, 30}},
							Until: cursorio.TextOffset{Byte: 52, LineColumn: cursorio.TextLineColumn{1, 36}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 53, LineColumn: cursorio.TextLineColumn{1, 37}},
							Until: cursorio.TextOffset{Byte: 57, LineColumn: cursorio.TextLineColumn{1, 41}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("band"),
						Object:    xsdobject.String("Four Parts Water"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 207, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 107, LineColumn: cursorio.TextLineColumn{2, 37}},
							Until: cursorio.TextOffset{Byte: 113, LineColumn: cursorio.TextLineColumn{2, 43}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 114, LineColumn: cursorio.TextLineColumn{2, 44}},
							Until: cursorio.TextOffset{Byte: 130, LineColumn: cursorio.TextLineColumn{2, 60}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("nationality"),
						Object:    xsdobject.String("British"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 207, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 167, LineColumn: cursorio.TextLineColumn{3, 24}},
							Until: cursorio.TextOffset{Byte: 180, LineColumn: cursorio.TextLineColumn{3, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 181, LineColumn: cursorio.TextLineColumn{3, 38}},
							Until: cursorio.TextOffset{Byte: 188, LineColumn: cursorio.TextLineColumn{3, 45}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/4",
			Snippet: `<div itemscope>
 <img itemprop="image" src="google-logo.png" alt="Google">
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("image"),
						Object:    rdf.IRI("google-logo.png"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 81, LineColumn: cursorio.TextLineColumn{2, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 31, LineColumn: cursorio.TextLineColumn{1, 15}},
							Until: cursorio.TextOffset{Byte: 38, LineColumn: cursorio.TextLineColumn{1, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 43, LineColumn: cursorio.TextLineColumn{1, 27}},
							Until: cursorio.TextOffset{Byte: 60, LineColumn: cursorio.TextLineColumn{1, 44}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/5",
			Snippet: `<h1 itemscope>
 <data itemprop="product-id" value="9678AOU879">The Instigator 2000</data>
</h1>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("product-id"),
						Object:    xsdobject.String("9678AOU879"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 95, LineColumn: cursorio.TextLineColumn{2, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 31, LineColumn: cursorio.TextLineColumn{1, 16}},
							Until: cursorio.TextOffset{Byte: 43, LineColumn: cursorio.TextLineColumn{1, 28}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 50, LineColumn: cursorio.TextLineColumn{1, 35}},
							Until: cursorio.TextOffset{Byte: 62, LineColumn: cursorio.TextLineColumn{1, 47}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/6",
			Snippet: `<div itemscope itemtype="http://schema.org/Product">
 <span itemprop="name">Panasonic White 60L Refrigerator</span>
 <img src="panasonic-fridge-60l-white.jpg" alt="">
  <div itemprop="aggregateRating"
       itemscope itemtype="http://schema.org/AggregateRating">
   <meter itemprop="ratingValue" min=0 value=3.5 max=5>Rated 3.5/5</meter>
   (based on <span itemprop="reviewCount">11</span> customer reviews)
  </div>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/Product"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 424, LineColumn: cursorio.TextLineColumn{8, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 15, LineColumn: cursorio.TextLineColumn{0, 15}},
							Until: cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{0, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 24, LineColumn: cursorio.TextLineColumn{0, 24}},
							Until: cursorio.TextOffset{Byte: 51, LineColumn: cursorio.TextLineColumn{0, 51}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Panasonic White 60L Refrigerator"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 424, LineColumn: cursorio.TextLineColumn{8, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 69, LineColumn: cursorio.TextLineColumn{1, 16}},
							Until: cursorio.TextOffset{Byte: 75, LineColumn: cursorio.TextLineColumn{1, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 76, LineColumn: cursorio.TextLineColumn{1, 23}},
							Until: cursorio.TextOffset{Byte: 108, LineColumn: cursorio.TextLineColumn{1, 55}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("aggregateRating"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 424, LineColumn: cursorio.TextLineColumn{8, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 183, LineColumn: cursorio.TextLineColumn{3, 16}},
							Until: cursorio.TextOffset{Byte: 200, LineColumn: cursorio.TextLineColumn{3, 33}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 169, LineColumn: cursorio.TextLineColumn{3, 2}},
							Until: cursorio.TextOffset{Byte: 417, LineColumn: cursorio.TextLineColumn{7, 8}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/AggregateRating"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 169, LineColumn: cursorio.TextLineColumn{3, 2}},
							Until: cursorio.TextOffset{Byte: 417, LineColumn: cursorio.TextLineColumn{7, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 218, LineColumn: cursorio.TextLineColumn{4, 17}},
							Until: cursorio.TextOffset{Byte: 226, LineColumn: cursorio.TextLineColumn{4, 25}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 227, LineColumn: cursorio.TextLineColumn{4, 26}},
							Until: cursorio.TextOffset{Byte: 262, LineColumn: cursorio.TextLineColumn{4, 61}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("ratingValue"),
						Object: rdf.Literal{
							LexicalForm: "3.5",
							// TODO confirm meter behavior being a decimal?
							Datatype: rdf.IRI("http://www.w3.org/2001/XMLSchema#decimal"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 169, LineColumn: cursorio.TextLineColumn{3, 2}},
							Until: cursorio.TextOffset{Byte: 417, LineColumn: cursorio.TextLineColumn{7, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 283, LineColumn: cursorio.TextLineColumn{5, 19}},
							Until: cursorio.TextOffset{Byte: 296, LineColumn: cursorio.TextLineColumn{5, 32}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 309, LineColumn: cursorio.TextLineColumn{5, 45}},
							Until: cursorio.TextOffset{Byte: 312, LineColumn: cursorio.TextLineColumn{5, 48}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("reviewCount"),
						Object:    xsdobject.String("11"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 169, LineColumn: cursorio.TextLineColumn{3, 2}},
							Until: cursorio.TextOffset{Byte: 417, LineColumn: cursorio.TextLineColumn{7, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 367, LineColumn: cursorio.TextLineColumn{6, 28}},
							Until: cursorio.TextOffset{Byte: 380, LineColumn: cursorio.TextLineColumn{6, 41}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 381, LineColumn: cursorio.TextLineColumn{6, 42}},
							Until: cursorio.TextOffset{Byte: 383, LineColumn: cursorio.TextLineColumn{6, 44}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/7",
			Snippet: `<div itemscope>
 I was born on <time itemprop="birthday" datetime="2009-05-10">May 10th 2009</time>.
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("birthday"),
						Object: rdf.Literal{
							LexicalForm: "2009-05-10",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#date"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 107, LineColumn: cursorio.TextLineColumn{2, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 46, LineColumn: cursorio.TextLineColumn{1, 30}},
							Until: cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{1, 40}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 66, LineColumn: cursorio.TextLineColumn{1, 50}},
							Until: cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{1, 62}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/8",
			Snippet: `<div itemscope>
 <p>Name: <span itemprop="name">Amanda</span></p>
 <p>Band: <span itemprop="band" itemscope> <span itemprop="name">Jazz Band</span> (<span itemprop="size">12</span> players)</span></p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Amanda"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 207, LineColumn: cursorio.TextLineColumn{3, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 41, LineColumn: cursorio.TextLineColumn{1, 25}},
							Until: cursorio.TextOffset{Byte: 47, LineColumn: cursorio.TextLineColumn{1, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 48, LineColumn: cursorio.TextLineColumn{1, 32}},
							Until: cursorio.TextOffset{Byte: 54, LineColumn: cursorio.TextLineColumn{1, 38}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("band"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 207, LineColumn: cursorio.TextLineColumn{3, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 91, LineColumn: cursorio.TextLineColumn{2, 25}},
							Until: cursorio.TextOffset{Byte: 97, LineColumn: cursorio.TextLineColumn{2, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 76, LineColumn: cursorio.TextLineColumn{2, 10}},
							Until: cursorio.TextOffset{Byte: 196, LineColumn: cursorio.TextLineColumn{2, 130}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Jazz Band"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 76, LineColumn: cursorio.TextLineColumn{2, 10}},
							Until: cursorio.TextOffset{Byte: 196, LineColumn: cursorio.TextLineColumn{2, 130}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 124, LineColumn: cursorio.TextLineColumn{2, 58}},
							Until: cursorio.TextOffset{Byte: 130, LineColumn: cursorio.TextLineColumn{2, 64}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 131, LineColumn: cursorio.TextLineColumn{2, 65}},
							Until: cursorio.TextOffset{Byte: 140, LineColumn: cursorio.TextLineColumn{2, 74}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("size"),
						Object:    xsdobject.String("12"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 76, LineColumn: cursorio.TextLineColumn{2, 10}},
							Until: cursorio.TextOffset{Byte: 196, LineColumn: cursorio.TextLineColumn{2, 130}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 164, LineColumn: cursorio.TextLineColumn{2, 98}},
							Until: cursorio.TextOffset{Byte: 170, LineColumn: cursorio.TextLineColumn{2, 104}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 171, LineColumn: cursorio.TextLineColumn{2, 105}},
							Until: cursorio.TextOffset{Byte: 173, LineColumn: cursorio.TextLineColumn{2, 107}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/9",
			Snippet: `<div itemscope id="amanda" itemref="a b"></div>
<p id="a">Name: <span itemprop="name">Amanda</span></p>
<div id="b" itemprop="band" itemscope itemref="c"></div>
<div id="c">
 <p>Band: <span itemprop="name">Jazz Band</span></p>
 <p>Size: <span itemprop="size">12</span> players</p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Amanda"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 47, LineColumn: cursorio.TextLineColumn{0, 47}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 79, LineColumn: cursorio.TextLineColumn{1, 31}},
							Until: cursorio.TextOffset{Byte: 85, LineColumn: cursorio.TextLineColumn{1, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 86, LineColumn: cursorio.TextLineColumn{1, 38}},
							Until: cursorio.TextOffset{Byte: 92, LineColumn: cursorio.TextLineColumn{1, 44}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("band"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 47, LineColumn: cursorio.TextLineColumn{0, 47}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 125, LineColumn: cursorio.TextLineColumn{2, 21}},
							Until: cursorio.TextOffset{Byte: 131, LineColumn: cursorio.TextLineColumn{2, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 104, LineColumn: cursorio.TextLineColumn{2, 0}},
							Until: cursorio.TextOffset{Byte: 160, LineColumn: cursorio.TextLineColumn{2, 56}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Jazz Band"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 104, LineColumn: cursorio.TextLineColumn{2, 0}},
							Until: cursorio.TextOffset{Byte: 160, LineColumn: cursorio.TextLineColumn{2, 56}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 199, LineColumn: cursorio.TextLineColumn{4, 25}},
							Until: cursorio.TextOffset{Byte: 205, LineColumn: cursorio.TextLineColumn{4, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 206, LineColumn: cursorio.TextLineColumn{4, 32}},
							Until: cursorio.TextOffset{Byte: 215, LineColumn: cursorio.TextLineColumn{4, 41}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("size"),
						Object:    xsdobject.String("12"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 104, LineColumn: cursorio.TextLineColumn{2, 0}},
							Until: cursorio.TextOffset{Byte: 160, LineColumn: cursorio.TextLineColumn{2, 56}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 252, LineColumn: cursorio.TextLineColumn{5, 25}},
							Until: cursorio.TextOffset{Byte: 258, LineColumn: cursorio.TextLineColumn{5, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 259, LineColumn: cursorio.TextLineColumn{5, 32}},
							Until: cursorio.TextOffset{Byte: 261, LineColumn: cursorio.TextLineColumn{5, 34}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/10",
			Snippet: `<div itemscope>
 <p>Flavors in my favorite ice cream:</p>
 <ul>
  <li itemprop="flavor">Lemon sorbet</li>
  <li itemprop="flavor">Apricot sorbet</li>
 </ul>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("flavor"),
						Object:    xsdobject.String("Lemon sorbet"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 163, LineColumn: cursorio.TextLineColumn{6, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 79, LineColumn: cursorio.TextLineColumn{3, 15}},
							Until: cursorio.TextOffset{Byte: 87, LineColumn: cursorio.TextLineColumn{3, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 88, LineColumn: cursorio.TextLineColumn{3, 24}},
							Until: cursorio.TextOffset{Byte: 100, LineColumn: cursorio.TextLineColumn{3, 36}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("flavor"),
						Object:    xsdobject.String("Apricot sorbet"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 163, LineColumn: cursorio.TextLineColumn{6, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 121, LineColumn: cursorio.TextLineColumn{4, 15}},
							Until: cursorio.TextOffset{Byte: 129, LineColumn: cursorio.TextLineColumn{4, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 130, LineColumn: cursorio.TextLineColumn{4, 24}},
							Until: cursorio.TextOffset{Byte: 144, LineColumn: cursorio.TextLineColumn{4, 38}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/11",
			Snippet: `<div itemscope>
 <span itemprop="favorite-color favorite-fruit">orange</span>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("favorite-color"),
						Object:    xsdobject.String("orange"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{2, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 32, LineColumn: cursorio.TextLineColumn{1, 16}},
							Until: cursorio.TextOffset{Byte: 63, LineColumn: cursorio.TextLineColumn{1, 47}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 64, LineColumn: cursorio.TextLineColumn{1, 48}},
							Until: cursorio.TextOffset{Byte: 70, LineColumn: cursorio.TextLineColumn{1, 54}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("favorite-fruit"),
						Object:    xsdobject.String("orange"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{2, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 32, LineColumn: cursorio.TextLineColumn{1, 16}},
							Until: cursorio.TextOffset{Byte: 63, LineColumn: cursorio.TextLineColumn{1, 47}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 64, LineColumn: cursorio.TextLineColumn{1, 48}},
							Until: cursorio.TextOffset{Byte: 70, LineColumn: cursorio.TextLineColumn{1, 54}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/12a",
			Snippet: `<figure>
 <img src="castle.jpeg">
 <figcaption><span itemscope><span itemprop="name">The Castle</span></span> (1986)</figcaption>
</figure>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("The Castle"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 47, LineColumn: cursorio.TextLineColumn{2, 13}},
							Until: cursorio.TextOffset{Byte: 109, LineColumn: cursorio.TextLineColumn{2, 75}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{2, 44}},
							Until: cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{2, 50}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 85, LineColumn: cursorio.TextLineColumn{2, 51}},
							Until: cursorio.TextOffset{Byte: 95, LineColumn: cursorio.TextLineColumn{2, 61}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.2/12b",
			Snippet: `<span itemscope><meta itemprop="name" content="The Castle"></span>
<figure>
 <img src="castle.jpeg">
 <figcaption>The Castle (1986)</figcaption>
</figure>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("The Castle"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 66, LineColumn: cursorio.TextLineColumn{0, 66}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 31, LineColumn: cursorio.TextLineColumn{0, 31}},
							Until: cursorio.TextOffset{Byte: 37, LineColumn: cursorio.TextLineColumn{0, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 46, LineColumn: cursorio.TextLineColumn{0, 46}},
							Until: cursorio.TextOffset{Byte: 58, LineColumn: cursorio.TextLineColumn{0, 58}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.3",
			Snippet: `<section itemscope itemtype="https://example.org/animals#cat">
 <h1 itemprop="name">Hedral</h1>
 <p itemprop="desc">Hedral is a male american domestic
 shorthair, with a fluffy black fur with white paws and belly.</p>
 <img itemprop="img" src="hedral.jpeg" alt="" title="Hedral, age 18 months">
</section>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://example.org/animals#cat"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 305, LineColumn: cursorio.TextLineColumn{5, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 19, LineColumn: cursorio.TextLineColumn{0, 19}},
							Until: cursorio.TextOffset{Byte: 27, LineColumn: cursorio.TextLineColumn{0, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 28, LineColumn: cursorio.TextLineColumn{0, 28}},
							Until: cursorio.TextOffset{Byte: 61, LineColumn: cursorio.TextLineColumn{0, 61}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Hedral"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 305, LineColumn: cursorio.TextLineColumn{5, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 77, LineColumn: cursorio.TextLineColumn{1, 14}},
							Until: cursorio.TextOffset{Byte: 83, LineColumn: cursorio.TextLineColumn{1, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{1, 21}},
							Until: cursorio.TextOffset{Byte: 90, LineColumn: cursorio.TextLineColumn{1, 27}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("desc"),
						Object:    xsdobject.String("Hedral is a male american domestic\n shorthair, with a fluffy black fur with white paws and belly."),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 305, LineColumn: cursorio.TextLineColumn{5, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 109, LineColumn: cursorio.TextLineColumn{2, 13}},
							Until: cursorio.TextOffset{Byte: 115, LineColumn: cursorio.TextLineColumn{2, 19}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 116, LineColumn: cursorio.TextLineColumn{2, 20}},
							Until: cursorio.TextOffset{Byte: 213, LineColumn: cursorio.TextLineColumn{3, 62}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("img"),
						Object:    rdf.IRI("hedral.jpeg"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 305, LineColumn: cursorio.TextLineColumn{5, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 233, LineColumn: cursorio.TextLineColumn{4, 15}},
							Until: cursorio.TextOffset{Byte: 238, LineColumn: cursorio.TextLineColumn{4, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 243, LineColumn: cursorio.TextLineColumn{4, 25}},
							Until: cursorio.TextOffset{Byte: 256, LineColumn: cursorio.TextLineColumn{4, 38}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.4",
			Snippet: `<dl itemscope
    itemtype="https://vocab.example.net/book"
    itemid="urn:isbn:0-330-34032-8">
 <dt>Title
 <dd itemprop="title">The Reality Dysfunction
 <dt>Author
 <dd itemprop="author">Peter F. Hamilton
 <dt>Publication date
 <dd><time itemprop="pubdate" datetime="1996-01-26">26 January 1996</time>
</dl>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("urn:isbn:0-330-34032-8"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://vocab.example.net/book"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 71, LineColumn: cursorio.TextLineColumn{2, 11}},
							Until: cursorio.TextOffset{Byte: 95, LineColumn: cursorio.TextLineColumn{2, 35}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 18, LineColumn: cursorio.TextLineColumn{1, 4}},
							Until: cursorio.TextOffset{Byte: 26, LineColumn: cursorio.TextLineColumn{1, 12}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 27, LineColumn: cursorio.TextLineColumn{1, 13}},
							Until: cursorio.TextOffset{Byte: 59, LineColumn: cursorio.TextLineColumn{1, 45}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("urn:isbn:0-330-34032-8"),
						Predicate: rdf.IRI("title"),
						Object:    xsdobject.String("The Reality Dysfunction\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 71, LineColumn: cursorio.TextLineColumn{2, 11}},
							Until: cursorio.TextOffset{Byte: 95, LineColumn: cursorio.TextLineColumn{2, 35}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 122, LineColumn: cursorio.TextLineColumn{4, 14}},
							Until: cursorio.TextOffset{Byte: 129, LineColumn: cursorio.TextLineColumn{4, 21}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 130, LineColumn: cursorio.TextLineColumn{4, 22}},
							Until: cursorio.TextOffset{Byte: 155, LineColumn: cursorio.TextLineColumn{5, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("urn:isbn:0-330-34032-8"),
						Predicate: rdf.IRI("author"),
						Object:    xsdobject.String("Peter F. Hamilton\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 71, LineColumn: cursorio.TextLineColumn{2, 11}},
							Until: cursorio.TextOffset{Byte: 95, LineColumn: cursorio.TextLineColumn{2, 35}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 180, LineColumn: cursorio.TextLineColumn{6, 14}},
							Until: cursorio.TextOffset{Byte: 188, LineColumn: cursorio.TextLineColumn{6, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 189, LineColumn: cursorio.TextLineColumn{6, 23}},
							Until: cursorio.TextOffset{Byte: 208, LineColumn: cursorio.TextLineColumn{7, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   rdf.IRI("urn:isbn:0-330-34032-8"),
						Predicate: rdf.IRI("pubdate"),
						Object: rdf.Literal{
							LexicalForm: "1996-01-26",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#date"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 71, LineColumn: cursorio.TextLineColumn{2, 11}},
							Until: cursorio.TextOffset{Byte: 95, LineColumn: cursorio.TextLineColumn{2, 35}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 249, LineColumn: cursorio.TextLineColumn{8, 20}},
							Until: cursorio.TextOffset{Byte: 258, LineColumn: cursorio.TextLineColumn{8, 29}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 268, LineColumn: cursorio.TextLineColumn{8, 39}},
							Until: cursorio.TextOffset{Byte: 280, LineColumn: cursorio.TextLineColumn{8, 51}},
						},
					},
				},
			},
		},
		{
			Name: "5.1.5/2",
			Snippet: `<section itemscope itemtype="https://example.org/animals#cat">
 <h1 itemprop="name https://example.com/fn">Hedral</h1>
 <p itemprop="desc">Hedral is a male American domestic
 shorthair, with a fluffy <span
 itemprop="https://example.com/color">black</span> fur with <span
 itemprop="https://example.com/color">white</span> paws and belly.</p>
 <img itemprop="img" src="hedral.jpeg" alt="" title="Hedral, age 18 months">
</section>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://example.org/animals#cat"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 430, LineColumn: cursorio.TextLineColumn{7, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 19, LineColumn: cursorio.TextLineColumn{0, 19}},
							Until: cursorio.TextOffset{Byte: 27, LineColumn: cursorio.TextLineColumn{0, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 28, LineColumn: cursorio.TextLineColumn{0, 28}},
							Until: cursorio.TextOffset{Byte: 61, LineColumn: cursorio.TextLineColumn{0, 61}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Hedral"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 430, LineColumn: cursorio.TextLineColumn{7, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 77, LineColumn: cursorio.TextLineColumn{1, 14}},
							Until: cursorio.TextOffset{Byte: 106, LineColumn: cursorio.TextLineColumn{1, 43}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 107, LineColumn: cursorio.TextLineColumn{1, 44}},
							Until: cursorio.TextOffset{Byte: 113, LineColumn: cursorio.TextLineColumn{1, 50}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("https://example.com/fn"),
						Object:    xsdobject.String("Hedral"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 430, LineColumn: cursorio.TextLineColumn{7, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 77, LineColumn: cursorio.TextLineColumn{1, 14}},
							Until: cursorio.TextOffset{Byte: 106, LineColumn: cursorio.TextLineColumn{1, 43}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 107, LineColumn: cursorio.TextLineColumn{1, 44}},
							Until: cursorio.TextOffset{Byte: 113, LineColumn: cursorio.TextLineColumn{1, 50}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("desc"),
						Object:    xsdobject.String("Hedral is a male American domestic\n shorthair, with a fluffy black fur with white paws and belly."),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 430, LineColumn: cursorio.TextLineColumn{7, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 132, LineColumn: cursorio.TextLineColumn{2, 13}},
							Until: cursorio.TextOffset{Byte: 138, LineColumn: cursorio.TextLineColumn{2, 19}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 139, LineColumn: cursorio.TextLineColumn{2, 20}},
							Until: cursorio.TextOffset{Byte: 338, LineColumn: cursorio.TextLineColumn{5, 66}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("https://example.com/color"),
						Object:    xsdobject.String("black"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 430, LineColumn: cursorio.TextLineColumn{7, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 216, LineColumn: cursorio.TextLineColumn{4, 10}},
							Until: cursorio.TextOffset{Byte: 243, LineColumn: cursorio.TextLineColumn{4, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 244, LineColumn: cursorio.TextLineColumn{4, 38}},
							Until: cursorio.TextOffset{Byte: 249, LineColumn: cursorio.TextLineColumn{4, 43}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("https://example.com/color"),
						Object:    xsdobject.String("white"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 430, LineColumn: cursorio.TextLineColumn{7, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 282, LineColumn: cursorio.TextLineColumn{5, 10}},
							Until: cursorio.TextOffset{Byte: 309, LineColumn: cursorio.TextLineColumn{5, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 310, LineColumn: cursorio.TextLineColumn{5, 38}},
							Until: cursorio.TextOffset{Byte: 315, LineColumn: cursorio.TextLineColumn{5, 43}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("img"),
						Object:    rdf.IRI("hedral.jpeg"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 430, LineColumn: cursorio.TextLineColumn{7, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 358, LineColumn: cursorio.TextLineColumn{6, 15}},
							Until: cursorio.TextOffset{Byte: 363, LineColumn: cursorio.TextLineColumn{6, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 368, LineColumn: cursorio.TextLineColumn{6, 25}},
							Until: cursorio.TextOffset{Byte: 381, LineColumn: cursorio.TextLineColumn{6, 38}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.2/1a",
			Snippet: `<dl itemscope itemtype="https://md.example.com/loco
                        https://md.example.com/lighting">
 <dt>Name:
 <dd itemprop="name">Tank Locomotive (DB 80)
 <dt>Product code:
 <dd itemprop="product-code">33041
 <dt>Scale:
 <dd itemprop="scale">HO
 <dt>Digital:
 <dd itemprop="digital">Delta
</dl>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://md.example.com/loco"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 306, LineColumn: cursorio.TextLineColumn{10, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 14, LineColumn: cursorio.TextLineColumn{0, 14}},
							Until: cursorio.TextOffset{Byte: 22, LineColumn: cursorio.TextLineColumn{0, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{0, 23}},
							Until: cursorio.TextOffset{Byte: 108, LineColumn: cursorio.TextLineColumn{1, 56}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://md.example.com/lighting"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 306, LineColumn: cursorio.TextLineColumn{10, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 14, LineColumn: cursorio.TextLineColumn{0, 14}},
							Until: cursorio.TextOffset{Byte: 22, LineColumn: cursorio.TextLineColumn{0, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{0, 23}},
							Until: cursorio.TextOffset{Byte: 108, LineColumn: cursorio.TextLineColumn{1, 56}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Tank Locomotive (DB 80)\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 306, LineColumn: cursorio.TextLineColumn{10, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 135, LineColumn: cursorio.TextLineColumn{3, 14}},
							Until: cursorio.TextOffset{Byte: 141, LineColumn: cursorio.TextLineColumn{3, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 142, LineColumn: cursorio.TextLineColumn{3, 21}},
							Until: cursorio.TextOffset{Byte: 167, LineColumn: cursorio.TextLineColumn{4, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("product-code"),
						Object:    xsdobject.String("33041\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 306, LineColumn: cursorio.TextLineColumn{10, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 199, LineColumn: cursorio.TextLineColumn{5, 14}},
							Until: cursorio.TextOffset{Byte: 213, LineColumn: cursorio.TextLineColumn{5, 28}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 214, LineColumn: cursorio.TextLineColumn{5, 29}},
							Until: cursorio.TextOffset{Byte: 221, LineColumn: cursorio.TextLineColumn{6, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("scale"),
						Object:    xsdobject.String("HO\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 306, LineColumn: cursorio.TextLineColumn{10, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 246, LineColumn: cursorio.TextLineColumn{7, 14}},
							Until: cursorio.TextOffset{Byte: 253, LineColumn: cursorio.TextLineColumn{7, 21}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 254, LineColumn: cursorio.TextLineColumn{7, 22}},
							Until: cursorio.TextOffset{Byte: 258, LineColumn: cursorio.TextLineColumn{8, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("digital"),
						Object:    xsdobject.String("Delta\n"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 306, LineColumn: cursorio.TextLineColumn{10, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 285, LineColumn: cursorio.TextLineColumn{9, 14}},
							Until: cursorio.TextOffset{Byte: 294, LineColumn: cursorio.TextLineColumn{9, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 295, LineColumn: cursorio.TextLineColumn{9, 24}},
							Until: cursorio.TextOffset{Byte: 301, LineColumn: cursorio.TextLineColumn{10, 0}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.2/1b",
			Snippet: `<dl itemscope itemtype="https://md.example.com/track
                        https://md.example.com/lighting">
 <dt>Name:
 <dd itemprop="name">Turnout Lantern Kit
 <dt>Product code:
 <dd itemprop="product-code">74470
 <dt>Purpose:
 <dd>For retrofitting 2 <span itemprop="track-type">C</span> Track
 turnouts. <meta itemprop="scale" content="HO">
</dl>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://md.example.com/track"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 351, LineColumn: cursorio.TextLineColumn{9, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 14, LineColumn: cursorio.TextLineColumn{0, 14}},
							Until: cursorio.TextOffset{Byte: 22, LineColumn: cursorio.TextLineColumn{0, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{0, 23}},
							Until: cursorio.TextOffset{Byte: 109, LineColumn: cursorio.TextLineColumn{1, 56}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://md.example.com/lighting"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 351, LineColumn: cursorio.TextLineColumn{9, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 14, LineColumn: cursorio.TextLineColumn{0, 14}},
							Until: cursorio.TextOffset{Byte: 22, LineColumn: cursorio.TextLineColumn{0, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{0, 23}},
							Until: cursorio.TextOffset{Byte: 109, LineColumn: cursorio.TextLineColumn{1, 56}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Turnout Lantern Kit\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 351, LineColumn: cursorio.TextLineColumn{9, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 136, LineColumn: cursorio.TextLineColumn{3, 14}},
							Until: cursorio.TextOffset{Byte: 142, LineColumn: cursorio.TextLineColumn{3, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 143, LineColumn: cursorio.TextLineColumn{3, 21}},
							Until: cursorio.TextOffset{Byte: 164, LineColumn: cursorio.TextLineColumn{4, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("product-code"),
						Object:    xsdobject.String("74470\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 351, LineColumn: cursorio.TextLineColumn{9, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 196, LineColumn: cursorio.TextLineColumn{5, 14}},
							Until: cursorio.TextOffset{Byte: 210, LineColumn: cursorio.TextLineColumn{5, 28}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 211, LineColumn: cursorio.TextLineColumn{5, 29}},
							Until: cursorio.TextOffset{Byte: 218, LineColumn: cursorio.TextLineColumn{6, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("track-type"),
						Object:    xsdobject.String("C"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 351, LineColumn: cursorio.TextLineColumn{9, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 270, LineColumn: cursorio.TextLineColumn{7, 39}},
							Until: cursorio.TextOffset{Byte: 282, LineColumn: cursorio.TextLineColumn{7, 51}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 283, LineColumn: cursorio.TextLineColumn{7, 52}},
							Until: cursorio.TextOffset{Byte: 284, LineColumn: cursorio.TextLineColumn{7, 53}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("scale"),
						Object:    xsdobject.String("HO"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 351, LineColumn: cursorio.TextLineColumn{9, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 324, LineColumn: cursorio.TextLineColumn{8, 26}},
							Until: cursorio.TextOffset{Byte: 331, LineColumn: cursorio.TextLineColumn{8, 33}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 340, LineColumn: cursorio.TextLineColumn{8, 42}},
							Until: cursorio.TextOffset{Byte: 344, LineColumn: cursorio.TextLineColumn{8, 46}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.2/1c",
			Snippet: `<dl itemscope itemtype="https://md.example.com/passengers">
 <dt>Name:
 <dd itemprop="name">Express Train Passenger Car (DB Am 203)
 <dt>Product code:
 <dd itemprop="product-code">8710
 <dt>Scale:
 <dd itemprop="scale">Z
</dl>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("https://md.example.com/passengers"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 226, LineColumn: cursorio.TextLineColumn{7, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 14, LineColumn: cursorio.TextLineColumn{0, 14}},
							Until: cursorio.TextOffset{Byte: 22, LineColumn: cursorio.TextLineColumn{0, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{0, 23}},
							Until: cursorio.TextOffset{Byte: 58, LineColumn: cursorio.TextLineColumn{0, 58}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Express Train Passenger Car (DB Am 203)\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 226, LineColumn: cursorio.TextLineColumn{7, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 85, LineColumn: cursorio.TextLineColumn{2, 14}},
							Until: cursorio.TextOffset{Byte: 91, LineColumn: cursorio.TextLineColumn{2, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 92, LineColumn: cursorio.TextLineColumn{2, 21}},
							Until: cursorio.TextOffset{Byte: 133, LineColumn: cursorio.TextLineColumn{3, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("product-code"),
						Object:    xsdobject.String("8710\n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 226, LineColumn: cursorio.TextLineColumn{7, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 165, LineColumn: cursorio.TextLineColumn{4, 14}},
							Until: cursorio.TextOffset{Byte: 179, LineColumn: cursorio.TextLineColumn{4, 28}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 180, LineColumn: cursorio.TextLineColumn{4, 29}},
							Until: cursorio.TextOffset{Byte: 186, LineColumn: cursorio.TextLineColumn{5, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("scale"),
						Object:    xsdobject.String("Z\n"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 226, LineColumn: cursorio.TextLineColumn{7, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 211, LineColumn: cursorio.TextLineColumn{6, 14}},
							Until: cursorio.TextOffset{Byte: 218, LineColumn: cursorio.TextLineColumn{6, 21}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 219, LineColumn: cursorio.TextLineColumn{6, 22}},
							Until: cursorio.TextOffset{Byte: 221, LineColumn: cursorio.TextLineColumn{7, 0}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.3/1a",
			Snippet: `<div itemscope>
 <p itemprop="a">1</p>
 <p itemprop="a">2</p>
 <p itemprop="b">test</p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 29, LineColumn: cursorio.TextLineColumn{1, 13}},
							Until: cursorio.TextOffset{Byte: 32, LineColumn: cursorio.TextLineColumn{1, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 33, LineColumn: cursorio.TextLineColumn{1, 17}},
							Until: cursorio.TextOffset{Byte: 34, LineColumn: cursorio.TextLineColumn{1, 18}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 52, LineColumn: cursorio.TextLineColumn{2, 13}},
							Until: cursorio.TextOffset{Byte: 55, LineColumn: cursorio.TextLineColumn{2, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{2, 17}},
							Until: cursorio.TextOffset{Byte: 57, LineColumn: cursorio.TextLineColumn{2, 18}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("b"),
						Object:    xsdobject.String("test"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 75, LineColumn: cursorio.TextLineColumn{3, 13}},
							Until: cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{3, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 79, LineColumn: cursorio.TextLineColumn{3, 17}},
							Until: cursorio.TextOffset{Byte: 83, LineColumn: cursorio.TextLineColumn{3, 21}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.3/1b",
			Snippet: `<div itemscope>
 <p itemprop="b">test</p>
 <p itemprop="a">1</p>
 <p itemprop="a">2</p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("b"),
						Object:    xsdobject.String("test"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 29, LineColumn: cursorio.TextLineColumn{1, 13}},
							Until: cursorio.TextOffset{Byte: 32, LineColumn: cursorio.TextLineColumn{1, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 33, LineColumn: cursorio.TextLineColumn{1, 17}},
							Until: cursorio.TextOffset{Byte: 37, LineColumn: cursorio.TextLineColumn{1, 21}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 55, LineColumn: cursorio.TextLineColumn{2, 13}},
							Until: cursorio.TextOffset{Byte: 58, LineColumn: cursorio.TextLineColumn{2, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 59, LineColumn: cursorio.TextLineColumn{2, 17}},
							Until: cursorio.TextOffset{Byte: 60, LineColumn: cursorio.TextLineColumn{2, 18}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{3, 13}},
							Until: cursorio.TextOffset{Byte: 81, LineColumn: cursorio.TextLineColumn{3, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 82, LineColumn: cursorio.TextLineColumn{3, 17}},
							Until: cursorio.TextOffset{Byte: 83, LineColumn: cursorio.TextLineColumn{3, 18}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.3/1c",
			Snippet: `<div itemscope>
 <p itemprop="a">1</p>
 <p itemprop="b">test</p>
 <p itemprop="a">2</p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 29, LineColumn: cursorio.TextLineColumn{1, 13}},
							Until: cursorio.TextOffset{Byte: 32, LineColumn: cursorio.TextLineColumn{1, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 33, LineColumn: cursorio.TextLineColumn{1, 17}},
							Until: cursorio.TextOffset{Byte: 34, LineColumn: cursorio.TextLineColumn{1, 18}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("b"),
						Object:    xsdobject.String("test"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 52, LineColumn: cursorio.TextLineColumn{2, 13}},
							Until: cursorio.TextOffset{Byte: 55, LineColumn: cursorio.TextLineColumn{2, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{2, 17}},
							Until: cursorio.TextOffset{Byte: 60, LineColumn: cursorio.TextLineColumn{2, 21}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{4, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 78, LineColumn: cursorio.TextLineColumn{3, 13}},
							Until: cursorio.TextOffset{Byte: 81, LineColumn: cursorio.TextLineColumn{3, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 82, LineColumn: cursorio.TextLineColumn{3, 17}},
							Until: cursorio.TextOffset{Byte: 83, LineColumn: cursorio.TextLineColumn{3, 18}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.3/1d",
			Snippet: `<div id="x">
 <p itemprop="a">1</p>
</div>
<div itemscope itemref="x">
 <p itemprop="b">test</p>
 <p itemprop="a">2</p>
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 43, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 126, LineColumn: cursorio.TextLineColumn{6, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 26, LineColumn: cursorio.TextLineColumn{1, 13}},
							Until: cursorio.TextOffset{Byte: 29, LineColumn: cursorio.TextLineColumn{1, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 30, LineColumn: cursorio.TextLineColumn{1, 17}},
							Until: cursorio.TextOffset{Byte: 31, LineColumn: cursorio.TextLineColumn{1, 18}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("b"),
						Object:    xsdobject.String("test"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 43, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 126, LineColumn: cursorio.TextLineColumn{6, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{4, 13}},
							Until: cursorio.TextOffset{Byte: 87, LineColumn: cursorio.TextLineColumn{4, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 88, LineColumn: cursorio.TextLineColumn{4, 17}},
							Until: cursorio.TextOffset{Byte: 92, LineColumn: cursorio.TextLineColumn{4, 21}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("a"),
						Object:    xsdobject.String("2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 43, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 126, LineColumn: cursorio.TextLineColumn{6, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 110, LineColumn: cursorio.TextLineColumn{5, 13}},
							Until: cursorio.TextOffset{Byte: 113, LineColumn: cursorio.TextLineColumn{5, 16}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 114, LineColumn: cursorio.TextLineColumn{5, 17}},
							Until: cursorio.TextOffset{Byte: 115, LineColumn: cursorio.TextLineColumn{5, 18}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.5/1",
			Snippet: `<!DOCTYPE HTML>
<html lang="en">
 <head>
  <title>Photo gallery</title>
 </head>
 <body>
  <h1>My photos</h1>
  <figure itemscope itemtype="http://n.whatwg.org/work" itemref="licenses">
   <img itemprop="work" src="images/house.jpeg" alt="A white house, boarded up, sits in a forest.">
   <figcaption itemprop="title">The house I found.</figcaption>
  </figure>
  <figure itemscope itemtype="http://n.whatwg.org/work" itemref="licenses">
   <img itemprop="work" src="images/mailbox.jpeg" alt="Outside the house is a mailbox. It has a leaflet inside.">
   <figcaption itemprop="title">The mailbox.</figcaption>
  </figure>
  <footer>
   <p id="licenses">All images licensed under the <a itemprop="license"
   href="http://www.opensource.org/licenses/mit-license.php">MIT
   license</a>.</p>
  </footer>
 </body>
</html>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://n.whatwg.org/work"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 112, LineColumn: cursorio.TextLineColumn{7, 2}},
							Until: cursorio.TextOffset{Byte: 361, LineColumn: cursorio.TextLineColumn{10, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 130, LineColumn: cursorio.TextLineColumn{7, 20}},
							Until: cursorio.TextOffset{Byte: 138, LineColumn: cursorio.TextLineColumn{7, 28}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 139, LineColumn: cursorio.TextLineColumn{7, 29}},
							Until: cursorio.TextOffset{Byte: 165, LineColumn: cursorio.TextLineColumn{7, 55}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("license"),
						Object:    rdf.IRI("http://www.opensource.org/licenses/mit-license.php"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 112, LineColumn: cursorio.TextLineColumn{7, 2}},
							Until: cursorio.TextOffset{Byte: 361, LineColumn: cursorio.TextLineColumn{10, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 695, LineColumn: cursorio.TextLineColumn{16, 62}},
							Until: cursorio.TextOffset{Byte: 704, LineColumn: cursorio.TextLineColumn{16, 71}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 713, LineColumn: cursorio.TextLineColumn{17, 8}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{17, 60}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("work"),
						Object:    rdf.IRI("images/house.jpeg"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 112, LineColumn: cursorio.TextLineColumn{7, 2}},
							Until: cursorio.TextOffset{Byte: 361, LineColumn: cursorio.TextLineColumn{10, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 203, LineColumn: cursorio.TextLineColumn{8, 17}},
							Until: cursorio.TextOffset{Byte: 209, LineColumn: cursorio.TextLineColumn{8, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 214, LineColumn: cursorio.TextLineColumn{8, 28}},
							Until: cursorio.TextOffset{Byte: 233, LineColumn: cursorio.TextLineColumn{8, 47}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("title"),
						Object:    xsdobject.String("The house I found."),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 112, LineColumn: cursorio.TextLineColumn{7, 2}},
							Until: cursorio.TextOffset{Byte: 361, LineColumn: cursorio.TextLineColumn{10, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 310, LineColumn: cursorio.TextLineColumn{9, 24}},
							Until: cursorio.TextOffset{Byte: 317, LineColumn: cursorio.TextLineColumn{9, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 318, LineColumn: cursorio.TextLineColumn{9, 32}},
							Until: cursorio.TextOffset{Byte: 336, LineColumn: cursorio.TextLineColumn{9, 50}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://n.whatwg.org/work"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 364, LineColumn: cursorio.TextLineColumn{11, 2}},
							Until: cursorio.TextOffset{Byte: 621, LineColumn: cursorio.TextLineColumn{14, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 382, LineColumn: cursorio.TextLineColumn{11, 20}},
							Until: cursorio.TextOffset{Byte: 390, LineColumn: cursorio.TextLineColumn{11, 28}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 391, LineColumn: cursorio.TextLineColumn{11, 29}},
							Until: cursorio.TextOffset{Byte: 417, LineColumn: cursorio.TextLineColumn{11, 55}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("license"),
						Object:    rdf.IRI("http://www.opensource.org/licenses/mit-license.php"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 364, LineColumn: cursorio.TextLineColumn{11, 2}},
							Until: cursorio.TextOffset{Byte: 621, LineColumn: cursorio.TextLineColumn{14, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 695, LineColumn: cursorio.TextLineColumn{16, 62}},
							Until: cursorio.TextOffset{Byte: 704, LineColumn: cursorio.TextLineColumn{16, 71}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 713, LineColumn: cursorio.TextLineColumn{17, 8}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{17, 60}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("work"),
						Object:    rdf.IRI("images/mailbox.jpeg"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 364, LineColumn: cursorio.TextLineColumn{11, 2}},
							Until: cursorio.TextOffset{Byte: 621, LineColumn: cursorio.TextLineColumn{14, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 455, LineColumn: cursorio.TextLineColumn{12, 17}},
							Until: cursorio.TextOffset{Byte: 461, LineColumn: cursorio.TextLineColumn{12, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 466, LineColumn: cursorio.TextLineColumn{12, 28}},
							Until: cursorio.TextOffset{Byte: 487, LineColumn: cursorio.TextLineColumn{12, 49}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("title"),
						Object:    xsdobject.String("The mailbox."),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 364, LineColumn: cursorio.TextLineColumn{11, 2}},
							Until: cursorio.TextOffset{Byte: 621, LineColumn: cursorio.TextLineColumn{14, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 576, LineColumn: cursorio.TextLineColumn{13, 24}},
							Until: cursorio.TextOffset{Byte: 583, LineColumn: cursorio.TextLineColumn{13, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 584, LineColumn: cursorio.TextLineColumn{13, 32}},
							Until: cursorio.TextOffset{Byte: 596, LineColumn: cursorio.TextLineColumn{13, 44}},
						},
					},
				},
			},
		},
		{
			Name: "5.2.6/1",
			Snippet: `<p itemscope></p> <!-- this is an item (with no properties and no type) -->
<svg itemscope></svg> <!-- this is not, it's just an SVG svg element with an invalid unknown attribute -->`,
			Expected: encodingtest.TripleStatementList{
				// TODO possible to differentiate property-less item?
			},
		},
		{
			Name: "5.3.1.2/1",
			Snippet: `<section id="jack" itemscope itemtype="http://microformats.org/profile/hcard">
 <h1 itemprop="fn">
  <span itemprop="n" itemscope>
   <span itemprop="given-name">Jack</span>
   <span itemprop="family-name">Bauer</span>
  </span>
 </h1>
 <img itemprop="photo" alt="" src="jack-bauer.jpg">
 <p itemprop="org" itemscope>
  <span itemprop="organization-name">Counter-Terrorist Unit</span>
  (<span itemprop="organization-unit">Los Angeles Division</span>)
 </p>
 <p>
  <span itemprop="adr" itemscope>
   <span itemprop="street-address">10201 W. Pico Blvd.</span><br>
   <span itemprop="locality">Los Angeles</span>,
   <span itemprop="region">CA</span>
   <span itemprop="postal-code">90064</span><br>
   <span itemprop="country-name">United States</span><br>
  </span>
  <span itemprop="geo">34.052339;-118.410623</span>
 </p>
 <h2>Assorted Contact Methods</h2>
 <ul>
  <li itemprop="tel" itemscope>
   <span itemprop="value">+1 (310) 597 3781</span> <span itemprop="type">work</span>
   <meta itemprop="type" content="voice">
  </li>
  <li><a itemprop="url" href="https://en.wikipedia.org/wiki/Jack_Bauer">I'm on Wikipedia</a>
  so you can leave a message on my user talk page.</li>
  <li><a itemprop="url" href="http://www.jackbauerfacts.com/">Jack Bauer Facts</a></li>
  <li itemprop="email"><a href="mailto:j.bauer@la.ctu.gov.invalid">j.bauer@la.ctu.gov.invalid</a></li>
  <li itemprop="tel" itemscope>
   <span itemprop="value">+1 (310) 555 3781</span> <span>
   <meta itemprop="type" content="cell">mobile phone</span>
  </li>
 </ul>
 <ins datetime="2008-07-20 21:00:00+01:00">
  <meta itemprop="rev" content="2008-07-20 21:00:00+01:00">
  <p itemprop="tel" itemscope><strong>Update!</strong>
  My new <span itemprop="type">home</span> phone number is
  <span itemprop="value">01632 960 123</span>.</p>
 </ins>
</section>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://microformats.org/profile/hcard"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 29, LineColumn: cursorio.TextLineColumn{0, 29}},
							Until: cursorio.TextOffset{Byte: 37, LineColumn: cursorio.TextLineColumn{0, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 38, LineColumn: cursorio.TextLineColumn{0, 38}},
							Until: cursorio.TextOffset{Byte: 77, LineColumn: cursorio.TextLineColumn{0, 77}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("fn"),
						Object:    xsdobject.String("\n  \n   Jack\n   Bauer\n  \n "),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 93, LineColumn: cursorio.TextLineColumn{1, 14}},
							Until: cursorio.TextOffset{Byte: 97, LineColumn: cursorio.TextLineColumn{1, 18}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 98, LineColumn: cursorio.TextLineColumn{1, 19}},
							Until: cursorio.TextOffset{Byte: 230, LineColumn: cursorio.TextLineColumn{6, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("n"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 116, LineColumn: cursorio.TextLineColumn{2, 17}},
							Until: cursorio.TextOffset{Byte: 119, LineColumn: cursorio.TextLineColumn{2, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 101, LineColumn: cursorio.TextLineColumn{2, 2}},
							Until: cursorio.TextOffset{Byte: 228, LineColumn: cursorio.TextLineColumn{5, 9}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("given-name"),
						Object:    xsdobject.String("Jack"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 101, LineColumn: cursorio.TextLineColumn{2, 2}},
							Until: cursorio.TextOffset{Byte: 228, LineColumn: cursorio.TextLineColumn{5, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 149, LineColumn: cursorio.TextLineColumn{3, 18}},
							Until: cursorio.TextOffset{Byte: 161, LineColumn: cursorio.TextLineColumn{3, 30}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 162, LineColumn: cursorio.TextLineColumn{3, 31}},
							Until: cursorio.TextOffset{Byte: 166, LineColumn: cursorio.TextLineColumn{3, 35}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("family-name"),
						Object:    xsdobject.String("Bauer"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 101, LineColumn: cursorio.TextLineColumn{2, 2}},
							Until: cursorio.TextOffset{Byte: 228, LineColumn: cursorio.TextLineColumn{5, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 192, LineColumn: cursorio.TextLineColumn{4, 18}},
							Until: cursorio.TextOffset{Byte: 205, LineColumn: cursorio.TextLineColumn{4, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 206, LineColumn: cursorio.TextLineColumn{4, 32}},
							Until: cursorio.TextOffset{Byte: 211, LineColumn: cursorio.TextLineColumn{4, 37}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("photo"),
						Object:    rdf.IRI("jack-bauer.jpg"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 251, LineColumn: cursorio.TextLineColumn{7, 15}},
							Until: cursorio.TextOffset{Byte: 258, LineColumn: cursorio.TextLineColumn{7, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 270, LineColumn: cursorio.TextLineColumn{7, 34}},
							Until: cursorio.TextOffset{Byte: 286, LineColumn: cursorio.TextLineColumn{7, 50}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("org"),
						Object:    testingBnode.MapBlankNodeIdentifier("b2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 301, LineColumn: cursorio.TextLineColumn{8, 13}},
							Until: cursorio.TextOffset{Byte: 306, LineColumn: cursorio.TextLineColumn{8, 18}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 289, LineColumn: cursorio.TextLineColumn{8, 1}},
							Until: cursorio.TextOffset{Byte: 457, LineColumn: cursorio.TextLineColumn{11, 5}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("organization-name"),
						Object:    xsdobject.String("Counter-Terrorist Unit"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 289, LineColumn: cursorio.TextLineColumn{8, 1}},
							Until: cursorio.TextOffset{Byte: 457, LineColumn: cursorio.TextLineColumn{11, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 335, LineColumn: cursorio.TextLineColumn{9, 17}},
							Until: cursorio.TextOffset{Byte: 354, LineColumn: cursorio.TextLineColumn{9, 36}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 355, LineColumn: cursorio.TextLineColumn{9, 37}},
							Until: cursorio.TextOffset{Byte: 377, LineColumn: cursorio.TextLineColumn{9, 59}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("organization-unit"),
						Object:    xsdobject.String("Los Angeles Division"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 289, LineColumn: cursorio.TextLineColumn{8, 1}},
							Until: cursorio.TextOffset{Byte: 457, LineColumn: cursorio.TextLineColumn{11, 5}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 403, LineColumn: cursorio.TextLineColumn{10, 18}},
							Until: cursorio.TextOffset{Byte: 422, LineColumn: cursorio.TextLineColumn{10, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 423, LineColumn: cursorio.TextLineColumn{10, 38}},
							Until: cursorio.TextOffset{Byte: 443, LineColumn: cursorio.TextLineColumn{10, 58}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("adr"),
						Object:    testingBnode.MapBlankNodeIdentifier("b3"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 480, LineColumn: cursorio.TextLineColumn{13, 17}},
							Until: cursorio.TextOffset{Byte: 485, LineColumn: cursorio.TextLineColumn{13, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 465, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{19, 9}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("street-address"),
						Object:    xsdobject.String("10201 W. Pico Blvd."),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 465, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{19, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 515, LineColumn: cursorio.TextLineColumn{14, 18}},
							Until: cursorio.TextOffset{Byte: 531, LineColumn: cursorio.TextLineColumn{14, 34}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 532, LineColumn: cursorio.TextLineColumn{14, 35}},
							Until: cursorio.TextOffset{Byte: 551, LineColumn: cursorio.TextLineColumn{14, 54}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("locality"),
						Object:    xsdobject.String("Los Angeles"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 465, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{19, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 581, LineColumn: cursorio.TextLineColumn{15, 18}},
							Until: cursorio.TextOffset{Byte: 591, LineColumn: cursorio.TextLineColumn{15, 28}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 592, LineColumn: cursorio.TextLineColumn{15, 29}},
							Until: cursorio.TextOffset{Byte: 603, LineColumn: cursorio.TextLineColumn{15, 40}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("region"),
						Object:    xsdobject.String("CA"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 465, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{19, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 630, LineColumn: cursorio.TextLineColumn{16, 18}},
							Until: cursorio.TextOffset{Byte: 638, LineColumn: cursorio.TextLineColumn{16, 26}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 639, LineColumn: cursorio.TextLineColumn{16, 27}},
							Until: cursorio.TextOffset{Byte: 641, LineColumn: cursorio.TextLineColumn{16, 29}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("postal-code"),
						Object:    xsdobject.String("90064"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 465, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{19, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 667, LineColumn: cursorio.TextLineColumn{17, 18}},
							Until: cursorio.TextOffset{Byte: 680, LineColumn: cursorio.TextLineColumn{17, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{17, 32}},
							Until: cursorio.TextOffset{Byte: 686, LineColumn: cursorio.TextLineColumn{17, 37}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("country-name"),
						Object:    xsdobject.String("United States"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 465, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 765, LineColumn: cursorio.TextLineColumn{19, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 716, LineColumn: cursorio.TextLineColumn{18, 18}},
							Until: cursorio.TextOffset{Byte: 730, LineColumn: cursorio.TextLineColumn{18, 32}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 731, LineColumn: cursorio.TextLineColumn{18, 33}},
							Until: cursorio.TextOffset{Byte: 744, LineColumn: cursorio.TextLineColumn{18, 46}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("geo"),
						Object:    xsdobject.String("34.052339;-118.410623"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 783, LineColumn: cursorio.TextLineColumn{20, 17}},
							Until: cursorio.TextOffset{Byte: 788, LineColumn: cursorio.TextLineColumn{20, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 789, LineColumn: cursorio.TextLineColumn{20, 23}},
							Until: cursorio.TextOffset{Byte: 810, LineColumn: cursorio.TextLineColumn{20, 44}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("tel"),
						Object:    testingBnode.MapBlankNodeIdentifier("b4"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 880, LineColumn: cursorio.TextLineColumn{24, 15}},
							Until: cursorio.TextOffset{Byte: 885, LineColumn: cursorio.TextLineColumn{24, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 867, LineColumn: cursorio.TextLineColumn{24, 2}},
							Until: cursorio.TextOffset{Byte: 1031, LineColumn: cursorio.TextLineColumn{27, 7}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b4"),
						Predicate: rdf.IRI("value"),
						Object:    xsdobject.String("+1 (310) 597 3781"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 867, LineColumn: cursorio.TextLineColumn{24, 2}},
							Until: cursorio.TextOffset{Byte: 1031, LineColumn: cursorio.TextLineColumn{27, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 915, LineColumn: cursorio.TextLineColumn{25, 18}},
							Until: cursorio.TextOffset{Byte: 922, LineColumn: cursorio.TextLineColumn{25, 25}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 923, LineColumn: cursorio.TextLineColumn{25, 26}},
							Until: cursorio.TextOffset{Byte: 940, LineColumn: cursorio.TextLineColumn{25, 43}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b4"),
						Predicate: rdf.IRI("type"),
						Object:    xsdobject.String("work"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 867, LineColumn: cursorio.TextLineColumn{24, 2}},
							Until: cursorio.TextOffset{Byte: 1031, LineColumn: cursorio.TextLineColumn{27, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 963, LineColumn: cursorio.TextLineColumn{25, 66}},
							Until: cursorio.TextOffset{Byte: 969, LineColumn: cursorio.TextLineColumn{25, 72}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 970, LineColumn: cursorio.TextLineColumn{25, 73}},
							Until: cursorio.TextOffset{Byte: 974, LineColumn: cursorio.TextLineColumn{25, 77}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b4"),
						Predicate: rdf.IRI("type"),
						Object:    xsdobject.String("voice"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 867, LineColumn: cursorio.TextLineColumn{24, 2}},
							Until: cursorio.TextOffset{Byte: 1031, LineColumn: cursorio.TextLineColumn{27, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1000, LineColumn: cursorio.TextLineColumn{26, 18}},
							Until: cursorio.TextOffset{Byte: 1006, LineColumn: cursorio.TextLineColumn{26, 24}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1015, LineColumn: cursorio.TextLineColumn{26, 33}},
							Until: cursorio.TextOffset{Byte: 1022, LineColumn: cursorio.TextLineColumn{26, 40}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("url"),
						Object:    rdf.IRI("https://en.wikipedia.org/wiki/Jack_Bauer"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1050, LineColumn: cursorio.TextLineColumn{28, 18}},
							Until: cursorio.TextOffset{Byte: 1055, LineColumn: cursorio.TextLineColumn{28, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1061, LineColumn: cursorio.TextLineColumn{28, 29}},
							Until: cursorio.TextOffset{Byte: 1103, LineColumn: cursorio.TextLineColumn{28, 71}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("url"),
						Object:    rdf.IRI("http://www.jackbauerfacts.com/"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1199, LineColumn: cursorio.TextLineColumn{30, 18}},
							Until: cursorio.TextOffset{Byte: 1204, LineColumn: cursorio.TextLineColumn{30, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1210, LineColumn: cursorio.TextLineColumn{30, 29}},
							Until: cursorio.TextOffset{Byte: 1242, LineColumn: cursorio.TextLineColumn{30, 61}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("email"),
						Object:    xsdobject.String("j.bauer@la.ctu.gov.invalid"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1284, LineColumn: cursorio.TextLineColumn{31, 15}},
							Until: cursorio.TextOffset{Byte: 1291, LineColumn: cursorio.TextLineColumn{31, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1292, LineColumn: cursorio.TextLineColumn{31, 23}},
							Until: cursorio.TextOffset{Byte: 1366, LineColumn: cursorio.TextLineColumn{31, 97}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("tel"),
						Object:    testingBnode.MapBlankNodeIdentifier("b5"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1387, LineColumn: cursorio.TextLineColumn{32, 15}},
							Until: cursorio.TextOffset{Byte: 1392, LineColumn: cursorio.TextLineColumn{32, 20}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1374, LineColumn: cursorio.TextLineColumn{32, 2}},
							Until: cursorio.TextOffset{Byte: 1529, LineColumn: cursorio.TextLineColumn{35, 7}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b5"),
						Predicate: rdf.IRI("value"),
						Object:    xsdobject.String("+1 (310) 555 3781"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1374, LineColumn: cursorio.TextLineColumn{32, 2}},
							Until: cursorio.TextOffset{Byte: 1529, LineColumn: cursorio.TextLineColumn{35, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1422, LineColumn: cursorio.TextLineColumn{33, 18}},
							Until: cursorio.TextOffset{Byte: 1429, LineColumn: cursorio.TextLineColumn{33, 25}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1430, LineColumn: cursorio.TextLineColumn{33, 26}},
							Until: cursorio.TextOffset{Byte: 1447, LineColumn: cursorio.TextLineColumn{33, 43}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b5"),
						Predicate: rdf.IRI("type"),
						Object:    xsdobject.String("cell"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1374, LineColumn: cursorio.TextLineColumn{32, 2}},
							Until: cursorio.TextOffset{Byte: 1529, LineColumn: cursorio.TextLineColumn{35, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1480, LineColumn: cursorio.TextLineColumn{34, 18}},
							Until: cursorio.TextOffset{Byte: 1486, LineColumn: cursorio.TextLineColumn{34, 24}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1495, LineColumn: cursorio.TextLineColumn{34, 33}},
							Until: cursorio.TextOffset{Byte: 1501, LineColumn: cursorio.TextLineColumn{34, 39}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("rev"),
						Object:    xsdobject.String("2008-07-20 21:00:00+01:00"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1598, LineColumn: cursorio.TextLineColumn{38, 17}},
							Until: cursorio.TextOffset{Byte: 1603, LineColumn: cursorio.TextLineColumn{38, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1612, LineColumn: cursorio.TextLineColumn{38, 31}},
							Until: cursorio.TextOffset{Byte: 1639, LineColumn: cursorio.TextLineColumn{38, 58}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("tel"),
						Object:    testingBnode.MapBlankNodeIdentifier("b6"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 1824, LineColumn: cursorio.TextLineColumn{43, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1655, LineColumn: cursorio.TextLineColumn{39, 14}},
							Until: cursorio.TextOffset{Byte: 1660, LineColumn: cursorio.TextLineColumn{39, 19}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1643, LineColumn: cursorio.TextLineColumn{39, 2}},
							Until: cursorio.TextOffset{Byte: 1805, LineColumn: cursorio.TextLineColumn{41, 50}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b6"),
						Predicate: rdf.IRI("type"),
						Object:    xsdobject.String("home"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1643, LineColumn: cursorio.TextLineColumn{39, 2}},
							Until: cursorio.TextOffset{Byte: 1805, LineColumn: cursorio.TextLineColumn{41, 50}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1720, LineColumn: cursorio.TextLineColumn{40, 24}},
							Until: cursorio.TextOffset{Byte: 1726, LineColumn: cursorio.TextLineColumn{40, 30}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1727, LineColumn: cursorio.TextLineColumn{40, 31}},
							Until: cursorio.TextOffset{Byte: 1731, LineColumn: cursorio.TextLineColumn{40, 35}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b6"),
						Predicate: rdf.IRI("value"),
						Object:    xsdobject.String("01632 960 123"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1643, LineColumn: cursorio.TextLineColumn{39, 2}},
							Until: cursorio.TextOffset{Byte: 1805, LineColumn: cursorio.TextLineColumn{41, 50}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1772, LineColumn: cursorio.TextLineColumn{41, 17}},
							Until: cursorio.TextOffset{Byte: 1779, LineColumn: cursorio.TextLineColumn{41, 24}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1780, LineColumn: cursorio.TextLineColumn{41, 25}},
							Until: cursorio.TextOffset{Byte: 1793, LineColumn: cursorio.TextLineColumn{41, 38}},
						},
					},
				},
			},
		},
		{
			Name: "5.3.1.2/2",
			Snippet: `<address itemscope itemtype="http://microformats.org/profile/hcard">
 <strong itemprop="fn"><span itemprop="n" itemscope><span itemprop="given-name">Alfred</span>
 <span itemprop="family-name">Person</span></span></strong> <br>
 <span itemprop="adr" itemscope>
  <span itemprop="street-address">1600 Amphitheatre Parkway</span> <br>
  <span itemprop="street-address">Building 43, Second Floor</span> <br>
  <span itemprop="locality">Mountain View</span>,
   <span itemprop="region">CA</span> <span itemprop="postal-code">94043</span>
 </span>
</address>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://microformats.org/profile/hcard"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 553, LineColumn: cursorio.TextLineColumn{9, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 19, LineColumn: cursorio.TextLineColumn{0, 19}},
							Until: cursorio.TextOffset{Byte: 27, LineColumn: cursorio.TextLineColumn{0, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 28, LineColumn: cursorio.TextLineColumn{0, 28}},
							Until: cursorio.TextOffset{Byte: 67, LineColumn: cursorio.TextLineColumn{0, 67}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("fn"),
						Object:    xsdobject.String("Alfred\n Person"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 553, LineColumn: cursorio.TextLineColumn{9, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 87, LineColumn: cursorio.TextLineColumn{1, 18}},
							Until: cursorio.TextOffset{Byte: 91, LineColumn: cursorio.TextLineColumn{1, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 92, LineColumn: cursorio.TextLineColumn{1, 23}},
							Until: cursorio.TextOffset{Byte: 213, LineColumn: cursorio.TextLineColumn{2, 50}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("n"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 553, LineColumn: cursorio.TextLineColumn{9, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 107, LineColumn: cursorio.TextLineColumn{1, 38}},
							Until: cursorio.TextOffset{Byte: 110, LineColumn: cursorio.TextLineColumn{1, 41}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 92, LineColumn: cursorio.TextLineColumn{1, 23}},
							Until: cursorio.TextOffset{Byte: 213, LineColumn: cursorio.TextLineColumn{2, 50}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("given-name"),
						Object:    xsdobject.String("Alfred"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 92, LineColumn: cursorio.TextLineColumn{1, 23}},
							Until: cursorio.TextOffset{Byte: 213, LineColumn: cursorio.TextLineColumn{2, 50}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 136, LineColumn: cursorio.TextLineColumn{1, 67}},
							Until: cursorio.TextOffset{Byte: 148, LineColumn: cursorio.TextLineColumn{1, 79}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 149, LineColumn: cursorio.TextLineColumn{1, 80}},
							Until: cursorio.TextOffset{Byte: 155, LineColumn: cursorio.TextLineColumn{1, 86}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("family-name"),
						Object:    xsdobject.String("Person"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 92, LineColumn: cursorio.TextLineColumn{1, 23}},
							Until: cursorio.TextOffset{Byte: 213, LineColumn: cursorio.TextLineColumn{2, 50}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 179, LineColumn: cursorio.TextLineColumn{2, 16}},
							Until: cursorio.TextOffset{Byte: 192, LineColumn: cursorio.TextLineColumn{2, 29}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 193, LineColumn: cursorio.TextLineColumn{2, 30}},
							Until: cursorio.TextOffset{Byte: 199, LineColumn: cursorio.TextLineColumn{2, 36}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("adr"),
						Object:    testingBnode.MapBlankNodeIdentifier("b2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 553, LineColumn: cursorio.TextLineColumn{9, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 244, LineColumn: cursorio.TextLineColumn{3, 16}},
							Until: cursorio.TextOffset{Byte: 249, LineColumn: cursorio.TextLineColumn{3, 21}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 229, LineColumn: cursorio.TextLineColumn{3, 1}},
							Until: cursorio.TextOffset{Byte: 542, LineColumn: cursorio.TextLineColumn{8, 8}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("street-address"),
						Object:    xsdobject.String("1600 Amphitheatre Parkway"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 229, LineColumn: cursorio.TextLineColumn{3, 1}},
							Until: cursorio.TextOffset{Byte: 542, LineColumn: cursorio.TextLineColumn{8, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 278, LineColumn: cursorio.TextLineColumn{4, 17}},
							Until: cursorio.TextOffset{Byte: 294, LineColumn: cursorio.TextLineColumn{4, 33}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 295, LineColumn: cursorio.TextLineColumn{4, 34}},
							Until: cursorio.TextOffset{Byte: 320, LineColumn: cursorio.TextLineColumn{4, 59}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("street-address"),
						Object:    xsdobject.String("Building 43, Second Floor"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 229, LineColumn: cursorio.TextLineColumn{3, 1}},
							Until: cursorio.TextOffset{Byte: 542, LineColumn: cursorio.TextLineColumn{8, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 350, LineColumn: cursorio.TextLineColumn{5, 17}},
							Until: cursorio.TextOffset{Byte: 366, LineColumn: cursorio.TextLineColumn{5, 33}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 367, LineColumn: cursorio.TextLineColumn{5, 34}},
							Until: cursorio.TextOffset{Byte: 392, LineColumn: cursorio.TextLineColumn{5, 59}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("locality"),
						Object:    xsdobject.String("Mountain View"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 229, LineColumn: cursorio.TextLineColumn{3, 1}},
							Until: cursorio.TextOffset{Byte: 542, LineColumn: cursorio.TextLineColumn{8, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 422, LineColumn: cursorio.TextLineColumn{6, 17}},
							Until: cursorio.TextOffset{Byte: 432, LineColumn: cursorio.TextLineColumn{6, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 433, LineColumn: cursorio.TextLineColumn{6, 28}},
							Until: cursorio.TextOffset{Byte: 446, LineColumn: cursorio.TextLineColumn{6, 41}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("region"),
						Object:    xsdobject.String("CA"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 229, LineColumn: cursorio.TextLineColumn{3, 1}},
							Until: cursorio.TextOffset{Byte: 542, LineColumn: cursorio.TextLineColumn{8, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 473, LineColumn: cursorio.TextLineColumn{7, 18}},
							Until: cursorio.TextOffset{Byte: 481, LineColumn: cursorio.TextLineColumn{7, 26}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 482, LineColumn: cursorio.TextLineColumn{7, 27}},
							Until: cursorio.TextOffset{Byte: 484, LineColumn: cursorio.TextLineColumn{7, 29}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("postal-code"),
						Object:    xsdobject.String("94043"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 229, LineColumn: cursorio.TextLineColumn{3, 1}},
							Until: cursorio.TextOffset{Byte: 542, LineColumn: cursorio.TextLineColumn{8, 8}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 507, LineColumn: cursorio.TextLineColumn{7, 52}},
							Until: cursorio.TextOffset{Byte: 520, LineColumn: cursorio.TextLineColumn{7, 65}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 521, LineColumn: cursorio.TextLineColumn{7, 66}},
							Until: cursorio.TextOffset{Byte: 526, LineColumn: cursorio.TextLineColumn{7, 71}},
						},
					},
				},
			},
		},
		{
			Name: "5.3.1.2/3",
			Snippet: `<span itemscope itemtype="http://microformats.org/profile/hcard"
><span itemprop=fn><span itemprop="n" itemscope><span itemprop="given-name"
>George</span> <span itemprop="family-name">Washington</span></span
></span></span>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://microformats.org/profile/hcard"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 224, LineColumn: cursorio.TextLineColumn{3, 15}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 16, LineColumn: cursorio.TextLineColumn{0, 16}},
							Until: cursorio.TextOffset{Byte: 24, LineColumn: cursorio.TextLineColumn{0, 24}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 25, LineColumn: cursorio.TextLineColumn{0, 25}},
							Until: cursorio.TextOffset{Byte: 64, LineColumn: cursorio.TextLineColumn{0, 64}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("fn"),
						Object:    xsdobject.String("George Washington"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 224, LineColumn: cursorio.TextLineColumn{3, 15}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 81, LineColumn: cursorio.TextLineColumn{1, 16}},
							Until: cursorio.TextOffset{Byte: 83, LineColumn: cursorio.TextLineColumn{1, 18}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{1, 19}},
							Until: cursorio.TextOffset{Byte: 210, LineColumn: cursorio.TextLineColumn{3, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("n"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 224, LineColumn: cursorio.TextLineColumn{3, 15}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 99, LineColumn: cursorio.TextLineColumn{1, 34}},
							Until: cursorio.TextOffset{Byte: 102, LineColumn: cursorio.TextLineColumn{1, 37}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{1, 19}},
							Until: cursorio.TextOffset{Byte: 210, LineColumn: cursorio.TextLineColumn{3, 1}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("given-name"),
						Object:    xsdobject.String("George"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{1, 19}},
							Until: cursorio.TextOffset{Byte: 210, LineColumn: cursorio.TextLineColumn{3, 1}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 128, LineColumn: cursorio.TextLineColumn{1, 63}},
							Until: cursorio.TextOffset{Byte: 140, LineColumn: cursorio.TextLineColumn{1, 75}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 142, LineColumn: cursorio.TextLineColumn{2, 1}},
							Until: cursorio.TextOffset{Byte: 148, LineColumn: cursorio.TextLineColumn{2, 7}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("family-name"),
						Object:    xsdobject.String("Washington"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{1, 19}},
							Until: cursorio.TextOffset{Byte: 210, LineColumn: cursorio.TextLineColumn{3, 1}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 171, LineColumn: cursorio.TextLineColumn{2, 30}},
							Until: cursorio.TextOffset{Byte: 184, LineColumn: cursorio.TextLineColumn{2, 43}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 185, LineColumn: cursorio.TextLineColumn{2, 44}},
							Until: cursorio.TextOffset{Byte: 195, LineColumn: cursorio.TextLineColumn{2, 54}},
						},
					},
				},
			},
		},
		{
			Name: "5.3.2.2/1a",
			Snippet: `<body itemscope itemtype="http://microformats.org/profile/hcalendar#vevent">
 ...
 <h1 itemprop="summary">Bluesday Tuesday: Money Road</h1>
 ...
 <time itemprop="dtstart" datetime="2009-05-05T19:00:00Z">May 5th @ 7pm</time>
 (until <time itemprop="dtend" datetime="2009-05-05T21:00:00Z">9pm</time>)
 ...
 <a href="http://livebrum.co.uk/2009/05/05/bluesday-tuesday-money-road"
    rel="bookmark" itemprop="url">Link to this page</a>
 ...
 <p>Location: <span itemprop="location">The RoadHouse</span></p>
 ...
 <p><input type=button value="Add to Calendar"
           onclick="location = getCalendar(this)"></p>
 ...
 <meta itemprop="description" content="via livebrum.co.uk">
</body>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://microformats.org/profile/hcalendar#vevent"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{16, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 16, LineColumn: cursorio.TextLineColumn{0, 16}},
							Until: cursorio.TextOffset{Byte: 24, LineColumn: cursorio.TextLineColumn{0, 24}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 25, LineColumn: cursorio.TextLineColumn{0, 25}},
							Until: cursorio.TextOffset{Byte: 75, LineColumn: cursorio.TextLineColumn{0, 75}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("summary"),
						Object:    xsdobject.String("Bluesday Tuesday: Money Road"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{16, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 96, LineColumn: cursorio.TextLineColumn{2, 14}},
							Until: cursorio.TextOffset{Byte: 105, LineColumn: cursorio.TextLineColumn{2, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 106, LineColumn: cursorio.TextLineColumn{2, 24}},
							Until: cursorio.TextOffset{Byte: 134, LineColumn: cursorio.TextLineColumn{2, 52}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("dtstart"),
						Object: rdf.Literal{
							LexicalForm: "2009-05-05T19:00:00Z",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#dateTime"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{16, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 161, LineColumn: cursorio.TextLineColumn{4, 16}},
							Until: cursorio.TextOffset{Byte: 170, LineColumn: cursorio.TextLineColumn{4, 25}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 180, LineColumn: cursorio.TextLineColumn{4, 35}},
							Until: cursorio.TextOffset{Byte: 202, LineColumn: cursorio.TextLineColumn{4, 57}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("dtend"),
						Object: rdf.Literal{
							LexicalForm: "2009-05-05T21:00:00Z",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#dateTime"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{16, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 247, LineColumn: cursorio.TextLineColumn{5, 23}},
							Until: cursorio.TextOffset{Byte: 254, LineColumn: cursorio.TextLineColumn{5, 30}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 264, LineColumn: cursorio.TextLineColumn{5, 40}},
							Until: cursorio.TextOffset{Byte: 286, LineColumn: cursorio.TextLineColumn{5, 62}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("url"),
						Object:    rdf.IRI("http://livebrum.co.uk/2009/05/05/bluesday-tuesday-money-road"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{16, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 404, LineColumn: cursorio.TextLineColumn{8, 28}},
							Until: cursorio.TextOffset{Byte: 409, LineColumn: cursorio.TextLineColumn{8, 33}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 313, LineColumn: cursorio.TextLineColumn{7, 9}},
							Until: cursorio.TextOffset{Byte: 375, LineColumn: cursorio.TextLineColumn{7, 71}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("location"),
						Object:    xsdobject.String("The RoadHouse"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{16, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 466, LineColumn: cursorio.TextLineColumn{10, 29}},
							Until: cursorio.TextOffset{Byte: 476, LineColumn: cursorio.TextLineColumn{10, 39}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 477, LineColumn: cursorio.TextLineColumn{10, 40}},
							Until: cursorio.TextOffset{Byte: 490, LineColumn: cursorio.TextLineColumn{10, 53}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("description"),
						Object:    xsdobject.String("via livebrum.co.uk"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 681, LineColumn: cursorio.TextLineColumn{16, 7}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 630, LineColumn: cursorio.TextLineColumn{15, 16}},
							Until: cursorio.TextOffset{Byte: 643, LineColumn: cursorio.TextLineColumn{15, 29}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 652, LineColumn: cursorio.TextLineColumn{15, 38}},
							Until: cursorio.TextOffset{Byte: 672, LineColumn: cursorio.TextLineColumn{15, 58}},
						},
					},
				},
			},
		},
		{
			Name: "5.3.2.2/1b",
			Snippet: `<div itemscope itemtype="http://microformats.org/profile/hcalendar#vevent">
 <p>I'm going to
 <strong itemprop="summary">Bluesday Tuesday: Money Road</strong>,
 <time itemprop="dtstart" datetime="2009-05-05T19:00:00Z">May 5th at 7pm</time>
 to <time itemprop="dtend" datetime="2009-05-05T21:00:00Z">9pm</time>,
 at <span itemprop="location">The RoadHouse</span>!</p>
 <p><a href="http://livebrum.co.uk/2009/05/05/bluesday-tuesday-money-road"
       itemprop="url">See this event on livebrum.co.uk</a>.</p>
 <meta itemprop="description" content="via livebrum.co.uk">
</div>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://microformats.org/profile/hcalendar#vevent"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 572, LineColumn: cursorio.TextLineColumn{9, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 15, LineColumn: cursorio.TextLineColumn{0, 15}},
							Until: cursorio.TextOffset{Byte: 23, LineColumn: cursorio.TextLineColumn{0, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 24, LineColumn: cursorio.TextLineColumn{0, 24}},
							Until: cursorio.TextOffset{Byte: 74, LineColumn: cursorio.TextLineColumn{0, 74}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("summary"),
						Object:    xsdobject.String("Bluesday Tuesday: Money Road"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 572, LineColumn: cursorio.TextLineColumn{9, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 111, LineColumn: cursorio.TextLineColumn{2, 18}},
							Until: cursorio.TextOffset{Byte: 120, LineColumn: cursorio.TextLineColumn{2, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 121, LineColumn: cursorio.TextLineColumn{2, 28}},
							Until: cursorio.TextOffset{Byte: 149, LineColumn: cursorio.TextLineColumn{2, 56}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("dtstart"),
						Object: rdf.Literal{
							LexicalForm: "2009-05-05T19:00:00Z",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#dateTime"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 572, LineColumn: cursorio.TextLineColumn{9, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 176, LineColumn: cursorio.TextLineColumn{3, 16}},
							Until: cursorio.TextOffset{Byte: 185, LineColumn: cursorio.TextLineColumn{3, 25}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 195, LineColumn: cursorio.TextLineColumn{3, 35}},
							Until: cursorio.TextOffset{Byte: 217, LineColumn: cursorio.TextLineColumn{3, 57}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("dtend"),
						Object: rdf.Literal{
							LexicalForm: "2009-05-05T21:00:00Z",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#dateTime"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 572, LineColumn: cursorio.TextLineColumn{9, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 259, LineColumn: cursorio.TextLineColumn{4, 19}},
							Until: cursorio.TextOffset{Byte: 266, LineColumn: cursorio.TextLineColumn{4, 26}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 276, LineColumn: cursorio.TextLineColumn{4, 36}},
							Until: cursorio.TextOffset{Byte: 298, LineColumn: cursorio.TextLineColumn{4, 58}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("location"),
						Object:    xsdobject.String("The RoadHouse"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 572, LineColumn: cursorio.TextLineColumn{9, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 330, LineColumn: cursorio.TextLineColumn{5, 19}},
							Until: cursorio.TextOffset{Byte: 340, LineColumn: cursorio.TextLineColumn{5, 29}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 341, LineColumn: cursorio.TextLineColumn{5, 30}},
							Until: cursorio.TextOffset{Byte: 354, LineColumn: cursorio.TextLineColumn{5, 43}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("url"),
						Object:    rdf.IRI("http://livebrum.co.uk/2009/05/05/bluesday-tuesday-money-road"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 572, LineColumn: cursorio.TextLineColumn{9, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 458, LineColumn: cursorio.TextLineColumn{7, 16}},
							Until: cursorio.TextOffset{Byte: 463, LineColumn: cursorio.TextLineColumn{7, 21}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 379, LineColumn: cursorio.TextLineColumn{6, 12}},
							Until: cursorio.TextOffset{Byte: 441, LineColumn: cursorio.TextLineColumn{6, 74}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("description"),
						Object:    xsdobject.String("via livebrum.co.uk"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 572, LineColumn: cursorio.TextLineColumn{9, 6}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 522, LineColumn: cursorio.TextLineColumn{8, 16}},
							Until: cursorio.TextOffset{Byte: 535, LineColumn: cursorio.TextLineColumn{8, 29}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 544, LineColumn: cursorio.TextLineColumn{8, 38}},
							Until: cursorio.TextOffset{Byte: 564, LineColumn: cursorio.TextLineColumn{8, 58}},
						},
					},
				},
			},
		},
		{
			Name: "5.3.2.1/1",
			Snippet: `<figure itemscope itemtype="http://n.whatwg.org/work">
 <img itemprop="work" src="mypond.jpeg">
 <figcaption>
  <p><cite itemprop="title">My Pond</cite></p>
  <p><small>Licensed under the <a itemprop="license"
  href="https://creativecommons.org/licenses/by-sa/4.0/">Creative
  Commons Attribution-Share Alike 4.0 International License</a>
  and the <a itemprop="license"
  href="http://www.opensource.org/licenses/mit-license.php">MIT
  license</a>.</small>
 </figcaption>
</figure>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://n.whatwg.org/work"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 483, LineColumn: cursorio.TextLineColumn{11, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 18, LineColumn: cursorio.TextLineColumn{0, 18}},
							Until: cursorio.TextOffset{Byte: 26, LineColumn: cursorio.TextLineColumn{0, 26}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 27, LineColumn: cursorio.TextLineColumn{0, 27}},
							Until: cursorio.TextOffset{Byte: 53, LineColumn: cursorio.TextLineColumn{0, 53}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("work"),
						Object:    rdf.IRI("mypond.jpeg"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 483, LineColumn: cursorio.TextLineColumn{11, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 70, LineColumn: cursorio.TextLineColumn{1, 15}},
							Until: cursorio.TextOffset{Byte: 76, LineColumn: cursorio.TextLineColumn{1, 21}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 81, LineColumn: cursorio.TextLineColumn{1, 26}},
							Until: cursorio.TextOffset{Byte: 94, LineColumn: cursorio.TextLineColumn{1, 39}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("title"),
						Object:    xsdobject.String("My Pond"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 483, LineColumn: cursorio.TextLineColumn{11, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 130, LineColumn: cursorio.TextLineColumn{3, 20}},
							Until: cursorio.TextOffset{Byte: 137, LineColumn: cursorio.TextLineColumn{3, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 138, LineColumn: cursorio.TextLineColumn{3, 28}},
							Until: cursorio.TextOffset{Byte: 145, LineColumn: cursorio.TextLineColumn{3, 35}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("license"),
						Object:    rdf.IRI("https://creativecommons.org/licenses/by-sa/4.0/"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 483, LineColumn: cursorio.TextLineColumn{11, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 200, LineColumn: cursorio.TextLineColumn{4, 43}},
							Until: cursorio.TextOffset{Byte: 209, LineColumn: cursorio.TextLineColumn{4, 52}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 217, LineColumn: cursorio.TextLineColumn{5, 7}},
							Until: cursorio.TextOffset{Byte: 266, LineColumn: cursorio.TextLineColumn{5, 56}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("license"),
						Object:    rdf.IRI("http://www.opensource.org/licenses/mit-license.php"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 0, LineColumn: cursorio.TextLineColumn{0, 0}},
							Until: cursorio.TextOffset{Byte: 483, LineColumn: cursorio.TextLineColumn{11, 9}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 362, LineColumn: cursorio.TextLineColumn{7, 22}},
							Until: cursorio.TextOffset{Byte: 371, LineColumn: cursorio.TextLineColumn{7, 31}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 379, LineColumn: cursorio.TextLineColumn{8, 7}},
							Until: cursorio.TextOffset{Byte: 431, LineColumn: cursorio.TextLineColumn{8, 59}},
						},
					},
				},
			},
		},
		{
			Name: "5.4.1/1",
			Snippet: `<!DOCTYPE HTML>
<html lang="en">
<title>My Blog</title>
<article itemscope itemtype="http://schema.org/BlogPosting">
 <header>
  <h1 itemprop="headline">Progress report</h1>
  <p><time itemprop="datePublished" datetime="2013-08-29">today</time></p>
  <link itemprop="url" href="?comments=0">
 </header>
 <p>All in all, he's doing well with his swim lessons. The biggest thing was he had trouble
 putting his head in, but we got it down.</p>
 <section>
  <h1>Comments</h1>
  <article itemprop="comment" itemscope itemtype="http://schema.org/UserComments" id="c1">
   <link itemprop="url" href="#c1">
   <footer>
    <p>Posted by: <span itemprop="creator" itemscope itemtype="http://schema.org/Person">
     <span itemprop="name">Greg</span>
    </span></p>
    <p><time itemprop="commentTime" datetime="2013-08-29">15 minutes ago</time></p>
   </footer>
   <p>Ha!</p>
  </article>
  <article itemprop="comment" itemscope itemtype="http://schema.org/UserComments" id="c2">
   <link itemprop="url" href="#c2">
   <footer>
    <p>Posted by: <span itemprop="creator" itemscope itemtype="http://schema.org/Person">
     <span itemprop="name">Charlotte</span>
    </span></p>
    <p><time itemprop="commentTime" datetime="2013-08-29">5 minutes ago</time></p>
   </footer>
   <p>When you say "we got it down"...</p>
  </article>
 </section>
</article>`,
			Expected: encodingtest.TripleStatementList{
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/BlogPosting"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 1343, LineColumn: cursorio.TextLineColumn{34, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 75, LineColumn: cursorio.TextLineColumn{3, 19}},
							Until: cursorio.TextOffset{Byte: 83, LineColumn: cursorio.TextLineColumn{3, 27}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 84, LineColumn: cursorio.TextLineColumn{3, 28}},
							Until: cursorio.TextOffset{Byte: 115, LineColumn: cursorio.TextLineColumn{3, 59}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("headline"),
						Object:    xsdobject.String("Progress report"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 1343, LineColumn: cursorio.TextLineColumn{34, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 142, LineColumn: cursorio.TextLineColumn{5, 15}},
							Until: cursorio.TextOffset{Byte: 152, LineColumn: cursorio.TextLineColumn{5, 25}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 153, LineColumn: cursorio.TextLineColumn{5, 26}},
							Until: cursorio.TextOffset{Byte: 168, LineColumn: cursorio.TextLineColumn{5, 41}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("datePublished"),
						Object: rdf.Literal{
							LexicalForm: "2013-08-29",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#date"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 1343, LineColumn: cursorio.TextLineColumn{34, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 194, LineColumn: cursorio.TextLineColumn{6, 20}},
							Until: cursorio.TextOffset{Byte: 209, LineColumn: cursorio.TextLineColumn{6, 35}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 219, LineColumn: cursorio.TextLineColumn{6, 45}},
							Until: cursorio.TextOffset{Byte: 231, LineColumn: cursorio.TextLineColumn{6, 57}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("url"),
						Object:    rdf.IRI("?comments=0"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 1343, LineColumn: cursorio.TextLineColumn{34, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 266, LineColumn: cursorio.TextLineColumn{7, 17}},
							Until: cursorio.TextOffset{Byte: 271, LineColumn: cursorio.TextLineColumn{7, 22}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 277, LineColumn: cursorio.TextLineColumn{7, 28}},
							Until: cursorio.TextOffset{Byte: 290, LineColumn: cursorio.TextLineColumn{7, 41}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("comment"),
						Object:    testingBnode.MapBlankNodeIdentifier("b1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 1343, LineColumn: cursorio.TextLineColumn{34, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 492, LineColumn: cursorio.TextLineColumn{13, 20}},
							Until: cursorio.TextOffset{Byte: 501, LineColumn: cursorio.TextLineColumn{13, 29}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 474, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 879, LineColumn: cursorio.TextLineColumn{22, 12}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/UserComments"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 474, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 879, LineColumn: cursorio.TextLineColumn{22, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 512, LineColumn: cursorio.TextLineColumn{13, 40}},
							Until: cursorio.TextOffset{Byte: 520, LineColumn: cursorio.TextLineColumn{13, 48}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 521, LineColumn: cursorio.TextLineColumn{13, 49}},
							Until: cursorio.TextOffset{Byte: 553, LineColumn: cursorio.TextLineColumn{13, 81}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("url"),
						Object:    rdf.IRI("#c1"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 474, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 879, LineColumn: cursorio.TextLineColumn{22, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 581, LineColumn: cursorio.TextLineColumn{14, 18}},
							Until: cursorio.TextOffset{Byte: 586, LineColumn: cursorio.TextLineColumn{14, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 592, LineColumn: cursorio.TextLineColumn{14, 29}},
							Until: cursorio.TextOffset{Byte: 597, LineColumn: cursorio.TextLineColumn{14, 34}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("creator"),
						Object:    testingBnode.MapBlankNodeIdentifier("b2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 474, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 879, LineColumn: cursorio.TextLineColumn{22, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 644, LineColumn: cursorio.TextLineColumn{16, 33}},
							Until: cursorio.TextOffset{Byte: 653, LineColumn: cursorio.TextLineColumn{16, 42}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 629, LineColumn: cursorio.TextLineColumn{16, 18}},
							Until: cursorio.TextOffset{Byte: 751, LineColumn: cursorio.TextLineColumn{18, 11}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/Person"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 629, LineColumn: cursorio.TextLineColumn{16, 18}},
							Until: cursorio.TextOffset{Byte: 751, LineColumn: cursorio.TextLineColumn{18, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 664, LineColumn: cursorio.TextLineColumn{16, 53}},
							Until: cursorio.TextOffset{Byte: 672, LineColumn: cursorio.TextLineColumn{16, 61}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 673, LineColumn: cursorio.TextLineColumn{16, 62}},
							Until: cursorio.TextOffset{Byte: 699, LineColumn: cursorio.TextLineColumn{16, 88}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b2"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Greg"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 629, LineColumn: cursorio.TextLineColumn{16, 18}},
							Until: cursorio.TextOffset{Byte: 751, LineColumn: cursorio.TextLineColumn{18, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 721, LineColumn: cursorio.TextLineColumn{17, 20}},
							Until: cursorio.TextOffset{Byte: 727, LineColumn: cursorio.TextLineColumn{17, 26}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 728, LineColumn: cursorio.TextLineColumn{17, 27}},
							Until: cursorio.TextOffset{Byte: 732, LineColumn: cursorio.TextLineColumn{17, 31}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b1"),
						Predicate: rdf.IRI("commentTime"),
						Object: rdf.Literal{
							LexicalForm: "2013-08-29",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#date"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 474, LineColumn: cursorio.TextLineColumn{13, 2}},
							Until: cursorio.TextOffset{Byte: 879, LineColumn: cursorio.TextLineColumn{22, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 778, LineColumn: cursorio.TextLineColumn{19, 22}},
							Until: cursorio.TextOffset{Byte: 791, LineColumn: cursorio.TextLineColumn{19, 35}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 801, LineColumn: cursorio.TextLineColumn{19, 45}},
							Until: cursorio.TextOffset{Byte: 813, LineColumn: cursorio.TextLineColumn{19, 57}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b0"),
						Predicate: rdf.IRI("comment"),
						Object:    testingBnode.MapBlankNodeIdentifier("b3"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 56, LineColumn: cursorio.TextLineColumn{3, 0}},
							Until: cursorio.TextOffset{Byte: 1343, LineColumn: cursorio.TextLineColumn{34, 10}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 900, LineColumn: cursorio.TextLineColumn{23, 20}},
							Until: cursorio.TextOffset{Byte: 909, LineColumn: cursorio.TextLineColumn{23, 29}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 882, LineColumn: cursorio.TextLineColumn{23, 2}},
							Until: cursorio.TextOffset{Byte: 1320, LineColumn: cursorio.TextLineColumn{32, 12}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/UserComments"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 882, LineColumn: cursorio.TextLineColumn{23, 2}},
							Until: cursorio.TextOffset{Byte: 1320, LineColumn: cursorio.TextLineColumn{32, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 920, LineColumn: cursorio.TextLineColumn{23, 40}},
							Until: cursorio.TextOffset{Byte: 928, LineColumn: cursorio.TextLineColumn{23, 48}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 929, LineColumn: cursorio.TextLineColumn{23, 49}},
							Until: cursorio.TextOffset{Byte: 961, LineColumn: cursorio.TextLineColumn{23, 81}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("url"),
						Object:    rdf.IRI("#c2"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 882, LineColumn: cursorio.TextLineColumn{23, 2}},
							Until: cursorio.TextOffset{Byte: 1320, LineColumn: cursorio.TextLineColumn{32, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 989, LineColumn: cursorio.TextLineColumn{24, 18}},
							Until: cursorio.TextOffset{Byte: 994, LineColumn: cursorio.TextLineColumn{24, 23}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1000, LineColumn: cursorio.TextLineColumn{24, 29}},
							Until: cursorio.TextOffset{Byte: 1005, LineColumn: cursorio.TextLineColumn{24, 34}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("creator"),
						Object:    testingBnode.MapBlankNodeIdentifier("b4"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 882, LineColumn: cursorio.TextLineColumn{23, 2}},
							Until: cursorio.TextOffset{Byte: 1320, LineColumn: cursorio.TextLineColumn{32, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1052, LineColumn: cursorio.TextLineColumn{26, 33}},
							Until: cursorio.TextOffset{Byte: 1061, LineColumn: cursorio.TextLineColumn{26, 42}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1037, LineColumn: cursorio.TextLineColumn{26, 18}},
							Until: cursorio.TextOffset{Byte: 1164, LineColumn: cursorio.TextLineColumn{28, 11}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b4"),
						Predicate: rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
						Object:    rdf.IRI("http://schema.org/Person"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1037, LineColumn: cursorio.TextLineColumn{26, 18}},
							Until: cursorio.TextOffset{Byte: 1164, LineColumn: cursorio.TextLineColumn{28, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1072, LineColumn: cursorio.TextLineColumn{26, 53}},
							Until: cursorio.TextOffset{Byte: 1080, LineColumn: cursorio.TextLineColumn{26, 61}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1081, LineColumn: cursorio.TextLineColumn{26, 62}},
							Until: cursorio.TextOffset{Byte: 1107, LineColumn: cursorio.TextLineColumn{26, 88}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b4"),
						Predicate: rdf.IRI("name"),
						Object:    xsdobject.String("Charlotte"),
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1037, LineColumn: cursorio.TextLineColumn{26, 18}},
							Until: cursorio.TextOffset{Byte: 1164, LineColumn: cursorio.TextLineColumn{28, 11}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1129, LineColumn: cursorio.TextLineColumn{27, 20}},
							Until: cursorio.TextOffset{Byte: 1135, LineColumn: cursorio.TextLineColumn{27, 26}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1136, LineColumn: cursorio.TextLineColumn{27, 27}},
							Until: cursorio.TextOffset{Byte: 1145, LineColumn: cursorio.TextLineColumn{27, 36}},
						},
					},
				},
				encodingtest.TripleStatement{
					Triple: rdf.Triple{
						Subject:   testingBnode.MapBlankNodeIdentifier("b3"),
						Predicate: rdf.IRI("commentTime"),
						Object: rdf.Literal{
							LexicalForm: "2013-08-29",
							Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#date"),
						},
					},
					TextOffsets: encoding.StatementTextOffsets{
						encoding.SubjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 882, LineColumn: cursorio.TextLineColumn{23, 2}},
							Until: cursorio.TextOffset{Byte: 1320, LineColumn: cursorio.TextLineColumn{32, 12}},
						},
						encoding.PredicateStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1191, LineColumn: cursorio.TextLineColumn{29, 22}},
							Until: cursorio.TextOffset{Byte: 1204, LineColumn: cursorio.TextLineColumn{29, 35}},
						},
						encoding.ObjectStatementOffsets: cursorio.TextOffsetRange{
							From:  cursorio.TextOffset{Byte: 1214, LineColumn: cursorio.TextLineColumn{29, 45}},
							Until: cursorio.TextOffset{Byte: 1226, LineColumn: cursorio.TextLineColumn{29, 57}},
						},
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

			out, err := triples.CollectErr(NewDecoder(htmlDocument))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// if false {
			// 	w := devrdfioutil.NewWriter(os.Stderr, devrdfioutil.WriterOptions{
			// 		Source: []byte(testcase.Snippet),
			// 	})

			// 	for _, stmt := range out {
			// 		if err := w.PutStatement(nil, stmt); err != nil {
			// 			t.Fatalf("unexpected error: %v", err)
			// 		}
			// 	}

			// 	if err := w.Close(); err != nil {
			// 		t.Fatalf("unexpected error: %v", err)
			// 	}
			// }

			microdataLivingAssertEquals(t, testcase.Expected.AsTriples(), out)
		})
	}
}
