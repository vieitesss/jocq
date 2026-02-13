package explorer

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/buffer"
	"github.com/vieitesss/jocq/internal/tui/views"
)

type PaneID int

const (
	InputPane PaneID = iota
	InPane
	OutPane
)

type ExplorerModel struct {
	// Each viewport (future component) will have:
	// - width
	// - height
	// - content
	// - mode
	// - flags

	In    viewport.Model
	Out   viewport.Model
	Input textinput.Model

	focused PaneID

	Data  *buffer.Data
	ratio float32
	query string
	ready bool
}

func NewExplorerModel(data *buffer.Data) ExplorerModel {
	in := viewport.New(0, 0)
	out := viewport.New(0, 0)

	input := textinput.New()
	input.Focus()

	e := ExplorerModel{
		ratio:   0.5,
		In:      in,
		Out:     out,
		Input:   input,
		Data:    data,
		focused: InputPane,
	}

	e.setFocusedPane(InputPane)

	return e
}

func (e ExplorerModel) Init() tea.Cmd {
	cmds := tea.Batch(
		e.Input.Cursor.BlinkCmd(),
		views.FetchRawData(e.Data),
		views.FetchDecodedData(e.Data),
	)
	return cmds
}

func (e ExplorerModel) Update(msg tea.Msg) (views.View, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case views.DecodedDataFetchedMsg:
		e, cmd = e.handleDecodedDataFetchedMsg(msg)

	case views.RawDataFetchedMsg:
		e, cmd = e.handleRawDataFetchedMsg(msg)

	case tea.KeyMsg:
		e, cmd = e.handleKeyMsg(msg)

	case tea.MouseMsg:
		e, cmd = e.handleMouseMsg(msg)

	case tea.WindowSizeMsg:
		e, cmd = e.handleWindowSizeMsg(msg)
	}

	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e ExplorerModel) View() string {
	return e.ExplorerView()
}
