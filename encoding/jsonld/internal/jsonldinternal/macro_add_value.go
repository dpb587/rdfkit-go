package jsonldinternal

import (
	"github.com/dpb587/inspectjson-go/inspectjson"
)

type macroAddValue struct {
	Value   inspectjson.Value
	Object  *inspectjson.ObjectValue
	Key     string
	AsArray bool
}

func (m macroAddValue) Call() {
	// [spec // 1] If *as array* is `true` and the value of *key* in *object* does not exist or is not an array, set it to a new array containing any original value.

	if m.AsArray {
		originalValue, ok := m.Object.Members[m.Key]
		if !ok {
			m.Object.Members[m.Key] = inspectjson.ObjectMember{
				Name: inspectjson.StringValue{
					Value: m.Key,
				},
				Value: inspectjson.ArrayValue{},
			}
		} else if _, ok := originalValue.Value.(inspectjson.ArrayValue); !ok {
			m.Object.Members[m.Key] = inspectjson.ObjectMember{
				Name: originalValue.Name,
				Value: inspectjson.ArrayValue{
					Values: []inspectjson.Value{
						originalValue.Value,
					},
				},
			}
		}
	}

	// [spec // 2] If *value* is an array, then for each element *v* in *value*, use add value recursively to add *v* to *key* in *entry*.

	if arrayValue, ok := m.Value.(inspectjson.ArrayValue); ok {
		for _, v := range arrayValue.Values {
			macroAddValue{
				Value:   v,
				Object:  m.Object,
				Key:     m.Key,
				AsArray: m.AsArray,
			}.Call()
		}
	} else {

		// [spec // 3] Otherwise:

		// [spec // 3.1] If *key* is not an entry in *object*, add *value* as the value of *key* in *object*.

		originalValue, ok := m.Object.Members[m.Key]
		if !ok {
			m.Object.Members[m.Key] = inspectjson.ObjectMember{
				Name: inspectjson.StringValue{
					Value: m.Key,
				},
				Value: m.Value,
			}
		} else {

			// [spec // 3.2] Otherwise:

			// [spec // 3.2.1] If the *value* of the *key* entry in *object* is not an array, set it to a new array containing the original value.

			if _, ok := originalValue.Value.(inspectjson.ArrayValue); !ok {
				m.Object.Members[m.Key] = inspectjson.ObjectMember{
					Name: originalValue.Name,
					Value: inspectjson.ArrayValue{
						Values: []inspectjson.Value{
							originalValue.Value,
						},
					},
				}
			}

			// [spec // 3.2.2] Append *value* to the value of the *key* entry in *object*.

			objectMembersKeyValueArray := m.Object.Members[m.Key].Value.(inspectjson.ArrayValue)
			objectMembersKeyValueArray.Values = append(
				objectMembersKeyValueArray.Values,
				m.Value,
			)

			objectMembersKey := m.Object.Members[m.Key]
			objectMembersKey.Value = objectMembersKeyValueArray

			m.Object.Members[m.Key] = objectMembersKey
		}
	}
}
