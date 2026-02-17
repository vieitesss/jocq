package treevp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestViewShowsRelativeLineNumbers(t *testing.T) {
	m := New(18, 3)
	m.SetNodes(testNodes(6))
	m.CursorDownN(2)

	view := m.View()
	lines := strings.Split(view, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	plain0 := ansi.Strip(lines[0])
	plain1 := ansi.Strip(lines[1])
	plain2 := ansi.Strip(lines[2])

	if !strings.HasPrefix(plain0, "2 ") {
		t.Fatalf("expected first visible line to start with relative 2, got %q", plain0)
	}

	if !strings.HasPrefix(plain1, "1 ") {
		t.Fatalf("expected second visible line to start with relative 1, got %q", plain1)
	}

	if !strings.HasPrefix(plain2, "0 ") {
		t.Fatalf("expected cursor line to start with relative 0, got %q", plain2)
	}

	for i, line := range lines {
		if got := ansi.StringWidth(line); got != 18 {
			t.Fatalf("expected line %d width 18, got %d", i, got)
		}
	}
}
