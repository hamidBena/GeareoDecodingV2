// app/ui.go
package app

type MenuOption struct {
	Label string
	Key   string // e.g. "1", "d", etc — what the user types to pick it
}

type UI interface {
	ShowMenu(title string, options []MenuOption) string // returns chosen Key
	Prompt(question string) string
	Confirm(question string) bool

	ShowError(err error)
	ShowMessage(msg string)
	ShowList(title string, items []string)

	ClearScreen() // optional, can be a no-op for CLI
	ClickAny()    // optional, can be a no-op for CLI
}
