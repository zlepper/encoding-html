package html

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// Simple usecases

func TestLoadIntoSimpleStruct(t *testing.T) {
	type simpleStruct struct {
		Foo string `css:".foo"`
		Bar string `css:".bar"`
	}

	html := `<body><p class="foo">This is the value of foo</p><p class="bar">This is the value of bar</p></body>`

	var s simpleStruct

	err := NewDecoder(strings.NewReader(html)).Decode(&s)
	assert.NoError(t, err)

	assert.Equal(t, "This is the value of foo", s.Foo)
	assert.Equal(t, "This is the value of bar", s.Bar)
}

func TestLoadIntoStructWithSubStruct(t *testing.T) {
	type complexStruct struct {
		Foo string `css:".foo"`
		Bar struct {
			Baz string `css:".baz"`
		} `css:".bar"`
	}

	//language=html
	html := `<body>
<p class="foo">This is foo</p>
<div class="bar">
    <p class="baz">This is baz</p>
</div>
</body>`

	var c complexStruct
	err := NewDecoder(strings.NewReader(html)).Decode(&c)

	assert.NoError(t, err)
	assert.Equal(t, "This is foo", c.Foo)
	assert.Equal(t, "This is baz", c.Bar.Baz)
}

func TestLoadIntoStructWithSimpleSubSlice(t *testing.T) {
	type sliceStruct struct {
		Foo string   `css:".foo"`
		Bar []string `css:".bar"`
	}

	//language=html
	html := `<body>
<p class="foo">This is foo</p>
<div>
    <p class="bar">Bar 1</p>
    <p class="bar">Bar 2</p>
    <p class="bar">Bar 3</p>
</div>
</body>`

	var s sliceStruct
	err := NewDecoder(strings.NewReader(html)).Decode(&s)

	assert.NoError(t, err)
	assert.Equal(t, "This is foo", s.Foo)
	assert.Equal(t, "Bar 1", s.Bar[0])
	assert.Equal(t, "Bar 2", s.Bar[1])
	assert.Equal(t, "Bar 3", s.Bar[2])

}

func TestLoadIntoStructWithStructSlice(t *testing.T) {
	type sliceStruct struct {
		Foo string `css:".foo"`
		Bar []struct {
			Baz string `css:".baz"`
		} `css:".bar"`
	}

	//language=html
	html := `<body>
<p class="foo">This is a foo</p>
<div class="bar">
    <p class="baz">Bazinga</p>
</div>
<div class="bar">
    <p class="baz">Bazooka</p>
</div>
</body>`

	var c sliceStruct
	err := NewDecoder(strings.NewReader(html)).Decode(&c)

	assert.NoError(t, err)
	assert.Equal(t, "This is a foo", c.Foo)
	assert.Equal(t, "Bazinga", c.Bar[0].Baz)
	assert.Equal(t, "Bazooka", c.Bar[1].Baz)
}

func TestLoadFromAttribute(t *testing.T) {
	type attrStruct struct {
		Foo  string `css:".foo"`
		Link string `css:".foo" extract:"attr" attr:"href"`
	}

	//language=html
	html := `<body><a class="foo" href="http://github.com/zlepper/encoding-html">Link to encoding-html on GitHub</a></body>`

	var a attrStruct
	err := Unmarshal([]byte(html), &a)

	assert.NoError(t, err)
	assert.Equal(t, "Link to encoding-html on GitHub", a.Foo)
	assert.Equal(t, "http://github.com/zlepper/encoding-html", a.Link)
}

// Odd edgecases that I want to make sure are handled
func TestLoadIntoStructWithSlicePointers(t *testing.T) {
	return

	type sliceStruct struct {
		Foo string `css:".foo"`
		Bar []*struct {
			Baz string `css:".baz"`
		} `css:".bar"`
	}

	//language=html
	html := `<body>
<p class="foo">This is a foo</p>
<div class="bar">
    <p class="baz">Bazinga</p>
</div>
<div class="bar">
    <p class="baz">Bazooka</p>
</div>
</body>`

	var c sliceStruct
	err := NewDecoder(strings.NewReader(html)).Decode(&c)

	assert.NoError(t, err)
	assert.Equal(t, "This is a foo", c.Foo)
	assert.Equal(t, "Bazinga", c.Bar[0].Baz)
	assert.Equal(t, "Bazooka", c.Bar[1].Baz)
}

func TestCanDecodeAllStandardTypes(t *testing.T) {
	type allTypesStruct struct {
		Bool bool `css:".bool"`

		Byte byte `css:".byte"`

		Float32 float32 `css:".float32"`
		Float64 float64 `css:".float64"`

		Int   int   `css:".int"`
		Int8  int8  `css:".int8"`
		Int16 int16 `css:".int16"`
		Int32 int32 `css:".int32"`
		Int64 int64 `css:".int64"`

		String string `css:".string"`

		UInt   uint   `css:".uint"`
		UInt8  uint8  `css:".uint8"`
		UInt16 uint16 `css:".uint16"`
		UInt32 uint32 `css:".uint32"`
		UInt64 uint64 `css:".uint64"`
	}

	//language=html
	html := `<body>
<p class="bool">true</p>
<p class="byte">1</p>
<p class="float32">5.5</p>
<p class="float64">4.566666</p>
<p class="int">4</p>
<p class="int8">127</p>
<p class="int16">32767</p>
<p class="int32">2147483647</p>
<p class="int64">9223372036854775807</p>
<p class="string">Some aweasome string</p>
<p class="uint">4294967295</p>
<p class="uint8">255</p>
<p class="uint16">65535</p>
<p class="uint32">4294967295</p>
<p class="uint64">18446744073709551615</p>
</body>`

	var s allTypesStruct
	err := NewDecoder(strings.NewReader(html)).Decode(&s)

	assert.NoError(t, err)
	assert.Equal(t, true, s.Bool)
	assert.Equal(t, byte(1), s.Byte)
	assert.Equal(t, float32(5.5), s.Float32)
	assert.Equal(t, float64(4.566666), s.Float64)
	assert.Equal(t, int(4), s.Int)
	assert.Equal(t, int8(127), s.Int8)
	assert.Equal(t, int16(32767), s.Int16)
	assert.Equal(t, int32(2147483647), s.Int32)
	assert.Equal(t, int64(9223372036854775807), s.Int64)
	assert.Equal(t, "Some aweasome string", s.String)
	assert.Equal(t, uint(4294967295), s.UInt)
	assert.Equal(t, uint8(255), s.UInt8)
	assert.Equal(t, uint16(65535), s.UInt16)
	assert.Equal(t, uint32(4294967295), s.UInt32)
	assert.Equal(t, uint64(18446744073709551615), s.UInt64)
}

func TestDefaultValues(t *testing.T) {
	type allTypesStruct struct {
		Bool bool `css:".bool" default:"true"`

		Byte byte `css:".byte" default:"1"`

		Float32 float32 `css:".float32" default:"5.5"`
		Float64 float64 `css:".float64" default:"4.566666"`

		Int   int   `css:".int" default:"4"`
		Int8  int8  `css:".int8" default:"127"`
		Int16 int16 `css:".int16" default:"32767"`
		Int32 int32 `css:".int32" default:"2147483647"`
		Int64 int64 `css:".int64" default:"9223372036854775807"`

		String string `css:".string" default:"Some aweasome string"`

		UInt   uint   `css:".uint" default:"4294967295"`
		UInt8  uint8  `css:".uint8" default:"255"`
		UInt16 uint16 `css:".uint16" default:"65535"`
		UInt32 uint32 `css:".uint32" default:"4294967295"`
		UInt64 uint64 `css:".uint64" default:"18446744073709551615"`
	}

	//language=html
	html := `<body></body>`

	var s allTypesStruct
	err := NewDecoder(strings.NewReader(html)).Decode(&s)

	assert.NoError(t, err)
	assert.Equal(t, true, s.Bool)
	assert.Equal(t, byte(1), s.Byte)
	assert.Equal(t, float32(5.5), s.Float32)
	assert.Equal(t, float64(4.566666), s.Float64)
	assert.Equal(t, int(4), s.Int)
	assert.Equal(t, int8(127), s.Int8)
	assert.Equal(t, int16(32767), s.Int16)
	assert.Equal(t, int32(2147483647), s.Int32)
	assert.Equal(t, int64(9223372036854775807), s.Int64)
	assert.Equal(t, "Some aweasome string", s.String)
	assert.Equal(t, uint(4294967295), s.UInt)
	assert.Equal(t, uint8(255), s.UInt8)
	assert.Equal(t, uint16(65535), s.UInt16)
	assert.Equal(t, uint32(4294967295), s.UInt32)
	assert.Equal(t, uint64(18446744073709551615), s.UInt64)
}

func TestDefaultValuesForChildStructWhereParentIsMissing(t *testing.T) {
	type Case struct {
		Foo string `css:".foo" default:"Hello"`
		Bar struct {
			Baz string `css:".baz" default:"world"`
		} `css:".bar"`
	}

	//language=html
	html := `<body></body>`

	var c Case
	err := Unmarshal([]byte(html), &c)

	assert.NoError(t, err)
	assert.Equal(t, "Hello", c.Foo)
	assert.Equal(t, "world", c.Bar.Baz)
}

func TestDefaultValueWhenParsingFails(t *testing.T) {
	type Case struct {
		Int   int     `css:".int" default:"42"`
		Uint  uint    `css:".uint" default:"42"`
		Float float64 `css:".float" default:"42.5"`
		Bool  bool    `css:".bool" default:"true"`
	}

	//language=html
	html := `<body>
<p class="int">something</p>
<p class="uint">something</p>
<p class="float">something</p>
<p class="bool">something</p>
</body>`

	var c Case

	err := Unmarshal([]byte(html), &c)

	assert.NoError(t, err)
	assert.Equal(t, 42, c.Int)
	assert.Equal(t, uint(42), c.Uint)
	assert.Equal(t, 42.5, c.Float)
	assert.Equal(t, true, c.Bool)
}

func TestInvalidCSSSelector(t *testing.T) {
	type TestStruct struct {
		Foo string `css:"...foo"`
	}

	//language=html
	html := `<body><p class="foo">Foo</p></body>`

	var f TestStruct
	err := Unmarshal([]byte(html), &f)

	assert.Error(t, err)
	assert.Equal(t, "", f.Foo)
}

func TestErrorShouldHappenWhenExtractAttrButNoAttribute(t *testing.T) {
	type TestStruct struct {
		Foo string `css:".foo" extract:"attr"`
	}

	//language=html
	html := `<body><p class="foo" data-stuffs="Bar">Foo</p></body>`

	var f TestStruct
	err := Unmarshal([]byte(html), &f)

	assert.Error(t, err)
	assert.Equal(t, "", f.Foo)
}
