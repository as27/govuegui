package govuegui

import (
	"fmt"
	"log"
	"net/http"

	"github.com/as27/golib/js/vuejsdev"
	"github.com/as27/golib/js/vueresourcemin"
	"github.com/as27/golib/js/vueroutermin"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
)

// DefaultPathPrefix defines the prefix for the all gui specific tasks
var DefaultPathPrefix = "/govuegui"

// DefaultServerPort defines the port of the gui server, when using
// `govuegui.Serve()`
var DefaultServerPort = ":2700"

// NewRouter returns a router from the gorillatoolkit
// http://www.gorillatoolkit.org/pkg/mux
// The router already includes all the paths which are needed
// for the gui. It can be called like:
//   r := govuegui.NewRouter()
//   // Add you own routes
//   r.HandleFunc("/products/{key}", ProductHandler)
func NewRouter(g *Gui) *mux.Router {
	r := mux.NewRouter()
	r.Handle(g.PathPrefix+"/", g)
	r.Handle(g.PathPrefix+"/data", g)
	r.Handle(g.PathPrefix+"/data/ws", g)
	r.HandleFunc(g.PathPrefix+"/lib/vue.min.js", vuejsdev.Handler)
	r.HandleFunc(g.PathPrefix+"/lib/vue-router.min.js", vueroutermin.Handler)
	r.HandleFunc(g.PathPrefix+"/lib/vue-resource.min.js", vueresourcemin.Handler)
	r.Handle(g.PathPrefix+"/app.css", g)
	r.Handle(g.PathPrefix+"/custom.css", g)
	r.Handle(g.PathPrefix+"/app.js", g)
	// Allow the template to add files
	r.Handle(g.PathPrefix+"/files/{filename}", g)
	return r
}

// Serve wraps the http.ListenAndServe() function, but adds the
// routes for the gui.
func Serve(g *Gui) error {
	r := NewRouter(g)
	log.Println("Serving gvg on port: ", g.ServerPort)
	gURL := fmt.Sprintf("http://localhost%s%s/", g.ServerPort, g.PathPrefix)
	log.Println("Open your browser and go to ", gURL)
	open.Run(gURL)
	return http.ListenAndServe(g.ServerPort, r)
}
