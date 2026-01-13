package earltesting

import "github.com/dpb587/rdfkit-go/rdf"

type AssertionProfile struct {
	// Test is equivalent to earl:test
	Test rdf.ObjectValue

	// Node is the subject for the given earl:Assertion
	Node rdf.BlankNode

	// ResultNode is the object of the assertion's earl:result property.
	ResultNode rdf.BlankNode

	// TestingName is the name reported by Go testing's t.Name()
	TestingName string
}
