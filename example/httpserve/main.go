package main

import (
	"fmt"
	"log"
	"time"

	"github.com/as27/govuegui"
	"github.com/as27/govuegui/gui/bulma"
)

var logrHead = [][]string{
	[]string{"Time", "Text"},
}

var logr = logrHead

var gui = govuegui.NewGui(bulma.Template)

func myl(gui *govuegui.Gui, s ...interface{}) {
	ts := time.Now().Format(time.StampMilli)
	logr = append(logr, []string{ts, fmt.Sprintln(s...)})
	gui.Form("Log").Box("Log").Table("Log").Set(logr)
	gui.Form("Log").Box("Log").Table("Log").Update()
}
func main() {
	gui.Form("Log").Box("Log").Table("Log")
	gui.Form("Log").Box("Log").Button("Empty Log").Action(
		func(gui *govuegui.Gui) {
			logr = logrHead
			myl(gui, "Log cleared")
			err := gui.Update("Log")
			if err != nil {
				myl(gui, "Error when updating log", err)
			}
		})
	inputBox := gui.Form("Input").Box("Input")
	ix := inputBox.Input("x")
	inputBox.Input("x").Set(0)
	ix.SetLabel("X value")
	ix.Option("o1", "a", "b")
	ix.Option("class", "active", "float", "left")
	inputBox.Input("y").Set(0)
	inputBox.Input("n").Set(0)

	a := 123
	b := 200
	c := a + b
	//quitCounter := make(chan bool)
	go counter(gui)
	gui.Form("Table").Box("Table").Table("A Table").Set(
		[][]string{
			{"h1", "my header", "hij"},
			{"abc", "def", "hij"},
			{"abc", "def", "hij"},
			{"abc", "def", "hij"},
		})
	gui.Form("Table").Box("Table").Button("Add row").Action(
		func(gui *govuegui.Gui) {
			i := gui.Form("Table").Box("Table").Table("A Table").Get()
			t := i.([][]string)
			t = append(t, []string{"r", "b", "ch"})
			gui.Form("Table").Box("Table").Table("A Table").Set(t)

		})
	gui.Form("Sum").Box("Numbers").Input("A").Set(&a)

	gui.Form("Sum").Box("Numbers").Input("B").Set(&b)
	gui.Form("Sum").Box("Numbers").Input("A + B").Set(&c)
	gui.Form("Sum").Box("Numbers").Text("Result").Set(&c)
	gui.Form("Sum").Box("Numbers").Button("WS Update").Action(
		func(gui *govuegui.Gui) {
			err := gui.Update()
			fmt.Println("Gui Update...", err)
			myl(gui, "Gui Update...")
		})
	gui.Form("Sum").Box("Numbers").Button("A Plus 1").Action(
		func(gui *govuegui.Gui) {
			a++
			c = a + b
			fmt.Println("A++ called")
			myl(gui, "A++ called")
		})

	log.Fatal(govuegui.Serve(gui))
}

func counter(g *govuegui.Gui) {
	c := 1
	ncounter := g.Form("Counter").Box("Numbers").Input("NCounter")
	ncounter.Watch = false
	ncounter.Set(&c)
	quit := make(chan bool)
	spb := g.Form("Counter").Box("Numbers").Button("Start/Pause")
	spb.SetLabel("Pause")
	spb.Action(
		func(g *govuegui.Gui) {
			quit <- true
		})
	status := g.Form("Counter").Box("Numbers").Text("Status")
	for {
		select {
		case <-time.Tick(time.Second * 2):
			spb.SetLabel("Pauses")
			status.Set("Running")
			c++
			myl(gui, "c++", c)
			err := g.Update("NCounter", "Log", "Status")
			if err != nil {
				fmt.Println("--->", err)
			}
		case <-quit:
			start := make(chan bool)
			spb.SetLabel("Start")
			status.Set("Paused")
			g.Update()

			spb.Action(
				func(g *govuegui.Gui) {
					start <- true
					spb.SetLabel("Pause")
					status.Set("Waiting for next tick")
					g.Update()
				})
			select {
			case <-start:
				spb.Action(
					func(g *govuegui.Gui) {
						quit <- true
					})
			}

		}
	}
}
