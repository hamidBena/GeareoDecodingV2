package cli

import (
	"GDv2/app"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLI struct {
	reader *bufio.Reader
}

func New() *CLI {
	return &CLI{reader: bufio.NewReader(os.Stdin)}
}

func (c *CLI) ClearScreen() {
	fmt.Print("\033[H\033[J")
}

func (c *CLI) ClickAny() {
	fmt.Print("Press ENTER to continue...")
	c.readLine()
}

func (c *CLI) ShowMenu(title string, options []app.MenuOption) string {
	fmt.Println("\n== " + title + " ==")
	for _, opt := range options {
		fmt.Printf("  [%s] %s\n", opt.Key, opt.Label)
	}
	fmt.Print("> ")
	return c.readLine()
}

func (c *CLI) Prompt(question string) string {
	fmt.Print(question + ": ")
	return c.readLine()
}

func (c *CLI) Confirm(question string) bool {
	fmt.Print(question + " (y/n): ")
	answer := c.readLine()
	return strings.ToLower(answer) == "y"
}

func (c *CLI) readLine() string {
	line, _ := c.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func (c *CLI) ShowError(err error) {
	fmt.Println("Error:", err)
	c.ClickAny()
}

func (c *CLI) ShowMessage(msg string) {
	fmt.Println(msg)
}

func (c *CLI) ShowList(title string, items []string) {
	fmt.Println("\n== " + title + " ==")
	for i, item := range items {
		fmt.Printf("  %d. %s\n", i+1, item)
	}
}

/*
	ShowError(err error)
	ShowMessage(msg string)
	ShowList(title string, items []string)
} */
