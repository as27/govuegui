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

// Clear removes all elements from the box.
func (b *Box) Clear() {
	for _, el := range b.Elements {
		// Remove values from storage
		b.gui.Data.Remove(el.ID())
	}
	// Set Elements to empty struct
	b.Elements = []*Element{}
}

// Element adds a element of ElementType to the box.
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

// Input adds a simple html input field to the box. The field is watched, so
// every change is submited from the browser and is availiable inside
// your go app.
// If you add a action to this field, the action is called everytime
// the field changes.
func (b *Box) Input(id string) *Element {
	e := b.Element(id, INPUT)
	e.Watch = true
	return e
}

// Table adds a html table inside to the box. That the table is correct
// rendered you need to pass a [][]string into via the set method.
// The first slice will be rendered as table header.
//   table := [][]string{
//		{"name", "age", "country"}, // is rendered as header
//		{"Andreas", "27", "Germany"},
//		{"Bob", "22", "Austria"},
//		}
//   gui.Form("myForm").Box("Box1").Table("myTable").Set(table)
func (b *Box) Table(id string) *Element {
	return b.Element(id, TABLE)
}

// Textarea is analog like the input field, just as a html textarea.
// The field is watched so every change is submited.
func (b *Box) Textarea(id string) *Element {
	e := b.Element(id, TEXTAREA)
	e.Watch = true
	return e
}

// Text enables you to write text on the gui. HTML tags are allowed.
func (b *Box) Text(id string) *Element {
	return b.Element(id, TEXT)
}

// List renders a html list. That the list is correct rendered you need
// to pass a slice of string into it.
//   list := []string{"one", "two", "three"}
//   gui.Form("myForm").Box("Box1").List("myTable").Set(list)
func (b *Box) List(id string) *Element {
	return b.Element(id, LIST)
}

// Button adds a button to the box. Every button can hold a action, which
// is called, when the button is clicked.
func (b *Box) Button(id string) *Element {
	return b.Element(id, BUTTON)
}

// Dropdown adds a dropdown element to the box. All the entries of the dropdown
// list needs to be added via the Option() method.
//   dropdown := gui.Form("myForm").Box("Box1").Dropdown("MyDropdown")
//   dropdown.Option("key1", "Value 1")
//   dropdown.Option("key2", "Value 2")
//   dropdown.Option("key3", "Value 3")
// The first argument specifies the key, which will be returned.
func (b *Box) Dropdown(id string) *Element {
	e := b.Element(id, DROPDOWN)
	e.Watch = true
	return e
}
