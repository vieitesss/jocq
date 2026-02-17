package treevp

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/tree"
)

type Model struct {
	nodes  []tree.Node
	cursor int
	offset int

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

func (m *Model) SetNodes(nodes []tree.Node) {
	m.nodes = nodes
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
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	key := keyMsg.String()
	if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
		m.pushCountDigit(int(key[0] - '0'))
		return m, nil
	}

	if key != "j" && key != "k" {
		m.clearCount()
	}

	switch key {
	case "up", "k":
		m.CursorUpN(m.consumeCount(1, key == "k"))

	case "down", "j":
		m.CursorDownN(m.consumeCount(1, key == "j"))

	case "ctrl+u":
		m.PageUp()

	case "ctrl+d":
		m.PageDown()

	case "g", "home":
		m.GoToTop()

	case "G", "end":
		m.GoToBottom()
	}

	return m, nil
}

func (m *Model) CursorUp() {
	m.CursorUpN(1)
}

func (m *Model) CursorUpN(n int) {
	if len(m.nodes) == 0 {
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
	if len(m.nodes) == 0 {
		return
	}

	if n <= 0 {
		return
	}

	m.cursor = min(len(m.nodes)-1, m.cursor+n)
	m.ensureCursorVisible()
}

func (m *Model) PageUp() {
	if len(m.nodes) == 0 {
		return
	}

	delta := max(1, m.Height/2)
	m.cursor = max(0, m.cursor-delta)
	m.ensureCursorVisible()
}

func (m *Model) PageDown() {
	if len(m.nodes) == 0 {
		return
	}

	delta := max(1, m.Height/2)
	m.cursor = min(len(m.nodes)-1, m.cursor+delta)
	m.ensureCursorVisible()
}

func (m *Model) GoToTop() {
	if len(m.nodes) == 0 {
		return
	}

	m.cursor = 0
	m.ensureCursorVisible()
}

func (m *Model) GoToBottom() {
	if len(m.nodes) == 0 {
		return
	}

	m.cursor = len(m.nodes) - 1
	m.ensureCursorVisible()
}

func (m Model) CursorNode() (tree.Node, bool) {
	if len(m.nodes) == 0 {
		return tree.Node{}, false
	}

	return m.nodes[m.cursor], true
}

func (m Model) CursorPercent() int {
	if len(m.nodes) == 0 {
		return 0
	}

	if len(m.nodes) == 1 {
		return 100
	}

	percent := int(math.Round(float64(m.cursor) * 100 / float64(len(m.nodes)-1)))
	return max(0, min(100, percent))
}

func (m Model) PendingCount() (int, bool) {
	if !m.hasCount {
		return 0, false
	}

	return m.pendingCount, true
}

func (m *Model) pushCountDigit(digit int) {
	if !m.hasCount {
		m.hasCount = true
		m.pendingCount = 0
	}

	m.pendingCount = m.pendingCount*10 + digit
}

func (m *Model) consumeCount(fallback int, allow bool) int {
	if !allow || !m.hasCount {
		return fallback
	}

	count := m.pendingCount
	m.clearCount()
	return count
}

func (m *Model) clearCount() {
	m.pendingCount = 0
	m.hasCount = false
}

func (m *Model) ResetCount() {
	m.clearCount()
}

func (m *Model) ensureCursorVisible() {
	if len(m.nodes) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}

	m.cursor = max(0, min(len(m.nodes)-1, m.cursor))

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

	maxOffset := max(0, len(m.nodes)-m.Height)
	m.offset = max(0, min(maxOffset, m.offset))
}
