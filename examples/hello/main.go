package main

import (
	"fmt"
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
	guiCounter.Action(func(gui *govuegui.Gui) {
		fmt.Println("Counter changed -->", guiCounter.Get().(*int))
	})
	go count(&counter, guiCounter)
	myBox2 := myForm.Box("Box 2")
	myDropdown := myBox2.Dropdown("MyDropdown")
	myDropdown.Option("key1", "Value 1")
	myDropdown.Option("key2", "Value 2")
	myDropdown.Option("key3", "Value 3")
	myDropdown.Option("key4", "Value 4")
	myDropdown.Set("key3")
	myDropdown.Action(func(gui *govuegui.Gui) {
		fmt.Printf("%#v", myDropdown.Get())
	})
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
