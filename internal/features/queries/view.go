package queries

import (
	"github.com/charmbracelet/lipgloss"
)

var srcStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder())

func (q QueriesModel) QueriesView() string {
	return q.Src.View()
}
