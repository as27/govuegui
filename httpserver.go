package govuegui

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	"github.com/as27/golib/css/purecss"
	"github.com/as27/golib/js/vuejsdev"
	"github.com/as27/golib/js/vueresourcemin"
	"github.com/as27/golib/js/vueroutermin"
	"github.com/gorilla/mux"
)

// PathPrefix defines the prefix for the all gui specific tasks
var PathPrefix = "/govuegui"

var useRice = false

// ServerPort defines the port of the gui server, when using
// `govuegui.Serve()`
var ServerPort = ":2700"

// NewRouter returns a router from the gorillatoolkit
// http://www.gorillatoolkit.org/pkg/mux
// The router already includes all the paths which are needed
// for the gui. It can be called like:
//   r := govuegui.NewRouter()
//   // Add you own routes
//   r.HandleFunc("/products/{key}", ProductHandler)
func NewRouter(g *Gui) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix+"/", rootHandler)
	r.Handle(PathPrefix+"/data", g)
	r.HandleFunc(PathPrefix+"/lib/vue.min.js", vuejsdev.Handler)
	r.HandleFunc(PathPrefix+"/lib/vue-router.min.js", vueroutermin.Handler)
	r.HandleFunc(PathPrefix+"/lib/vue-resource.min.js", vueresourcemin.Handler)
	r.HandleFunc(PathPrefix+"/lib/pure.min.css", purecss.Handler)
	jsPrefix := PathPrefix + "/lib/"
	if useRice {
		box := rice.MustFindBox("lib")
		htmlFiles := http.StripPrefix(jsPrefix, http.FileServer(box.HTTPBox()))
		r.PathPrefix(jsPrefix).Handler(htmlFiles)
	} else {
		r.PathPrefix(jsPrefix).
			Handler(
				http.StripPrefix(jsPrefix,
					http.FileServer(http.Dir("lib"))))
	}

	return r
}

// Serve wraps the http.ListenAndServe() function, but adds the
// routes for the gui.
func Serve(g *Gui) error {
	r := NewRouter(g)
	log.Println("Serving gvg on port: ", ServerPort)
	return http.ListenAndServe(ServerPort, r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var templateString string
	fpath := filepath.Join("html", "index.html")
	_, err := os.Stat(fpath)
	if err == os.ErrNotExist {
		templateBox, err := rice.FindBox("html")
		if err != nil {
			log.Fatal(err)
		}
		templateString, err = templateBox.String("index.html")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tmpB, err := ioutil.ReadFile(fpath)
		if err != nil {
			log.Fatal(err)
		}
		templateString = string(tmpB)
	}

	tmplMessage, err := template.New("message").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]string)
	data["PathPrefix"] = PathPrefix
	tmplMessage.Execute(w, data)
}
