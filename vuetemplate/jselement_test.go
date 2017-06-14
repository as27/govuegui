package vuetemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{
			Got: JSElement{
				VUECOMPONENT,
				"var1",
				"{}",
			},
			Expect: `const var1 = Vue.component('var1', {});`,
		},
		{
			Got: JSElement{
				VUEAPP,
				"app",
				"{}",
			},
			Expect: `const app = new Vue({});`,
		},
		{
			Got: JSElement{
				VUEROUTER,
				"router",
				"[{}]",
			},
			Expect: `const router = new VueRouter([{}]);`,
		},
		{
			Got: JSElement{
				WEBSOCKET,
				"socket",
				"ws://server/ws",
			},
			Expect: `var socket = new WebSocket("ws://server/ws");`,
		},
	}
	for _, tc := range testCases {
		jse := tc.Got
		assert.Equal(t, clearString(tc.Expect), clearString(jse.String()))
	}
}
