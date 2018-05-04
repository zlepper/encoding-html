package html

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"io"
	"reflect"
	"strconv"
)

// Parses the html into nodes
func parseHtml(reader io.Reader) (*html.Node, error) {
	doc, err := html.Parse(reader)
	return doc, err
}

// A decoder for parsing the html into structs
type Decoder struct {
	reader io.Reader
}

// Creates a new decoder
func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{reader}
}

// Decodes the html into the given struct
// uses reflection
// Also doesn't stream content, so watch your memory.
func (d *Decoder) Decode(v interface{}) error {
	root, err := parseHtml(d.reader)
	if err != nil {
		return err
	}

	mainStruct := reflect.ValueOf(v)

	if mainStruct.Kind() != reflect.Ptr {
		return errors.New("v should be a pointer to a struct")
	}

	err = processPtr(mainStruct, []*html.Node{root}, "")

	return err
}

// Unmarshals the html in b into the given pointer in v
func Unmarshal(b []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(b)).Decode(v)
}

func processField(field reflectableField, node *html.Node) error {
	// Get the tag
	cssSelector, ok := field.field.Tag.Lookup(CSS)
	// Only parse if known tag
	if ok {
		nodes, err := getNodes(node, cssSelector)
		if err != nil {
			return err
		}

		switch {
		// For slice we should operate on all
		case field.isSlice():
			return processSlice(field.value, nodes, field.field.Tag)

			// For structs we should work recursive
		case field.isStruct():
			// Just take the first match
			if len(nodes) > 0 {
				return processStruct(field.value, nodes[0])
			} else {
				return processStruct(field.value, &html.Node{})
			}

			// For anything else, we should just set the value
		default:
			if len(nodes) > 0 {
				return setValue(field.value, nodes[0], field.field.Tag)
			} else {
				return setValue(field.value, &html.Node{}, field.field.Tag)
			}
		}
	}
	return nil
}

// At this point the field should be a slice
func processSlice(v reflect.Value, nodes []*html.Node, tag reflect.StructTag) error {
	if v.Kind() != reflect.Slice {
		return errors.New("field is not a slice")
	}

	// Make sure there is enough capacity to actually load all the nodes
	if v.Cap() < len(nodes) {
		newSlice := reflect.MakeSlice(v.Type(), len(nodes), len(nodes))
		reflect.Copy(newSlice, v)
		v.Set(newSlice)
	}
	if v.Len() < len(nodes) {
		v.SetLen(len(nodes))
	}

	var processSlicePart func(v reflect.Value, node *html.Node) error
	kind := reflect.TypeOf(v.Interface()).Elem().Kind()
	switch kind {
	case reflect.Struct:
		processSlicePart = processStruct
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		processSlicePart = func(v reflect.Value, node *html.Node) error {
			return setValue(v, node, tag)
		}
	case reflect.Ptr:
		processSlicePart = func(v reflect.Value, node *html.Node) error {
			return processPtr(v, []*html.Node{node}, "")
		}
	default:
		return errors.New(fmt.Sprintf("unknown field type: '%s'", kind.String()))
	}
	for index, node := range nodes {
		v := v.Index(index)
		err := processSlicePart(v, node)
		if err != nil {
			return err
		}
	}

	return nil
}

func processStruct(s reflect.Value, node *html.Node) error {
	fields := getFields(s)
	for _, field := range fields {
		err := processField(field, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func processPtr(ptr reflect.Value, nodes []*html.Node, tag reflect.StructTag) error {
	v := ptr.Elem()
	switch v.Kind() {
	case reflect.Struct:
		if len(nodes) > 0 {
			return processStruct(v, nodes[0])
		}
	case reflect.Slice:
		return processSlice(v, nodes, tag)
	default:
		return errors.New(fmt.Sprintf("unknown underlying pointer type: '%s'", v.Kind().String()))
	}
	return nil
}

func setValue(v reflect.Value, node *html.Node, tag reflect.StructTag) error {
	how := tag.Get(EXTRACT)
	if how == "" {
		how = TEXT
	}
	var text string

	switch how {
	case TEXT:
		text = getNodeText(node)
	case ATTRIBUTE:
		key := tag.Get(ATTRIBUTE)
		if key == "" {
			return errors.New("no attribute specified for attribute extraction")
		}
		text = getNodeAttribute(node, key)
	default:
		return errors.New(fmt.Sprintf("unknown how format: '%s'", how))
	}

	// In this case we should return to a default value and use that, if specified
	if text == "" {
		defaultValue, ok := tag.Lookup(DEFAULT)
		if ok {
			text = defaultValue
		}
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(text)
	case reflect.Bool:
		return trySetBool(v, text, tag)
	case reflect.Float32, reflect.Float64:
		return trySetFloat(v, text, tag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return trySetInt(v, text, tag)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return trySetUInt(v, text, tag)
	default:
		return errors.New(fmt.Sprintf("unknown value kind: '%s'", v.Kind().String()))
	}
	return nil
}

// Tries to set a bool
func trySetBool(v reflect.Value, text string, tag reflect.StructTag) error {
	b, err := strconv.ParseBool(text)
	if err != nil {
		defaultValue, ok := tag.Lookup(DEFAULT)
		if ok {
			b, err = strconv.ParseBool(defaultValue)
		}
		if err != nil {
			return err
		}
	}
	v.SetBool(b)
	return nil
}

func trySetFloat(v reflect.Value, text string, tag reflect.StructTag) error {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		defaultValue, ok := tag.Lookup(DEFAULT)
		if ok {
			f, err = strconv.ParseFloat(defaultValue, 64)
		}
		if err != nil {
			return err
		}
	}
	v.SetFloat(f)
	return nil
}
func trySetInt(v reflect.Value, text string, tag reflect.StructTag) error {
	i, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		defaultValue, ok := tag.Lookup(DEFAULT)
		if ok {
			i, err = strconv.ParseInt(defaultValue, 10, 64)
		}
		if err != nil {
			return err
		}
	}
	v.SetInt(i)
	return nil
}
func trySetUInt(v reflect.Value, text string, tag reflect.StructTag) error {
	i, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		defaultValue, ok := tag.Lookup(DEFAULT)
		if ok {
			i, err = strconv.ParseUint(defaultValue, 10, 64)
		}
		if err != nil {
			return err
		}
	}
	v.SetUint(i)
	return nil
}

// Grabs all the nodes matching the given selector
func getNodes(root *html.Node, cssSelector string) ([]*html.Node, error) {
	selector, err := cascadia.Compile(cssSelector)
	if err != nil {
		return nil, err
	}

	matches := selector.MatchAll(root)
	return matches, nil
}

// Gets all the text in the node
func getNodeText(node *html.Node) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			// Keep newlines and spaces, like jQuery
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(node)

	return buf.String()
}

// Gets the specified attribute from the node
// if the attribute cannot be found, returns an empty string
func getNodeAttribute(node *html.Node, attributeKey string) string {
	for _, attribute := range node.Attr {
		if attribute.Key == attributeKey {
			return attribute.Val
		}
	}
	return ""
}
