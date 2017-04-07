package main

import (
	"log"

	"github.com/as27/govuegui"
)

func main() {
	gui := govuegui.NewGui()
	gui.Form("Form1").Box("Box1").Input("Name").Set("Smith")
	gui.Form("Form1").Box("Box1").Input("Age").Set(27)
	log.Fatal(govuegui.Serve())
}
