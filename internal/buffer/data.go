package buffer

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Data struct {
	mu      sync.RWMutex // guards concurrent access to Data fields; also supports future async ingestion
	raw     [][]byte
	decoded []any // raw to json data
	bytes   uint64
}

func NewData() *Data {
	return &Data{
		raw:     make([][]byte, 0),
		decoded: make([]any, 0),
	}
}

// Append stores raw bytes and decoded JSON value.
func (d *Data) Append(raw []byte) error {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	d.mu.Lock()
	d.raw = append(d.raw, raw)
	d.decoded = append(d.decoded, v)
	d.bytes += uint64(len(raw))
	d.mu.Unlock()

	return nil
}

// This is a READ-ONLY function.
// DO NOT MUTATE THE RETURNED SLICES.
func (d *Data) Raw() [][]byte {
	d.mu.RLock()
	defer d.mu.RUnlock()

	l := len(d.raw)
	snapshot := d.raw[:l:l]

	return snapshot
}

func (d *Data) Decoded() []any {
	d.mu.RLock()
	defer d.mu.RUnlock()

	l := len(d.decoded)
	snapshot := d.decoded[:l:l]

	return snapshot
}

// This is a READ-ONLY function.
// DO NOT MUTATE THE RETURNED SLICES.
func (d *Data) RawRange(from, to int) [][]byte {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if from < 0 || to < 0 || from > to {
		return nil
	}

	l := len(d.raw)

	if from > l {
		return nil
	}

	if to > l {
		to = l
	}

	snapshot := d.raw[from:to:to]

	return snapshot
}
