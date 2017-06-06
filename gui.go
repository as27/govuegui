package govuegui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/as27/govuegui/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Gui groups different forms together.
type Gui struct {
	Title      string
	Forms      []*Form
	Data       *storage.Data
	UpdateData *storage.Data
	hub        *hub
	Actions    map[string]func(*Gui) `json:"-"`
}

// NewGui returns a pointer to a new instance of a gui
func NewGui() *Gui {
	return &Gui{
		Title:      "My govuigui app",
		hub:        newWebsocketHub(),
		Data:       storage.New(),
		UpdateData: storage.New(),
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
	router.HandleFunc(PathPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		var templateString string
		templateString = htmlTemplate
		tmplMessage, err := template.New("message").Parse(templateString)
		if err != nil {
			log.Fatal(err)
		}
		data := make(map[string]string)
		data["PathPrefix"] = PathPrefix
		data["Title"] = g.Title
		tmplMessage.Execute(w, data)
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
			newG := NewGui()
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
