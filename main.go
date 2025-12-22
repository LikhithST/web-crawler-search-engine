package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// CONFIGURATION
	const numWorkers = 10
	seedURL := "https://go.dev"
	const maxPages = 50                   // Stop after finding this many unique pages
	const targetDomain = "https://go.dev" // Only follow links starting with this

	// 1. Create Channels
	// Buffered channels prevent blocking if one side is slightly faster
	jobs := make(chan string, 100)
	results := make(chan PageData, 100)

	index := NewIndex()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run a goroutine to handle the save
	go func() {
		<-sigChan // Block until signal received
		fmt.Println("\nShutdown signal received. Saving index...")

		file, _ := os.Create("index.json")
		json.NewEncoder(file).Encode(index.store)

		fmt.Println("Index saved. Exiting.")
		os.Exit(0)
	}()

	// 2. Start Workers
	// We spawn 5 goroutines that will all listen to the SAME 'jobs' channel.
	// Go automatically distributes the work.
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// 3. Add the Seed
	jobs <- seedURL

	// 4. The Event Loop
	// We need to track how many items are "in flight" so we know when to stop.
	// Otherwise, the program runs forever or exits too early.
	activeRequests := 1

	// Using a dedicated visited map to avoid cycles (Phase 3 preview)
	visited := make(map[string]bool)
	visited[seedURL] = true

	for activeRequests > 0 {
		// We wait for a result from ANY worker
		pageData := <-results
		activeRequests-- // One worker just finished
		// If fetch failed, pageData.URL will be empty
		if pageData.URL == "" {
			continue
		}

		// 1. FEED THE BRAIN
		index.Add(pageData.URL, pageData.Content)

		// 3. Feed the Crawler (WITH GUARDRAILS)
		// Stop adding work if we hit our limit
		if len(visited) >= maxPages {
			continue
		}

		for _, link := range pageData.Links {
			if !visited[link] {
				if strings.HasPrefix(link, targetDomain) {
					visited[link] = true

					// Increment active count BEFORE sending to avoid race conditions
					activeRequests++

					// Send to worker
					// WARNING: This can block if the 'jobs' channel is full!
					// In production, run this in a goroutine: go func() { jobs <- link }()
					go func(l string) { jobs <- l }(link)
				}
			}
		}

		fmt.Printf("Queue Status: %d active workers\n", activeRequests)
	}

	// SEARCH INTERFACE
	fmt.Println("\nCrawling Complete!")
	fmt.Println("Indexed", len(index.store), "unique words.")

	// Simple REPL (Read-Eval-Print Loop)
	var query string
	for {
		fmt.Print("\nSearch (> to exit): ")
		fmt.Scanln(&query)
		if query == ">" {
			break
		}

		results, ok := index.store[strings.ToLower(query)]
		if !ok {
			fmt.Println("No results found.")
		} else {
			fmt.Println("Found on pages:")
			for _, url := range results {
				fmt.Println("-", url)
			}
		}
	}

}
