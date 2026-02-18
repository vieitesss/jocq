package treevp

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/tree"
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

func TestCountAlsoAppliesToArrowKeys(t *testing.T) {
	m := New(40, 5)
	m.SetNodes(testNodes(12))

	m, _ = m.Update(runeKey('5'))
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	if m.cursor != 5 {
		t.Fatalf("expected down arrow to consume count and move 5 lines, got %d", m.cursor)
	}

	m, _ = m.Update(runeKey('j'))
	if m.cursor != 6 {
		t.Fatalf("expected trailing j to move one line after count consumption, got %d", m.cursor)
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

func TestPendingCount(t *testing.T) {
	m := New(40, 5)
	m.SetNodes(testNodes(12))

	if count, ok := m.PendingCount(); ok || count != 0 {
		t.Fatalf("expected no pending count, got %d %v", count, ok)
	}

	m, _ = m.Update(runeKey('1'))
	m, _ = m.Update(runeKey('2'))
	if count, ok := m.PendingCount(); !ok || count != 12 {
		t.Fatalf("expected pending count 12, got %d %v", count, ok)
	}

	m, _ = m.Update(runeKey('j'))
	if count, ok := m.PendingCount(); ok || count != 0 {
		t.Fatalf("expected pending count to be consumed, got %d %v", count, ok)
	}
}

func TestToggleCollapseOnContainer(t *testing.T) {
	m := New(80, 8)
	m.SetNodes(tree.Flatten([]any{
		map[string]any{
			"a": map[string]any{"b": 1.0},
			"c": 2.0,
		},
	}))

	if len(m.visible) != len(m.nodes) {
		t.Fatalf("expected all nodes to be visible initially, got %d of %d", len(m.visible), len(m.nodes))
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	node, ok := m.CursorNode()
	if !ok {
		t.Fatalf("expected cursor node after collapsing")
	}

	if !node.Collapsed {
		t.Fatalf("expected root container to be collapsed")
	}

	if len(m.visible) != 1 {
		t.Fatalf("expected collapsed root to show a single line, got %d", len(m.visible))
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	node, ok = m.CursorNode()
	if !ok {
		t.Fatalf("expected cursor node after expanding")
	}

	if node.Collapsed {
		t.Fatalf("expected root container to be expanded")
	}

	if len(m.visible) != len(m.nodes) {
		t.Fatalf("expected all nodes visible after expanding, got %d of %d", len(m.visible), len(m.nodes))
	}
}

func TestNavigationSkipsCollapsedChildren(t *testing.T) {
	m := New(80, 8)
	m.SetNodes(tree.Flatten([]any{
		map[string]any{
			"a": map[string]any{"b": 1.0},
			"c": 2.0,
		},
	}))

	m.CursorDownN(1) // ".a" object
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m.CursorDown()

	node, ok := m.CursorNode()
	if !ok {
		t.Fatalf("expected cursor node after moving down")
	}

	if node.Path != ".c" {
		t.Fatalf("expected cursor to skip collapsed children and land on .c, got %q", node.Path)
	}
}

func TestToggleCollapseOnArray(t *testing.T) {
	m := New(80, 8)
	m.SetNodes(tree.Flatten([]any{
		[]any{1.0, 2.0, 3.0},
	}))

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	node, ok := m.CursorNode()
	if !ok {
		t.Fatalf("expected cursor node after collapsing")
	}

	if node.Type != tree.ArrayOpen || !node.Collapsed {
		t.Fatalf("expected collapsed root array node, got %+v", node)
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
