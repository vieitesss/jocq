package treevp

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/tree"
)

type Model struct {
	nodes   []tree.Node
	visible []int
	cursor  int
	offset  int

	pendingCount int
	hasCount     bool

	Width  int
	Height int
}

func New(width, height int) Model {
	return Model{
		Width:  max(0, width),
		Height: max(0, height),
	}
}

// SetNodes takes ownership of the provided slice.
// Callers must not mutate it after this call.
func (m *Model) SetNodes(nodes []tree.Node) {
	m.nodes = nodes
	m.rebuildVisible()
	m.cursor = 0
	m.offset = 0
	m.ensureCursorVisible()
}

func (m *Model) SetSize(width, height int) {
	m.Width = max(0, width)
	m.Height = max(0, height)
	m.ensureCursorVisible()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

func (m *Model) CursorUp() {
	m.CursorUpN(1)
}

func (m *Model) CursorUpN(n int) {
	if len(m.visible) == 0 {
		return
	}

	if n <= 0 {
		return
	}

	m.cursor = max(0, m.cursor-n)
	m.ensureCursorVisible()
}

func (m *Model) CursorDown() {
	m.CursorDownN(1)
}

func (m *Model) CursorDownN(n int) {
	if len(m.visible) == 0 {
		return
	}

	if n <= 0 {
		return
	}

	m.cursor = min(len(m.visible)-1, m.cursor+n)
	m.ensureCursorVisible()
}

func (m *Model) PageUp() {
	if len(m.visible) == 0 {
		return
	}

	delta := max(1, m.Height/2)
	m.cursor = max(0, m.cursor-delta)
	m.ensureCursorVisible()
}

func (m *Model) PageDown() {
	if len(m.visible) == 0 {
		return
	}

	delta := max(1, m.Height/2)
	m.cursor = min(len(m.visible)-1, m.cursor+delta)
	m.ensureCursorVisible()
}

func (m *Model) GoToTop() {
	if len(m.visible) == 0 {
		return
	}

	m.cursor = 0
	m.ensureCursorVisible()
}

func (m *Model) GoToBottom() {
	if len(m.visible) == 0 {
		return
	}

	m.cursor = len(m.visible) - 1
	m.ensureCursorVisible()
}

func (m Model) CursorNode() (tree.Node, bool) {
	if len(m.visible) == 0 {
		return tree.Node{}, false
	}

	nodeIndex := m.visible[m.cursor]
	return m.nodes[nodeIndex], true
}

func (m Model) CursorPercent() int {
	if len(m.visible) == 0 {
		return 0
	}

	if len(m.visible) == 1 {
		return 100
	}

	percent := int(math.Round(float64(m.cursor) * 100 / float64(len(m.visible)-1)))
	return max(0, min(100, percent))
}

func (m Model) PendingCount() (int, bool) {
	if !m.hasCount {
		return 0, false
	}

	return m.pendingCount, true
}

func (m *Model) ToggleCollapse() bool {
	nodeIndex, ok := m.cursorNodeIndex()
	if !ok {
		return false
	}

	node := m.nodes[nodeIndex]
	if !node.Collapsible {
		return false
	}

	m.nodes[nodeIndex].Collapsed = !node.Collapsed
	m.rebuildVisible()
	m.setCursorByNodeIndex(nodeIndex)
	m.ensureCursorVisible()
	return true
}

func (m *Model) ensureCursorVisible() {
	if len(m.visible) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}

	m.cursor = max(0, min(len(m.visible)-1, m.cursor))

	if m.Height <= 0 {
		m.offset = 0
		return
	}

	if m.cursor < m.offset {
		m.offset = m.cursor
	}

	if m.cursor >= m.offset+m.Height {
		m.offset = m.cursor - m.Height + 1
	}

	maxOffset := max(0, len(m.visible)-m.Height)
	m.offset = max(0, min(maxOffset, m.offset))
}

func (m Model) cursorNodeIndex() (int, bool) {
	if len(m.visible) == 0 {
		return 0, false
	}

	return m.visible[m.cursor], true
}

func (m *Model) setCursorByNodeIndex(nodeIndex int) {
	for i, visibleIndex := range m.visible {
		if visibleIndex == nodeIndex {
			m.cursor = i
			return
		}
	}

	// Clamp cursor to the last visible row when the previous node is no longer visible.
	m.cursor = min(m.cursor, max(0, len(m.visible)-1))
}

// rebuildVisible scans flattened nodes and builds a list of row indices to render.
// When it hits a collapsed container opener, it records that container depth and
// skips every deeper node until the matching close node at the same depth.
func (m *Model) rebuildVisible() {
	m.visible = m.visible[:0]
	skipDepth := -1 // -1 means no active collapsed range.

	for i, node := range m.nodes {
		if skipDepth >= 0 {
			if node.Depth > skipDepth {
				continue
			}

			if isContainerClose(node) && node.Depth == skipDepth {
				skipDepth = -1
				continue
			}
		}

		m.visible = append(m.visible, i)
		if isContainerOpen(node) && node.Collapsed {
			skipDepth = node.Depth
		}
	}
}

func isContainerOpen(node tree.Node) bool {
	return node.Type == tree.ObjectOpen || node.Type == tree.ArrayOpen
}

func isContainerClose(node tree.Node) bool {
	return node.Type == tree.ObjectClose || node.Type == tree.ArrayClose
}
