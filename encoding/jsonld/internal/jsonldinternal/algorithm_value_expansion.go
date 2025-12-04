package jsonldinternal

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
)

type algorithmValueExpansion struct {
	activeContext  *Context
	activeProperty string
	value          inspectjson.Value

	// [dpb] for tracking source offsets
	activePropertySourceOffsets *cursorio.TextOffsetRange

	processor *contextProcessor
}

func (vars algorithmValueExpansion) Call() *ExpandedObject {
	// [spec // 5.3.2 // 1] If the *active property* has a type mapping in *active context* that is `@id`, and the value is a string, return a new map containing a single entry where the key is @id and the value is the result IRI expanding value using true for document relative and false for vocab.

	activePropertyTermDefinition := vars.activeContext.TermDefinitions[vars.activeProperty]

	if activePropertyTermDefinition != nil && activePropertyTermDefinition.TypeMapping == ExpandedIRIasKeyword("@id") {
		expandedValue, err := algorithmIRIExpansion{
			value:            vars.value,
			documentRelative: true,
			vocab:            false,
			//
			activeContext: vars.activeContext,
		}.Call()
		if err != nil {
			panic(err)
		}

		// [dpb] this is not defined by the spec, but supports test case #t0088
		if rawValue, ok := expandedValue.(ExpandedIRIasRawValue); ok {
			return &ExpandedObject{
				Members: map[string]ExpandedValue{
					"@value": &ExpandedScalarPrimitive{
						Value: rawValue.Value,
					},
				},
				PropertySourceOffsets: vars.activePropertySourceOffsets,
			}
		}

		return &ExpandedObject{
			Members: map[string]ExpandedValue{
				"@id": expandedValue.NewPropertyValue(nil, vars.value.GetSourceOffsets()),
			},
			PropertySourceOffsets: vars.activePropertySourceOffsets,
		}
	}

	// [spec // 5.3.2 // 2] If *active property* has a type mapping in *active context* that is `@vocab`, and the *value* is a string, return a new map containing a single entry where the key is `@id` and the value is the result of IRI expanding *value* using `true` for *document relative*.

	if activePropertyTermDefinition != nil && activePropertyTermDefinition.TypeMapping == ExpandedIRIasKeyword("@vocab") {
		expandedValue, err := algorithmIRIExpansion{
			value:            vars.value,
			documentRelative: true,
			vocab:            true, // not explicitly mentioned by spec?
			//
			activeContext: vars.activeContext,
		}.Call()
		if err != nil {
			panic(err)
		}

		// [dpb] this is not defined by the spec, but supports test case #t0088
		if rawValue, ok := expandedValue.(ExpandedIRIasRawValue); ok {
			return &ExpandedObject{
				Members: map[string]ExpandedValue{
					"@value": &ExpandedScalarPrimitive{
						Value: rawValue.Value,
					},
				},
				PropertySourceOffsets: vars.activePropertySourceOffsets,
			}
		}

		return &ExpandedObject{
			Members: map[string]ExpandedValue{
				"@id": expandedValue.NewPropertyValue(nil, vars.value.GetSourceOffsets()),
			},
			PropertySourceOffsets: vars.activePropertySourceOffsets,
		}
	}

	// [spec // 5.3.2 // 3] Otherwise, initialize *result* to a map with an `@value` entry whose value is set to *value*.

	result := &ExpandedObject{
		Members: map[string]ExpandedValue{
			"@value": &ExpandedScalarPrimitive{
				Value: vars.value,
			},
		},
		PropertySourceOffsets: vars.activePropertySourceOffsets,
	}

	// [spec // 5.3.2 // 4] If *active property* has a type mapping in *active context*, other than `@id`, `@vocab`, or `@none`, add `@type` to *result* and set its value to the value associated with the type mapping.

	if activePropertyTermDefinition != nil && activePropertyTermDefinition.TypeMapping != nil && (activePropertyTermDefinition.TypeMapping != ExpandedIRIasKeyword("@id") &&
		activePropertyTermDefinition.TypeMapping != ExpandedIRIasKeyword("@vocab") &&
		activePropertyTermDefinition.TypeMapping != ExpandedIRIasKeyword("@none")) {

		result.Members["@type"] = activePropertyTermDefinition.TypeMapping.NewPropertyValue(
			nil,
			activePropertyTermDefinition.TypeMappingValue.GetSourceOffsets(),
		)
	} else {

		// [spec // 5.3.2 // 5] Otherwise, if *value* is a string:

		if _, ok := vars.value.(inspectjson.StringValue); ok {

			// [spec // 5.3.2 // 5.1] Initialize language to the language mapping for active property in active context, if any, otherwise to the default language of active context.

			var language inspectjson.Value

			if activePropertyTermDefinition != nil && activePropertyTermDefinition.LanguageMappingValue != nil {
				language = activePropertyTermDefinition.LanguageMappingValue
			} else if vars.activeContext.DefaultLanguageValue != nil {
				language = *vars.activeContext.DefaultLanguageValue
			}

			// [spec // 5.3.2 // 5.3] If language is not null, add @language to result with the value language.

			if _, ok := language.(inspectjson.StringValue); ok {
				result.Members["@language"] = &ExpandedScalarPrimitive{
					Value: language,
				}
			}

			// [spec // 5.3.2 // 5.2] Initialize direction to the direction mapping for active property in active context, if any, otherwise to the default base direction of active context.

			var direction inspectjson.Value

			if activePropertyTermDefinition != nil && activePropertyTermDefinition.DirectionMappingValue != nil {
				direction = activePropertyTermDefinition.DirectionMappingValue
			} else if vars.activeContext.DefaultDirectionValue != nil {
				direction = *vars.activeContext.DefaultDirectionValue
			}

			// [spec // 5.3.2 // 5.4] If direction is not null, add @direction to result with the value direction.

			if _, ok := direction.(inspectjson.StringValue); ok {
				result.Members["@direction"] = &ExpandedScalarPrimitive{
					Value: direction,
				}
			}
		}
	}

	return result
}
