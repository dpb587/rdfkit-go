package rdfdescriptionstruct_test

import (
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/x/rdfdescriptionstruct"
)

func TestMarshal_Basic(t *testing.T) {
	type TestStruct struct {
		Subject rdf.IRI `rdf:"s"`
		Type    rdf.IRI `rdf:"o,p=rdf:type"`
		Name    string  `rdf:"o,p=mf:name"`
	}

	input := TestStruct{
		Subject: "https://example.com/test1",
		Type:    "https://example.com/TestType",
		Name:    "Test Name",
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	resource := resources[0]

	if resource.GetResourceSubject() != rdf.IRI("https://example.com/test1") {
		t.Errorf("Subject mismatch: got %v, want %v", resource.GetResourceSubject(), "https://example.com/test1")
	}

	stmts := resource.GetResourceStatements()
	if len(stmts) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(stmts))
	}

	// Check Type statement
	typeStmt := stmts[0].(rdfdescription.ObjectStatement)
	if typeStmt.Predicate != rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type") {
		t.Errorf("Type predicate mismatch: got %v", typeStmt.Predicate)
	}
	if typeStmt.Object != rdf.IRI("https://example.com/TestType") {
		t.Errorf("Type object mismatch: got %v", typeStmt.Object)
	}

	// Check Name statement
	nameStmt := stmts[1].(rdfdescription.ObjectStatement)
	// mf: is not a standard prefix, so it won't be expanded
	expectedNamePred := rdf.IRI("mf:name")
	if nameStmt.Predicate != expectedNamePred {
		t.Errorf("Name predicate mismatch: got %v, want %v", nameStmt.Predicate, expectedNamePred)
	}
	expectedNameLiteral := xsdobject.String("Test Name").(rdf.Literal)
	nameLiteral, ok := nameStmt.Object.(rdf.Literal)
	if !ok {
		t.Fatalf("Name object is not a Literal: got %T", nameStmt.Object)
	}
	if nameLiteral.LexicalForm != expectedNameLiteral.LexicalForm || nameLiteral.Datatype != expectedNameLiteral.Datatype {
		t.Errorf("Name object mismatch: got %v, want %v", nameStmt.Object, expectedNameLiteral)
	}
}

func TestMarshal_CompactIRI(t *testing.T) {
	type TestStruct struct {
		Subject rdf.IRI `rdf:"s"`
		Type    rdf.IRI `rdf:"o,p=rdf:type"`
		Label   string  `rdf:"o,p=rdfs:label"`
	}

	input := TestStruct{
		Subject: "https://example.com/test",
		Type:    "https://example.com/Type",
		Label:   "Test Label",
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	resource := resources[0]

	// Check Label statement uses expanded IRI
	stmts := resource.GetResourceStatements()
	labelStmt := stmts[1].(rdfdescription.ObjectStatement)
	expectedPred := "http://www.w3.org/2000/01/rdf-schema#label"
	if labelStmt.Predicate != rdf.IRI(expectedPred) {
		t.Errorf("Label predicate mismatch: got %v, want %v", labelStmt.Predicate, expectedPred)
	}
}

func TestMarshal_NumericTypes(t *testing.T) {
	type TestStruct struct {
		Subject rdf.IRI `rdf:"s"`
		UInt8   uint8   `rdf:"o,p=https://example.com/uint8"`
		UInt16  uint16  `rdf:"o,p=https://example.com/uint16"`
		UInt32  uint32  `rdf:"o,p=https://example.com/uint32"`
		UInt64  uint64  `rdf:"o,p=https://example.com/uint64"`
		Int16   int16   `rdf:"o,p=https://example.com/int16"`
		Int32   int32   `rdf:"o,p=https://example.com/int32"`
		Int64   int64   `rdf:"o,p=https://example.com/int64"`
		Float32 float32 `rdf:"o,p=https://example.com/float32"`
		Float64 float64 `rdf:"o,p=https://example.com/float64"`
	}

	input := TestStruct{
		Subject: "https://example.com/test",
		UInt8:   255,
		UInt16:  65535,
		UInt32:  4294967295,
		UInt64:  18446744073709551615,
		Int16:   -32768,
		Int32:   -2147483648,
		Int64:   -9223372036854775808,
		Float32: 3.14,
		Float64: 2.718281828,
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	resource := resources[0]

	stmts := resource.GetResourceStatements()
	if len(stmts) != 9 {
		t.Fatalf("Expected 9 statements, got %d", len(stmts))
	}

	// Verify each statement has a literal object
	for i, stmt := range stmts {
		objStmt := stmt.(rdfdescription.ObjectStatement)
		if _, ok := objStmt.Object.(rdf.Literal); !ok {
			t.Errorf("Statement %d: expected Literal object, got %T", i, objStmt.Object)
		}
	}
}

func TestMarshal_Pointer(t *testing.T) {
	type TestStruct struct {
		Subject rdf.IRI `rdf:"s"`
		Name    *string `rdf:"o,p=mf:name"`
		Count   *int32  `rdf:"o,p=https://example.com/count"`
	}

	name := "Test"
	count := int32(42)

	input := TestStruct{
		Subject: "https://example.com/test",
		Name:    &name,
		Count:   &count,
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	resource := resources[0]

	stmts := resource.GetResourceStatements()
	if len(stmts) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(stmts))
	}
}

func TestMarshal_NilPointer(t *testing.T) {
	type TestStruct struct {
		Subject rdf.IRI `rdf:"s"`
		Name    *string `rdf:"o,p=mf:name"`
	}

	input := TestStruct{
		Subject: "https://example.com/test",
		Name:    nil,
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	resource := resources[0]

	// Nil pointer should produce no statement
	stmts := resource.GetResourceStatements()
	if len(stmts) != 0 {
		t.Fatalf("Expected 0 statements for nil pointer, got %d", len(stmts))
	}
}

func TestMarshal_Slice(t *testing.T) {
	type TestStruct struct {
		Subject rdf.IRI   `rdf:"s"`
		Types   []rdf.IRI `rdf:"o,p=rdf:type"`
	}

	input := TestStruct{
		Subject: "https://example.com/test",
		Types: []rdf.IRI{
			"https://example.com/Type1",
			"https://example.com/Type2",
			"https://example.com/Type3",
		},
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	resource := resources[0]

	stmts := resource.GetResourceStatements()
	if len(stmts) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(stmts))
	}

	// All should have same predicate
	for _, stmt := range stmts {
		objStmt := stmt.(rdfdescription.ObjectStatement)
		if objStmt.Predicate != rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type") {
			t.Errorf("Unexpected predicate: %v", objStmt.Predicate)
		}
	}
}

func TestMarshal_BlankNode(t *testing.T) {
	type TestStruct struct {
		Subject rdf.BlankNode `rdf:"s"`
		Name    string        `rdf:"o,p=mf:name"`
	}

	bn := rdf.NewBlankNode()
	input := TestStruct{
		Subject: bn,
		Name:    "Test",
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	resource := resources[0]

	if resource.GetResourceSubject() != bn {
		t.Errorf("Subject mismatch: got %v, want %v", resource.GetResourceSubject(), bn)
	}
}

func TestMarshal_Recursive(t *testing.T) {
	type NestedStruct struct {
		Subject rdf.IRI `rdf:"s"`
		Name    string  `rdf:"o,p=mf:name"`
	}

	type TestStruct struct {
		Subject rdf.IRI       `rdf:"s"`
		Type    rdf.IRI       `rdf:"o,p=rdf:type"`
		Nested  *NestedStruct `rdf:"o,p=https://example.com/nested"`
	}

	input := TestStruct{
		Subject: "https://example.com/test",
		Type:    "https://example.com/Type",
		Nested: &NestedStruct{
			Subject: "https://example.com/nested",
			Name:    "Nested Name",
		},
	}

	resources, err := rdfdescriptionstruct.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	// Should have 2 resources: main and nested
	if len(resources) != 2 {
		t.Fatalf("Expected 2 resources, got %d", len(resources))
	}

	mainResource := resources[0]
	nestedResource := resources[1]

	// Main resource should have 2 statements: Type and link to nested
	stmts := mainResource.GetResourceStatements()
	if len(stmts) != 2 {
		t.Fatalf("Expected 2 statements in main resource, got %d", len(stmts))
	}

	// First statement is Type
	stmt0 := stmts[0].(rdfdescription.ObjectStatement)
	if stmt0.Predicate != rdf.IRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type") {
		t.Errorf("First statement should be Type")
	}

	// Second statement links to nested resource
	nestedLinkStmt := stmts[1].(rdfdescription.ObjectStatement)
	if nestedLinkStmt.Predicate != rdf.IRI("https://example.com/nested") {
		t.Errorf("Second statement predicate mismatch: got %v", nestedLinkStmt.Predicate)
	}
	if nestedLinkStmt.Object != rdf.IRI("https://example.com/nested") {
		t.Errorf("Second statement object mismatch: got %v", nestedLinkStmt.Object)
	}

	// Nested resource should have 1 statement: Name
	nestedStmts := nestedResource.GetResourceStatements()
	if len(nestedStmts) != 1 {
		t.Fatalf("Expected 1 statement in nested resource, got %d", len(nestedStmts))
	}

	// Check nested's Name statement
	nestedNameStmt := nestedStmts[0].(rdfdescription.ObjectStatement)
	// mf: is not a standard prefix, so it won't be expanded
	expectedNamePred := rdf.IRI("mf:name")
	if nestedNameStmt.Predicate != expectedNamePred {
		t.Errorf("Nested statement predicate mismatch: got %v, want %v", nestedNameStmt.Predicate, expectedNamePred)
	}
}

func TestMarshal_NoSubject(t *testing.T) {
	type TestStruct struct {
		Name string `rdf:"o,p=mf:name"`
	}

	input := TestStruct{
		Name: "Test",
	}

	_, err := rdfdescriptionstruct.Marshal(input)
	if err == nil {
		t.Fatal("Expected error for struct without subject field")
	}
}

func TestMarshal_RoundTrip(t *testing.T) {
	type TestStruct struct {
		Subject rdf.IRI `rdf:"s"`
		Type    rdf.IRI `rdf:"o,p=rdf:type"`
		Name    string  `rdf:"o,p=mf:name"`
		Count   int32   `rdf:"o,p=https://example.com/count"`
	}

	original := TestStruct{
		Subject: "https://example.com/test",
		Type:    "https://example.com/TestType",
		Name:    "Test Name",
		Count:   42,
	}

	// Marshal to resource
	resources, err := rdfdescriptionstruct.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	// Unmarshal back to struct
	var result TestStruct
	err = rdfdescriptionstruct.UnmarshalResource(resources[0], &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// Compare
	if result.Subject != original.Subject {
		t.Errorf("Subject mismatch: got %v, want %v", result.Subject, original.Subject)
	}
	if result.Type != original.Type {
		t.Errorf("Type mismatch: got %v, want %v", result.Type, original.Type)
	}
	if result.Name != original.Name {
		t.Errorf("Name mismatch: got %v, want %v", result.Name, original.Name)
	}
	if result.Count != original.Count {
		t.Errorf("Count mismatch: got %v, want %v", result.Count, original.Count)
	}
}
