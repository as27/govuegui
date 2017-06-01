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

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/as27/govuegui/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ElementType defines the
type ElementType string

// Defining the allowed ElementTypes
const (
	INPUT    ElementType = "GVGINPUT"
	TEXTAREA             = "GVGTEXTAREA"
	SELECT               = "GVGSELECT"
	TEXT                 = "GVGTEXT"
	TABLE                = "GVGTABLE"
	LIST                 = "GVGLIST"
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

// Gui groups different forms together.
type Gui struct {
	Forms      []*Form
	Data       *storage.Data
	UpdateData *storage.Data
	hub        *hub
	Actions    map[string]func() `json:"-"`
}

// NewGui returns a pointer to a new instance of a gui
func NewGui() *Gui {
	return &Gui{
		hub:        newWebsocketHub(),
		Data:       storage.New(),
		UpdateData: storage.New(),
		Actions:    make(map[string]func()),
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
			Key: id,
			gui: g,
		}
		g.Forms = append(g.Forms, form)
	}
	return form
}

// ServeHTTP implements the http handler interface
func (g *Gui) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	//r.HandleFunc(PathPrefix+"/", rootHandler)
	prefix := PathPrefix + "/data"
	router.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		action := q.Get("action")
		if action != "" {
			a, ok := g.Actions[action]
			if ok {
				a()
			}
			return
		}
		if r.Method == "GET" {
			b, err := g.Marshal()
			if err != nil {
				log.Println(err)
			}
			w.Write(b)
		}
		if r.Method == "POST" {
			rbody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			newG := NewGui()
			err = json.Unmarshal(rbody, newG)
			if err != nil {
				log.Println(err)
			}
			err = g.Data.SetData(newG.Data)
			if err != nil {
				log.Println(err)
			}
		}
	})
	router.HandleFunc(prefix+"/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		defer g.hub.removeConnection(conn)
		g.hub.addConnection(conn)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("gui serve ws:", err)
			}
			fmt.Println(message)
		}
	})
	router.ServeHTTP(w, r)

}

func (g *Gui) Marshal() ([]byte, error) {
	return json.MarshalIndent(g, "", "  ")
}

// Update sends the values to the websocket. If a dataKey is specified
// just the data to the key is updated.
func (g *Gui) Update(dataKeys ...string) error {
	// Clear update data at the beginning to ensure that is no data from
	// the last update call.
	g.clearUpdateData()
	for _, key := range dataKeys {
		keyData := g.Data.Get(key)
		if keyData == nil {
			// Search for the element ID inside the data, when the key does
			// not exist the data is nil. Then new keys are searched and
			// updated
			availiableKeys := g.Data.GetKeys()
			var matchingKeys []string
			suffix := fmt.Sprintf("-%s", key)
			for _, ak := range availiableKeys {
				if strings.HasSuffix(ak, suffix) {
					matchingKeys = append(matchingKeys, ak)
				}
			}
			if len(matchingKeys) == 0 {
				return fmt.Errorf("key '%s' not found", key)
			}
			return g.Update(matchingKeys...)
		}
		err := g.UpdateData.Set(key, g.Data.Get(key))
		if err != nil {
			return err
		}
	}
	err := g.hub.writeJSON(g)
	return err
}

func (g *Gui) clearUpdateData() {
	g.UpdateData = storage.New()
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

/*type Button struct {
}

func (b *Button) Action(f func()) {

}*/
