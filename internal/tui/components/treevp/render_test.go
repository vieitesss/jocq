package treevp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/vieitesss/jocq/internal/tree"
)

func TestRenderLineDoesNotWrapLongValues(t *testing.T) {
	node := tree.Node{
		Type:   tree.KeyValue,
		Depth:  1,
		Key:    "about",
		Value:  strings.Repeat("a", 120),
		IsLast: true,
	}

	line := RenderLine(node, false, 40)
	if strings.Contains(line, "\n") {
		t.Fatalf("expected single visual line, got wrapped output")
	}

	if got := ansi.StringWidth(line); got != 40 {
		t.Fatalf("expected width 40, got %d", got)
	}
}

func TestRenderLineCollapsedContainerSummary(t *testing.T) {
	node := tree.Node{
		Type:      tree.ObjectOpen,
		Depth:     1,
		Key:       "user",
		Collapsed: true,
		IsLast:    true,
	}

	line := ansi.Strip(RenderLine(node, false, 80))
	if !strings.Contains(line, "+ ") {
		t.Fatalf("expected collapsed marker in line, got %q", line)
	}

	if !strings.Contains(line, "{...}") {
		t.Fatalf("expected collapsed object summary, got %q", line)
	}
}

func TestRenderLineCollapsedArraySummary(t *testing.T) {
	node := tree.Node{
		Type:      tree.ArrayOpen,
		Depth:     1,
		Collapsed: true,
		IsLast:    true,
	}

	line := ansi.Strip(RenderLine(node, false, 80))
	if !strings.Contains(line, "[...]") {
		t.Fatalf("expected collapsed array summary, got %q", line)
	}
}
