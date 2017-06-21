package govuegui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/as27/govuegui/storage"
)

// Gui groups different forms together.
type Gui struct {
	PathPrefix string `json:"-"`
	ServerPort string `json:"-"`
	Title      string
	Forms      []*Form
	Data       *storage.Data
	UpdateData *storage.Data
	Active     string // route to the active box
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
	if err != nil {
		return err
	}
	// After sending the data to the hub the Active route has to be cleared
	// that the route is not set again with the next update
	g.Active = ""
	return nil
}

func (g *Gui) clearUpdateData() {
	g.UpdateData = storage.New()
}
