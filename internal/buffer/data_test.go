package buffer

import "testing"

func TestAppendStoresRawAndDecoded(t *testing.T) {
	d := NewData()
	raw := []byte(`{"a":1}`)

	if err := d.Append(raw); err != nil {
		t.Fatalf("append returned unexpected error: %v", err)
	}

	gotRaw := d.Raw()
	if len(gotRaw) != 1 {
		t.Fatalf("expected one raw item, got %d", len(gotRaw))
	}
	if string(gotRaw[0]) != string(raw) {
		t.Fatalf("unexpected raw payload: %q", string(gotRaw[0]))
	}

	gotDecoded := d.Decoded()
	if len(gotDecoded) != 1 {
		t.Fatalf("expected one decoded item, got %d", len(gotDecoded))
	}

	obj, ok := gotDecoded[0].(map[string]any)
	if !ok {
		t.Fatalf("expected decoded object map, got %T", gotDecoded[0])
	}

	value, ok := obj["a"].(float64)
	if !ok {
		t.Fatalf("expected decoded value to be float64, got %T", obj["a"])
	}
	if value != 1 {
		t.Fatalf("expected decoded value 1, got %v", value)
	}
}

func TestAppendInvalidJSONReturnsErrorAndDoesNotMutateData(t *testing.T) {
	d := NewData()
	if err := d.Append([]byte(`{"ok":true}`)); err != nil {
		t.Fatalf("seed append returned unexpected error: %v", err)
	}

	if err := d.Append([]byte(`{"broken"`)); err == nil {
		t.Fatalf("expected append to fail for invalid json")
	}

	if got := len(d.Raw()); got != 1 {
		t.Fatalf("expected raw length to remain 1, got %d", got)
	}

	if got := len(d.Decoded()); got != 1 {
		t.Fatalf("expected decoded length to remain 1, got %d", got)
	}
}
