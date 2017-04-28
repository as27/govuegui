const PathPrefix = "/govuegui";
const Server = "localhost:2700"

const gvginput = Vue.component('gvginput',{
    template: `<input class="input" type="text" v-model="data.Data.data[element.id]">`,
    props: ['data', 'element']
})
const gvgtextarea = Vue.component('gvgtextarea',{
    template: `<textarea class="textarea" v-model="data.Data.data[element.id]"></textarea>`,
    props: ['data', 'element']
})
const gvgtext = Vue.component('gvgtext',{
    template: `<div class="text">{{data.Data.data[element.id]}}</div>`,
    props: ['data', 'element']
})
const gvgbutton = Vue.component('gvgbutton',{
    template: `<div><br><button class="button is-primary" @click="callAction">{{element.id}}</button><br></div>`,
    props: ['data', 'element'],
    methods:{
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
    }
})
const gvgelement = Vue.component('gvgelement',{
    template: `<div class="field"><label v-if="renderLabel" class="label">{{element.id}}</label>
    
    <component :is=element.type :element=element :data=data v-model="data.Data.data[element.id]"></component>
    
    </div>`,
    components: {
        GVGINPUT: gvginput,
        GVGTEXTAREA: gvgtextarea,
        GVGTEXT: gvgtext,
        GVGBUTTON: gvgbutton
    },
    computed: {
        renderLabel: function(){
            if (this.element.type != 'GVGBUTTON'){
                return true;
            }
            return false;
        }
    },
    props: {
        data: Object,
        element: {
            type: Object,
            default: function(){
                return {id:""}
            }
        }
    }
})

const gvgbox = Vue.component('gvgbox',{
    template: `<div><h2 class="subtitle">{{box.id}}</h2>
    <div class="gvgelement" v-for="element in box.elements">
    <gvgelement :element=element :data=data></gvgelement>
    </div>
    </div>`,
    props: {
        data: Object,
        box: {
            type: Object,
            default: function(){
                return {id:""}
            }
        }
    }
});

const gvgform = Vue.component('gvgform', {
    template: `<div><h1 class="title is-1">{{form.id}}</h1>
    <div class="box" v-for="box in form.Boxes">
    <gvgbox :box=box :data=data></gvgbox></div>
    <button class="button is-primary" @click="saveData">Submit</button>
    </div>`,
    data: function(){
        return {
            myForm:{id:''}
        }
    },
    methods: {
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
    },
    components: {
        gvgbox: gvgbox
    },
    props: {
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
    }
})

const gvgforms = Vue.component('gvgforms', {
    template: `<div class="columns">
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
        
        </div>`,

    data: function(){
        return {
        }
    },
    props:{
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
    } ,
    //['data', 'forms', 'formid'],
    components: {
        gvgform: gvgform
    }
})



const routes = [
    {
        path: '/',
        name: 'home',
        components: {
           default: gvgforms
        },
        props: {
            default: true
        },
        children: [
             {
                path: '/:formid',
                name: 'gvgform',
                component: gvgform,
                props: true
            }
        ]
    }
]

const router = new VueRouter({
    routes: routes
});

var socket = new WebSocket("ws://"+Server+PathPrefix+"/data/ws");
socket.onmessage = function(evt){
    var newData = JSON.parse(evt.data);
    //console.log(evt.data); //TODO: Remove in production
    console.log(evt);
    app.data = newData;
};

const app = new Vue({
    router,
    data: {
        data: {},
        dataNew: {},
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
}).$mount('#govuegui');