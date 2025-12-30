package rdfdescriptionstruct_test

import (
	"fmt"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/x/rdfdescriptionstruct"
)

// Example demonstrates basic usage of rdfdescriptionstruct.
func Example() {
	// Define a struct with rdf tags (using compact IRIs)
	type Person struct {
		ID   rdf.IRI `rdf:"s"`
		Name string  `rdf:"o,p=foaf:name"`
		Age  *int64  `rdf:"o,p=http://xmlns.com/foaf/0.1/age"`
	}

	// Create an RDF resource
	resource := &rdfdescription.SubjectResource{
		Subject: rdf.IRI("http://example.com/alice"),
		Statements: rdfdescription.StatementList{
			rdfdescription.ObjectStatement{
				Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
				Object: rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: "Alice",
				},
			},
			rdfdescription.ObjectStatement{
				Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/age"),
				Object: rdf.Literal{
					Datatype:    xsdiri.Integer_Datatype,
					LexicalForm: "30",
				},
			},
		},
	}

	// Unmarshal the resource into the struct
	var person Person
	err := rdfdescriptionstruct.UnmarshalResource(resource, &person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %s\n", person.ID)
	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("Age: %d\n", *person.Age)

	// Output:
	// ID: http://example.com/alice
	// Name: Alice
	// Age: 30
}

// Example_recursive demonstrates recursive resource unmarshaling.
func Example_recursive() {
	// Define structs with nested resources (using compact IRIs)
	type Address struct {
		City string `rdf:"o,p=schema:addressLocality"`
	}

	type Person struct {
		Name    string   `rdf:"o,p=foaf:name"`
		Address *Address `rdf:"o,p=schema:address"`
	}

	// Create a builder and add triples
	builder := rdfdescription.NewResourceListBuilder()

	builder.Add(rdf.Triple{
		Subject:   rdf.IRI("http://example.com/person1"),
		Predicate: rdf.IRI("http://xmlns.com/foaf/0.1/name"),
		Object: rdf.Literal{
			Datatype:    xsdiri.String_Datatype,
			LexicalForm: "Bob",
		},
	})
	builder.Add(rdf.Triple{
		Subject:   rdf.IRI("http://example.com/person1"),
		Predicate: rdf.IRI("http://schema.org/address"),
		Object:    rdf.IRI("http://example.com/addr1"),
	})
	builder.Add(rdf.Triple{
		Subject:   rdf.IRI("http://example.com/addr1"),
		Predicate: rdf.IRI("http://schema.org/addressLocality"),
		Object: rdf.Literal{
			Datatype:    xsdiri.String_Datatype,
			LexicalForm: "Seattle",
		},
	})

	// Create the main resource
	resource := &rdfdescription.SubjectResource{
		Subject:    rdf.IRI("http://example.com/person1"),
		Statements: builder.GetResourceStatements(rdf.IRI("http://example.com/person1")),
	}

	// Unmarshal with the builder for recursive unmarshaling
	var person Person
	err := rdfdescriptionstruct.Unmarshal(builder, resource, &person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("City: %s\n", person.Address.City)

	// Output:
	// Name: Bob
	// City: Seattle
}

// Example_marshal demonstrates converting a Go struct to an RDF resource.
func Example_marshal() {
	// Define a struct with rdf tags (using compact IRIs)
	type Person struct {
		ID   rdf.IRI `rdf:"s"`
		Type rdf.IRI `rdf:"o,p=rdf:type"`
		Name string  `rdf:"o,p=foaf:name"`
		Age  int32   `rdf:"o,p=foaf:age"`
	}

	// Create a Go struct
	person := Person{
		ID:   "http://example.com/alice",
		Type: "http://xmlns.com/foaf/0.1/Person",
		Name: "Alice",
		Age:  30,
	}

	// Marshal to RDF resource
	resources, err := rdfdescriptionstruct.Marshal(person)
	if err != nil {
		panic(err)
	}

	// Get the main resource (first in the list)
	resource := resources[0]

	fmt.Printf("Subject: %s\n", resource.GetResourceSubject())

	statements := resource.GetResourceStatements()
	fmt.Printf("Statements: %d\n", len(statements))

	for _, stmt := range statements {
		objStmt := stmt.(rdfdescription.ObjectStatement)
		// Display object value - handle Literals specially
		var objDisplay string
		if lit, ok := objStmt.Object.(rdf.Literal); ok {
			objDisplay = lit.LexicalForm
		} else {
			objDisplay = fmt.Sprintf("%v", objStmt.Object)
		}
		fmt.Printf("  %s = %s\n", objStmt.Predicate, objDisplay)
	}

	// Output:
	// Subject: http://example.com/alice
	// Statements: 3
	//   http://www.w3.org/1999/02/22-rdf-syntax-ns#type = http://xmlns.com/foaf/0.1/Person
	//   http://xmlns.com/foaf/0.1/name = Alice
	//   http://xmlns.com/foaf/0.1/age = 30
}

// Example_customPrefixes demonstrates how to use custom prefix mappings
// with the Unmarshaler.
func Example_customPrefixes() {
	// Define a struct with custom prefixes
	type Book struct {
		ID     rdf.IRI `rdf:"s"`
		Title  string  `rdf:"o,p=dc:title"`
		Author string  `rdf:"o,p=dc:creator"`
	}

	// Create an RDF resource with Dublin Core predicates
	resource := &rdfdescription.SubjectResource{
		Subject: rdf.IRI("http://example.com/book1"),
		Statements: rdfdescription.StatementList{
			rdfdescription.ObjectStatement{
				Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/title"),
				Object: rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: "The Go Programming Language",
				},
			},
			rdfdescription.ObjectStatement{
				Predicate: rdf.IRI("http://purl.org/dc/elements/1.1/creator"),
				Object: rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: "Alan Donovan",
				},
			},
		},
	}

	// Create custom prefix map for Dublin Core
	customPrefixes := map[string]rdf.IRI{
		"dc": "http://purl.org/dc/elements/1.1/",
	}

	// Create unmarshaler with custom prefixes
	config := rdfdescriptionstruct.NewUnmarshalerConfig().SetPrefixes(customPrefixes)
	unmarshaler := rdfdescriptionstruct.NewUnmarshaler(config)

	// Unmarshal using the custom unmarshaler
	var book Book
	err := unmarshaler.UnmarshalResource(resource, &book)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %s\n", book.ID)
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)

	// Output:
	// ID: http://example.com/book1
	// Title: The Go Programming Language
	// Author: Alan Donovan
}
