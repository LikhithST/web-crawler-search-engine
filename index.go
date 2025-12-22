package main

import (
	"strings"
	"sync"
)

// ThreadSafeIndex wraps the map with a mutex
type ThreadSafeIndex struct {
	store map[string][]string
	mu    sync.RWMutex // Read/Write Mutex
}

func NewIndex() *ThreadSafeIndex {
	return &ThreadSafeIndex{
		store: make(map[string][]string),
	}
}

func (idx *ThreadSafeIndex) Add(url string, content string) {
	// 1. Tokenize content (Split by spaces)
	words := strings.Fields(content)

	// LOCK (Write Lock)
	idx.mu.Lock()
	defer idx.mu.Unlock() // Unlock automatically when function finishes

	// 2. Add to map
	for _, word := range words {
		word = strings.ToLower(word)
		// Basic cleanup: remove punctuation (simple version)
		word = strings.Trim(word, ".,!?;:\"()")

		if len(word) < 3 {
			continue
		} // Skip tiny words like "is", "at"

		// Check if we already have this URL for this word to avoid duplicates
		// (A Set would be better here, but a slice check is fine for MVP)
		found := false
		for _, existingURL := range idx.store[word] {
			if existingURL == url {
				found = true
				break
			}
		}

		if !found {
			idx.store[word] = append(idx.store[word], url)
		}
	}
}

func (idx *ThreadSafeIndex) Search(query string) []string {
	// R-LOCK (Read Lock) - Multiple readers allowed at once!
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return idx.store[query]
}
