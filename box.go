package govuegui

import "fmt"

// Box is the way elements are grouped. Every Element
type Box struct {
	Key      string             `json:"id"`
	Options  map[string]*Option `json:"options"`
	gui      *Gui
	form     *Form
	Elements []*Element `json:"elements"`
}

// ID returns the id of the box
func (b *Box) ID() string {
	return fmt.Sprintf("%s-%s", b.form.ID(), b.Key)
}

func (b *Box) Option(opt string, values ...string) *Box {
	b.Options[opt] = &Option{
		Option: opt,
		Values: values,
	}
	return b
}

func (b *Box) Clear() {
	for _, el := range b.Elements {
		// Remove values from storage
		b.gui.Data.Remove(el.ID())
	}
	// Set Elements to empty struct
	b.Elements = []*Element{}
}

func (b *Box) Element(id string, inputType ElementType) *Element {
	var el *Element
	for _, e := range b.Elements {
		if e.Key == id {
			el = e
			break
		}
	}
	if el == nil {
		el = &Element{
			Key:       id,
			Label:     id,
			gui:       b.gui,
			box:       b,
			InputType: inputType,
			Watch:     false,
			Options:   make(map[string]*Option),
		}
		el.DataKey = el.ID()
		b.Elements = append(b.Elements, el)
	}
	return el
}

func (b *Box) Input(id string) *Element {
	e := b.Element(id, INPUT)
	e.Watch = true
	return e
}

func (b *Box) Table(id string) *Element {
	return b.Element(id, TABLE)
}

func (b *Box) Textarea(id string) *Element {
	e := b.Element(id, TEXTAREA)
	e.Watch = true
	return e
}

func (b *Box) Text(id string) *Element {
	return b.Element(id, TEXT)
}

func (b *Box) List(id string) *Element {
	return b.Element(id, LIST)
}

func (b *Box) Button(id string) *Element {
	return b.Element(id, BUTTON)
}

func (b *Box) Dropdown(id string) *Element {
	e := b.Element(id, DROPDOWN)
	e.Watch = true
	return e
}
