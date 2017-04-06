package govuegui

import "fmt"

// Element represents a simple html element
type Element struct {
	id        string
	gui       *Gui
	inputType ElementType
	options   []*Option
}

func NewElement(id string, gui *Gui, inputType ElementType) *Element {
	return &Element{
		id:        id,
		gui:       gui,
		inputType: inputType,
	}
}

// ID returns the id of the element
func (e *Element) ID() string {
	return e.id
}

func (e *Element) optionsName() string {
	return fmt.Sprintf("element-options-%s", e.ID())
}

// Option sets the given values as option
func (e *Element) Option(opt string, values ...string) {
	o := e.getOption(opt)
	if o != nil {
		o.Values = values
	} else {
		newOption := Option{
			Option: opt,
			Values: values,
		}
		e.gui.Data.Set(e.optionsName(), append(e.options, &newOption))
	}

}

func (e *Element) getOption(opt string) *Option {
	opts, err := e.gui.Data.GetWithErrors(e.optionsName())
	if err != nil {
		return nil
	}
	for _, o := range opts.([]*Option) {
		if o.Option == opt {
			return o
		}
	}
	return nil
}

func (e *Element) Set(i interface{}) error {
	return e.gui.Data.Set(e.id, i)
}

func (e *Element) Get() interface{} {
	return e.gui.Data.Get(e.id)
}
