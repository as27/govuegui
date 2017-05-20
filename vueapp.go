package govuegui

import (
	"fmt"
	"net/http"

	"bytes"

	"github.com/as27/govuegui/vuetemplate"
)

func vueappHandler(w http.ResponseWriter, r *http.Request) {
	serverVar := "localhost" + ServerPort
	vuetemplate.NewJSElement(vuetemplate.CONSTANT, "PathPrefix", PathPrefix).WriteTo(w)
	vuetemplate.NewJSElement(vuetemplate.CONSTANT, "Server", serverVar).WriteTo(w)
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

	comp = vuetemplate.NewComponent("gvgtable")
	comp.Template = `<div class="text">
    <table class="table is-narrow">
    <thead>
    <tr><th v-for="cell in data.Data.data[element.id][0]">{{cell}}</th></tr>
    </thead>
    <tr v-for="(row,index) in data.Data.data[element.id]" v-if="index > 0">
    <td v-for="cell in row">{{cell}}</td>
    </tr>
    </table>
    </div>`
	comp.Props = "['data', 'element']"
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgbutton")
	comp.Template = `<div><br><button class="button is-primary" @click="callAction">{{element.label}}</button><br></div>`
	comp.Props = "['data', 'element']"
	comp.Methods = `{
        callAction: function(){
            this.$http.get(PathPrefix+"/data?action="+this.element.id).then(
                res => {
                    this.$http.get(PathPrefix + "/data").then(
                        (res) => {
                            app.data = res.body;
                            for (var i = 0; i < app.data.Forms.length; i++) {
                                app.forms[app.data.Forms[i].id] = app.data.Forms[i];
                            }
                        }
                    );
                },res=>{console.log("There is an error")}
            );
        }
    }`
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgelement")
	comp.Template = `<div class="field"><label v-if="renderLabel" class="label">{{element.label}}</label>
    <component :is=element.type :element=element :data=data v-model="data.Data.data[element.id]"></component>
    </div>`
	comp.Props = "['data', 'element']"
	comp.Components = `{
        GVGINPUT: gvginput,
        GVGTEXTAREA: gvgtextarea,
        GVGTEXT: gvgtext,
        GVGTABLE: gvgtable,
        GVGBUTTON: gvgbutton }`
	comp.Computed = `{
        renderLabel: function(){
            if (this.element.type != 'GVGBUTTON'){
                return true;
            }
            return false;
        } }`
	comp.Props = `{
        data: Object,
        element: {
            type: Object,
            default: function(){
                return {
                    id:"",
                    label:""
                }
            }
        }
    }`
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgbox")
	comp.Template = `<div><h2 class="subtitle">{{box.id}}</h2>
    <div class="gvgelement" v-for="element in box.elements">
    <gvgelement :element=element :data=data></gvgelement>
    </div>
    </div>`
	comp.Props = `{
        data: Object,
        box: {
            type: Object,
            default: function(){
                return {id:""}
            }
        }
    }`
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgform")
	comp.Template = `<div><h1 class="title is-1">{{form.id}}</h1>
    <div class="box" v-for="box in form.Boxes">
    <gvgbox :box=box :data=data></gvgbox></div>
    <button class="button is-primary" @click="saveData">Submit</button>
    </div>`
	comp.Data = "{myForm:{id:''}}"
	comp.Methods = `{
        saveData: function () {
            this.$http.post(PathPrefix + "/data", this.data).then(
                res => {
                    console.log("post ready");
                    this.$http.get(PathPrefix + "/data").then(
                        (res) => {
                            app.data = res.body;
                            for (var i = 0; i < app.data.Forms.length; i++) {
                                app.forms[app.data.Forms[i].id] = app.data.Forms[i];
                            }
                        }
                    );
                },res=>{console.log("There is an error")}
            );
        }
    }`
	comp.Components = "{gvgbox: gvgbox}"
	comp.Props = `{
        data: Object,
        formid: String,
        form: {
            type: Object,
            default: function(){
                return {
                    id: "",
                    Boxes: [{id:""}]
                }
            }
        }
    }`
	comp.WriteTo(w)

	comp = vuetemplate.NewComponent("gvgforms")
	comp.Template = `<div class="columns">
        <div class="column is-one-quarter">
        <aside class="menu">
        <p class="menu-label">
            Forms
        </p> 
            <ul class="menu-list">
            <li  v-for="form in data.Forms">
                <router-link 
                    active-class="is-active"
                    :to="{name: 'gvgform', params: { formid: form.id}}">
                    {{form.id}}</router-link>
            </li>
            </ul>
        </aside>
        </div>
            <div class="column">
            <router-view :data=data :form=forms[formid] :formid=formid></router-view>
            </div>
        
        </div>`
	comp.Data = "{}"
	comp.Props = `{
        data: Object,
        formid: {
            type: String,
            default: "defObj"
        },
        forms: {
            type: Object,
            default: {
                defObj: {id: "defObj"}
            }
        }
    }`
	comp.Components = `{
        gvgform: gvgform
    }`
	comp.WriteTo(w)
	route := vuetemplate.NewVue()
	route.Path = "/"
	route.Name = "home"
	route.Components = "{default: gvgforms}"
	route.Props = "{default: true}"
	route.Children = `[
             {
                path: '/:formid',
                name: 'gvgform',
                component: gvgform,
                props: true
            }
        ]`
	routes := []vuetemplate.Vue{route}
	router := vuetemplate.NewRouter("router", routes)
	router.WriteTo(w)
	ws := vuetemplate.NewJSElement(
		vuetemplate.WEBSOCKET,
		"socket",
		fmt.Sprintf("ws://%s%s/data/ws", serverVar, PathPrefix),
	)
	ws.WriteTo(w)
	b := bytes.NewBufferString(`socket.onmessage = function(evt){
    var newData = JSON.parse(evt.data);
    var updateAll = true;
        for(var dataKey in newData.UpdateData.data){
            app.data.Data.data[dataKey] = newData.UpdateData.data[dataKey];
            updateAll = false;
        }
    if (updateAll){
    app.data = newData;
    }
};	`)
	w.Write(b.Bytes())

	app := vuetemplate.NewJSElement(
		vuetemplate.VUEAPP,
		"app",
		`{
    router,
    data: {
        data: {},
        forms: {}
    },
    methods: {
        fetchData: function () {
            this.$http.get(PathPrefix+"/data").then(
                (res) => {
                    this.data = res.body;
                    for (var i = 0;i < this.data.Forms.length;i++){
                        this.forms[this.data.Forms[i].id] = this.data.Forms[i];
                    }
                }
            );
        },
        saveData: function () {
            this.$http.post('/collection/' + this.cid + '/' + this.iid, this.item.data);
        }
    },
    created: function () {
        this.fetchData()
    },
}`)
	app.WriteTo(w)
	b = bytes.NewBufferString("app.$mount('#govuegui');")
	w.Write(b.Bytes())
}

var htmlTemplate = `<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    
   
    <script src="{{ .PathPrefix }}/lib/vue.min.js"></script>
    <script src="{{ .PathPrefix }}/lib/vue-router.min.js"></script>
    <script src="{{ .PathPrefix }}/lib/vue-resource.min.js"></script>

     <link rel="stylesheet" type="text/css" href="{{ .PathPrefix }}/lib/bulma.css" >
    
    </head>
    <body class="page-grid">

        <div id="govuegui" class="container">
            <router-view :data=data :forms=forms ></router-view>
             
        </div>
        <script src="{{ .PathPrefix }}/app.js"></script>
    </body>
</html>`
