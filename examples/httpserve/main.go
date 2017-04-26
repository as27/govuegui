package main

import (
	"fmt"
	"log"

	"github.com/as27/govuegui"
)

func main() {
	gui := govuegui.NewGui()
	/*gui.Form("Form1").Box("Box1").Input("Name").Set("Smith")
	gui.Form("Form1").Box("Box1").Input("Age").Set(27)
	gui.Form("Form1").Box("Box1").Input("Age").Option("class", "active", "int")
	b1 := gui.Form("Form1").Box("Box1")
	b1.Textarea("Area").Set("This is the text of the textarea")
	gui.Form("Form1").Box("Box2").Input("Comment").Set("This is a comment.")
	gui.Form("Form2").Box("B2").Input("Title").Set("Mr. Andersson")
	addressForm := gui.Form("Adress form")
	addressForm.Box("Name").Input("First Name")
	addressForm.Box("Name").Input("Last Name")
	addressForm.Box("Private").Input("Street")
	addressForm.Box("Private").Input("City")*/
	inputBox := gui.Form("Input").Box("Input")
	inputBox.Input("x").Set(0)
	inputBox.Input("y").Set(0)
	inputBox.Input("n").Set(0)
	resultBox := gui.Form("Input").Box("Result")

	a := 123
	b := 200
	c := a + b
	gui.Form("Sum").Box("Numbers").Input("A").Set(&a)
	gui.Form("Sum").Box("Numbers").Input("B").Set(&b)
	gui.Form("Sum").Box("Numbers").Input("A + B").Set(&c)
	gui.Form("Sum").Box("Numbers").Text("Result").Set(&c)
	gui.Form("Sum").Box("Numbers").Button("A Plus 1").Action(
		func() {
			a++
			c = a + b
			fmt.Println("A++ called")
		})
	gui.Form("Sum").Box("Numbers").Button("B Plus 1").Action(
		func() {
			b++
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
