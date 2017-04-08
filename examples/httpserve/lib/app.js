const PathPrefix = "/govuegui";

const gvgforms = Vue.component('gvgforms',{
    template: `<div>
        <ul>
        
        </ul>
        </div>`,
    props: {
        data: Object
    }
})

const routes = [
    {
        path: '/',
        name: 'home',
        components: {
           default: gvgforms
        },
        props: [],
        children: [
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