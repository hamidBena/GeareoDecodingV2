package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("DocTabs Container")

	tab1 := container.NewVBox()

	tab1.Add(widget.NewLabel("This is the content of Tab 1"))

	wid1 := widget.NewLabel("name: ")
	val1 := widget.NewEntry()
	wid2 := widget.NewLabel("age: ")
	val2 := widget.NewEntry()
	tab1.Add(container.New(layout.NewFormLayout(), wid1, val1, wid2, val2))

	tabs := container.NewDocTabs(
		container.NewTabItem("Doc 1", widget.NewLabel("Content of document 1")),
		container.NewTabItem("Doc 2", widget.NewForm(
			widget.NewFormItem("Name", widget.NewEntry()),
			widget.NewFormItem("Age", widget.NewEntry()),
		)),
		container.NewTabItem("Doc 3", tab1),
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()
}

/*func main() {
	userInterface := cli.New()
	app.Run(userInterface)
}*/

/*func main() {
	raw, err := os.ReadFile("C:\\Users\\CRT PC\\AppData\\LocalLow\\WitWeld\\Geareo\\projects\\all\\data.json")
	if err != nil {
		fmt.Println("read error:", err)
		os.Exit(1)
	}

	var sf model.SaveFile
	if err := json.Unmarshal(raw, &sf); err != nil {
		fmt.Println("unmarshal error:", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(sf, "", "  ")
	if err != nil {
		fmt.Println("marshal error:", err)
		os.Exit(1)
	}

	err = os.WriteFile("output.json", out, 0644)
	if err != nil {
		fmt.Println("write error:", err)
		os.Exit(1)
	}

	fmt.Println("Successfully wrote output.json")
}*/
