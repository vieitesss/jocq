package queries

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type QueriesModel struct {
	Src viewport.Model

	width  int
	height int
	text   string
}

func NewQueriesModel() QueriesModel {
	return QueriesModel{}
}

func (q QueriesModel) Init() tea.Cmd {
	return nil
}

func (q QueriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return q, tea.Quit

		case "enter":
			q.text += "\n"

		case "backspace":
			q.text = q.text[:len(q.text)-1]

		default:
			q.text += msg.String()
		}

		q.Src.SetContent(q.text)

	case tea.WindowSizeMsg:
		q.Src = viewport.New(msg.Width, msg.Height)
		q.Src.Style = srcStyle
	}

	return q, nil
}

func (q QueriesModel) View() string {
	return q.QueriesView()
}
