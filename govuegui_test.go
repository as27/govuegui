package govuegui

import (
	"reflect"
	"testing"
)

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

func TestSetGet(t *testing.T) {
	gui := NewGui()
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
