package govuegui

import (
	"reflect"
	"testing"
)

func TestGuiForm(t *testing.T) {
	gui := NewGui(GuiTemplate{})
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
	form := &Form{
		Key:     "123",
		Options: make(map[string]*Option),
	}
	b1 := form.Box("box1")
	if b1 != form.Box("box1") {
		t.Error("box1 has not been created correctly.")
	}
	testOpts := []string{"val1", "val2", "val3"}
	form.Option("f1option", testOpts...)
	if !reflect.DeepEqual(
		form.Options["f1option"].Values,
		testOpts,
	) {
		t.Errorf("Form Options not set correct\nExp: %v\nGot: %v",
			testOpts,
			form.Options["f1option"].Values,
		)
	}
}

func TestBox(t *testing.T) {
	gui := NewGui(GuiTemplate{})
	b1 := gui.Form("F1").Box("b1")
	b1.Input("I1")
	if b1.Elements[0].ID() != "F1-b1-I1" {
		t.Errorf("Input field not added to Box1!\nGot: %s", b1.Elements[0].ID())
	}
	b1.Option("title", "This is my title")
	if b1.Options["title"].Values[0] != "This is my title" {
		t.Error("Option was not set correct at Box1")
	}
}

func TestSetGet(t *testing.T) {
	gui := NewGui(GuiTemplate{})
	testString := "Value of a string"
	err := gui.Form("myForm").Box("Box1").Textarea("t1").Set(testString)
	if err != nil {
		t.Error(err)
	}
	got := gui.Form("myForm").Box("Box1").Textarea("t1").Get()
	if !reflect.DeepEqual(got, testString) {
		t.Errorf("Another value expected.\nGot: %s\nExp: %s",
			got,
			testString,
		)
	}
}
