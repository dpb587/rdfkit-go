package pvspecification

import (
	"strings"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

const (
	ShorthandTextAttributeName_Required  = "required"
	ShorthandTextAttributeName_Value     = "value"
	ShorthandTextAttributeName_Name      = "name"
	ShorthandTextAttributeName_Readonly  = "readonly"
	ShorthandTextAttributeName_Multiple  = "multiple"
	ShorthandTextAttributeName_MinLength = "minlength"
	ShorthandTextAttributeName_MaxLength = "maxlength"
	ShorthandTextAttributeName_Pattern   = "pattern"
	ShorthandTextAttributeName_Min       = "min"
	ShorthandTextAttributeName_Max       = "max"
	ShorthandTextAttributeName_Step      = "step"
)

var ShorthandTextAttributeMappings = map[string]rdf.IRI{
	ShorthandTextAttributeName_Required:  schemairi.ValueRequired_Property,
	ShorthandTextAttributeName_Value:     schemairi.DefaultValue_Property,
	ShorthandTextAttributeName_Name:      schemairi.ValueName_Property,
	ShorthandTextAttributeName_Readonly:  schemairi.ReadonlyValue_Property,
	ShorthandTextAttributeName_Multiple:  schemairi.MultipleValues_Property,
	ShorthandTextAttributeName_MinLength: schemairi.ValueMinLength_Property,
	ShorthandTextAttributeName_MaxLength: schemairi.ValueMaxLength_Property,
	ShorthandTextAttributeName_Pattern:   schemairi.ValuePattern_Property,
	ShorthandTextAttributeName_Min:       schemairi.MinValue_Property,
	ShorthandTextAttributeName_Max:       schemairi.MaxValue_Property,
	ShorthandTextAttributeName_Step:      schemairi.StepValue_Property,
}

//

type ShorthandTextAttribute struct {
	Name  string
	Value *string
}

// GetPredicate returns the predicate mapping (based on [ShorthandTextAttributeMappings]) for the attribute's name. This
// mapping is case-insensitive and converted to lower case.
func (psse ShorthandTextAttribute) GetPredicate() (rdf.IRI, bool) {
	predicate, known := ShorthandTextAttributeMappings[strings.ToLower(psse.Name)]
	if !known {
		return "", false
	}

	return predicate, true
}

//

type ShorthandText struct {
	Attributes []ShorthandTextAttribute
}

func (pst ShorthandText) AsText() string {
	var entryStrings []string

	for _, entry := range pst.Attributes {
		if entry.Value != nil {
			entryStrings = append(entryStrings, entry.Name+"="+*entry.Value)
		} else {
			entryStrings = append(entryStrings, entry.Name)
		}
	}

	return strings.Join(entryStrings, " ")
}

// AsPropertyValueSpecification converts the shorthand text into a resource of type PropertyValueSpecification.
//
// Attributes which do not have a predicate mapping are excluded. Attributes which have a nil value will use an IRI of
// True, and all other values are added as a literal with a Text datatype.
func (pst ShorthandText) AsPropertyValueSpecification() rdfdescription.AnonResource {
	res := rdfdescription.AnonResource{
		Statements: rdfdescription.StatementList{
			rdfdescription.ObjectStatement{
				Predicate: rdfiri.Type_Property,
				Object:    schemairi.PropertyValueSpecification_Class,
			},
		},
	}

	for _, entry := range pst.Attributes {
		predicate, ok := entry.GetPredicate()
		if !ok {
			continue
		}

		if entry.Value == nil {
			// true iff nil; empty string as True for a Text-expected property would not resolve later
			res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
				Predicate: predicate,
				Object:    schemairi.True_Boolean,
			})
		} else {
			// avoiding value type resolution and validation here; caller should implement, if important
			res.Statements = append(res.Statements, rdfdescription.ObjectStatement{
				Predicate: predicate,
				Object: rdf.Literal{
					Datatype:    schemairi.Text_Class,
					LexicalForm: *entry.Value,
				},
			})
		}
	}

	return res
}
