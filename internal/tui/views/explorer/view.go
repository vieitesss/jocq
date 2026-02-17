package explorer

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/vieitesss/jocq/internal/tui/theme"
)

const viewportChromeHeight = 3

var (
	paneStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	paneStyleBlur   = paneStyle.BorderForeground(lipgloss.Color(theme.Gray))
	paneStyleFocus  = paneStyle.BorderForeground(lipgloss.Color(theme.Pink))
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Gray))
	titleStyleFocus = titleStyle.Foreground(lipgloss.Color(theme.Pink))
	titleMetaStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GrayMuted))
	titleMetaFocus  = titleMetaStyle.Foreground(lipgloss.Color(theme.Pink))
)

func (e ExplorerModel) viewportHeight(height, inputHeight, helpHeight int) int {
	return max(0, height-inputHeight-viewportChromeHeight-helpHeight)
}

func (e ExplorerModel) ExplorerView() string {
	if !e.ready {
		return "Loading..."
	}

	inLabel := "Raw JSON [Tab]"
	outLabel := "Query Result [S+Tab]"
	inFocused := false
	outFocused := false
	inPaneStyle := paneStyleBlur
	outPaneStyle := paneStyleBlur

	switch e.focused {
	case InPane:
		inLabel = "Raw JSON"
		outLabel = "Query Result [Tab]"
		inFocused = true
		inPaneStyle = paneStyleFocus

	case OutPane:
		outLabel = "Query Result"
		inLabel = "Raw JSON [S+Tab]"
		outFocused = true
		outPaneStyle = paneStyleFocus
	}

	inCount, inHasCount := e.In.PendingCount()
	inTitle := renderPaneTitle(inLabel, sourcePaneMeta(inCount, inHasCount, e.In.CursorPercent()), e.In.Width, inFocused)
	outTitle := renderPaneTitle(outLabel, percentMeta(viewportPercent(e.Out)), e.Out.Width, outFocused)

	inPaneHeight := e.In.Height + viewportChromeHeight
	outPaneHeight := e.Out.Height + viewportChromeHeight

	paneIn := inPaneStyle.
		Width(e.In.Width).
		Height(e.In.Height + 1).
		MaxHeight(inPaneHeight).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			inTitle,
			e.In.View(),
		))
	paneOut := outPaneStyle.
		Width(e.Out.Width).
		Height(e.Out.Height + 1).
		MaxHeight(outPaneHeight).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			outTitle,
			fitContentWidth(e.Out.View(), e.Out.Width),
		))
	inputView := fitContentWidth(e.Input.View(), e.width)
	helpView := fitContentWidth(e.help.View(e.keys), e.width)

	v := lipgloss.JoinVertical(lipgloss.Top,
		inputView,
		lipgloss.JoinHorizontal(lipgloss.Top,
			paneIn,
			paneOut,
		),
		helpView,
	)

	return v
}

func renderPaneTitle(label, meta string, width int, focused bool) string {
	if width <= 0 {
		return ""
	}

	leftStyle := titleStyle
	rightStyle := titleMetaStyle
	if focused {
		leftStyle = titleStyleFocus
		rightStyle = titleMetaFocus
	}

	left := leftStyle.Render(label)
	right := rightStyle.Render(meta)

	return joinEdge(left, right, width)
}

func sourcePaneMeta(count int, hasCount bool, percent int) string {
	percentPart := percentMeta(percent)
	if !hasCount {
		return percentPart
	}

	return fmt.Sprintf("%d  •  %s", count, percentPart)
}

func percentMeta(percent int) string {
	return fmt.Sprintf("%3d%%", clampPercent(percent))
}

func joinEdge(left, right string, width int) string {
	if width <= 0 {
		return ""
	}

	leftWidth := ansi.StringWidth(left)
	rightWidth := ansi.StringWidth(right)

	if leftWidth+1+rightWidth <= width {
		return left + strings.Repeat(" ", width-leftWidth-rightWidth) + right
	}

	if rightWidth >= width {
		return ansi.Truncate(right, width, "")
	}

	leftBudget := max(0, width-rightWidth-1)
	left = ansi.Truncate(left, leftBudget, "")
	padding := max(1, width-ansi.StringWidth(left)-rightWidth)

	return left + strings.Repeat(" ", padding) + right
}

func viewportPercent(v viewport.Model) int {
	totalLines := v.TotalLineCount()
	if totalLines == 0 {
		return 0
	}

	if totalLines <= max(1, v.Height) {
		return 100
	}

	scrollPercent := v.ScrollPercent()
	if math.IsNaN(scrollPercent) {
		return 0
	}

	return clampPercent(int(math.Round(scrollPercent * 100)))
}

func clampPercent(percent int) int {
	if percent < 0 {
		return 0
	}

	if percent > 100 {
		return 100
	}

	return percent
}

func fitContentWidth(content string, width int) string {
	if width <= 0 || content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		clipped := ansi.Truncate(line, width, "")
		if padding := width - ansi.StringWidth(clipped); padding > 0 {
			clipped += strings.Repeat(" ", padding)
		}
		lines[i] = clipped
	}

	return strings.Join(lines, "\n")
}
