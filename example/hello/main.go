package main

import (
	"fmt"
	"log"
	"time"

	"github.com/as27/govuegui"
	"github.com/as27/govuegui/gui/photon"
)

func main() {
	gui := govuegui.NewGui(photon.Template)
	gui.Title = "Hello gui!"
	myForm := gui.Form("myForm")
	myBox1 := myForm.Box("Box 1")
	counter := 0
	guiCounter := myBox1.Text("Counter")
	guiCounter.Set(&counter)
	guiCounter.Action(func(gui *govuegui.Gui) {
		fmt.Println("Counter changed -->", *guiCounter.Get().(*int))
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
	myInput := myBox2.Input("MyInput")
	myInput.Set("MyInput value")
	myInput.Action(func(gui *govuegui.Gui) {
		fmt.Println(myInput.ID(), "myInput -->", myInput.Get())
		//gui.Update(myInput.ID())
	})

	gui.Form("Test").Box("Test").Button("Start").Action(func(gui *govuegui.Gui) {
		txt := "Das ist der Text<br>"
		gui.Form("Test").Box("Test").Text("Testtxt").Set(txt)
		gui.Update()
		time.Sleep(2 * time.Second)
		gui.Form("Test").Box("Test").Text("Testtxt").Set(txt)
		txt = txt + "Noch eine Zeile<br>"
		gui.Form("Test").Box("Test").Text("Testtxt").Set(txt)
		gui.Update()
		time.Sleep(2 * time.Second)
		txt = txt + "Noch eine Zeile<br>"
		gui.Form("Test").Box("Test").Text("Testtxt").Set(txt)
		gui.Update()
		time.Sleep(2 * time.Second)
		txt = txt + "Noch eine Zeile<br>"
		gui.Form("Test").Box("Test").Text("Testtxt").Set(txt)
		gui.Update()
		time.Sleep(2 * time.Second)
		txt = txt + "Noch eine Zeile<br>"
		gui.Update()
	})
	gui.Form("Test").Box("Test").Text("Testtxt").Set("")
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
