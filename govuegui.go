// Package govuegui provides a simple gui, which can be used via a
// http server inside the browser. There are three different abstractions
// to build the gui. Every level gets a identifier as a string.
//
// First level is the Form. Every Form has one submit button.
//
// Inside a form every element is grouped into a Box. Each Form can
// hold as many Boxes as wanted.
//
//
//
// The api let's you define everything on a very simple way:
//   Form("abc").Box("cde").Input("name").Value("myvalue")
//   Form("abc").Box("cde").Input("name").BindString(myString)
//   Form("abc").Box("cde").Textarea("name2").Value("myvalue")
//   Form("abc").Box("cde").Select("name2").Option("val1", "Value1")
//   Form("abc").Box("cde").Select("name2").Option("val2", "Value2")
//   Form("abc").Box("cde").Select("name2").Option("val3", "Value3")
//   Form("abc").Box("cde").Each(func(){})
package govuegui

// ElementType defines the
type ElementType int

type dataStorage interface {
	Get(string) interface{}
	Set(string, interface{}) error
}

// Defining the allowed ElementTypes
const (
	INPUT ElementType = iota
	TEXTAREA
	SELECT
)

// option holds the one option of a element
type Option struct {
	Option string
	Values []string
}

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

func (e *Element) Option(opt string, values ...string) {
	o := e.getOption(opt)
	if o != nil {
		o.Values = values
	} else {
		newOption := Option{
			Option: opt,
			Values: values,
		}
		e.options = append(e.options, &newOption)
	}

}

func (e *Element) getOption(opt string) *Option {
	for _, o := range e.options {
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

// Box is the way elements are grouped. Every Element
type Box struct {
	id       string
	gui      *Gui
	Elements []*Element
}

// ID returns the id of the box
func (b *Box) ID() string {
	return b.id
}

func (b *Box) Input(id string) *Element {
	return NewElement(id, b.gui, INPUT)
}

func (b *Box) Textarea(id string) *Element {
	return NewElement(id, b.gui, TEXTAREA)
}

// Form wrapps one ore more Boxes
type Form struct {
	id    string
	gui   *Gui
	Boxes []*Box
}

// ID returns the id of the form
func (f *Form) ID() string {
	return f.id
}

// Box returns the pointer to the box with the given id. If there
// is no box with that id, a new one is created.
func (f *Form) Box(id string) *Box {
	var box *Box
	for _, b := range f.Boxes {
		if b.ID() == id {
			box = b
			break
		}
	}
	if box == nil {
		box = &Box{
			id:  id,
			gui: f.gui,
		}
		f.Boxes = append(f.Boxes, box)
	}
	return box
}

// Gui groups different forms together.
type Gui struct {
	Forms []*Form
	Data  dataStorage
}

// NewGui returns a pointer to a new instance of a gui
func NewGui(ds dataStorage) *Gui {
	return &Gui{
		Data: ds,
	}
}

// Form returns the pointer to a form. If the id exists the existing
// Form is used.
func (g *Gui) Form(id string) *Form {
	// Find Form
	var form *Form
	for _, f := range g.Forms {
		if f.ID() == id {
			form = f
			break
		}
	}
	if form == nil {
		form = &Form{
			id:  id,
			gui: g,
		}
		g.Forms = append(g.Forms, form)
	}
	return form
}
