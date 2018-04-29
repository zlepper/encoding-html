package html

import (
	"reflect"
)

type reflectableField struct {
	field reflect.StructField
	value reflect.Value
}

// Checks if the underlying field is a slice
func (f reflectableField) isSlice() bool {
	return f.field.Type.Kind() == reflect.Slice
}

// Checks if the underlying field is a struct
func (f reflectableField) isStruct() bool {
	return f.field.Type.Kind() == reflect.Struct
}

// Gets the underlying kind of the slice, e.g. for []int -> int
func (f reflectableField) getUnderlyingKind() reflect.Kind {
	return f.value.Type().Kind()
}

// Gets the value of the tag, or a default value if the tag
// is not assigned, or is empty
func (f reflectableField) getTagOrDefault(key, d string) string {
	value, ok := f.field.Tag.Lookup(key)
	if !ok || value == "" {
		return d
	}
	return value
}

// Gets all the fields that can be changed
func getFields(v reflect.Value) []reflectableField {
	t := v.Type()
	fields := make([]reflectableField, 0)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			fields = append(fields, reflectableField{t.Field(i), field})
		}
	}
	return fields
}
