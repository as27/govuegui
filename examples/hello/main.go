package main

import (
	"log"
	"time"

	"github.com/as27/govuegui"
)

func main() {
	gui := govuegui.NewGui()
	gui.Title = "Hello gui!"
	myForm := gui.Form("myForm")
	myBox1 := myForm.Box("Box 1")
	counter := 0
	guiCounter := myBox1.Text("Counter")
	guiCounter.Set(&counter)
	go count(&counter, guiCounter)
	myBox2 := myForm.Box("Box 2")
	myDropdown := myBox2.Dropdown("MyDropdown")
	myDropdown.Option("key1", "Value 1")
	myDropdown.Option("key2", "Value 2")
	myDropdown.Option("key3", "Value 3")
	myDropdown.Option("key4", "Value 4")
	myDropdown.Set("key3")
	log.Fatal(govuegui.Serve(gui))
}

func count(counter *int, formElement *govuegui.Element) {
	for {
		select {
		case <-time.Tick(time.Second * 2):
			*counter++
			formElement.Update()
		}
	}
}
