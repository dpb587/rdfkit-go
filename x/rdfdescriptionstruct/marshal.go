package rdfdescriptionstruct

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdtype"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

// Marshal converts a struct into an rdfdescription.ResourceList.
// The struct fields must have rdf tags indicating how to map to RDF data.
// Returns a ResourceList containing the main resource and any nested resources.
func Marshal(from any) (rdfdescription.ResourceList, error) {
	resources, err := marshalResources(from)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// marshalResources performs the actual marshaling work, collecting all resources.
func marshalResources(from any) (rdfdescription.ResourceList, error) {
	allResources := make(rdfdescription.ResourceList, 0)
	mainResource, nestedResources, err := marshalResource(from)
	if err != nil {
		return nil, err
	}
	allResources = append(allResources, mainResource)
	allResources = append(allResources, nestedResources...)
	return allResources, nil
}

// marshalResource performs the actual marshaling work.
// Returns the main resource, any nested resources, and an error.
func marshalResource(from any) (*rdfdescription.SubjectResource, rdfdescription.ResourceList, error) {
	// from must be a struct or pointer to struct
	fromValue := reflect.ValueOf(from)
	if fromValue.Kind() == reflect.Ptr {
		if fromValue.IsNil() {
			return nil, nil, fmt.Errorf("from must not be nil")
		}
		fromValue = fromValue.Elem()
	}

	if fromValue.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("from must be a struct or pointer to struct, got %s", fromValue.Kind())
	}

	fromType := fromValue.Type()

	var subject rdf.SubjectValue
	var statements []rdfdescription.Statement
	var nestedResources rdfdescription.ResourceList

	// Iterate over struct fields
	for i := 0; i < fromValue.NumField(); i++ {
		field := fromType.Field(i)
		fieldValue := fromValue.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get rdf tag
		tag := field.Tag.Get("rdf")
		if tag == "" {
			continue
		}

		// Parse tag
		parsedTag, err := parseTag(tag, defaultPrefixMap)
		if err != nil {
			return nil, nil, &InvalidTagError{Tag: tag, Err: err}
		}

		// Handle subject tag
		if parsedTag.kind == "s" {
			subjectValue, err := marshalSubject(fieldValue)
			if err != nil {
				return nil, nil, &UnmarshalError{Field: field.Name, Err: err}
			}
			subject = subjectValue
			continue
		}

		// Handle object tag
		if parsedTag.kind == "o" {
			objectStatements, objectNestedResources, err := marshalObject(subject, parsedTag.predicate, fieldValue)
			if err != nil {
				return nil, nil, &UnmarshalError{Field: field.Name, Err: err}
			}
			statements = append(statements, objectStatements...)
			nestedResources = append(nestedResources, objectNestedResources...)
		}
	}

	if subject == nil {
		return nil, nil, fmt.Errorf("no subject field found (must have rdf:\"s\" tag)")
	}

	return &rdfdescription.SubjectResource{
		Subject:    subject,
		Statements: statements,
	}, nestedResources, nil
}

// marshalSubject converts a field value to a SubjectValue.
func marshalSubject(fieldValue reflect.Value) (rdf.SubjectValue, error) {
	// Handle pointer types
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return nil, fmt.Errorf("subject field is nil")
		}
		fieldValue = fieldValue.Elem()
	}

	// Check if it's already a SubjectValue type
	if fieldValue.Type().AssignableTo(reflect.TypeOf((*rdf.SubjectValue)(nil)).Elem()) {
		return fieldValue.Interface().(rdf.SubjectValue), nil
	}

	// Check specific types
	switch v := fieldValue.Interface().(type) {
	case rdf.IRI:
		return v, nil
	case rdf.BlankNode:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported subject type: %T", v)
	}
}

// marshalObject converts a field value to one or more statements.
// Returns statements for this field and any nested resources.
func marshalObject(subject rdf.SubjectValue, predicate rdf.IRI, fieldValue reflect.Value) ([]rdfdescription.Statement, rdfdescription.ResourceList, error) {
	if subject == nil {
		return nil, nil, fmt.Errorf("subject must be set before marshaling objects")
	}

	// Handle pointer types
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			// Nil pointers produce no statements
			return nil, nil, nil
		}
		fieldValue = fieldValue.Elem()
	}

	// Handle slices
	if fieldValue.Kind() == reflect.Slice {
		var statements []rdfdescription.Statement
		var nestedResources rdfdescription.ResourceList
		for i := 0; i < fieldValue.Len(); i++ {
			stmts, nested, err := marshalObjectValue(subject, predicate, fieldValue.Index(i))
			if err != nil {
				return nil, nil, err
			}
			statements = append(statements, stmts...)
			nestedResources = append(nestedResources, nested...)
		}
		return statements, nestedResources, nil
	}

	// Handle single value
	return marshalObjectValue(subject, predicate, fieldValue)
}

// marshalObjectValue converts a single value to a statement.
// Returns statements for this value and any nested resources.
func marshalObjectValue(subject rdf.SubjectValue, predicate rdf.IRI, fieldValue reflect.Value) ([]rdfdescription.Statement, rdfdescription.ResourceList, error) {
	// Handle pointer types
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return nil, nil, nil
		}
		fieldValue = fieldValue.Elem()
	}

	// Check if this is a struct (custom resource type) that should be recursively marshaled
	if fieldValue.Kind() == reflect.Struct {
		// Check if it's one of the known RDF types that should NOT be recursively marshaled
		fieldType := fieldValue.Type()
		if fieldType == reflect.TypeOf(rdf.Literal{}) {
			// It's an rdf.Literal, handle as ObjectValue
			objectValue, err := marshalToObjectValue(fieldValue)
			if err != nil {
				return nil, nil, err
			}
			return []rdfdescription.Statement{
				rdfdescription.ObjectStatement{
					Predicate: predicate,
					Object:    objectValue,
				},
			}, nil, nil
		}

		// For other struct types, recursively marshal as a nested resource
		nestedResource, deeperNested, err := marshalResource(fieldValue.Interface())
		if err != nil {
			return nil, nil, fmt.Errorf("marshal nested resource: %w", err)
		}

		// Create statement linking to the nested resource
		statement := rdfdescription.ObjectStatement{
			Predicate: predicate,
			Object:    nestedResource.Subject,
		}

		// Return the linking statement, the nested resource itself, and any deeper nested resources
		allNested := append(rdfdescription.ResourceList{nestedResource}, deeperNested...)
		return []rdfdescription.Statement{statement}, allNested, nil
	}

	// Convert field value to ObjectValue
	objectValue, err := marshalToObjectValue(fieldValue)
	if err != nil {
		return nil, nil, err
	}

	return []rdfdescription.Statement{
		rdfdescription.ObjectStatement{
			Predicate: predicate,
			Object:    objectValue,
		},
	}, nil, nil
}

// marshalToObjectValue converts a field value to an ObjectValue.
func marshalToObjectValue(fieldValue reflect.Value) (rdf.ObjectValue, error) {
	// Handle pointer types
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return nil, fmt.Errorf("cannot marshal nil pointer")
		}
		fieldValue = fieldValue.Elem()
	}

	// Check if it's already an ObjectValue type
	if fieldValue.Type().AssignableTo(reflect.TypeOf((*rdf.ObjectValue)(nil)).Elem()) {
		return fieldValue.Interface().(rdf.ObjectValue), nil
	}

	// Check specific RDF types
	switch v := fieldValue.Interface().(type) {
	case rdf.IRI:
		return v, nil
	case rdf.BlankNode:
		return v, nil
	case rdf.Literal:
		return v, nil
	}

	// Handle builtin types - convert to Literal
	return marshalBuiltinToLiteral(fieldValue)
}

// marshalBuiltinToLiteral creates a Literal from a builtin Go type.
func marshalBuiltinToLiteral(fieldValue reflect.Value) (rdf.ObjectValue, error) {
	switch fieldValue.Kind() {
	case reflect.String:
		return xsdtype.String(fieldValue.String()).AsObjectValue(), nil

	case reflect.Uint8:
		return xsdtype.UnsignedByte(uint8(fieldValue.Uint())).AsObjectValue(), nil
	case reflect.Uint16:
		return xsdtype.UnsignedShort(uint16(fieldValue.Uint())).AsObjectValue(), nil
	case reflect.Uint32:
		return xsdtype.UnsignedInt(uint32(fieldValue.Uint())).AsObjectValue(), nil
	case reflect.Uint64:
		return xsdtype.UnsignedLong(fieldValue.Uint()).AsObjectValue(), nil

	case reflect.Int16:
		return xsdtype.Short(int16(fieldValue.Int())).AsObjectValue(), nil
	case reflect.Int32:
		return xsdtype.Int(int32(fieldValue.Int())).AsObjectValue(), nil
	case reflect.Int64:
		// Use xsd:integer for int64
		lexical := strconv.FormatInt(fieldValue.Int(), 10)
		return rdf.Literal{
			Datatype:    xsdiri.Integer_Datatype,
			LexicalForm: lexical,
		}, nil

	case reflect.Float32:
		return xsdtype.Float(float32(fieldValue.Float())).AsObjectValue(), nil
	case reflect.Float64:
		return xsdtype.Double(fieldValue.Float()).AsObjectValue(), nil

	default:
		return rdf.Literal{}, fmt.Errorf("unsupported builtin type: %s", fieldValue.Kind())
	}
}
