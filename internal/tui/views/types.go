package views

import tea "github.com/charmbracelet/bubbletea"

type View interface {
	Init() tea.Cmd
	Update(tea.Msg) (View, tea.Cmd)
	View() string
}
