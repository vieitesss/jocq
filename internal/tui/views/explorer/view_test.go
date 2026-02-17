package explorer

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/vieitesss/jocq/internal/tree"
)

func TestFitContentWidthTruncatesAndPads(t *testing.T) {
	content := "short\nthis is a very long line"
	got := fitContentWidth(content, 10)
	lines := strings.Split(got, "\n")

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	for i, line := range lines {
		if width := ansi.StringWidth(line); width != 10 {
			t.Fatalf("expected line %d width 10, got %d", i, width)
		}
	}
}

func TestFitContentWidthHandlesEmptyAndZeroWidth(t *testing.T) {
	if got := fitContentWidth("", 10); got != "" {
		t.Fatalf("expected empty content to stay empty, got %q", got)
	}

	if got := fitContentWidth("abc", 0); got != "" {
		t.Fatalf("expected zero width to return empty, got %q", got)
	}
}

func TestExplorerViewDoesNotOverflowWindowWidth(t *testing.T) {
	const (
		windowWidth  = 80
		windowHeight = 24
	)

	e := NewExplorerModel(nil)
	e.ready = true
	e.help.Width = windowWidth
	e.resizeViewports(windowWidth, windowHeight)
	e.In.SetNodes([]tree.Node{{Type: tree.ArrayElement, Value: strings.Repeat("a", 200)}})
	e.Out.SetContent(strings.Repeat("b", 500) + "\n" + strings.Repeat("c", 500))

	view := e.ExplorerView()
	lines := strings.Split(view, "\n")
	inputHeight := lipgloss.Height(e.Input.View())
	helpHeight := lipgloss.Height(e.help.View(e.keys))
	paneEnd := len(lines) - helpHeight
	if paneEnd < inputHeight {
		t.Fatalf("invalid layout boundaries: input=%d help=%d lines=%d", inputHeight, helpHeight, len(lines))
	}

	for i, line := range lines[inputHeight:paneEnd] {
		if width := ansi.StringWidth(line); width > windowWidth {
			t.Fatalf("expected line %d to fit width %d, got %d: %q", i, windowWidth, width, ansi.Strip(line))
		}
	}
}

func TestExplorerViewShowsPendingCountInSourceTitle(t *testing.T) {
	e := NewExplorerModel(nil)
	e.ready = true
	e.help.Width = 80
	e.resizeViewports(80, 24)
	e.In.SetNodes([]tree.Node{{Type: tree.ArrayElement, Value: "x"}})

	e.In, _ = e.In.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	e.In, _ = e.In.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})

	view := ansi.Strip(e.ExplorerView())
	if !strings.Contains(view, "12  •") {
		t.Fatalf("expected source title to include pending count metadata, got %q", view)
	}
}
