package govuegui

import "fmt"

// Box is the way elements are grouped. Every Element
type Box struct {
	Key      string    `json:"id"`
	Options  []*Option `json:"options"`
	gui      *Gui
	form     *Form
	Elements []*Element `json:"elements"`
}

// ID returns the id of the box
func (b *Box) ID() string {
	return fmt.Sprintf("%s-%s", b.form.ID(), b.Key)
}

func (b *Box) Option(opt string, values ...string) {
	addOption(b, opt, values...)
}

func (b *Box) Clear() {
	for _, el := range b.Elements {
		// Remove values from storage
		b.gui.Data.Remove(el.ID())
	}
	// Set Elements to empty struct
	b.Elements = []*Element{}
}

func (b *Box) getOption(opt string) *Option {
	return getOption(opt, b.Options)
}

func (b *Box) appendOption(o *Option) {
	b.Options = append(b.Options, o)
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
		}
		el.DataKey = el.ID()
		b.Elements = append(b.Elements, el)
	}
	return el
}

func (b *Box) Input(id string) *Element {
	return b.Element(id, INPUT)
}

func (b *Box) Table(id string) *Element {
	return b.Element(id, TABLE)
}

func (b *Box) Textarea(id string) *Element {
	return b.Element(id, TEXTAREA)
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
	return b.Element(id, DROPDOWN)
}
