package main

import (
	"io"
	"net/url" // NEW: We need this to parse URLs
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks now accepts the baseURL to resolve relative paths
func ExtractLinks(baseURL string, body io.Reader) []string {
	var links []string

	// Parse the base URL once so we can reuse it
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil // If the base URL is broken, we can't do anything
	}

	z := html.NewTokenizer(body)

	for {
		tokenType := z.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := cleanLink(base, attr.Val)
						if link != "" {
							links = append(links, link)
						}
					}
				}
			}
		}
	}
}

// cleanLink now takes the base URL object to resolve paths
func cleanLink(base *url.URL, link string) string {
	link = strings.TrimSpace(link)
	if link == "" || strings.HasPrefix(link, "javascript:") || strings.HasPrefix(link, "#") {
		return ""
	}

	// Parse the link found in the HTML
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// The Magic: ResolveReference combines the Base URL with the relative link
	// Example: Base("https://go.dev/blog/") + Link("../about") = "https://go.dev/about"
	return base.ResolveReference(u).String()
}

// ExtractText gets the visible content from the HTML body
func ExtractText(body io.Reader) string {
	z := html.NewTokenizer(body)
	var contentBuilder strings.Builder

	for {
		tokenType := z.Next()
		switch tokenType {
		case html.ErrorToken:
			return contentBuilder.String()

		case html.TextToken:
			text := strings.TrimSpace(z.Token().Data)
			if len(text) > 0 {
				// Determine if we should ignore this text (css, js)
				// For a simple MVP, we accept everything.
				// In production, you'd check parent tags to avoid <script> content.
				contentBuilder.WriteString(text + " ")
			}
		}
	}
}
