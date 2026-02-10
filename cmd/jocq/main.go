package main

import (
	_ "embed"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/tui"
)

//go:embed example.json
var example string

func main() {
	_, err := tea.NewProgram(tui.NewApp(), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
