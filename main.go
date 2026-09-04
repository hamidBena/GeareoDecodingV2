package main

import (
	UI "GDv2/ui/fyne"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Geareo Save Editor")

	// Initially show an open-file view.
	myWindow.SetContent(UI.NewOpenFileView(myWindow))

	myWindow.Resize(fyne.NewSize(1000, 800))
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
