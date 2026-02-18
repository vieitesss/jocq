package ingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vieitesss/jocq/internal/buffer"
)

func TestScanValidJSONAppendsToBuffer(t *testing.T) {
	f := tempFileWithContent(t, `{"a":1}`)
	defer f.Close()

	data := buffer.NewData()
	s := NewScanner(f, data)

	if err := s.Scan(); err != nil {
		t.Fatalf("scan returned unexpected error: %v", err)
	}

	if got := len(data.Raw()); got != 1 {
		t.Fatalf("expected one raw entry, got %d", got)
	}
	if got := len(data.Decoded()); got != 1 {
		t.Fatalf("expected one decoded entry, got %d", got)
	}
}

func TestScanInvalidJSONReturnsErrorAndDoesNotAppend(t *testing.T) {
	f := tempFileWithContent(t, `{"a":`)
	defer f.Close()

	data := buffer.NewData()
	s := NewScanner(f, data)

	err := s.Scan()
	if err == nil {
		t.Fatalf("expected scan error for invalid json")
	}

	if !strings.Contains(err.Error(), "invalid json") {
		t.Fatalf("expected invalid json error, got: %v", err)
	}

	if !strings.Contains(err.Error(), filepath.Base(f.Name())) {
		t.Fatalf("expected error to include input filename, got: %v", err)
	}

	if got := len(data.Raw()); got != 0 {
		t.Fatalf("expected no raw entries after failed scan, got %d", got)
	}
	if got := len(data.Decoded()); got != 0 {
		t.Fatalf("expected no decoded entries after failed scan, got %d", got)
	}
}

func tempFileWithContent(t *testing.T, content string) *os.File {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "scanner-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}

	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("seek temp file: %v", err)
	}

	return f
}
