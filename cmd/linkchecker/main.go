package main

import (
	"fmt"
	"linkchecker/pkg/checker"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}

	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch %s: %s\n", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	links := checker.ExtractLinks(resp.Body)

	linkStatuses := make([]checker.LinkStatus, len(links))

	for i, link := range links {
		linkStatuses[i] = checker.CheckLink(link)
	}

	for _, status := range linkStatuses {
		fmt.Printf("Link: %s Status: %s\n", status.Link, status.Status)
	}
}
