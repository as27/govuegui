package govuegui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/as27/govuegui/vuetemplate"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

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
	router.HandleFunc(g.PathPrefix+"/files/{filename}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars["filename"]
		fmt.Println("--->", filename)
		fileHandler, ok := g.template.Files[filename]
		if !ok {
			http.NotFound(w, r)
			return
		}
		fileHandler(w, r)
	})
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
