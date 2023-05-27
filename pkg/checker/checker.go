package checker

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex
var visited = make(map[string]bool)

func CheckLink(link string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(link)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch %s: %s\n", link, err)
		return
	}
	defer resp.Body.Close()

	filename := filepath.Join("statuses", strconv.Itoa(resp.StatusCode)+".txt")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(link + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to file: %s\n", err)
	}
}

func ExtractLinks(body io.Reader) []string {
	z := html.NewTokenizer(body)
	var links []string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						// Skip URLs with "#"
						if strings.Contains(a.Val, "#") {
							continue
						}
						// Remove trailing slashes
						link := strings.TrimRight(a.Val, "/")
						links = append(links, link)
					}
				}
			}
		}
	}
}

func CheckLinksRecursively(baseURL *url.URL, currentURL *url.URL, wg *sync.WaitGroup) {
	defer wg.Done()

	mu.Lock()
	if visited[currentURL.String()] {
		mu.Unlock()
		return
	}

	resp, err := http.Get(currentURL.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch %s: %s\n", currentURL, err)
		mu.Unlock()
		return
	}
	defer resp.Body.Close()

	visited[currentURL.String()] = true
	mu.Unlock()

	wgLink := &sync.WaitGroup{}
	wgLink.Add(1)
	go CheckLink(currentURL.String(), wgLink)

	if resp.StatusCode == http.StatusOK {
		links := ExtractLinks(resp.Body)

		for _, link := range links {
			absoluteURL := resolveURL(baseURL, link)
			parsedURL, err := url.Parse(absoluteURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse URL: %s\n", err)
				continue
			}
			// Only visit links that have the same hostname as the base URL
			if parsedURL.Hostname() == baseURL.Hostname() {
				wg.Add(1)
				go CheckLinksRecursively(baseURL, parsedURL, wg)
			}
		}
	}

	wgLink.Wait()
}

func resolveURL(baseURL *url.URL, link string) string {
	relativeURL, err := url.Parse(link)
	if err != nil {
		return ""
	}
	absoluteURL := baseURL.ResolveReference(relativeURL)
	return absoluteURL.String()
}
