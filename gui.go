package govuegui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/as27/govuegui/storage"
	"github.com/as27/govuegui/vuetemplate"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Gui groups different forms together.
type Gui struct {
	PathPrefix string `json:"-"`
	ServerPort string `json:"-"`
	Title      string
	Forms      []*Form
	Data       *storage.Data
	UpdateData *storage.Data
	hub        *hub
	template   GuiTemplate
	Actions    map[string]func(*Gui) `json:"-"`
}

// NewGui returns a pointer to a new instance of a gui
func NewGui(t GuiTemplate) *Gui {
	return &Gui{
		PathPrefix: DefaultPathPrefix,
		ServerPort: DefaultServerPort,
		Title:      "My govuigui app",
		hub:        newWebsocketHub(),
		Data:       storage.New(),
		UpdateData: storage.New(),
		template:   t,
		Actions:    make(map[string]func(*Gui)),
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
			Key:     id,
			gui:     g,
			Options: make(map[string]*Option),
		}
		g.Forms = append(g.Forms, form)
	}
	return form
}

// ServeHTTP implements the http handler interface
func (g *Gui) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	//r.HandleFunc(PathPrefix+"/", rootHandler)
	prefix := g.PathPrefix + "/data"
	router.HandleFunc(g.PathPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		var templateString string
		templateString = htmlTemplate
		tmplMessage, err := template.New("message").Parse(templateString)
		if err != nil {
			log.Fatal(err)
		}
		// Define the variables for the template
		data := make(map[string]string)
		data["PathPrefix"] = g.PathPrefix
		data["Title"] = g.Title
		data["Body"] = g.template.Body
		tmplMessage.Execute(w, data)
	})
	router.HandleFunc(g.PathPrefix+"/app.css", g.template.CSSHandler)
	router.HandleFunc(g.PathPrefix+"/custom.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		b := bytes.NewBufferString(g.template.CustomCSS)
		w.Write(b.Bytes())
	})
	router.HandleFunc(g.PathPrefix+"/app.js", func(w http.ResponseWriter, r *http.Request) {
		serverVar := "localhost" + g.ServerPort
		vuetemplate.NewJSElement(vuetemplate.CONSTANT, "PathPrefix", g.PathPrefix).WriteTo(w)
		vuetemplate.NewJSElement(vuetemplate.CONSTANT, "Server", serverVar).WriteTo(w)
		vuetemplate.NewJSElement(vuetemplate.CONSTANT, "appTitle", g.Title).WriteTo(w)
		for _, vueComp := range VueComponents {
			var comp *vuetemplate.Component
			if vueComp.CompFunc != nil {
				comp = vueComp.CompFunc(g.template)
			} else {
				tString, err := getTemplateFromElementType(vueComp.ElementType, g.template)
				if err != nil {
					panic(err)
				}
				comp = vuegvgdefaultelement(vueComp.Name, tString)
			}
			comp.WriteTo(w)
		}
		/*vuegvgdefaultelement("gvginput", g.template.GvgInput).WriteTo(w)
		vuegvgdefaultelement("gvgtextarea", g.template.GvgTextarea).WriteTo(w)
		vuegvgdefaultelement("gvgtext", g.template.GvgText).WriteTo(w)
		vuegvgdefaultelement("gvgtable", g.template.GvgTable).WriteTo(w)
		vuegvgdefaultelement("gvgdropdown", g.template.GvgDropdown).WriteTo(w)
		vuegvgdefaultelement("gvglist", g.template.GvgList).WriteTo(w)
		vuegvgbutton(g.template).WriteTo(w)*/
		vuegvgelement(g.template).WriteTo(w)
		vuegvgbox(g.template).WriteTo(w)
		vuegvgform(g.template).WriteTo(w)
		vuegvgforms(g.template).WriteTo(w)
		b := bytes.NewBufferString(vueappstring)
		w.Write(b.Bytes())
	})
	router.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		action := q.Get("action")
		if action != "" {
			a, ok := g.Actions[action]
			if ok {
				a(g)
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
			newG := NewGui(g.template)
			err = json.Unmarshal(rbody, newG)
			if err != nil {
				log.Println(err)
			}
			err = g.Data.SetData(newG.Data)
			if err != nil {
				log.Println(err)
			}
			g.Update()
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

// Marshal wraps the json marshal method.
func (g *Gui) Marshal() ([]byte, error) {
	return json.MarshalIndent(g, "", "  ")
}

// Update sends the values of the gui to the websocket. The dataKeys are
// defining, which field are going to be updated. It is also possible to
// call that function without arguments, then everything will be updated.
// The update of all fields has to be used carefully, because that update
// overwrites also user inputs, which are not submited.
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
