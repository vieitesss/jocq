package query

import (
	"strings"
	"sync"

	"github.com/itchyny/gojq"
)

type cachedProgram struct {
	code *gojq.Code
	err  error
}

var (
	programCacheMu sync.RWMutex
	programCache   = make(map[string]cachedProgram)
)

func Execute(query string, data []any) ([]any, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return data, nil
	}

	code, err := getProgram(trimmed)
	if err != nil {
		return nil, err
	}

	result := make([]any, 0)
	for _, d := range data {
		iter := code.Run(d)
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

func getProgram(query string) (*gojq.Code, error) {
	programCacheMu.RLock()
	if cached, ok := programCache[query]; ok {
		programCacheMu.RUnlock()
		return cached.code, cached.err
	}
	programCacheMu.RUnlock()

	q, err := gojq.Parse(query)
	if err != nil {
		programCacheMu.Lock()
		if _, exists := programCache[query]; !exists {
			programCache[query] = cachedProgram{err: err}
		}
		programCacheMu.Unlock()
		return nil, err
	}

	code, err := gojq.Compile(q)
	programCacheMu.Lock()
	if cached, ok := programCache[query]; ok {
		programCacheMu.Unlock()
		return cached.code, cached.err
	}
	programCache[query] = cachedProgram{code: code, err: err}
	programCacheMu.Unlock()

	return code, err
}
