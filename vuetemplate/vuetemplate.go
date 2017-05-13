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
	"strings"
	"text/template"
)

var helperFunc = template.FuncMap{
	"function":   func(s string) string { return fmt.Sprintf("function(){\nreturn %s\n}", s) },
	"backquotes": func(s string) string { return fmt.Sprintf("`%s`", s) },
}

// Vue defines all possible JS objects used with vue.js. There are not only
// the core elements availiable. Route and component properties are added here,
// too.
type Vue struct {
	Template string // vue template also used inside components
	Data     string
	Props    string // for handling values inside components
	Children string // used inside components
	Computed string
	Methods  string
	Watch    string
	Path     string // just used inside routes

}

func (v *Vue) WriteTo(w io.Writer) (int64, error) {
	b := &bytes.Buffer{}
	t := template.Must(template.New("vue").Funcs(helperFunc).Parse(vueTemplate))
	b.Write([]byte("{"))
	t.Execute(b, v)
	s := strings.TrimRight(b.String(), "\t\n ,") + "}"
	n, err := w.Write([]byte(s))
	return int64(n), err
}

const vueTemplate = `{{with .Template}}template: {{backquotes .}},{{end}}
	 {{with .Data}}data: {{function .}},{{end}}
	 {{with .Props}}props: {{.}},{{end}}
	 {{with .Children}}children: {{.}},{{end}}
	 {{with .Computed}}computed: {{.}},{{end}}
	 {{with .Methods}}methods: {{.}},{{end}}
	 {{with .Watch}}watch: {{.}},{{end}}
	 {{with .Path}}path: {{.}},{{end}}`

// Component is used for vuejs components
type Component struct {
	Vue
	Name string
}

// NewComponent creates a component and returns the pointer
func NewComponent(name string) *Component {
	return &Component{
		Name: name,
	}
}

// WriteTo implements the WriterTo interface. It takes a io.Writer and
// writes the js block into the writer.
func (c *Component) WriteTo(w io.Writer) (int64, error) {
	jse := NewJSElement(VUECOMPONENT, c.Name, "")
	// JSElement implements the io.Writer
	n, err := c.Vue.WriteTo(jse)
	if err != nil {
		return int64(n), err
	}
	n, err = jse.WriteTo(w)
	return int64(n), err
}

type Router struct {
	Vue
	Name string
}

func NewRouter(name string) *Router {
	return &Router{
		Name: name,
	}
}

func (r *Router) WriteTo(w io.Writer) (int64, error) {
	jse := NewJSElement(VUEROUTER, r.Name, "")
	n, err := r.Vue.WriteTo(jse)
	if err != nil {
		return int64(n), err
	}
	n, err = jse.WriteTo(w)
	return int64(n), err
}
