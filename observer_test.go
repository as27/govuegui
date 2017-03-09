package govuegui

import (
	"testing"
	"time"
)

func TestObserver(t *testing.T) {
	myString := "OldValue"
	testString := "-"
	o := NewObserver()
	o.RefreshTime = time.Millisecond
	o.AddString("myString", &myString)
	o.Start()
	o.Subscribe(func(event, key, oldval string) {
		testString = oldval
		//fmt.Println(event, key, oldval)
	})
	tests := []struct {
		value  string
		expect string
	}{
		{"NewValue", "OldValue"},
		{"AnotherNewValue", "NewValue"},
	}
	for _, test := range tests {
		myString = test.value
		// Test needs to wait becuase the observer needs some time
		time.Sleep(time.Millisecond * 300)
		if testString != test.expect {
			t.Errorf("Expect:'%s' Got:'%s'", test.expect, testString)
		}
	}
	o.Stop()
}
