package explorer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vieitesss/jocq/internal/query"
	"github.com/vieitesss/jocq/internal/tui/theme"
	"github.com/vieitesss/jocq/internal/tui/views"
)

func (e ExplorerModel) handleDecodedDataFetchedMsg(msg views.DecodedDataFetchedMsg) (ExplorerModel, tea.Cmd) {
	toRender, err := query.Execute(e.query, msg.Content)
	if err != nil {
		e.Out.SetContent(fmt.Sprintf("Error: %v", err))
		return e, nil
	}

	content, err := json.MarshalIndent(toRender, "", "  ")
	if err != nil {
		e.Out.SetContent(fmt.Sprintf("Error: %v", err))
		return e, nil
	}

	e.Out.SetContent(string(content))

	return e, nil
}

func (e ExplorerModel) handleRawDataFetchedMsg(msg views.RawDataFetchedMsg) (ExplorerModel, tea.Cmd) {
	var totalLen int
	for _, line := range msg.Content {
		totalLen += len(line)
	}

	var builder strings.Builder
	builder.Grow(totalLen)
	for _, line := range msg.Content {
		builder.Write(line)
	}

	e.In.SetContent(builder.String())

	return e, nil
}

func (e ExplorerModel) handleKeyMsg(msg tea.KeyMsg) (ExplorerModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch {
	case key.Matches(msg, e.keys.Quit):
		return e, tea.Quit

	case key.Matches(msg, e.keys.NextPane):
		e.focusNextPane()

	case key.Matches(msg, e.keys.PrevPane):
		e.focusPrevPane()

	case key.Matches(msg, e.keys.RunQuery):
		e.query = e.Input.Value()
		cmds = append(cmds, views.FetchDecodedData(e.Data))

	case key.Matches(msg, e.keys.ToggleHelp):
		e.help.ShowAll = !e.help.ShowAll
		e.resizeViewports(e.width, e.height)
	}

	cmd = e.updateFocusedPane(msg)
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e ExplorerModel) handleMouseMsg(msg tea.MouseMsg) (ExplorerModel, tea.Cmd) {
	cmd := e.updateFocusedPane(msg)
	return e, cmd
}

func (e ExplorerModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) (ExplorerModel, tea.Cmd) {
	e.width = msg.Width
	e.height = msg.Height
	e.help.Width = msg.Width
	e.resizeViewports(msg.Width, msg.Height)

	e.ready = true

	return e, nil
}

func (e *ExplorerModel) resizeViewports(width, height int) {
	inWidth := int(float32(width) * e.ratio)
	outWidth := width - inWidth
	inputHeight := lipgloss.Height(e.Input.View())
	helpHeight := lipgloss.Height(e.help.View(e.keys))
	viewportHeight := e.viewportHeight(height, inputHeight, helpHeight)

	e.In.Height = viewportHeight
	e.In.Width = max(0, inWidth-2)

	e.Out.Height = viewportHeight
	e.Out.Width = max(0, outWidth-2)

	e.Input.Width = width
}

func (e *ExplorerModel) updateFocusedPane(msg tea.Msg) tea.Cmd {
	switch e.focused {
	case InputPane:
		var cmd tea.Cmd
		e.Input, cmd = e.Input.Update(msg)
		return cmd

	case InPane:
		var cmd tea.Cmd
		e.In, cmd = e.In.Update(msg)
		return cmd

	case OutPane:
		var cmd tea.Cmd
		e.Out, cmd = e.Out.Update(msg)
		return cmd
	}

	return nil
}

func (e *ExplorerModel) setFocusedPane(pane PaneID) {
	if pane < InputPane || pane > OutPane {
		return
	}

	e.focused = pane
	e.keys.SetFocusMode(pane)
	e.updateInputPlaceholder()
	if e.ready {
		e.resizeViewports(e.width, e.height)
	}

	if pane == InputPane {
		e.Input.Focus()
		e.setInputColor(theme.Pink)
		return
	}

	e.Input.Blur()
	e.setInputColor(theme.Gray)
}

func (e *ExplorerModel) setInputColor(color string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	e.Input.PromptStyle = style
	e.Input.TextStyle = style
	e.Input.PlaceholderStyle = style
	e.Input.Cursor.Style = style
}

func (e *ExplorerModel) updateInputPlaceholder() {
	switch e.focused {
	case InPane:
		e.Input.Placeholder = "[S+Tab]"
	case OutPane:
		e.Input.Placeholder = "[Tab]"
	default:
		e.Input.Placeholder = ""
	}
}

func (e *ExplorerModel) focusPrevPane() {
	if e.focused == InputPane {
		e.setFocusedPane(OutPane)
		return
	}

	e.setFocusedPane(e.focused - 1)
}

func (e *ExplorerModel) focusNextPane() {
	if e.focused == OutPane {
		e.setFocusedPane(InputPane)
		return
	}

	e.setFocusedPane(e.focused + 1)
}
