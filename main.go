package main

import (
	UI "GDv2/ui/fyne"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Geareo Save Editor")

	myWindow.SetContent(UI.NewOpenFileView(myWindow))

	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}
