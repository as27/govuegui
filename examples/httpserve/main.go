package main

import (
	"fmt"
	"log"
	"time"

	"github.com/as27/govuegui"
)

func main() {
	gui := govuegui.NewGui()

	inputBox := gui.Form("Input").Box("Input")
	inputBox.Input("x").Set(0)
	inputBox.Input("y").Set(0)
	inputBox.Input("n").Set(0)
	resultBox := gui.Form("Input").Box("Result")

	a := 123
	b := 200
	c := a + b
	//quitCounter := make(chan bool)
	go counter(gui)
	gui.Form("Sum").Box("Numbers").Input("A").Set(&a)
	gui.Form("Sum").Box("Numbers").Input("B").Set(&b)
	gui.Form("Sum").Box("Numbers").Input("A + B").Set(&c)
	gui.Form("Sum").Box("Numbers").Text("Result").Set(&c)
	gui.Form("Sum").Box("Numbers").Button("WS Update").Action(
		func() {
			err := gui.Update()
			fmt.Println("Gui Update...", err)
		})
	gui.Form("Sum").Box("Numbers").Button("A Plus 1").Action(
		func() {
			a++
			c = a + b
			fmt.Println("A++ called")
		})

	gui.CB = func() {
		//a = gui.Form("Sum").Box("Numbers").Input("A").Get().(int)
		//gui.Form("Sum").Box("Numbers").Input("A + B").Set(a)
		fmt.Println("a wird gesetzt: ", a)
		//d := gui.Form("Sum").Box("Numbers").Input("A").Get()
		c = a + b
		//gui.Form("Sum").Box("Numbers").Input("A + B").Set(a + b)
		n := inputBox.Input("n")
		resultBox.Clear()
		for i := 1; i <= n.Get().(int); i++ {
			name := fmt.Sprintf("n=%d: (x+y)*n", i)
			x := inputBox.Input("x").Get().(int)
			y := inputBox.Input("y").Get().(int)
			resultBox.Text(name).Set((x + y) * i)
		}
	}
	log.Fatal(govuegui.Serve(gui))
}

func counter(g *govuegui.Gui) {
	c := 1
	g.Form("Counter").Box("Numbers").Input("NCounter").Set(&c)
	quit := make(chan bool)
	g.Form("Counter").Box("Numbers").Button("Start/Pause").Action(
		func() {
			quit <- true
		})
	status := g.Form("Counter").Box("Numbers").Text("Status")
	for {
		select {
		case <-time.Tick(time.Second * 2):
			status.Set("Running")
			c++
			err := g.Update()
			if err != nil {
				fmt.Println("--->", err)
			}
		case <-quit:
			start := make(chan bool)
			status.Set("Paused")
			g.Update()
			g.Form("Sum").Box("Numbers").Button("Start/Pause").Action(
				func() {
					start <- true
					status.Set("Waiting for next tick")
					g.Update()
				})
			select {
			case <-start:
				g.Form("Sum").Box("Numbers").Button("Start/Pause").Action(
					func() {
						quit <- true
					})
			}

		}
	}
}
