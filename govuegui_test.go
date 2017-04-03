package govuegui

import "testing"
import "reflect"

func TestGuiForm(t *testing.T) {
	gui := NewGui()
	testCases := []struct {
		formid  string
		boxesid []string
	}{
		{
			"form1",
			[]string{"box1"},
		},
		{
			"form2",
			[]string{"box21", "box22", "box23"},
		},
	}
	for _, tc := range testCases {
		f := gui.Form(tc.formid)
		for _, b := range tc.boxesid {
			// Adding boxes
			gui.Form(tc.formid).Box(b)
		}
		if f.ID() != tc.formid {
			t.Errorf("Failed generating Forms.\nExpected: %s Got: %s", tc.formid, f.ID())
		}
		if len(f.Boxes) != len(tc.boxesid) {
			t.Errorf("Form: %s\nNot all boxes are created!", tc.formid)
		}
	}

}

func TestForm(t *testing.T) {
	form := &Form{id: "123"}
	b1 := form.Box("box1")
	if b1 != form.Box("box1") {
		t.Error("box1 has not been created correctly.")
	}
}

func TestElements(t *testing.T) {
	box := &Box{id: "testbox"}
	inputElement := box.Input("inputid")
	if inputElement.inputType != INPUT {
		t.Error("Wrong ElementType!")
	}
	// Add options
	inputElement.Option("class", "form")
	inputElement.Option("title", "myElement")
	inputElement.Option("class", "active", "form")
	o := inputElement.GetOption("class")
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
	if textareaElement.inputType != TEXTAREA {
		t.Error("Wrong ElementType!")
	}
}
