package treevp

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vieitesss/jocq/internal/tui/theme"
)

var (
	relativeLineNumberStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GrayMuted))
	relativeLineNumberCursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Pink))
)

func (m Model) View() string {
	if m.Height <= 0 {
		return ""
	}

	gutterWidth := m.gutterWidth()
	separator := " "
	separatorWidth := 1
	if gutterWidth == 0 {
		separator = ""
		separatorWidth = 0
	}

	contentWidth := max(0, m.Width-gutterWidth-separatorWidth)

	lines := make([]string, 0, m.Height)

	if len(m.nodes) > 0 {
		start := max(0, min(m.offset, len(m.nodes)-1))
		end := min(len(m.nodes), start+m.Height)

		for i := start; i < end; i++ {
			content := RenderLine(m.nodes[i], i == m.cursor, contentWidth)
			lines = append(lines, m.renderLineWithGutter(i, content, gutterWidth, separator))
		}
	}

	for len(lines) < m.Height {
		lines = append(lines, strings.Repeat(" ", max(0, m.Width)))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderLineWithGutter(index int, content string, width int, separator string) string {
	if width <= 0 {
		return content
	}

	relative := index - m.cursor
	if relative < 0 {
		relative *= -1
	}

	label := strconv.Itoa(relative)
	if len(label) > width {
		label = label[len(label)-width:]
	} else if len(label) < width {
		label = strings.Repeat(" ", width-len(label)) + label
	}

	style := relativeLineNumberStyle
	if relative == 0 {
		style = relativeLineNumberCursorStyle
	}

	return style.Render(label) + separator + content
}

func (m Model) gutterWidth() int {
	if m.Width <= 1 || len(m.nodes) == 0 {
		return 0
	}

	width := digitCount(max(1, len(m.nodes)-1))
	maxAllowed := max(0, m.Width-1)
	return min(width, maxAllowed)
}

func digitCount(value int) int {
	if value < 0 {
		value *= -1
	}

	count := 1
	for value >= 10 {
		count++
		value /= 10
	}

	return count
}
