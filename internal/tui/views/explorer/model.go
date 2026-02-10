package explorer

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	text  string
	ratio float32
}

func NewExplorerModel() ExplorerModel {
	in := viewport.New(0, 0)
	in.Style = roundedBorder

	out := viewport.New(0, 0)
	out.Style = roundedBorder

	input := textinput.New()
	input.Focus()
	return ExplorerModel{
		ratio: 0.5,
		In:    in,
		Out:   out,
		Input: input,
	}
}

func (e ExplorerModel) Init() tea.Cmd {
	return e.Input.Cursor.BlinkCmd()
}

func (e ExplorerModel) Update(msg tea.Msg) (views.View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return e, tea.Quit

		case "enter":
			e.Input.Reset()

		case "backspace":
			if len(e.text) > 0 {
				e.text = e.text[:len(e.text)-1]
			}

		default:
			e.text += msg.String()
		}

	case tea.WindowSizeMsg:
		e.In.Height = e.viewportHeight(msg.Height)
		e.In.Width = e.viewportWidth(msg.Width)

		e.Out.Height = e.viewportHeight(msg.Height)
		e.Out.Width = e.viewportWidth(msg.Width)

		e.Input.Width = msg.Width
	}

	var cmd tea.Cmd
	e.Input, cmd = e.Input.Update(msg)

	return e, cmd
}

func (e ExplorerModel) View() string {
	return e.ExplorerView()
}

func (e ExplorerModel) viewportWidth(width int) int {
	w := float32(width) * e.ratio
	return int(w)
}

func (e ExplorerModel) viewportHeight(height int) int {
	return height - lipgloss.Height(e.Input.View())
}
