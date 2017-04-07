package govuegui

import (
	"io/ioutil"
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/as27/golib/css/purecss"
	"github.com/as27/golib/js/vuejsdev"
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
//   r.HandleFunc("/products/{key}", ProductHandler)
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix+"/lib/vue.min.js", vuejsdev.Handler)
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
func Serve() error {
	r := NewRouter()
	log.Println("Serving gvg on port: ", ServerPort)
	return http.ListenAndServe(ServerPort, r)
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("js/app.js")
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	w.Write(b)
}
