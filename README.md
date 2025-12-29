# web-crawler-search-engine: Concurrent Web Crawler & Search Engine

**web-crawler-search-engine** is a high-performance, multi-threaded web crawler and search engine built entirely in Go. It demonstrates advanced concurrency patterns to crawl web pages in parallel, index their content, and provide an instant keyword search interface.

## 🚀 Features

* **Concurrent Worker Pool:** Uses Go Channels and Goroutines to crawl multiple pages simultaneously.
* **Custom HTML Tokenizer:** Efficiently streams and parses HTML to extract links and text without loading entire DOM trees into memory.
* **Thread-Safe Indexing:** Implements `sync.RWMutex` to safely manage concurrent writes to the inverted index.
* **Inverted Index Search:** Provides $O(1)$ lookup time for keyword searches.
* **Graceful Shutdown:** Captures `SIGINT` (Ctrl+C) signals to serialize and save the index to `index.json` before exiting.
* **Domain Guardrails:** Restricts crawling to a specific domain (e.g., `go.dev`) to prevent scope creep.

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **Concurrency:** Goroutines, Buffered Channels, `sync.WaitGroup`, `sync.RWMutex`
* **Parsing:** `golang.org/x/net/html` (Tokenizer API)
* **Data Storage:** In-Memory Map (persisted to JSON)

## 📦 Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/likhithst/web-crawler-search-engine.git
    cd web-crawler-search-engine
    ```

2.  **Install dependencies:**
    ```bash
    go mod init web-crawler-search-engine
    go get golang.org/x/net/html
    ```

## 🏃‍♂️ Usage

1.  **Run the Crawler:**
    ```bash
    go run .
    ```

2.  **Monitor Progress:**
    The terminal will show the crawler status as it visits pages:
    ```text
    Status: 5 pages visited, 12 pending...
    Status: 10 pages visited, 24 pending...
    ```

3.  **Search:**
    Once the crawl hits the page limit (default: 50) or the queue empties, the search prompt appears.
    ```text
    Crawling Complete!
    Indexed 4500 unique words.

    Search (> to exit): concurrency
    Found on pages:
    - [https://go.dev/doc/effective_go](https://go.dev/doc/effective_go)
    - [https://go.dev/tour/concurrency](https://go.dev/tour/concurrency)
    ```

4.  **Save & Exit:**
    Press `>` to exit, or `Ctrl+C` at any time to save the current index to `index.json`.

## 📂 Project Structure

* `main.go`: The entry point. Manages the orchestration, user interface, and signal handling.
* `crawler.go`: Contains the `worker` logic, HTTP request handling, and user-agent settings.
* `tokenizer.go`: Custom logic to parse HTML streams, extract clean text, and resolve relative URLs.
* `index.go` (or inside main): Defines the `ThreadSafeIndex` struct and thread-safe methods (`Add`, `Search`).
* `structs.go`: Defines data models like `PageData`.

## 🧠 Key Concepts Demonstrated

### 1. The Worker Pool Pattern
Instead of spawning a Goroutine for every single URL (which could exhaust system resources), web-crawler-search-engine uses a fixed pool of workers (e.g., 5-10) that consume from a shared `jobs` channel. This ensures predictable resource usage.

### 2. Race Condition Prevention
Since multiple workers try to write to the `index` map simultaneously, the application uses a `sync.RWMutex`.
* **Writes:** specific keywords obtain a `Lock()` (Exclusive).
* **Reads:** Searches use `RLock()` (Shared), allowing multiple users to search at once without blocking.

### 3. Stream Processing
Rather than using heavy DOM-parsing libraries, the `tokenizer.go` uses `io.Reader` streams. This minimizes memory allocation, allowing the crawler to handle large pages efficiently.

<!-- ## 🔮 Future Improvements

* [ ] Implement `robots.txt` parsing to respect site crawling policies.
* [ ] Add a persistent database (SQLite or PostgreSQL) instead of JSON.
* [ ] Implement a ranking algorithm (TF-IDF) for better search relevance.

---
*Built as a Capstone Project to master Go Concurrency patterns.* -->