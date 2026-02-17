package treevp

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/tui/views/explorer/tree"
)

func TestCountedDownMotionWithJK(t *testing.T) {
	m := New(40, 5)
	m.SetNodes(testNodes(12))

	m, _ = m.Update(runeKey('5'))
	m, _ = m.Update(runeKey('j'))

	if m.cursor != 5 {
		t.Fatalf("expected cursor at 5, got %d", m.cursor)
	}
}

func TestCountedUpMotionWithJK(t *testing.T) {
	m := New(40, 5)
	m.SetNodes(testNodes(12))
	m.GoToBottom()

	m, _ = m.Update(runeKey('3'))
	m, _ = m.Update(runeKey('k'))

	if m.cursor != 8 {
		t.Fatalf("expected cursor at 8, got %d", m.cursor)
	}
}

func TestCountOnlyAppliesToJK(t *testing.T) {
	m := New(40, 5)
	m.SetNodes(testNodes(12))

	m, _ = m.Update(runeKey('5'))
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	if m.cursor != 1 {
		t.Fatalf("expected down arrow to move one line, got %d", m.cursor)
	}

	m, _ = m.Update(runeKey('j'))
	if m.cursor != 2 {
		t.Fatalf("expected trailing j to move one line, got %d", m.cursor)
	}
}

func TestCountClearsOnNonJKKey(t *testing.T) {
	m := New(40, 5)
	m.SetNodes(testNodes(12))
	m.CursorDownN(6)

	m, _ = m.Update(runeKey('4'))
	m, _ = m.Update(runeKey('g'))
	m, _ = m.Update(runeKey('j'))

	if m.cursor != 1 {
		t.Fatalf("expected cursor at 1 after g then j, got %d", m.cursor)
	}
}

func TestCursorPercent(t *testing.T) {
	var empty Model
	if got := empty.CursorPercent(); got != 0 {
		t.Fatalf("expected empty cursor percent to be 0, got %d", got)
	}

	one := New(40, 5)
	one.SetNodes(testNodes(1))
	if got := one.CursorPercent(); got != 100 {
		t.Fatalf("expected one-line cursor percent to be 100, got %d", got)
	}

	middle := New(40, 5)
	middle.SetNodes(testNodes(11))
	middle.CursorDownN(5)
	if got := middle.CursorPercent(); got != 50 {
		t.Fatalf("expected middle cursor percent to be 50, got %d", got)
	}
}

func runeKey(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func testNodes(count int) []tree.Node {
	nodes := make([]tree.Node, count)
	for i := range nodes {
		nodes[i] = tree.Node{Type: tree.ArrayElement, Value: i}
	}

	return nodes
}
