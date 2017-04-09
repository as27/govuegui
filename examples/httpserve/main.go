package main

import (
	"log"

	"github.com/as27/govuegui"
)

func main() {
	gui := govuegui.NewGui()
	gui.Form("Form1").Box("Box1").Input("Name").Set("Smith")
	gui.Form("Form1").Box("Box1").Input("Age").Set(27)
	gui.Form("Form1").Box("Box1").Input("Age").Option("class", "active", "int")
	b1 := gui.Form("Form1").Box("Box1")
	b1.Textarea("Area").Set("This is the text of the textarea")
	gui.Form("Form1").Box("Box2").Input("Comment").Set("This is a comment.")
	gui.Form("Form2").Box("B2").Input("Title").Set("Mr. Andersson")
	log.Fatal(govuegui.Serve(gui))
}
