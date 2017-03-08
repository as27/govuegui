package govuegui

import (
	"fmt"
	"testing"
)

func TestAddString(t *testing.T) {

}

func TestObserver(t *testing.T) {
	myString := ""
	outerString := "Outer"
	o := NewObserver()
	o.AddString(&myString)
	o.Start()
	o.Subscribe(func(event, key, oldval string) {
		fmt.Println(outerString)
		fmt.Println(event, key, oldval)
	})
	myString = "Me"
	//time.Sleep(time.Second)
	myString = "You"
	//time.Sleep(time.Second)
	o.Stop()
}
