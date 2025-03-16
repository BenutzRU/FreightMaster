package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Тест Fyne")
	w.SetContent(container.NewVBox(
		widget.NewLabel("Если ты видишь это окно, Fyne работает!"),
	))
	w.ShowAndRun()
}
