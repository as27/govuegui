package govuegui

import "testing"

func TestObserver(t *testing.T) {
	myString := ""
	o := NewObserver()
	o.AddString(&myString)
	o.Start()
	myString = "Me"
	o.Stop()
}
