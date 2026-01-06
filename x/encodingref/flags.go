package encodingref

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// UnmarshalFlags unmarshals flag strings into a struct based on field names.
// The d parameter must be a pointer to a struct.
// Flags are in the format "{PATH}" or "{PATH}={VALUE}" where PATH is a dot-separated
// traversal path based on field names with lowercase first letter.
func UnmarshalFlags(d any, flags []string) error {
	rv := reflect.ValueOf(d)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("d must be a pointer, got %T", d)
	}
	if rv.IsNil() {
		return fmt.Errorf("d must not be nil")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("d must point to a struct, got %T", d)
	}

	for _, flag := range flags {
		path, value, hasValue := parseFlag(flag)
		if err := setField(rv, path, value, hasValue); err != nil {
			return fmt.Errorf("flag %q: %w", flag, err)
		}
	}

	return nil
}

func parseFlag(flag string) (path []string, value string, hasValue bool) {
	parts := strings.SplitN(flag, "=", 2)
	path = strings.Split(parts[0], ".")
	if len(parts) == 2 {
		value = parts[1]
		hasValue = true
	}
	return
}

func setField(rv reflect.Value, path []string, value string, hasValue bool) error {
	if len(path) == 0 {
		return fmt.Errorf("empty path")
	}

	// Traverse to the target field
	for i, segment := range path {
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				// Initialize nil pointer
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}

		if rv.Kind() != reflect.Struct {
			return fmt.Errorf("cannot traverse into non-struct type %v at path segment %d", rv.Type(), i)
		}

		field, err := findFieldByName(rv, segment)
		if err != nil {
			return err
		}

		// If this is not the last segment, continue traversing
		if i < len(path)-1 {
			rv = field
			continue
		}

		// Last segment - set the value
		return setFieldValue(field, value, hasValue)
	}

	return nil
}

func findFieldByName(rv reflect.Value, name string) (reflect.Value, error) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		// Convert field name to lowercase first letter
		fieldName := toLowerFirst(field.Name)
		if fieldName == name {
			return rv.Field(i), nil
		}
	}

	return reflect.Value{}, fmt.Errorf("field %q not found", name)
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func setFieldValue(field reflect.Value, value string, hasValue bool) error {
	if !field.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	fieldType := field.Type()
	isPtr := fieldType.Kind() == reflect.Ptr

	if !hasValue {
		// No value provided - initialize empty or set bool to true
		if isPtr {
			elemType := fieldType.Elem()
			if elemType.Kind() == reflect.Bool {
				// Special case: *bool should be set to true
				trueVal := true
				field.Set(reflect.ValueOf(&trueVal))
			} else {
				// Initialize empty value
				field.Set(reflect.New(elemType))
			}
		} else if fieldType.Kind() == reflect.Bool {
			// Special case: bool should be set to true
			field.SetBool(true)
		} else {
			// Initialize with zero value
			field.Set(reflect.Zero(fieldType))
		}
		return nil
	}

	// Value provided - parse according to type
	targetType := fieldType
	if isPtr {
		targetType = fieldType.Elem()
	}

	parsedValue, err := parseValue(value, targetType)
	if err != nil {
		return fmt.Errorf("failed to parse value: %w", err)
	}

	if isPtr {
		// Create pointer to parsed value
		ptrVal := reflect.New(targetType)
		ptrVal.Elem().Set(parsedValue)
		field.Set(ptrVal)
	} else {
		field.Set(parsedValue)
	}

	return nil
}

func parseValue(value string, targetType reflect.Type) (reflect.Value, error) {
	// Check for json.RawMessage first
	if targetType == reflect.TypeOf(json.RawMessage{}) {
		return reflect.ValueOf(json.RawMessage(value)), nil
	}

	switch targetType.Kind() {
	case reflect.String:
		// Handle string and custom string types
		result := reflect.New(targetType).Elem()
		result.SetString(value)
		return result, nil

	case reflect.Slice:
		if targetType.Elem().Kind() == reflect.Uint8 {
			// []byte
			return reflect.ValueOf([]byte(value)), nil
		}

	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetBool(parsed)
		return result, nil

	case reflect.Int:
		parsed, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetInt(parsed)
		return result, nil

	case reflect.Int8:
		parsed, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetInt(parsed)
		return result, nil

	case reflect.Int16:
		parsed, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetInt(parsed)
		return result, nil

	case reflect.Int32:
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetInt(parsed)
		return result, nil

	case reflect.Int64:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetInt(parsed)
		return result, nil

	case reflect.Uint:
		parsed, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetUint(parsed)
		return result, nil

	case reflect.Uint8:
		parsed, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetUint(parsed)
		return result, nil

	case reflect.Uint16:
		parsed, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetUint(parsed)
		return result, nil

	case reflect.Uint32:
		parsed, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetUint(parsed)
		return result, nil

	case reflect.Uint64:
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		result := reflect.New(targetType).Elem()
		result.SetUint(parsed)
		return result, nil
	}

	return reflect.Value{}, fmt.Errorf("unsupported type %v", targetType)
}
