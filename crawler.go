package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// This function loops FOREVER, waiting for work.
func worker(id int, jobs <-chan string, results chan<- PageData) {
	client := &http.Client{
		Timeout: 10 * time.Second, // BEST PRACTICE: Never wait forever
	}

	for url := range jobs {
		fmt.Printf("Worker %d processing %s\n", id, url)

		// 1. Fetch the page
		req, _ := http.NewRequest("GET", url, nil)

		// Identify yourself!
		req.Header.Set("User-Agent", "GoSpider-Bot/1.0 (+https://github.com/likhithst/web-crawler-search-engine)")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Error:", err)
			results <- PageData{} // Send empty result on error
			continue
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// 2. Extract links (Using your Phase 1 function!)
		// Note: We use the URL itself as base to resolve relative links
		links := ExtractLinks(url, bytes.NewReader(bodyBytes))

		text := ExtractText(bytes.NewReader(bodyBytes))

		// 3. Send results back
		results <- PageData{
			URL:     url,
			Links:   links,
			Content: text,
		}

		// Optional: Be polite to websites
		time.Sleep(500 * time.Millisecond)
	}
}
