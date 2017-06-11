package photon

import (
	"net/http"

	"github.com/as27/govuegui"
	"github.com/as27/govuegui/gui/photon/src/pkg/photoncss"
	"github.com/as27/govuegui/gui/photon/src/pkg/photonwoff"
)

var Template = govuegui.GuiTemplate{
	CSSHandler: photoncss.Handler,
	CustomCSS: `#govuegui{
        max-width:1200px;
        margin: auto;
    }`,
	Files: map[string]func(w http.ResponseWriter, r *http.Request){
		"photon-entypo.woff": photonwoff.Handler,
	},
	Body: `<body> 
        <div id="govuegui" class="window">
	    <header class="toolbar toolbar-header">
		<h1 class="title">{{appTitle}}</h1>

        <div class="toolbar-actions">
            <div class="btn-group">
                <router-link
                    v-for="form in data.Forms"
                    active-class="active"
                    class="btn btn-default"
                    tag="button"
                    :to="{name: 'gvgform', params: { formid: form.id}}">
                    <span class="icon icon-menu icon-text"></span>
                    {{form.id}}</router-link>
            </div>
        </div>

		</header>	
            <router-view :data=data :forms=forms ></router-view>
        <footer class="toolbar toolbar-footer">
		<h1 class="title">
        <strong>govuigui</strong> 
        by <a href="https://as27.github.io/" target="_blank">Andreas Schr&ouml;pfer</a>
      | <a href="https://github.com/as27/govuegui" target="_blank">govuigui github page</a>
		</h1></footer>
        </div>
    </body>
`,
	GvgForms: `
    <div class="window-content">
            <router-view :data=data :form=forms[formid] :formid=formid></router-view>
     </div>
  `,

	GvgForm: `
     <div class="pane-group">
      <div class="pane-sm sidebar">
        <nav class="nav-group">
            <h5 class="nav-group-title">Boxes</h5>
            <router-link
                v-for="box in form.Boxes"
                active-class="active"
                class="nav-group-item"
                :to="{ name: 'gvgbox', params: { boxid: box.id}}">
                <span class="icon icon-home"></span>
                {{box.id}}</router-link>
        </nav>
      </div>
      <div class="pane">
        <gvgbox :box=myBox :data=data></gvgbox>
        <button class="btn btn-large btn-primary" @click="saveData">Submit</button>
      </div>
    </div>
    `,

	GvgBox: `<form class="padded-more">
    <h2>{{box.id}}</h2>
            <gvgelement 
                v-for="element in box.elements"
                :element=element 
                :data=data></gvgelement>
    </form>
   `,

	GvgElement: `
    <div class="form-group">
        <label v-if="renderLabel" class="label">{{element.label}}</label>
        <component :is=element.type :element=element :data=data v-model="data.Data.data[element.id]"></component>
    </div>`,
	GvgButton: `<div><br><button class="btn btn-large btn-primary" @click="callAction">{{element.label}}</button><br></div>`,

	GvgList: `<div class="text">
   <ul>
   <li v-for="litem in data.Data.data[element.id]">{{litem}}</li>
   </ul> 
    </div>`,
	GvgDropdown: `<div>
    <select v-model="data.Data.data[element.id]" class="form-control">
    <option v-for="oitem in element.options" v-bind:value="oitem.Option">{{oitem.Values[0]}}</option>
    </select>
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

	GvgTextarea: `<textarea class="form-control" v-model="data.Data.data[element.id]"></textarea>`,

	GvgInput: `<input class="form-control" type="text" v-model="data.Data.data[element.id]">`,
}
