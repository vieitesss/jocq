package treevp

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/vieitesss/jocq/internal/tree"
	"github.com/vieitesss/jocq/internal/tui/theme"
)

const indentSize = 2

var (
	lineStyle        = lipgloss.NewStyle()
	bracketStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.JSONBracket))
	keyStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.JSONKey))
	stringStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.JSONString))
	numberStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.JSONNumber))
	boolStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.JSONBool))
	nullStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.JSONNull))
	punctuationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GrayMuted))
)

func RenderLine(node tree.Node, isCursor bool, width int) string {
	renderer := lineRenderer{cursor: isCursor}

	var builder strings.Builder
	builder.Grow(max(32, width))
	builder.WriteString(renderer.plain(strings.Repeat(" ", max(0, node.Depth*indentSize))))

	switch node.Type {
	case tree.ObjectOpen:
		writeCollapseMarker(&builder, renderer, node.Collapsed)
		writeKeyPrefix(&builder, renderer, node.Key)
		if node.Collapsed {
			builder.WriteString(renderer.render(bracketStyle, "{...}"))
		} else {
			builder.WriteString(renderer.render(bracketStyle, "{"))
		}

	case tree.ObjectClose:
		builder.WriteString(renderer.render(bracketStyle, "}"))

	case tree.ArrayOpen:
		writeCollapseMarker(&builder, renderer, node.Collapsed)
		writeKeyPrefix(&builder, renderer, node.Key)
		if node.Collapsed {
			builder.WriteString(renderer.render(bracketStyle, "[...]"))
		} else {
			builder.WriteString(renderer.render(bracketStyle, "["))
		}

	case tree.ArrayClose:
		builder.WriteString(renderer.render(bracketStyle, "]"))

	case tree.KeyValue:
		writeKeyPrefix(&builder, renderer, node.Key)
		builder.WriteString(renderScalar(renderer, node.Value))

	case tree.ArrayElement:
		builder.WriteString(renderScalar(renderer, node.Value))
	}

	if !node.IsLast {
		builder.WriteString(renderer.render(punctuationStyle, ","))
	}

	line := builder.String()
	if width > 0 {
		line = ansi.Truncate(line, width, "")
		if padding := width - ansi.StringWidth(line); padding > 0 {
			line += renderer.plain(strings.Repeat(" ", padding))
		}
	}

	return line
}

type lineRenderer struct {
	cursor bool
}

func (r lineRenderer) render(style lipgloss.Style, value string) string {
	if r.cursor {
		style = style.Background(lipgloss.Color(theme.CursorLineBg))
	}

	return style.Render(value)
}

func (r lineRenderer) plain(value string) string {
	return r.render(lineStyle, value)
}

func writeKeyPrefix(builder *strings.Builder, renderer lineRenderer, key string) {
	if key == "" {
		return
	}

	builder.WriteString(renderer.render(keyStyle, jsonString(key)))
	builder.WriteString(renderer.render(punctuationStyle, ": "))
}

func writeCollapseMarker(builder *strings.Builder, renderer lineRenderer, collapsed bool) {
	if collapsed {
		builder.WriteString(renderer.render(punctuationStyle, "+ "))
		return
	}

	builder.WriteString(renderer.render(punctuationStyle, "- "))
}

func renderScalar(renderer lineRenderer, value any) string {
	switch value := value.(type) {
	case string:
		return renderer.render(stringStyle, jsonString(value))

	case bool:
		if value {
			return renderer.render(boolStyle, "true")
		}
		return renderer.render(boolStyle, "false")

	case nil:
		return renderer.render(nullStyle, "null")

	case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return renderer.render(numberStyle, jsonValue(value))

	default:
		return renderer.render(punctuationStyle, jsonValue(value))
	}
}

func jsonValue(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		return strconv.Quote("<invalid>")
	}

	return string(b)
}

func jsonString(value string) string {
	b, err := json.Marshal(value)
	if err != nil {
		return strconv.Quote(value)
	}

	return string(b)
}
