package explorer

import (
	"github.com/charmbracelet/lipgloss"
)

var paneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

var paneStyleBlur = paneStyle.BorderForeground(lipgloss.Color("240"))
var paneStyleFocus = paneStyle.BorderForeground(lipgloss.Color("205"))

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("240"))
var titleStyleFocus = titleStyle.Foreground(lipgloss.Color("205"))
var hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

func (e ExplorerModel) viewportHeight(height int) int {
	return max(0, height-lipgloss.Height(e.Input.View())-4)
}

func (e ExplorerModel) ExplorerView() string {
	if !e.ready {
		return "Loading..."
	}

	inTitle := titleStyle.Render("Raw JSON [Tab]")
	outTitle := titleStyle.Render("Query Result [S+Tab]")
	inPaneStyle := paneStyleBlur
	outPaneStyle := paneStyleBlur

	switch e.focused {
	case InPane:
		inTitle = titleStyleFocus.Render("Raw JSON")
		outTitle = titleStyle.Render("Query Result [Tab]")
		inPaneStyle = paneStyleFocus

	case OutPane:
		outTitle = titleStyleFocus.Render("Query Result")
		inTitle = titleStyle.Render("Raw JSON [S+Tab]")
		outPaneStyle = paneStyleFocus
	}

	paneIn := inPaneStyle.Width(e.In.Width + 2).Render(lipgloss.JoinVertical(lipgloss.Left,
		inTitle,
		e.In.View(),
	))
	paneOut := outPaneStyle.Width(e.Out.Width + 2).Render(lipgloss.JoinVertical(lipgloss.Left,
		outTitle,
		e.Out.View(),
	))

	v := lipgloss.JoinVertical(lipgloss.Top,
		e.Input.View(),
		hintStyle.Render("Tab: cycle focus (Query -> Raw JSON -> Query Result)"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			paneIn,
			paneOut,
		),
	)

	return v
}
