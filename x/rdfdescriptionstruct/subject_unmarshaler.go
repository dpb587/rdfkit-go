package rdfdescriptionstruct

import (
	"reflect"

	"github.com/dpb587/rdfkit-go/rdf"
)

// unmarshalSubject unmarshals a subject value into a struct field.
func unmarshalSubject(subject rdf.SubjectValue, fieldValue reflect.Value, fieldType reflect.Type) error {
	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		if subject == nil {
			// Set nil pointer
			fieldValue.Set(reflect.Zero(fieldType))
			return nil
		}

		// Create new pointer value
		ptr := reflect.New(fieldType.Elem())
		if err := unmarshalSubject(subject, ptr.Elem(), fieldType.Elem()); err != nil {
			return err
		}
		fieldValue.Set(ptr)
		return nil
	}

	// Handle non-pointer types
	switch fieldType {
	case reflect.TypeOf((*rdf.SubjectValue)(nil)).Elem():
		// rdf.SubjectValue interface - accept anything
		fieldValue.Set(reflect.ValueOf(subject))
		return nil

	case reflect.TypeOf(rdf.IRI("")):
		// rdf.IRI
		iri, ok := subject.(rdf.IRI)
		if !ok {
			return &TypeMismatchError{Expected: "rdf.IRI", Got: reflect.TypeOf(subject).String()}
		}
		fieldValue.Set(reflect.ValueOf(iri))
		return nil

	case reflect.TypeOf((*rdf.BlankNode)(nil)).Elem():
		// rdf.BlankNode interface
		bn, ok := subject.(rdf.BlankNode)
		if !ok {
			return &TypeMismatchError{Expected: "rdf.BlankNode", Got: reflect.TypeOf(subject).String()}
		}
		fieldValue.Set(reflect.ValueOf(bn))
		return nil

	default:
		return &TypeMismatchError{
			Expected: "rdf.SubjectValue, rdf.IRI, or rdf.BlankNode",
			Got:      fieldType.String(),
		}
	}
}
