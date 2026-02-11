package ingest

import (
	"fmt"
	"io"
	"os"

	"github.com/vieitesss/jocq/internal/buffer"
)

const SIZE_THRESHOLD = 100 * 1000 * 1000 // 100 MB

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
		return fmt.Errorf("The file %s is empty", s.input.Name())
	}

	if size <= SIZE_THRESHOLD {
		d := make([]byte, int(size))
		_, err = io.ReadFull(s.input, d)
		if err != nil {
			return err
		}
		s.data.Append(d)
	} else {
		return fmt.Errorf("File too big! (%d > %d)", size, SIZE_THRESHOLD)
	}

	return nil
}
