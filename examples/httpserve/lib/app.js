const PathPrefix = "/govuegui";

const gvgforms = Vue.component('gvgforms',{
    template: `<div>
        <ul>
        <li v-for="form in data.Forms">
            <router-link 
                :to="{name: 'gvgform', params: { formid: form.id}}">
                {{form.id}}</router-link>
        </li>
        </ul>
        <gvgform :data=data :formid=formid></gvgform>
        </div>`,
        props: ['data', 'formid']
})

const gvgform = Vue.component('gvgform', {
    template: `<div>Hi</div>`,
    computed: {
        form: function () {
            forms = this.data.Forms
            for (index = 0; index < forms.length; ++index) {
                if (forms[index].formid == this.formid){
                    return forms[index];
                }
            }
        }
    },
    props: ['data', 'formid']
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
        data: {}
    },
    methods: {
        fetchData: function () {
            this.$http.get(PathPrefix+"/data").then(
                (res) => {
                    this.data = res.body;
                    this.dataLoaded = true;
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