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
	"time"
)

type Link struct {
	URL    string
	Status int
}

type LinkCheckResult struct {
	Origin string
	Link   Link
}

const StatusesDir = "statuses"

var mu sync.Mutex
var visited = make(map[string]bool)
var linkCheckResults = make([]LinkCheckResult, 0)
var linkCheckBuffer = make([]LinkCheckResult, 0)
var BufferSize = 1
var brokenLinks = make(map[string]bool)
var ProcessedURLs = 0
var didLinkCheckingStart bool
var delayBetweenRequests = time.Second * 2

var client = &http.Client{
	Timeout: time.Second * 10,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// Adjust User-Agent for redirects
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, "+
			"like Gecko) Chrome/92.0.4515.131 Safari/537.36")
		return nil
	},
}

func StartLinkChecking() {
	didLinkCheckingStart = true
}

func DidLinkCheckingStart() bool {
	return didLinkCheckingStart
}

func IncrementProcessedURLs() {
	mu.Lock()
	defer mu.Unlock()
	ProcessedURLs++
}

func WriteResultsToFileBuffer() {
	mu.Lock()
	defer mu.Unlock()

	// Check if the directory exists, if not, try to create it
	if _, err := os.Stat(StatusesDir); os.IsNotExist(err) {
		err := os.MkdirAll(StatusesDir, 0755)
		if err != nil {
			fmt.Printf("Folder %s couldn't be created, please create it manually\n", StatusesDir)
			return
		}
	}

	for _, result := range linkCheckBuffer {
		filename := filepath.Join(StatusesDir, strconv.Itoa(result.Link.Status)+".txt")
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", err)
			continue
		}

		var output string
		if result.Link.Status == 404 {
			output = result.Origin + " -> " + result.Link.URL + "\n"
		} else {
			output = result.Link.URL + "\n"
		}

		if _, err := file.WriteString(output); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to file: %s\n", err)
		}
		file.Close()
	}

	// Empty the buffer
	linkCheckBuffer = make([]LinkCheckResult, 0)
}

func CheckLink(link string) Link {
	resp, err := client.Get(link) // using client.Get instead of http.Get
	if err != nil {
		return Link{
			URL:    link,
			Status: 0,
		}
	}
	defer resp.Body.Close()

	return Link{
		URL:    link,
		Status: resp.StatusCode,
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

func CheckLinksRecursively(baseURL *url.URL, currentURL *url.URL, origin string, wg *sync.WaitGroup) {
	defer wg.Done()

	mu.Lock()
	if visited[currentURL.String()] && !brokenLinks[currentURL.String()] {
		mu.Unlock()
		return
	}
	visited[currentURL.String()] = true
	mu.Unlock()

	fmt.Printf("Fetching %s\n", currentURL.String())

	resp, err := client.Get(currentURL.String()) // using client.Get instead of http.Get
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch %s: %s\n", currentURL, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Finished fetching %s\n", currentURL.String())

	status := resp.StatusCode

	// Create a new link check result
	result := LinkCheckResult{
		Origin: origin,
		Link:   Link{URL: currentURL.String(), Status: status},
	}

	mu.Lock()
	linkCheckResults = append(linkCheckResults, result)
	linkCheckBuffer = append(linkCheckBuffer, result)
	mu.Unlock()

	if len(linkCheckBuffer) >= BufferSize {
		WriteResultsToFileBuffer()
	}

	if status == 404 {
		mu.Lock()
		brokenLinks[currentURL.String()] = true
		mu.Unlock()
	}

	IncrementProcessedURLs()

	StartLinkChecking()

	if status == http.StatusOK {
		links := ExtractLinks(resp.Body)

		for _, link := range links {
			absoluteURL := resolveURL(baseURL, link)
			parsedURL, err := url.Parse(absoluteURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse URL: %s\n", err)
				continue
			}
			time.Sleep(delayBetweenRequests)
			if parsedURL.Hostname() == baseURL.Hostname() {
				wg.Add(1)
				go CheckLinksRecursively(baseURL, parsedURL, currentURL.String(), wg)
			}
			time.Sleep(delayBetweenRequests)
		}
	}
}

func resolveURL(baseURL *url.URL, link string) string {
	relativeURL, err := url.Parse(link)
	if err != nil {
		return ""
	}
	absoluteURL := baseURL.ResolveReference(relativeURL)
	return absoluteURL.String()
}

func ProcessLinkCheckResults() {
	uniqueLinks := make(map[string]bool)
	for _, result := range linkCheckResults {
		if !uniqueLinks[result.Link.URL] {
			result.Link.Status = CheckLink(result.Link.URL).Status
			uniqueLinks[result.Link.URL] = true
		}
	}

	// Write remaining results in the buffer
	if len(linkCheckBuffer) > 0 {
		WriteResultsToFileBuffer()
	}
}
