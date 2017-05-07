package vuetemplate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVue(t *testing.T) {
	// v := Vue{}
	// assert.Equal(t, "ab", v.String())
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
			}`,
		},
	}
	for _, tc := range testCases {
		jse := tc.Got
		assert.Equal(t, tc.Expect, jse.String())
	}
}
func clearString(s string) string {
	return strings.Trim(s, "\n\t")
}
func TestComponent(t *testing.T) {
}
