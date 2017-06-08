// Package govuegui provides a simple gui, which can be used via a
// http server inside the browser. There are three different elements
// for building the gui. Every level gets a identifier as a string.
//
// First level is the Form. Every Form has one submit button. A form
// can have boxes. It depends on the implementation how the boxes are
// rendered.
//
// Inside of a box the fields are defined.
//
// The api let's you define everything on a very simple way:
//   Form("abc").Box("cde").Input("name").Value("myvalue")
//   Form("abc").Box("cde").Input("name").BindString(myString)
//   Form("abc").Box("cde").Textarea("name2").Value("myvalue")
//   Form("abc").Box("cde").Select("name2").Option("val1", "Value1")
//   Form("abc").Box("cde").Select("name2").Option("val2", "Value2")
//   Form("abc").Box("cde").Select("name2").Option("val3", "Value3")
//   Form("abc").Box("cde").Button("Click me").Action(func(gui *govuegui.Gui){
//	    // do something when the Button is clicked
//     })
//
// The gui will be show up inside the browser of the user as a webapp. The
// app uses vuejs with a websocket connection.
package govuegui

// ElementType defines the
type ElementType string

// Defining the allowed ElementTypes. The value of each type is used inside
// the vueapp.
const (
	INPUT    ElementType = "GVGINPUT"
	TEXTAREA             = "GVGTEXTAREA"
	SELECT               = "GVGSELECT"
	TEXT                 = "GVGTEXT"
	TABLE                = "GVGTABLE"
	LIST                 = "GVGLIST"
	DROPDOWN             = "GVGDROPDOWN"
	BUTTON               = "GVGBUTTON"
)

// Option holds the one option of a element
type Option struct {
	Option string
	Values []string
}

// Form is the first level of grouping everything.
type Form struct {
	Key     string             `json:"id"`
	Options map[string]*Option `json:"options"`
	gui     *Gui
	Boxes   []*Box
}

// ID returns the id of the form
func (f *Form) ID() string {
	return f.Key
}

func (f *Form) Option(opt string, values ...string) *Form {
	f.Options[opt] = &Option{
		Option: opt,
		Values: values,
	}
	return f

}

// Box returns the pointer to the box with the given id. If there
// is no box with that id, a new one is created.
func (f *Form) Box(id string) *Box {
	var box *Box
	for _, b := range f.Boxes {
		if b.Key == id {
			box = b
			break
		}
	}
	if box == nil {
		box = &Box{
			Key:     id,
			form:    f,
			gui:     f.gui,
			Options: make(map[string]*Option),
		}
		f.Boxes = append(f.Boxes, box)
	}
	return box
}
