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

func TestExecuteEvictsLeastRecentlyUsedProgramWhenCapacityExceeded(t *testing.T) {
	resetProgramCacheWithCapacity(2)
	t.Cleanup(resetProgramCache)

	data := []any{map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}}

	if _, err := Execute(".a", data); err != nil {
		t.Fatalf("unexpected error executing .a: %v", err)
	}
	if _, err := Execute(".b", data); err != nil {
		t.Fatalf("unexpected error executing .b: %v", err)
	}
	if _, err := Execute(".c", data); err != nil {
		t.Fatalf("unexpected error executing .c: %v", err)
	}

	if _, ok := getCachedProgram(".a"); ok {
		t.Fatalf("expected .a to be evicted as least recently used entry")
	}

	if _, ok := getCachedProgram(".b"); !ok {
		t.Fatalf("expected .b to remain cached")
	}
	if _, ok := getCachedProgram(".c"); !ok {
		t.Fatalf("expected .c to remain cached")
	}
}

func TestExecuteCacheHitRefreshesRecency(t *testing.T) {
	resetProgramCacheWithCapacity(2)
	t.Cleanup(resetProgramCache)

	data := []any{map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}}

	if _, err := Execute(".a", data); err != nil {
		t.Fatalf("unexpected error executing .a: %v", err)
	}
	if _, err := Execute(".b", data); err != nil {
		t.Fatalf("unexpected error executing .b: %v", err)
	}
	// Refresh .a so .b becomes the least recently used entry.
	if _, err := Execute(".a", data); err != nil {
		t.Fatalf("unexpected error executing .a (refresh): %v", err)
	}
	if _, err := Execute(".c", data); err != nil {
		t.Fatalf("unexpected error executing .c: %v", err)
	}

	if _, ok := getCachedProgram(".b"); ok {
		t.Fatalf("expected .b to be evicted after .a recency refresh")
	}
	if _, ok := getCachedProgram(".a"); !ok {
		t.Fatalf("expected .a to stay cached after recency refresh")
	}
	if _, ok := getCachedProgram(".c"); !ok {
		t.Fatalf("expected .c to be cached")
	}
}

func resetProgramCache() {
	programCacheMu.Lock()
	programCache = newLRUProgramCache(defaultProgramCacheCapacity)
	programCacheMu.Unlock()
}

func resetProgramCacheWithCapacity(capacity int) {
	programCacheMu.Lock()
	programCache = newLRUProgramCache(capacity)
	programCacheMu.Unlock()
}

func programCacheSize() int {
	programCacheMu.Lock()
	defer programCacheMu.Unlock()
	return programCache.len()
}

func getCachedProgram(query string) (cachedProgram, bool) {
	programCacheMu.Lock()
	defer programCacheMu.Unlock()

	elem, ok := programCache.items[query]
	if !ok {
		return cachedProgram{}, false
	}

	item := elem.Value.(cacheItem)
	return item.value, true
}
