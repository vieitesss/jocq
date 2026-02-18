package query

import "testing"

func TestExecuteCachesCompiledProgram(t *testing.T) {
	resetProgramCache()

	data := []any{
		map[string]any{"a": 1.0},
	}

	got, err := Execute(".a", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected one result, got %d", len(got))
	}

	if size := programCacheSize(); size != 1 {
		t.Fatalf("expected one cached program, got %d", size)
	}

	cached, ok := getCachedProgram(".a")
	if !ok || cached.code == nil || cached.err != nil {
		t.Fatalf("expected compiled cached program, got %+v", cached)
	}

	got, err = Execute(".a", data)
	if err != nil {
		t.Fatalf("unexpected error on second execute: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected one result on second execute, got %d", len(got))
	}
	if size := programCacheSize(); size != 1 {
		t.Fatalf("expected cache size to stay at one, got %d", size)
	}
}

func TestExecuteCachesParseError(t *testing.T) {
	resetProgramCache()

	data := []any{map[string]any{"a": 1.0}}

	if _, err := Execute(".[", data); err == nil {
		t.Fatalf("expected parse error")
	}

	if size := programCacheSize(); size != 1 {
		t.Fatalf("expected one cached parse error, got %d", size)
	}

	cached, ok := getCachedProgram(".[")
	if !ok || cached.code != nil || cached.err == nil {
		t.Fatalf("expected cached error entry, got %+v", cached)
	}

	if _, err := Execute(".[", data); err == nil {
		t.Fatalf("expected parse error on second execute")
	}
	if size := programCacheSize(); size != 1 {
		t.Fatalf("expected cache size to stay at one, got %d", size)
	}
}

func resetProgramCache() {
	programCacheMu.Lock()
	programCache = make(map[string]cachedProgram)
	programCacheMu.Unlock()
}

func programCacheSize() int {
	programCacheMu.RLock()
	defer programCacheMu.RUnlock()
	return len(programCache)
}

func getCachedProgram(query string) (cachedProgram, bool) {
	programCacheMu.RLock()
	defer programCacheMu.RUnlock()
	cached, ok := programCache[query]
	return cached, ok
}
