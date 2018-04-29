package html

import (
	"fmt"
	"reflect"
	"testing"
)

// Used for mocking around with the reflect library,
// trying to figure out what is going on

type TestData struct {
	Title       string `css:".title"`
	Description string `css:".description"`
	SubSection  []struct {
		Field1 int `css:".field-1"`
		Field2 int `css:".field-2"`
	} `css:".subsection"`
}

const css = "css"

func TestReflectionExperimentation(t *testing.T) {
	d := TestData{}

	ty := reflect.TypeOf(d)

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		tag, ok := field.Tag.Lookup(css)
		if ok {
			fmt.Printf("Field %s had css tag with value %s and was type %v\n", field.Name, tag, field.Type)
		}
		if field.Type.Kind() == reflect.Slice {
			fmt.Println(field.Type.Elem().Kind() == reflect.Struct)
		}
	}
}
