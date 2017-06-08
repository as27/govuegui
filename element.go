package govuegui

import "fmt"

// Element represents a simple html element. The element is not defined directly. It always
// is initialized via a Box type.
//   // The element of the type INPUT is created via the Input() method of the Box type.
//   inputElement := gui.Form("myForm").Box("myBox").Input("Input")
type Element struct {
	Key       string `json:"-"`
	DataKey   string `json:"id"`
	Label     string `json:"label"`
	Watch     bool   `json:"watch"`
	gui       *Gui
	box       *Box
	InputType ElementType        `json:"type"`
	Options   map[string]*Option `json:"options"`
}

// ID returns the id of the element
func (e *Element) ID() string {
	return fmt.Sprintf("%s-%s", e.box.ID(), e.Key)
}

// SetLabel changes the label of the element inside the gui. The default value of the label
// is always the id.
func (e *Element) SetLabel(l string) {
	e.Label = l
}

// Option sets the given values as option
func (e *Element) Option(opt string, values ...string) *Element {
	e.Options[opt] = &Option{
		Option: opt,
		Values: values,
	}
	return e
}

// Set takes a value for rendering inside a element of the gui.
func (e *Element) Set(i interface{}) error {
	return e.gui.Data.Set(e.ID(), i)
}

// Get returns a value out from the gui.
func (e *Element) Get() interface{} {
	return e.gui.Data.Get(e.ID())
}

// Update is the method to let the gui send the value from the gui server
// to the browser.
func (e *Element) Update() *Element {
	e.gui.Update(e.ID())
	return e
}

// Action takes a callback function. For input fields that function
// is called when the value changes.
func (e *Element) Action(f func(*Gui)) *Element {
	e.gui.Actions[e.ID()] = f
	return e
}
