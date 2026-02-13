package query

import (
	"strings"

	"github.com/itchyny/gojq"
)

func Execute(query string, data []any) ([]any, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return []any{}, nil
	}

	q, err := gojq.Parse(trimmed)
	if err != nil {
		return nil, err
	}

	result := make([]any, 0)
	for _, d := range data {
		iter := q.Run(d)
		for {
			v, ok := iter.Next()
			if !ok {
				// No more objects
				break
			}
			if err, ok := v.(error); ok {
				if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
					// clean stop
					break
				}
				// real error
				return nil, err
			}
			result = append(result, v)
		}
	}

	return result, nil
}
