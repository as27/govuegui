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
	notification := make(chan stringValue)
	o.Subscribe(notification)
	tests := []struct {
		value  string
		expect string
	}{
		{"NewValue", "OldValue"},
		{"AnotherNewValue", "NewValue"},
	}
	for _, test := range tests {
		time.Sleep(time.Millisecond * 40)
		myString = test.value

		select {
		case n := <-notification:
			//fmt.Printf("%v", n)
			testString = n.value

		}

		// Test needs to wait becuase the observer needs some time
		time.Sleep(time.Millisecond * 30)
		if testString != test.expect {
			t.Errorf("Expect:'%s' Got:'%s'", test.expect, testString)
		}
	}
	o.Stop()
}
