package rdfdescriptionstruct

import (
	"fmt"
	"reflect"

	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

// Unmarshaler provides unmarshaling with custom configuration.
type Unmarshaler struct {
	prefixes iri.PrefixMapper
}

// NewUnmarshaler creates a new Unmarshaler with optional configuration.
func NewUnmarshaler(options ...UnmarshalerOption) *Unmarshaler {
	u := &Unmarshaler{
		prefixes: defaultPrefixMap,
	}

	for _, opt := range options {
		opt.apply(u)
	}

	return u
}

// UnmarshalResource unmarshals an rdfdescription.Resource into a struct.
// The struct fields must have rdf tags indicating how to map RDF data.
func (u *Unmarshaler) UnmarshalResource(from rdfdescription.Resource, to any) error {
	return u.Unmarshal(nil, from, to)
}

// Unmarshal unmarshals an rdfdescription.Resource into a struct using a ResourceListBuilder.
// The builder is required for recursive unmarshaling of resource-type object values.
func (u *Unmarshaler) Unmarshal(builder *rdfdescription.ResourceListBuilder, from rdfdescription.Resource, to any) error {
	return u.unmarshalResource(builder, from, to)
}

// unmarshalResourceWithPrefixes performs the actual unmarshaling work with a custom prefix map.
func (u *Unmarshaler) unmarshalResource(builder *rdfdescription.ResourceListBuilder, resource rdfdescription.Resource, to any) error {
	// to must be a pointer to a struct
	toValue := reflect.ValueOf(to)
	if toValue.Kind() != reflect.Ptr {
		return fmt.Errorf("to must be a pointer, got %s", toValue.Kind())
	}

	toValue = toValue.Elem()
	if toValue.Kind() != reflect.Struct {
		return fmt.Errorf("to must be a pointer to struct, got pointer to %s", toValue.Kind())
	}

	toType := toValue.Type()

	// Iterate over struct fields
	for i := 0; i < toValue.NumField(); i++ {
		field := toType.Field(i)
		fieldValue := toValue.Field(i)

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
		tagInfo, err := parseTag(tag, u.prefixes)
		if err != nil {
			return &UnmarshalError{
				Field: field.Name,
				Err:   &InvalidTagError{Tag: tag, Err: err},
			}
		}

		if tagInfo == nil {
			continue
		}

		// Handle based on tag kind
		switch tagInfo.kind {
		case "s":
			// Subject field
			subject := resource.GetResourceSubject()
			if err := unmarshalSubject(subject, fieldValue, field.Type); err != nil {
				return &UnmarshalError{Field: field.Name, Err: err}
			}

		case "o":
			// Object field
			if err := u.unmarshalObject(builder, resource, tagInfo.predicate, fieldValue, field.Type); err != nil {
				return &UnmarshalError{Field: field.Name, Err: err}
			}

		default:
			return &UnmarshalError{
				Field: field.Name,
				Err:   fmt.Errorf("unknown tag kind: %s", tagInfo.kind),
			}
		}
	}

	return nil
}
