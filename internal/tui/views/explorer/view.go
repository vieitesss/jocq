package explorer

import (
	"github.com/charmbracelet/lipgloss"
)

var roundedBorder = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder())

func (e ExplorerModel) ExplorerView() string {
	v := lipgloss.JoinVertical(lipgloss.Top,
		e.Input.View(),
		lipgloss.JoinHorizontal(lipgloss.Top,
			e.In.View(),
			e.Out.View(),
		),
	)

	return v
}
