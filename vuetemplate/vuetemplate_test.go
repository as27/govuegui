package vuetemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompontent(t *testing.T) {
	c := NewComponent("test")
	c.Template = "<div>abc</div>"
	b := &bytes.Buffer{}
	c.WriteTo(b)
	expect := `const test = Vue.component('test', {
		template: ` + "`" + "<div>abc</div>`" + `});`
	assert.Equal(t, clearString(expect), clearString(b.String()))
}

func TestVue(t *testing.T) {
	v := Vue{}
	v.Data = `{
		val1: "value",
		int1: 1
	}`
	v.Template = `<div>abc</div>`
	expect := `{
		template: ` + "`" + `<div>abc</div>` + "`" + `, 
		data: function(){
			return {
				val1: "value",
				int1: 1
			}
		}
	}`
	b := &bytes.Buffer{}
	v.WriteTo(b)
	assert.Equal(t, clearString(expect), clearString(b.String()))
}
func TestJSElement(t *testing.T) {
	testCases := []struct {
		Got    JSElement
		Expect string
	}{
		{
			Got: JSElement{
				CONSTANT,
				"var1",
				"val1",
			},
			Expect: `const var1 = "val1";`,
		},
		{
			Got: JSElement{
				VARIABLE,
				"var1",
				"val1",
			},
			Expect: `var var1 = "val1";`,
		},
		{
			Got: JSElement{
				LETSTMT,
				"var1",
				"val1",
			},
			Expect: `let var1 = "val1";`,
		},
		{
			Got: JSElement{
				FUNCTION,
				"var1",
				"return \"a\"",
			},
			Expect: `const var1 = function() {
				return "a";
			};`,
		},
	}
	for _, tc := range testCases {
		jse := tc.Got
		assert.Equal(t, clearString(tc.Expect), clearString(jse.String()))
	}
}
func clearString(s string) string {
	r := strings.NewReplacer(
		"\n", "",
		"\t", "",
		"  ", "")
	return r.Replace(s)
}
func TestComponent(t *testing.T) {
}
