package govuegui

import (
	"net/http"

	"github.com/as27/govuegui/vuetemplate"
)

func vueappHandler(w http.ResponseWriter, r *http.Request) {
	vuetemplate.NewJSElement(vuetemplate.CONSTANT, "PathPrefix", PathPrefix).WriteTo(w)
	vuetemplate.NewJSElement(vuetemplate.CONSTANT, "Server", "localhost"+ServerPort).WriteTo(w)
	comp := vuetemplate.NewComponent("gvginput")
	comp.Template = `<input class="input" type="text" v-model="data.Data.data[element.id]">`
	comp.Props = "['data', 'element']"
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgtextarea")
	comp.Template = `<textarea class="textarea" v-model="data.Data.data[element.id]"></textarea>`
	comp.Props = "['data', 'element']"
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgtext")
	comp.Template = `<div class="text">{{data.Data.data[element.id]}}</div>`
	comp.Props = "['data', 'element']"
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgbutton")
	comp.Template = `<div><br><button class="button is-primary" @click="callAction">{{element.label}}</button><br></div>`
	comp.Props = "['data', 'element']"
	comp.WriteTo(w)

}
