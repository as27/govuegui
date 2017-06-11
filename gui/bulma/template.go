package bulma

import (
	"net/http"

	"github.com/as27/golib/css/bulma"
	"github.com/as27/govuegui"
	"github.com/as27/govuegui/gui/bulma/src/pkg/fontawesomecss"
	"github.com/as27/govuegui/gui/bulma/src/pkg/fontawesomewoff"
)

var Template = govuegui.GuiTemplate{
	CSSHandler: bulma.Handler,
	Files: map[string]func(w http.ResponseWriter, r *http.Request){
		"font-awesome.css":         fontawesomecss.Handler,
		"fontawesome-webfont.woff": fontawesomewoff.Handler,
	},
	HeadAdd: `<link rel="stylesheet" href="files/font-awesome.css">`,
	Body: `<body class="page-grid">
        <div id="govuegui" class="container">
    <div class="container">
   <div class=""><h1 class="title is-1">{{appTitle}}</h1></div> 
    </div>
            <router-view :data=data :forms=forms ></router-view>
             
        </div>
    <section class="section">
        <footer class="footer">
  <div class="container">
   <div class="content has-text-centered">
      <p>
        <strong>govuigui</strong> 
        by <a href="https://as27.github.io/" target="_blank">Andreas Schr&ouml;pfer</a>
      </p>
      <p><a href="https://github.com/as27/govuegui" target="_blank">govuigui github page</a></p>
    </div>
  </div>
</footer>
    </section>
       
    </body>
`,
	GvgForms: `<div class="columns">
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

	GvgForm: `<div>
    <div class="tabs">
    <ul>
    <router-link v-for="box in form.Boxes"
        active-class="is-active"
        tag="li"
        :to="{ name: 'gvgbox', params: { boxid: box.id}}">
       <a> {{box.id}}</a>
    </router-link>
    </ul>
    </div>
    <div class="box"><gvgbox :box=myBox :data=data></gvgbox></div>
    <button class="button is-primary" @click="saveData">Submit</button>
    </div>`,

	GvgBox: `<div><h2 class="subtitle">{{box.id}}</h2>
    <div class="gvgelement" v-for="element in box.elements">
    <gvgelement :element=element :data=data></gvgelement>
    </div>
    </div>`,

	GvgElement: `<div class="field"><label v-if="renderLabel" class="label">{{element.label}}</label>
    <component :is=element.type :element=element :data=data v-model="data.Data.data[element.id]"></component>
    </div>`,
	GvgButton: `<div><br><button class="button is-primary" @click="callAction">{{element.label}}</button><br></div>`,

	GvgList: `<div class="text">
   <ul>
   <li v-for="litem in data.Data.data[element.id]">{{litem}}</li>
   </ul> 
    </div>`,
	GvgDropdown: `<div class="field">
    <p class="control">
    <span class="select">
    <select v-model="data.Data.data[element.id]">
    <option v-for="oitem in element.options" v-bind:value="oitem.Option">{{oitem.Values[0]}}</option>
    </select></span>
    </p>
    </div>`,

	GvgTable: `<div class="text">
    <table class="table is-narrow">
    <thead>
    <tr><th v-for="cell in data.Data.data[element.id][0]">{{cell}}</th></tr>
    </thead>
    <tr v-for="(row,index) in data.Data.data[element.id]" v-if="index > 0">
    <td v-for="cell in row">{{cell}}</td>
    </tr>
    </table>
    </div>`,
	GvgText: `<div class="text" v-html="data.Data.data[element.id]"></div>`,

	GvgTextarea: `<textarea class="textarea" v-model="data.Data.data[element.id]"></textarea>`,

	GvgInput: `<input class="input" type="text" v-model="data.Data.data[element.id]">`,
}
