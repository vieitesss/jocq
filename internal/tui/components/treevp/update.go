package treevp

import tea "github.com/charmbracelet/bubbletea"

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	if isCountDigit(key) {
		m.pushCountDigit(int(key[0] - '0'))
		return m, nil
	}

	if !isLineMotion(key) {
		m.clearCount()
	}

	switch key {
	case "up", "k":
		m.CursorUpN(m.consumeCount(1))

	case "down", "j":
		m.CursorDownN(m.consumeCount(1))

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

func (m *Model) pushCountDigit(digit int) {
	if !m.hasCount {
		m.hasCount = true
		m.pendingCount = 0
	}

	m.pendingCount = m.pendingCount*10 + digit
}

func (m *Model) consumeCount(fallback int) int {
	if !m.hasCount {
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

func isCountDigit(key string) bool {
	return len(key) == 1 && key[0] >= '0' && key[0] <= '9'
}

func isLineMotion(key string) bool {
	return key == "j" || key == "k" || key == "up" || key == "down"
}
