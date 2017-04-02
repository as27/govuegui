package govuegui

// Option for an element is represented by a function, which takes
// strings as input. For example Class() could be a function where
// css classes can added to a element.
// Class("active","box")
type Option func(vals ...string)

// Element represents a simple html element
type Element struct {
	id      string
	Key     string
	Name    string
	Options []Option
}

// Box is the way elements are grouped. Every Element
type Box struct {
	id       string
	Elements []*Element
}

// ID returns the id of the box
func (b *Box) ID() string {
	return b.id
}

// Form wrapps one ore more Boxes
type Form struct {
	id    string
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
		box = &Box{id: id}
		f.Boxes = append(f.Boxes, box)
	}
	return box
}

// Gui groups different forms together.
type Gui struct {
	Forms []*Form
}

// NewGui returns a pointer to a new instance of a gui
func NewGui() *Gui {
	return &Gui{}
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
		form = &Form{id: id}
		g.Forms = append(g.Forms, form)
	}
	return form
}

// Form("abc").Box("cde").Input("name","Name").Value("myvalue")
// Form("abc").Box("cde").Textarea("name2","Name").Value("myvalue")
// Form("abc").Box("cde").Select("name2","Name").Option("myvalue")
