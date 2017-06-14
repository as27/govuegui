package govuegui

import (
	"github.com/as27/govuegui/vuetemplate"
)

type VueComponent struct {
	Name        string
	ElementType ElementType
	CompFunc    func(GuiTemplate) *vuetemplate.Component
}

var VueComponents = []VueComponent{
	{"gvginput", INPUT, nil},
	{"gvgtextarea", TEXTAREA, nil},
	{"gvgtext", TEXT, nil},
	{"gvgtable", TABLE, nil},
	{"gvgdropdown", DROPDOWN, nil},
	{"gvglist", LIST, nil},
	{"gvgbutton", BUTTON, vuegvgbutton},
}

func vuegvgdefaultelement(name, template string) *vuetemplate.Component {
	comp := vuetemplate.NewComponent(name)
	comp.Template = template
	comp.Props = "['data', 'element']"
	return comp
}
func vuegvgbutton(t GuiTemplate) *vuetemplate.Component {
	comp := vuetemplate.NewComponent("gvgbutton")
	comp.Template = t.GvgButton
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
	return comp
}

func vuegvgelement(t GuiTemplate) *vuetemplate.Component {
	comp := vuetemplate.NewComponent("gvgelement")
	comp.Template = t.GvgElement
	comp.Data = `{updateRuns: false}`
	comp.Props = "['data', 'element']"
	comp.Components = `{
        GVGINPUT: gvginput,
        GVGTEXTAREA: gvgtextarea,
        GVGTEXT: gvgtext,
        GVGTABLE: gvgtable,
        GVGLIST: gvglist,
        GVGDROPDOWN: gvgdropdown,
        GVGBUTTON: gvgbutton }`
	comp.Watch = `{
        datastring: function(val, oldVal){
            if (this.updateRuns===false && this.element.watch===true){
                //console.log("-------------------");
                //console.log("val: "+val);
                //console.log("oldVal: "+oldVal);
                //console.log(this.element.id);
                this.$root.saveData();
                this.$root.callAction(this.element.id);
            }
        }}`
	comp.Computed = `{
        datastring: function(){
            return JSON.stringify(this.data.Data.data[this.element.id])
        },
        renderLabel: function(){
            if (this.element.type != 'GVGBUTTON'){
                return true;
            }
            return false;
        } }`
	comp.BeforeUpdate = `function(){
        this.updateRuns = true;
    }`
	comp.Updated = `function(){
        this.updateRuns = false;
    }`
	comp.BeforeMount = `function(){
        this.updateRuns = true;
        console.log("mounted called");
    }`
	comp.Props = `{
        data: Object,
        element: {
            type: Object,
            default: function(){
                return {
                    id:"",
                    label:"",
                    options:[{Option:"",Values:[""]}]
                }
            }
        }
    }`
	return comp
}

func vuegvgbox(t GuiTemplate) *vuetemplate.Component {
	comp := vuetemplate.NewComponent("gvgbox")
	comp.Template = t.GvgBox
	comp.Props = `{
        data: Object,
        box: {
            type: Object,
            default: function(){
                return {id:""}
            }
        }
    }`

	return comp
}

func vuegvgform(t GuiTemplate) *vuetemplate.Component {
	gvgform := vuetemplate.NewComponent("gvgform")
	gvgform.Template = t.GvgForm
	gvgform.Data = `{
        myForm:{id:''},
        myBox:{id:''},
    }`
	gvgform.Methods = `{
        saveData: function () {
            console.log("saveData called");
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
        },
        getBox: function (){
            if (typeof this.boxid === "undefined"){
                this.myBox = this.form.Boxes[0]
                return
            };
            for (var i = 0; i < this.form.Boxes.length;i++){
                myBox = this.form.Boxes[i]
                if (myBox.id===this.boxid){
                    this.myBox = myBox;
                    return;
                }
            }
        }
    }`
	gvgform.Watch = `{
        '$route': 'getBox'
    }`
	gvgform.Mounted = `function() {
        this.getBox();
    }`
	// gvgform.Computed = "{}"
	gvgform.Components = "{gvgbox: gvgbox}"
	gvgform.Props = `{
        data: Object,
        formid: String,
        boxid: String,
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
	return gvgform
}

func vuegvgforms(t GuiTemplate) *vuetemplate.Component {
	comp := vuetemplate.NewComponent("gvgforms")
	comp.Template = t.GvgForms
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
	return comp
}

var vueappstring = `/*----------------------*/
const router = new VueRouter({routes: [{
name: 'home', 
props: {default: true}, 
children: [
             {
                path: '/:formid',
                name: 'gvgform',
                component: gvgform,
                props: true,
                children: [
                    {
                        path: '/:formid/:boxid',
                        name: 'gvgbox',
                        component: gvgbox,
                        props: true
                    }
                ]
            }
        ], 
components: {default: gvgforms}, 
path: '/'}]});
/*----------------------*/
var socket = new WebSocket("ws://"+Server+PathPrefix+"/data/ws");
/*----------------------*/
socket.onmessage = function(evt){
    var newData = JSON.parse(evt.data);
    var updateAll = true;
        for(var dataKey in newData.UpdateData.data){
            app.data.Data.data[dataKey] = newData.UpdateData.data[dataKey];
            updateAll = false;
        }
    if (updateAll){
    app.data = newData;
    }
};	const app = new Vue({
    router,
    data: {
        appTitle: appTitle,
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
            console.log("saveData called");
            this.$http.post(PathPrefix + "/data", this.data).then(
                res => {
                    console.log("SaveData success");
                },res=>{console.log("There is an error")}
            );
        },
        callAction: function(id) {
            this.$http.get(PathPrefix+"/data?action="+id);
        }
    },
    created: function () {
        this.fetchData()
    },
});
/*----------------------*/
app.$mount('#govuegui');`

var htmlTemplate = `<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    
   
    <script src="{{ .PathPrefix }}/lib/vue.min.js"></script>
    <script src="{{ .PathPrefix }}/lib/vue-router.min.js"></script>
    <script src="{{ .PathPrefix }}/lib/vue-resource.min.js"></script>

     <link rel="stylesheet" type="text/css" href="{{ .PathPrefix }}/app.css" >
     <link rel="stylesheet" type="text/css" href="{{ .PathPrefix }}/custom.css" >
     {{ .HeadAdd }}
   <title>{{ .Title }}</title>
    </head>
    {{ .Body }} 
    <script src="{{ .PathPrefix }}/app.js"></script>
   </html>`
