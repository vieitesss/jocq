package explorer

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vieitesss/jocq/internal/buffer"
	"github.com/vieitesss/jocq/internal/tui/theme"
	"github.com/vieitesss/jocq/internal/tui/views"
	"github.com/vieitesss/jocq/internal/tui/views/explorer/treevp"
)

type PaneID int

const (
	InputPane PaneID = iota
	InPane
	OutPane
)

type ExplorerModel struct {
	// The source JSON.
	In treevp.Model

	// The resulting JSON after executing a query.
	Out viewport.Model

	Input textinput.Model
	Data  *buffer.Data

	// The amount of terminal width (0 <= ratio <= 1) the In viewport has to take.
	ratio float32

	focused      PaneID
	query        string
	ready        bool
	sourceLoaded bool
	help         help.Model
	keys         KeyMap
	width        int
	height       int
}

func NewExplorerModel(data *buffer.Data) ExplorerModel {
	in := treevp.New(0, 0)
	out := viewport.New(0, 0)

	input := textinput.New()
	input.Focus()
	helpModel := help.New()
	helpModel.Styles.ShortKey = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Gray))
	helpModel.Styles.FullKey = helpModel.Styles.ShortKey
	helpModel.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GrayMuted))
	helpModel.Styles.FullDesc = helpModel.Styles.ShortDesc
	helpModel.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GrayMuted))
	helpModel.Styles.FullSeparator = helpModel.Styles.ShortSeparator
	helpModel.Styles.Ellipsis = helpModel.Styles.ShortSeparator
	helpModel.ShortSeparator = "  •  "
	helpModel.FullSeparator = "     "

	e := ExplorerModel{
		ratio:   0.5,
		In:      in,
		Out:     out,
		Input:   input,
		Data:    data,
		help:    helpModel,
		keys:    NewKeyMap(),
		focused: InputPane,
	}

	e.setFocusedPane(InputPane)

	return e
}

func (e ExplorerModel) Init() tea.Cmd {
	cmds := tea.Batch(
		e.Input.Cursor.BlinkCmd(),
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
