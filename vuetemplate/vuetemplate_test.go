package vuetemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompontent(t *testing.T) {
	c := NewComponent("test")
	c.Template = "<div>abc</div>"
	b := &bytes.Buffer{}
	c.WriteTo(b)
	expect := `const test = Vue.component('test', {
		template: ` + "`" + "<div>abc</div>`" + `});`
	assert.Equal(t, clearString(expect), clearString(b.String()))
}

func TestRouter(t *testing.T) {
	v := NewVue()
	v.Path = "/"
	routes := []Vue{v}
	r := NewRouter("router", routes)
	b := &bytes.Buffer{}
	r.WriteTo(b)
	expect := `const router = VueRouter({
		router: [
			{
				path: '/'
			}
			]
			});`
	assert.Equal(t, clearString(expect), clearString(b.String()))
}

func TestVue(t *testing.T) {
	v := Vue{}
	v.Data = `{
		val1: "value",
		int1: 1
	}`
	v.Path = "/abc/"
	v.Template = `<div>abc</div>`
	expect := `{
		template: ` + "`" + `<div>abc</div>` + "`" + `, 
		data: function(){
			return {
				val1: "value",
				int1: 1
			}
		}, 
		path: '/abc/'
	}`
	b := &bytes.Buffer{}
	v.WriteTo(b)
	assert.Equal(t, clearString(expect), clearString(b.String()))
}

func clearString(s string) string {
	r := strings.NewReplacer(
		"\n", "",
		"\t", "",
		"  ", "")
	return r.Replace(s)
}
func TestComponent(t *testing.T) {
}
