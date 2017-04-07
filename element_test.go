package govuegui

import (
	"reflect"
	"testing"
)

func TestElements(t *testing.T) {
	gui := NewGui()
	box := gui.Form("myForm").Box("testbox")
	inputElement := box.Input("inputid")
	if inputElement.InputType != INPUT {
		t.Error("Wrong ElementType!")
	}
	// Add options
	inputElement.Option("class", "form")
	inputElement.Option("title", "myElement")
	inputElement.Option("class", "active", "form")
	o := inputElement.getOption("class")
	if o == nil {
		t.Fatal("Option `class` not found!")
	}
	exp := []string{"active", "form"}
	if !reflect.DeepEqual(o.Values, exp) {
		t.Errorf("Got the wrong values!\nExp: %v\nGot: %v",
			exp,
			o.Values,
		)
	}
	//inputElement.Class("active", "form")
	textareaElement := box.Textarea("textid")
	if textareaElement.InputType != TEXTAREA {
		t.Error("Wrong ElementType!")
	}
}
