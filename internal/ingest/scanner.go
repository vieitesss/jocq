package ingest

import (
	"fmt"
	"io"
	"os"

	"github.com/vieitesss/jocq/internal/buffer"
)

const SizeThreshold = 100 * 1000 * 1000 // 100 MB

type Scanner struct {
	input *os.File
	data  *buffer.Data
}

func (s Scanner) inputSize() (int64, error) {
	info, err := s.input.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func NewScanner(input *os.File, data *buffer.Data) Scanner {
	return Scanner{
		input: input,
		data:  data,
	}
}

// Reads input and writes into data
func (s Scanner) Scan() error {
	size, err := s.inputSize()
	if err != nil {
		return err
	}

	if size == 0 {
		return fmt.Errorf("file %s is empty", s.input.Name())
	}

	if size <= SizeThreshold {
		d := make([]byte, int(size))
		_, err = io.ReadFull(s.input, d)
		if err != nil {
			return err
		}
		s.data.Append(d)
	} else {
		return fmt.Errorf("file too big (%d > %d)", size, SizeThreshold)
	}

	return nil
}
