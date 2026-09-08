package customWidgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ItemCard(title string, description string, icon fyne.Resource, onClick func()) *fyne.Container {

	card := widget.NewCard(title, description, nil)
	if icon != nil {
		card.Image = &canvas.Image{Resource: icon, FillMode: canvas.ImageFillContain}
	}

	content := widget.NewButton("Build", onClick)

	content.OnTapped = onClick

	card.SetContent(content)

	return container.NewBorder(card, nil, nil, nil)
}
