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
func ParseHtml(reader io.Reader) (*html.Node, error) {
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
	root, err := ParseHtml(d.reader)
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

func processField(field reflectableField, node *html.Node) error {
	// Get the tag
	cssSelector, ok := field.field.Tag.Lookup(CSS)
	// Only parse if known tag
	if ok {
		how := field.getTagOrDefault(EXTRACT, TEXT)
		nodes, err := getNodes(node, cssSelector)
		if err != nil {
			return err
		}

		switch {
		// For slice we should operate on all
		case field.isSlice():
			return processSlice(field.value, nodes, how)

			// For structs we should work recursive
		case field.isStruct():
			// Just take the first match
			if len(nodes) > 0 {
				return processStruct(field.value, nodes[0])
			}

			// For anything else, we should just set the value
		default:
			if len(nodes) > 0 {
				return setValue(field.value, nodes[0], how)
			}
		}
	}
	return nil
}

// At this point the field should be a slice
func processSlice(v reflect.Value, nodes []*html.Node, how string) error {
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
			return setValue(v, node, how)
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
		fmt.Printf("Processing field %s for struct", field.field.Name)
		err := processField(field, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func processPtr(ptr reflect.Value, nodes []*html.Node, how string) error {
	v := ptr.Elem()
	println(ptr.CanInterface())
	switch v.Kind() {
	case reflect.Struct:
		if len(nodes) > 0 {
			return processStruct(v, nodes[0])
		}
	case reflect.Slice:
		return processSlice(v, nodes, how)
	default:
		return errors.New(fmt.Sprintf("unknown underlying pointer type: '%s'", v.Kind().String()))
	}
	return nil
}

func setValue(v reflect.Value, node *html.Node, how string) error {
	var text string
	switch how {
	case TEXT:
		text = getNodeText(node)
	default:
		return errors.New(fmt.Sprintf("unknown how format: '%s'", how))
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(text)
	case reflect.Bool:
		b, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	default:
		return errors.New(fmt.Sprintf("unknown value kind: '%s'", v.Kind().String()))
	}
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
