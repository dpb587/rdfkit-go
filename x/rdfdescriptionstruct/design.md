---
title: Struct Tags for rdfdescription
---

Go supports "tags" being defined on fields for structs which can be used by arbitrary marshal/unmarshal services. This implementation allows a user to unmarshal a single or list of structs based on an `rdfdescription.Resource` or `rdfdescription.ResourceList`.

The available function definitions are:

* `Unmarshal(builder rdfdescription.ResourceListBuilder, from rdfdescription.Resource, to any) error`
* `UnmarshalResource(from rdfdescription.Resource, to any) error`
* `Marshal(from any) (rdfdescription.ResourceList, error)`

For additional customizations, a `Unmarshaler` may be created with variadic arguments of `UnmarshalerOption` type. The only implementation of a `UnmarshalerOption` is `UnmarshalerConfig` which has the following methods:

* `SetPrefixes(iriutil.PrefixMap)` - override the default RDFa Initial Context (Widely Used) prefixes.

# Tags

Struct tags for the `rdf` namespace must be one of the following syntaxes.

* `s` - indicates the subject of the current resource scope. Potential value behaviors are defined in the "Tagged Subject" section.
* `o,p={PREDICATE}` - indicates the object value of a statement. Potential value behaviors are defined in the "Tagged Object Scalar" section.

The `PREDICATE` token may be either an IRI or a compact IRI (e.g. `rdf:type`). By default, the rdfacontext.WidelyUsedInitialContext() prefixes are used for compact IRIs.

## Tagged Subject

A tagged subject value must be a scalar of one of the following types.

* `rdf.SubjectValue` - the most generic type and always acceptable
* `rdf.IRI` - it is an error if the incoming value is not of type `rdf.IRI`; pointer values are supported.
* `rdf.BlankNode` - it is an error if the incoming value is not of the type `rdf.BlankNode`

## Tagged Object Scalar

A tagged object value may be a scalar or slice of one of the following types.

* `rdf.ObjectValue` - the most generic type and always acceptable
* `rdf.IRI` - it is an error if the incoming value is not of type `rdf.IRI`; pointer values are supported
* `rdf.BlankNode` - it is an error if the incoming value is not of type `rdf.BlankNode`
* `rdf.Literal` - it is an error if the incoming value is not of type `rdf.Literal`; pointer values are supported

Additionally, the following builtin types are supported for incoming `rdf.Literal` value types with special value casting conditions of the lexical form based on its the literal data type; otherwise it is an error. Pointer values are supported.

* `string` - lexical form for `xsd:string`
* `uint8` - strconv.ParseInt of lexical form for `xsd:unsignedByte`
* `uint16` - strconv.ParseInt of lexical form for `xsd:unsignedShort`
* `uint32` - strconv.ParseInt of lexical form for `xsd:unsignedInt`
* `uint64` - strconv.ParseInt of lexical form for `xsd:unsignedLong`
* `int16` - strconv.ParseInt of lexical form for `xsd:short`
* `int32` - strconv.ParseInt of lexical form for `xsd:int`
* `int64` - strconv.ParseInt of lexical form for `xsd:integer`, `xsd:long`
* `float32` - strconv.ParseFloat of lexical form for `xsd:float`
* `float64` - strconv.ParseFloat of lexical form for `xsd:decimal`, `xsd:double`

Additionally, the generic `Collection[T]` type is supported for RDF list traversal.

* `Collection[T]` - if the object value is an IRI or Blank Node, check if it represents a resource from the scope's ResourceListBuilder. If so, and if it offers `rdf:first` and `rdf:rest` properties, follow the [List](https://www.w3.org/TR/rdf11-mt/#rdf-collections) of values, process each list value in the context of the field type `T` (i.e. per "Tagged Object Scalar" section). The `Collection[T]` type is defined as `type Collection[T any] []T`.

Any other type is assumed to be a custom struct which represents a Resource that can be recursively unmarshalled. Custom structs can be unmarshaled from:

* **ObjectStatement** - when the object is an IRI or BlankNode that references a resource with statements in the ResourceListBuilder
* **AnonResourceStatement** - when the object is an inline blank node with embedded properties (no separate subject in the ResourceListBuilder)

# Examples

## RDFC10 Example

Given the following Turtle description.

```ttl
@base <https://example.com/> .
@prefix : <manifest#> .
@prefix rdf:  <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix mf:   <http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#> .
@prefix rdfc: <https://w3c.github.io/rdf-canon/tests/vocab#> .
@prefix rdft: <http://www.w3.org/ns/rdftest#> .

:test004m a rdfc:RDFC10MapTest;
  mf:name "bnode plus embed w/subject (map test)";
  rdfc:computationalComplexity "low";
  rdft:approval rdft:Approved;
  mf:action <rdfc10/test004-in.nq>;
  mf:result <rdfc10/test004-rdfc10map.json>;
  .
```

Given the following user-defined struct and tags.

```go
type Example struct {
  ID                      rdf.SubjectValue `rdf:"s"`
  Type                    rdf.IRI          `rdf:"o,p=rdf:type"`
  Name                    string           `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#name"`
  ComputationalComplexity *string          `rdf:"o,p=https://w3c.github.io/rdf-canon/tests/vocab#computationalComplexity"`
  Approval                rdf.IRI          `rdf:"o,p=http://www.w3.org/ns/rdftest#approval"`
  Action                  rdf.IRI          `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#action"`
  Result                  rdf.IRI          `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#reasult"`
}
```

The following result is expected.

```go
result := Example{
  ID:                      rdf.IRI("manifest#test004m"),
  Type:                    rdf.IRI("https://w3c.github.io/rdf-canon/tests/vocab#RDFC10MapTest"),
  Name:                    "bnode plus embed w/subject (map test)",
  ComputationalComplexity: ptr.Value("low"),
  Approval:                rdf.IRI("http://www.w3.org/ns/rdftest#Approved"),
  Action:                  rdf.IRI("https://example.com/rdfc10/test004-in.nq"),
  Result:                  rdf.IRI("https://example.com/rdfc10/test004-rdfc10map.json"),
}
```

## Jelly Example

Given the following Turtle description.

```ttl
PREFIX jellyt: <https://w3id.org/jelly/dev/tests/vocab#>
PREFIX mf:     <http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#>
PREFIX rdfs:   <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdft:   <http://www.w3.org/ns/rdftest#>

BASE           <https://w3id.org/jelly/dev/tests/rdf/from_jelly/>

<triples_rdf_1_1/pos_015> a jellyt:TestPositive, jellyt:TestRdfFromJelly ;
    mf:name "Four (4) frames, the first frame is empty. Prefix table disabled." ;
    rdft:approval rdft:Proposed ;
    mf:requires jellyt:requirementPhysicalTypeTriples ;
    mf:action <triples_rdf_1_1/pos_015/in.jelly> ;
    mf:result ( 
        <triples_rdf_1_1/pos_015/out_000.nt>
        <triples_rdf_1_1/pos_015/out_001.nt>
        <triples_rdf_1_1/pos_015/out_002.nt>
        <triples_rdf_1_1/pos_015/out_003.nt>
    ) .
```

Given the following user-defined struct and tags.

```go
type Example struct {
  Type     []rdf.IRI          `rdf:"o,p=rdf:type"`
  Name     string             `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#name"`
  Requires []rdf.IRI          `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#requires"`
  Action   rdf.IRI            `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#action"`
  Result   Collection[rdf.IRI] `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#result"`
}
```

The following result is expected.

```go
result := Example{
  Type:                    []rdf.IRI{
    "https://w3id.org/jelly/dev/tests/vocab#TestPositive",
    "https://w3id.org/jelly/dev/tests/vocab#TestRdfFromJelly",
  },
  Name:                    "Four (4) frames, the first frame is empty. Prefix table disabled.",
  Requires:                []rdf.IRI{
    "https://w3id.org/jelly/dev/tests/vocab#requirementPhysicalTypeTriples",
  }
  Action:                  rdf.IRI("https://w3id.org/jelly/dev/tests/rdf/from_jelly/triples_rdf_1_1/pos_015/in.jelly"),
  Result:                  Collection[rdf.IRI]{
    "https://w3id.org/jelly/dev/tests/rdf/from_jelly/triples_rdf_1_1/pos_015/out_000.nt",
    "https://w3id.org/jelly/dev/tests/rdf/from_jelly/triples_rdf_1_1/pos_015/out_001.nt",
    "https://w3id.org/jelly/dev/tests/rdf/from_jelly/triples_rdf_1_1/pos_015/out_002.nt",
    "https://w3id.org/jelly/dev/tests/rdf/from_jelly/triples_rdf_1_1/pos_015/out_003.nt",
  },
}
```

## EARL Example

Given the following Turtle example.

```ttl
PREFIX dawg: <http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#>
PREFIX dc:   <http://purl.org/dc/elements/1.1/>
PREFIX dct:  <http://purl.org/dc/terms/>
PREFIX doap: <http://usefulinc.com/ns/doap#>
PREFIX earl: <http://www.w3.org/ns/earl#>
PREFIX foaf: <http://xmlns.com/foaf/0.1/>
PREFIX rdf:  <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX rdft: <http://www.w3.org/ns/rdftest#>
PREFIX xsd:  <http://www.w3.org/2001/XMLSchema#>

[ rdf:type         earl:Assertion ;
  earl:assertedBy  <http://jena.apache.org/#jena> ;
  earl:mode        earl:automatic ;
  earl:result      [ rdf:type      earl:TestResult ;
                     dc:date       "2021-12-18+00:00"^^xsd:date ;
                     earl:outcome  earl:passed
                   ] ;
  earl:subject     <http://jena.apache.org/#jena> ;
  earl:test        <https://w3c.github.io/rdf-star/tests/trig/eval#trig-star-2>
] .
```

Given the following user-defined struct and tags.

```go
type Assertion struct {
	Subject    rdf.ObjectValue `rdf:"o,p=http://www.w3.org/ns/earl#subject"`
	Test       rdf.ObjectValue `rdf:"o,p=http://www.w3.org/ns/earl#test"`
	AssertedBy rdf.ObjectValue `rdf:"o,p=http://www.w3.org/ns/earl#assertedBy"`
	Mode       rdf.ObjectValue `rdf:"o,p=http://www.w3.org/ns/earl#mode"`
	Result     *TestResult     `rdf:"o,p=http://www.w3.org/ns/earl#result"`
}

type TestResult struct {
	Outcome rdf.IRI     `rdf:"o,p=http://www.w3.org/ns/earl#outcome"`
	Date    rdf.Literal `rdf:"o,p=http://purl.org/dc/elements/1.1/date"`
}
```

The following result is expected.

```go
result := Assertion{
  Subject: rdf.IRI("http://jena.apache.org/#jena"),
  Test: rdf.IRI("https://w3c.github.io/rdf-star/tests/trig/eval#trig-star-2"),
  AssertedBy: rdf.IRI("http://jena.apache.org/#jena"),
  Mode: rdf.IRI("http://www.w3.org/ns/earl#automatic"),
  Result: &TestResult{
    Outcome: rdf.IRI("http://www.w3.org/ns/earl#passed"),
    Date: rdf.Literal{
      LexicalForm: "2021-12-18+00:00",
      Datatype: rdf.IRI("http://www.w3.org/2001/XMLSchema#date"),
    },
  },
}
```
