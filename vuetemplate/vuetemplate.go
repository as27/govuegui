/*Package vuetemplate allows to serve vue.js apps over a go api.
The abstraction works over different elements:
 * JSType defines the different statements, which are used inside JS
 * JSElement is a full JavaScript statement for example `var v1 = "val";`
 * Vue is the definition of the vue object
 * Component defines a vue component
*/
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
)

// JSElement represents the different variable declarations
// of JS.
type JSElement struct {
	JSType  JSType
	VarName string
	Value   string
}

func NewJSElement(t JSType, name, value string) JSElement {
	return JSElement{
		JSType:  t,
		VarName: name,
		Value:   value,
	}
}

// String creates a JS line for the element
func (jse JSElement) String() string {
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
func (jse JSElement) WriteTo(w io.Writer) (int64, error) {
	b := bytes.NewBufferString(jse.String())
	n, err := w.Write(b.Bytes())
	return int64(n), err
}

type Vue struct {
	JSElement
	Options map[string]string
}

type Component struct {
	Vue
	Name string
}
