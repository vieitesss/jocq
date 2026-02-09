package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/root"
)

func main() {
	_, err := tea.NewProgram(root.NewRoot(), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
