package rdfdescriptionstruct

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

// isCollectionType checks if a type is a Collection[T] type.
// Collection types are detected by checking if the type name starts with "Collection[".
func isCollectionType(t reflect.Type) bool {
	if t.Kind() != reflect.Slice {
		return false
	}
	// Check if the type name indicates it's a Collection type
	// This works because Collection[T] is defined as type Collection[T any] []T
	// and Go's reflection will show the type name as "Collection[...]"
	typeName := t.Name()
	return len(typeName) > 11 && typeName[:11] == "Collection["
}

// unmarshalObject unmarshals object value(s) into a struct field.
func (u *Unmarshaler) unmarshalObject(builder *rdfdescription.ResourceListBuilder, resource rdfdescription.Resource, predicate rdf.IRI, fieldValue reflect.Value, fieldType reflect.Type) error {
	statements := resource.GetResourceStatements().GroupByPredicate()
	objects := statements[predicate]

	if len(objects) == 0 {
		// No values for this predicate
		return nil
	}

	// Check if field is a Collection type
	isCollection := isCollectionType(fieldType)

	// Determine if field is a slice
	if fieldType.Kind() == reflect.Slice {
		return u.unmarshalObjectSlice(builder, objects, fieldValue, fieldType, isCollection)
	}

	// Single value
	if len(objects) > 1 {
		return fmt.Errorf("multiple values found for non-slice field")
	}

	stmt := objects[0]

	// Handle ObjectStatement (normal case: object is IRI, BlankNode, or Literal)
	if objStmt, ok := stmt.(rdfdescription.ObjectStatement); ok {
		return u.unmarshalObjectValue(builder, objStmt.Object, fieldValue, fieldType)
	}

	// Handle AnonResourceStatement (inline blank node with nested properties)
	if anonStmt, ok := stmt.(rdfdescription.AnonResourceStatement); ok {
		// For nested structs, unmarshal the anonymous resource directly
		return u.unmarshalAnonResourceToStruct(builder, anonStmt.AnonResource, fieldValue, fieldType)
	}

	return fmt.Errorf("expected ObjectStatement or AnonResourceStatement, got %T", stmt)
}

// unmarshalObjectSlice unmarshals multiple object values into a slice field.
func (u *Unmarshaler) unmarshalObjectSlice(builder *rdfdescription.ResourceListBuilder, objects rdfdescription.StatementList, fieldValue reflect.Value, fieldType reflect.Type, isCollection bool) error {
	elemType := fieldType.Elem()
	slice := reflect.MakeSlice(fieldType, 0, len(objects))

	for _, stmt := range objects {
		// Handle ObjectStatement
		if objStmt, ok := stmt.(rdfdescription.ObjectStatement); ok {
			// Check if we should expand rdf:List (for Collection types)
			if isCollection && builder != nil {
				// Try to expand as an rdf:List
				expanded, err := u.expandRDFList(builder, objStmt.Object, elemType)
				if err != nil {
					return err
				}
				if expanded.IsValid() {
					// Successfully expanded as a list, append all values
					for i := 0; i < expanded.Len(); i++ {
						slice = reflect.Append(slice, expanded.Index(i))
					}
					continue
				}
			}

			// Not a list or expansion disabled, process as single value
			elemValue := reflect.New(elemType).Elem()
			if err := u.unmarshalObjectValue(builder, objStmt.Object, elemValue, elemType); err != nil {
				return err
			}

			slice = reflect.Append(slice, elemValue)
			continue
		}

		// Handle AnonResourceStatement (might be an RDF list structure)
		if anonStmt, ok := stmt.(rdfdescription.AnonResourceStatement); ok {
			if isCollection && builder != nil {
				// The anonymous resource might be an rdf:List structure
				// Check if this looks like an rdf:List node (has rdf:first and rdf:rest)
				anonStatements := anonStmt.AnonResource.GetResourceStatements()
				stmtsByPred := anonStatements.GroupByPredicate()

				if len(stmtsByPred[rdfiri.First_Property]) > 0 && len(stmtsByPred[rdfiri.Rest_Property]) > 0 {
					// This is an rdf:List node wrapped as AnonResource
					// Traverse it manually since we don't have its subject
					expanded, err := u.expandRDFListFromAnon(builder, anonStmt.AnonResource, elemType)
					if err != nil {
						return err
					}
					if expanded.IsValid() {
						for i := 0; i < expanded.Len(); i++ {
							slice = reflect.Append(slice, expanded.Index(i))
						}
						continue
					}
				}
			}

			// Not a list - could be a nested struct
			if !isCollection {
				elemValue := reflect.New(elemType).Elem()
				if err := u.unmarshalAnonResourceToStruct(builder, anonStmt.AnonResource, elemValue, elemType); err != nil {
					return err
				}
				slice = reflect.Append(slice, elemValue)
				continue
			}

			// Collection type but not an RDF list and no builder
			if isCollection && builder == nil {
				return fmt.Errorf("expected ObjectStatement (Collection type requires ResourceListBuilder - use UnmarshalBuilder instead of Unmarshal)")
			}
			return fmt.Errorf("expected ObjectStatement for non-list value (got AnonResourceStatement)")
		}

		// Unknown statement type
		if isCollection && builder == nil {
			return fmt.Errorf("expected ObjectStatement (Collection type requires ResourceListBuilder - use UnmarshalBuilder instead of Unmarshal)")
		}
		return fmt.Errorf("expected ObjectStatement or AnonResourceStatement, got %T", stmt)
	}

	fieldValue.Set(slice)
	return nil
}

// unmarshalObjectValue unmarshals a single object value.
func (u *Unmarshaler) unmarshalObjectValue(builder *rdfdescription.ResourceListBuilder, object rdf.ObjectValue, fieldValue reflect.Value, fieldType reflect.Type) error {
	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		if object == nil {
			fieldValue.Set(reflect.Zero(fieldType))
			return nil
		}

		ptr := reflect.New(fieldType.Elem())
		if err := u.unmarshalObjectValue(builder, object, ptr.Elem(), fieldType.Elem()); err != nil {
			return err
		}
		fieldValue.Set(ptr)
		return nil
	}

	// Handle RDF types
	switch fieldType {
	case reflect.TypeOf((*rdf.ObjectValue)(nil)).Elem():
		fieldValue.Set(reflect.ValueOf(object))
		return nil

	case reflect.TypeOf(rdf.IRI("")):
		iri, ok := object.(rdf.IRI)
		if !ok {
			return &TypeMismatchError{Expected: "rdf.IRI", Got: reflect.TypeOf(object).String()}
		}
		fieldValue.Set(reflect.ValueOf(iri))
		return nil

	case reflect.TypeOf((*rdf.BlankNode)(nil)).Elem():
		bn, ok := object.(rdf.BlankNode)
		if !ok {
			return &TypeMismatchError{Expected: "rdf.BlankNode", Got: reflect.TypeOf(object).String()}
		}
		fieldValue.Set(reflect.ValueOf(bn))
		return nil

	case reflect.TypeOf(rdf.Literal{}):
		lit, ok := object.(rdf.Literal)
		if !ok {
			return &TypeMismatchError{Expected: "rdf.Literal", Got: reflect.TypeOf(object).String()}
		}
		fieldValue.Set(reflect.ValueOf(lit))
		return nil
	}

	// Handle builtin Go types (must be from Literal)
	lit, ok := object.(rdf.Literal)
	if ok {
		return unmarshalLiteralToBuiltin(lit, fieldValue, fieldType)
	}

	// For any other type, if the object is an IRI or BlankNode,
	// treat it as a custom resource type that can be recursively unmarshaled
	var subject rdf.SubjectValue
	switch obj := object.(type) {
	case rdf.IRI:
		subject = obj
	case rdf.BlankNode:
		subject = obj
	default:
		return &TypeMismatchError{Expected: "rdf.IRI or rdf.BlankNode for custom resource type", Got: reflect.TypeOf(object).String()}
	}

	if builder == nil {
		return fmt.Errorf("builder required for resource unmarshaling (custom struct type %s)", fieldType.String())
	}

	// Get the resource for this subject
	subResource := &rdfdescription.SubjectResource{
		Subject:    subject,
		Statements: builder.GetSubjectStatements(subject),
	}

	// Recursively unmarshal into the field
	return u.unmarshalResource(builder, subResource, fieldValue.Addr().Interface())
}

// unmarshalAnonResourceToStruct unmarshals an AnonResource into a struct field.
// This is used when an object value is an inline blank node (AnonResourceStatement)
// rather than a reference to a blank node or IRI (ObjectStatement).
func (u *Unmarshaler) unmarshalAnonResourceToStruct(builder *rdfdescription.ResourceListBuilder, anon rdfdescription.AnonResource, fieldValue reflect.Value, fieldType reflect.Type) error {
	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		ptr := reflect.New(fieldType.Elem())
		if err := u.unmarshalAnonResourceToStruct(builder, anon, ptr.Elem(), fieldType.Elem()); err != nil {
			return err
		}
		fieldValue.Set(ptr)
		return nil
	}

	// The field type should be a struct for nested resources
	if fieldType.Kind() != reflect.Struct {
		return fmt.Errorf("cannot unmarshal AnonResource into non-struct type %s", fieldType.String())
	}

	// Recursively unmarshal the anonymous resource as a Resource
	return u.unmarshalResource(builder, anon, fieldValue.Addr().Interface())
}

// unmarshalLiteralToBuiltin unmarshals a Literal to a builtin Go type.
func unmarshalLiteralToBuiltin(lit rdf.Literal, fieldValue reflect.Value, fieldType reflect.Type) error {
	switch fieldType.Kind() {
	case reflect.String:
		if lit.Datatype != xsdiri.String_Datatype {
			return &TypeMismatchError{Expected: "xsd:string", Got: string(lit.Datatype)}
		}
		fieldValue.SetString(lit.LexicalForm)
		return nil

	case reflect.Uint8:
		if lit.Datatype != xsdiri.UnsignedByte_Datatype {
			return &TypeMismatchError{Expected: "xsd:unsignedByte", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseUint(lit.LexicalForm, 10, 8)
		if err != nil {
			return fmt.Errorf("parse uint8: %w", err)
		}
		fieldValue.SetUint(val)
		return nil

	case reflect.Uint16:
		if lit.Datatype != xsdiri.UnsignedShort_Datatype {
			return &TypeMismatchError{Expected: "xsd:unsignedShort", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseUint(lit.LexicalForm, 10, 16)
		if err != nil {
			return fmt.Errorf("parse uint16: %w", err)
		}
		fieldValue.SetUint(val)
		return nil

	case reflect.Uint32:
		if lit.Datatype != xsdiri.UnsignedInt_Datatype {
			return &TypeMismatchError{Expected: "xsd:unsignedInt", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseUint(lit.LexicalForm, 10, 32)
		if err != nil {
			return fmt.Errorf("parse uint32: %w", err)
		}
		fieldValue.SetUint(val)
		return nil

	case reflect.Uint64:
		if lit.Datatype != xsdiri.UnsignedLong_Datatype {
			return &TypeMismatchError{Expected: "xsd:unsignedLong", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseUint(lit.LexicalForm, 10, 64)
		if err != nil {
			return fmt.Errorf("parse uint64: %w", err)
		}
		fieldValue.SetUint(val)
		return nil

	case reflect.Int16:
		if lit.Datatype != xsdiri.Short_Datatype {
			return &TypeMismatchError{Expected: "xsd:short", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseInt(lit.LexicalForm, 10, 16)
		if err != nil {
			return fmt.Errorf("parse int16: %w", err)
		}
		fieldValue.SetInt(val)
		return nil

	case reflect.Int32:
		if lit.Datatype != xsdiri.Int_Datatype {
			return &TypeMismatchError{Expected: "xsd:int", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseInt(lit.LexicalForm, 10, 32)
		if err != nil {
			return fmt.Errorf("parse int32: %w", err)
		}
		fieldValue.SetInt(val)
		return nil

	case reflect.Int64:
		if lit.Datatype != xsdiri.Integer_Datatype && lit.Datatype != xsdiri.Long_Datatype {
			return &TypeMismatchError{Expected: "xsd:integer or xsd:long", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseInt(lit.LexicalForm, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int64: %w", err)
		}
		fieldValue.SetInt(val)
		return nil

	case reflect.Float32:
		if lit.Datatype != xsdiri.Float_Datatype {
			return &TypeMismatchError{Expected: "xsd:float", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseFloat(lit.LexicalForm, 32)
		if err != nil {
			return fmt.Errorf("parse float32: %w", err)
		}
		fieldValue.SetFloat(val)
		return nil

	case reflect.Float64:
		if lit.Datatype != xsdiri.Decimal_Datatype && lit.Datatype != xsdiri.Double_Datatype {
			return &TypeMismatchError{Expected: "xsd:decimal or xsd:double", Got: string(lit.Datatype)}
		}
		val, err := strconv.ParseFloat(lit.LexicalForm, 64)
		if err != nil {
			return fmt.Errorf("parse float64: %w", err)
		}
		fieldValue.SetFloat(val)
		return nil

	default:
		return &TypeMismatchError{Expected: "supported builtin type", Got: fieldType.String()}
	}
}

// expandRDFList attempts to expand an rdf:List into a slice of values.
// Returns nil if the object is not an rdf:List or is rdf:nil.
func (u *Unmarshaler) expandRDFList(builder *rdfdescription.ResourceListBuilder, object rdf.ObjectValue, elemType reflect.Type) (reflect.Value, error) {
	// Object must be IRI or BlankNode to potentially be an rdf:List
	var subject rdf.SubjectValue
	switch obj := object.(type) {
	case rdf.IRI:
		subject = obj
	case rdf.BlankNode:
		subject = obj
	default:
		// Not a subject type, can't be a list
		return reflect.Value{}, nil
	}

	// Check if it's rdf:nil
	if iri, ok := subject.(rdf.IRI); ok && iri == rdfiri.Nil_List {
		// Empty list
		return reflect.Value{}, nil
	}

	// Get the resource for this subject
	statements := builder.GetSubjectStatements(subject)
	if len(statements) == 0 {
		// No statements, not a list
		return reflect.Value{}, nil
	}

	// Group statements by predicate
	stmtsByPred := statements.GroupByPredicate()

	// Check if this looks like an rdf:List structure
	// A list node should have rdf:first and rdf:rest properties
	// (Some RDF parsers don't generate explicit rdf:type rdf:List statements)
	hasFirst := len(stmtsByPred[rdfiri.First_Property]) > 0
	hasRest := len(stmtsByPred[rdfiri.Rest_Property]) > 0

	if !hasFirst || !hasRest {
		// Not an rdf:List structure
		return reflect.Value{}, nil
	}

	// It's an rdf:List, traverse it
	result := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
	current := subject

	for {
		// Check for rdf:nil
		if iri, ok := current.(rdf.IRI); ok && iri == rdfiri.Nil_List {
			break
		}

		// Get statements for current node
		currentStmts := builder.GetSubjectStatements(current).GroupByPredicate()

		// Get rdf:first
		firstStmts := currentStmts[rdfiri.First_Property]
		if len(firstStmts) == 0 {
			return reflect.Value{}, fmt.Errorf("rdf:List node missing rdf:first")
		}
		if len(firstStmts) > 1 {
			return reflect.Value{}, fmt.Errorf("rdf:List node has multiple rdf:first values")
		}

		firstObjStmt, ok := firstStmts[0].(rdfdescription.ObjectStatement)
		if !ok {
			return reflect.Value{}, fmt.Errorf("expected ObjectStatement for rdf:first, got %T", firstStmts[0])
		}

		// Unmarshal the first value
		elemValue := reflect.New(elemType).Elem()
		if err := u.unmarshalObjectValue(builder, firstObjStmt.Object, elemValue, elemType); err != nil {
			return reflect.Value{}, fmt.Errorf("unmarshal rdf:List value: %w", err)
		}
		result = reflect.Append(result, elemValue)

		// Get rdf:rest
		restStmts := currentStmts[rdfiri.Rest_Property]
		if len(restStmts) == 0 {
			return reflect.Value{}, fmt.Errorf("rdf:List node missing rdf:rest")
		}
		if len(restStmts) > 1 {
			return reflect.Value{}, fmt.Errorf("rdf:List node has multiple rdf:rest values")
		}

		restObjStmt, ok := restStmts[0].(rdfdescription.ObjectStatement)
		if !ok {
			return reflect.Value{}, fmt.Errorf("expected ObjectStatement for rdf:rest, got %T", restStmts[0])
		}

		// Move to next node
		switch rest := restObjStmt.Object.(type) {
		case rdf.IRI:
			current = rest
		case rdf.BlankNode:
			current = rest
		default:
			return reflect.Value{}, fmt.Errorf("rdf:rest must be IRI or BlankNode, got %T", restObjStmt.Object)
		}
	}

	return result, nil
}

// expandRDFListFromAnon expands an RDF list starting from an AnonResource.
// This handles the case where the list head is wrapped as an AnonResourceStatement
// instead of an ObjectStatement (which happens when the blank node has only 1 reference).
func (u *Unmarshaler) expandRDFListFromAnon(builder *rdfdescription.ResourceListBuilder, anon rdfdescription.AnonResource, elemType reflect.Type) (reflect.Value, error) {
	result := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
	currentAnon := &anon

	for {
		if currentAnon == nil {
			break
		}

		currentStmts := currentAnon.GetResourceStatements().GroupByPredicate()

		// Get rdf:first
		firstStmts := currentStmts[rdfiri.First_Property]
		if len(firstStmts) == 0 {
			return reflect.Value{}, fmt.Errorf("rdf:List node missing rdf:first")
		}
		if len(firstStmts) > 1 {
			return reflect.Value{}, fmt.Errorf("rdf:List node has multiple rdf:first values")
		}

		firstObjStmt, ok := firstStmts[0].(rdfdescription.ObjectStatement)
		if !ok {
			return reflect.Value{}, fmt.Errorf("expected ObjectStatement for rdf:first, got %T", firstStmts[0])
		}

		// Unmarshal the first value
		elemValue := reflect.New(elemType).Elem()
		if err := u.unmarshalObjectValue(builder, firstObjStmt.Object, elemValue, elemType); err != nil {
			return reflect.Value{}, fmt.Errorf("unmarshal rdf:List value: %w", err)
		}
		result = reflect.Append(result, elemValue)

		// Get rdf:rest
		restStmts := currentStmts[rdfiri.Rest_Property]
		if len(restStmts) == 0 {
			return reflect.Value{}, fmt.Errorf("rdf:List node missing rdf:rest")
		}
		if len(restStmts) > 1 {
			return reflect.Value{}, fmt.Errorf("rdf:List node has multiple rdf:rest values")
		}

		// Check what type of statement rest is
		switch restStmt := restStmts[0].(type) {
		case rdfdescription.ObjectStatement:
			// rest points to an IRI (rdf:nil) or BlankNode
			if iri, ok := restStmt.Object.(rdf.IRI); ok && iri == rdfiri.Nil_List {
				// End of list
				return result, nil
			} else if bn, ok := restStmt.Object.(rdf.BlankNode); ok {
				// Continue with next node - need to get its statements from builder
				nextStmts := builder.GetSubjectStatements(bn)
				if len(nextStmts) == 0 {
					return reflect.Value{}, fmt.Errorf("blank node in rdf:rest has no statements")
				}
				// But we can't easily create an AnonResource from statements...
				// We need to use expandRDFList instead since we have a subject now
				remaining, err := u.expandRDFList(builder, bn, elemType)
				if err != nil {
					return reflect.Value{}, err
				}
				if !remaining.IsValid() {
					return reflect.Value{}, fmt.Errorf("failed to expand remaining list from blank node")
				}
				// Append all remaining values
				for i := 0; i < remaining.Len(); i++ {
					result = reflect.Append(result, remaining.Index(i))
				}
				return result, nil
			} else {
				return reflect.Value{}, fmt.Errorf("rdf:rest must be rdf:nil, IRI, or BlankNode, got %T", restStmt.Object)
			}
		case rdfdescription.AnonResourceStatement:
			// rest points to another anonymous resource (next list node)
			next := restStmt.AnonResource
			currentAnon = &next
		default:
			return reflect.Value{}, fmt.Errorf("unexpected statement type for rdf:rest: %T", restStmt)
		}
	}

	return result, nil
}
