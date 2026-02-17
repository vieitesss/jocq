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
