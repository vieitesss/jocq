package buffer

import (
	"encoding/json"
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

// Add data both to the raw and decoded arrays
func (d *Data) Append(raw []byte) {
	// What will be added to "decoded"
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		// It's preferred to have the raw data, even if
		// the decoded data is not successfully shown.
		v = nil
	}

	d.mu.Lock()
	d.raw = append(d.raw, raw)
	d.decoded = append(d.decoded, v)
	d.bytes += uint64(len(raw))
	d.mu.Unlock()
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

// This is a READ-ONLY function.
// DO NOT MUTATE THE RETURNED SLICES.
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
