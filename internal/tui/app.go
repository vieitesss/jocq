package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/buffer"
	"github.com/vieitesss/jocq/internal/tui/views"
	"github.com/vieitesss/jocq/internal/tui/views/explorer"
)

type ViewID int

const (
	ExplorerView ViewID = iota
)

type AppModel struct {
	// The views
	ExplorerModel explorer.ExplorerModel

	Active views.View
	Views  map[ViewID]views.View
}

func NewApp(data *buffer.Data) AppModel {
	em := explorer.NewExplorerModel(data)

	views := make(map[ViewID]views.View, 1)
	views[ExplorerView] = em

	return AppModel{
		Active:        views[ExplorerView],
		Views:         views,
		ExplorerModel: em,
	}
}

func (a AppModel) Init() tea.Cmd {
	cmds := tea.Batch(
		tea.WindowSize(),
		a.Active.Init(),
	)

	return cmds
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	a.Active, cmd = a.Active.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a AppModel) View() string {
	return a.Active.View()
}
