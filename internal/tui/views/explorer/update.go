package explorer

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/query"
	"github.com/vieitesss/jocq/internal/tui/views"
)

func (e ExplorerModel) handleDecodedDataFetchedMsg(msg views.DecodedDataFetchedMsg) (ExplorerModel, tea.Cmd) {
	toRender, err := query.Execute(e.query, msg.Content)
	if err != nil {
		e.Out.SetContent(fmt.Sprint(err))
		return e, nil
	}

	content, err := json.MarshalIndent(toRender, "", "  ")
	if err != nil {
		e.Out.SetContent(fmt.Sprint(err))
		return e, nil
	}

	e.Out.SetContent(string(content))
	e.panes[OutPane] = e.Out

	return e, nil
}

func (e ExplorerModel) handleRawDataFetchedMsg(msg views.RawDataFetchedMsg) (ExplorerModel, tea.Cmd) {
	lines := []string{}
	for _, line := range msg.Content {
		lines = append(lines, string(line))
	}

	e.In.SetContent(strings.Join(lines, ""))
	e.panes[InPane] = e.In

	return e, nil
}

func (e ExplorerModel) handleKeyMsg(msg tea.KeyMsg) (ExplorerModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg.String() {
	case "ctrl+c":
		return e, tea.Quit

	case "tab":
		e.focusNextPane()

	case "shift+tab":
		e.focusPrevPane()

	case "enter":
		if e.focused == InputPane {
			e.query = e.Input.Value()
			cmds = append(cmds, views.FetchDecodedData(e.Data))
		}
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
	inWidth := int(float32(msg.Width) * e.ratio)
	outWidth := msg.Width - inWidth

	e.In.Height = e.viewportHeight(msg.Height)
	e.In.Width = max(0, inWidth-2)

	e.Out.Height = e.viewportHeight(msg.Height)
	e.Out.Width = max(0, outWidth-2)

	e.Input.Width = msg.Width
	e.panes[InputPane] = e.Input
	e.panes[InPane] = e.In
	e.panes[OutPane] = e.Out

	e.ready = true

	return e, nil
}

func (e *ExplorerModel) updateFocusedPane(msg tea.Msg) tea.Cmd {
	switch e.focused {
	case InputPane:
		var cmd tea.Cmd
		e.Input, cmd = e.Input.Update(msg)
		e.panes[InputPane] = e.Input
		return cmd

	case InPane:
		var cmd tea.Cmd
		e.In, cmd = e.In.Update(msg)
		e.panes[InPane] = e.In
		return cmd

	case OutPane:
		var cmd tea.Cmd
		e.Out, cmd = e.Out.Update(msg)
		e.panes[OutPane] = e.Out
		return cmd
	}

	return nil
}

func (e *ExplorerModel) setFocusedPane(pane PaneID) {
	if _, ok := e.panes[pane]; !ok {
		return
	}

	e.focused = pane

	if pane == InputPane {
		e.Input.Focus()
		e.panes[InputPane] = e.Input
		return
	}

	e.Input.Blur()
	e.panes[InputPane] = e.Input
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
