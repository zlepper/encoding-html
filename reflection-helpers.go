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
