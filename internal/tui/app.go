package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/tui/views"
	"github.com/vieitesss/jocq/internal/tui/views/explorer"
)

type ViewID int

const (
	ExplorerView ViewID = iota
)

type AppModel struct {
	active        views.View
	views         map[ViewID]views.View
	ExplorerModel explorer.ExplorerModel
}

func NewApp() AppModel {
	views := make(map[ViewID]views.View, 1)
	views[ExplorerView] = explorer.NewExplorerModel()

	return AppModel{
		active: views[ExplorerView],
	}
}

func (AppModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	a.active, cmd = a.active.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a AppModel) View() string {
	return a.active.View()
}
