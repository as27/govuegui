const PathPrefix = "/govuegui";
const gvgelement = Vue.component('gvgelement',{
    template: `<div>{{element.id}}</div>`,
    props: {
        element: {
            type: Object,
            default: function(){
                return {id:""}
            }
        }
    }
})

const gvgbox = Vue.component('gvgbox',{
    template: `<div>{{box.id}}
    <div class="gvgelement" v-for="element in box.elements">
    <gvgelement :element=element></gvgelement>
    </div>
    </div>`,
    props: {
        box: {
            type: Object,
            default: function(){
                return {id:""}
            }
        }
    }
});

const gvgform = Vue.component('gvgform', {
    template: `<div><h1>{{form.id}}</h1>
    <div class="box" v-for="box in form.Boxes">
    <gvgbox :box=box></gvgbox></div>
    
    </div>`,
    data: function(){
        return {
            myForm:{id:''}
        }
    },
    methods: {
       
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
    template: `<div>
        <div id="menu">
            <div class="pure-menu">
            <a class="pure-menu-heading" href="#">Company</a>
            <ul class="pure-menu-list">
            <li class="pure-menu-item" v-for="form in data.Forms">
                <router-link 
                    class="pure-menu-link"
                    :to="{name: 'gvgform', params: { formid: form.id}}">
                    {{form.id}}</router-link>
            </li>
            </ul>
            </div>
            <div id="main">
            <router-view :data=data :form=forms[formid] :formid=formid></router-view>
            </div>
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

const app = new Vue({
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
}).$mount('#govuegui');