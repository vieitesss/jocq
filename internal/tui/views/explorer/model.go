package explorer

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vieitesss/jocq/internal/buffer"
	"github.com/vieitesss/jocq/internal/tui/views"
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

	Data  *buffer.Data
	Ratio float32
}

func NewExplorerModel(data *buffer.Data) ExplorerModel {
	in := viewport.New(0, 0)
	in.Style = roundedBorder

	out := viewport.New(0, 0)
	out.Style = roundedBorder

	input := textinput.New()
	input.Focus()
	return ExplorerModel{
		Ratio: 0.5,
		In:    in,
		Out:   out,
		Input: input,
		Data:  data,
	}
}

func (e ExplorerModel) Init() tea.Cmd {
	cmds := tea.Batch(
		e.Input.Cursor.BlinkCmd(),
		views.FetchRawData(e.Data),
	)
	return cmds
}

func (e ExplorerModel) Update(msg tea.Msg) (views.View, tea.Cmd) {
	switch msg := msg.(type) {
	case views.RawDataFetchedMsg:
		lines := []string{}

		for _, line := range msg.Content {
			lines = append(lines, string(line))
		}

		e.In.SetContent(strings.Join(lines, ""))

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return e, tea.Quit

		case "enter":
			e.Input.Reset()
		}

	case tea.WindowSizeMsg:
		inWidth := int(float32(msg.Width) * e.Ratio)
		outWidth := msg.Width - inWidth

		e.In.Height = e.viewportHeight(msg.Height)
		e.In.Width = inWidth

		e.Out.Height = e.viewportHeight(msg.Height)
		e.Out.Width = outWidth

		e.Input.Width = msg.Width
	}

	var cmd tea.Cmd
	e.Input, cmd = e.Input.Update(msg)

	return e, cmd
}

func (e ExplorerModel) View() string {
	return e.ExplorerView()
}

func (e ExplorerModel) viewportHeight(height int) int {
	return height - lipgloss.Height(e.Input.View())
}
