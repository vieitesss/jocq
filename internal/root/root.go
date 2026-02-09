package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/features/queries"
)

type ScreenId int

const (
	QueriesScreen ScreenId = iota
)

type RootModel struct {
	active       ScreenId
	queriesModel queries.QueriesModel
}

func NewRoot() RootModel {
	return RootModel{
		active:       QueriesScreen,
		queriesModel: queries.NewQueriesModel(),
	}
}

func (r RootModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (r RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		model tea.Model
	)

	switch r.active {
	case QueriesScreen:
		model, cmd = r.queriesModel.Update(msg)
	}

	cmds = append(cmds, cmd)

	switch model := model.(type) {
	case queries.QueriesModel:
		r.queriesModel = model
	}

	return r, tea.Batch(cmds...)
}

func (r RootModel) View() string {
	switch r.active {
	case QueriesScreen:
		return r.queriesModel.View()
	default:
		return "I don't know what screen to show you :("
	}
}
