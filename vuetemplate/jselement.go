package vuetemplate

import (
	"bytes"
	"fmt"
	"io"
)

type JSType int

const (
	CONSTANT JSType = iota
	VARIABLE
	LETSTMT
	FUNCTION
	WEBSOCKET
	VUECOMPONENT
	VUEAPP
	VUEROUTER
)

// JSElement represents the different variable declarations
// of JS.
type JSElement struct {
	JSType  JSType
	VarName string
	Value   string
}

func NewJSElement(t JSType, name, value string) *JSElement {
	return &JSElement{
		JSType:  t,
		VarName: name,
		Value:   value,
	}
}

// String creates a JS line for the element
func (jse *JSElement) String() string {
	var def = ""
	switch jse.JSType {
	default:
		def = "const"
	case CONSTANT:
		def = "const"
	case VARIABLE:
		def = "var"
	case LETSTMT:
		def = "let"
	case FUNCTION:
		return fmt.Sprintf("const %s = function() {\n%s;\n};",
			jse.VarName,
			jse.Value,
		)
	case WEBSOCKET:
		return fmt.Sprintf("var %s = new WebSocket(\"%s\");",
			jse.VarName,
			jse.Value,
		)
	case VUECOMPONENT:
		return fmt.Sprintf("const %s = Vue.component('%s', %s);",
			jse.VarName,
			jse.VarName,
			jse.Value,
		)
	case VUEAPP:
		return fmt.Sprintf("const %s = Vue(%s);",
			jse.VarName,
			jse.Value,
		)
	case VUEROUTER:
		return fmt.Sprintf("const %s = new VueRouter(%s);",
			jse.VarName,
			jse.Value,
		)
	}
	return fmt.Sprintf("%s %s = \"%s\";",
		def,
		jse.VarName,
		jse.Value,
	)
}

// WriteTo implements the io.WriterTo interface by wrapping the String()
// function. WriteTo makes it easier to serve the data inside of a http
// handler.
func (jse *JSElement) WriteTo(w io.Writer) (int64, error) {
	b := bytes.NewBufferString(jse.String())
	n, err := w.Write(b.Bytes())
	return int64(n), err
}

// Write implements the io.Writer. The write method writes everything into
// the jse.Value
func (jse *JSElement) Write(p []byte) (n int, err error) {
	b := bytes.NewBufferString(jse.Value)
	n, err = b.Write(p)
	jse.Value = b.String()
	return n, err
}
func (jse JSElement) Create(w io.Writer, wt io.WriterTo) {
}
