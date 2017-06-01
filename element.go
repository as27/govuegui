package govuegui

import "fmt"

// Element represents a simple html element
type Element struct {
	Key       string `json:"-"`
	DataKey   string `json:"id"`
	Label     string `json:"label"`
	gui       *Gui
	box       *Box
	InputType ElementType `json:"type"`
	Options   []*Option   `json:"options"`
}

// ID returns the id of the element
func (e *Element) ID() string {
	return fmt.Sprintf("%s-%s", e.box.ID(), e.Key)
}

func (e *Element) SetLabel(l string) {
	e.Label = l
}

// Option sets the given values as option
func (e *Element) Option(opt string, values ...string) {
	addOption(e, opt, values...)
}

func (e *Element) getOption(opt string) *Option {
	return getOption(opt, e.Options)
}

func (e *Element) appendOption(o *Option) {
	e.Options = append(e.Options, o)
}

func (e *Element) Set(i interface{}) error {
	return e.gui.Data.Set(e.ID(), i)
}

func (e *Element) Get() interface{} {
	return e.gui.Data.Get(e.ID())
}

func (e *Element) Update() *Element {
	e.gui.Update(e.ID())
	return e
}

func (e *Element) Action(f func()) *Element {
	e.gui.Actions[e.ID()] = f
	return e
}
