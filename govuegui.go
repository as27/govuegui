// Package govuegui provides a simple gui, which can be used via a
// http server inside the browser. There are three different elements
// for building the gui. Every level gets a identifier as a string.
//
// First level is the Form. Every Form has one submit button.
//
// Inside a form every element is grouped into a Box. Each Form can
// hold as many Boxes as wanted.
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

type optioner interface {
	Option(string, ...string)
	getOption(string) *Option
	appendOption(*Option)
}

func addOption(o optioner, opt string, values ...string) {
	op := o.getOption(opt)
	if op != nil {
		op.Values = values
	} else {
		newOption := Option{
			Option: opt,
			Values: values,
		}
		o.appendOption(&newOption)
	}
}

func getOption(opt string, opts []*Option) *Option {
	for _, o := range opts {
		if o.Option == opt {
			return o
		}
	}
	return nil
}

// Form wrapps one ore more Boxes
type Form struct {
	Key     string    `json:"id"`
	Options []*Option `json:"options"`
	gui     *Gui
	Boxes   []*Box
}

// ID returns the id of the form
func (f *Form) ID() string {
	return f.Key
}

func (f *Form) Option(opt string, values ...string) {
	addOption(f, opt, values...)
}

func (f *Form) getOption(opt string) *Option {
	return getOption(opt, f.Options)
}

func (f *Form) appendOption(o *Option) {
	f.Options = append(f.Options, o)
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
			Key:  id,
			form: f,
			gui:  f.gui,
		}
		f.Boxes = append(f.Boxes, box)
	}
	return box
}

/*type Button struct {
}

func (b *Button) Action(f func()) {

}*/
