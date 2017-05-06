package vuetemplate

import "fmt"

type JSType int

const (
	CONSTANT JSType = iota
	VARIABLE
)

// JSElement represents the different variable declarations
// of JS.
type JSElement struct {
	JSType  JSType
	VarName string
	Value   string
}

// String creates a JS line for the element
func (jse JSElement) String() string {
	var def = ""
	switch jse.JSType {
	case CONSTANT:
		def = "const"
	case VARIABLE:
		def = "var"
	}
	return fmt.Sprintf("%s %s = \"%s\";",
		def,
		jse.VarName,
		jse.Value,
	)
}

// WriteTo implements the io.WriterTo interface
func (jse JSElement) WriteTo(w io.Writer) (int64, error){
	b := bytes.NewBufferString(jse.String())
	return w.Write(b.Bytes())

type Vue struct {
	El       string
	Data     string
	Computed string
	Methods  string
}

type Component struct {
	Vue
	Name string
}

func v() string {
	fmt.Println()
	return "a"
}

var a = v()
