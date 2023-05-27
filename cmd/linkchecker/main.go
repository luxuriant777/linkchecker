package main

import (
	"flag"
	"fmt"
	"linkchecker/pkg/checker"
	"net/http" // Add this for http.Get
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s <url>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nThe <url> argument must be a URL in the following formats:\n")
	fmt.Fprintf(os.Stderr, "\nexample.com\nwww.example.com\nhttp://example.com\nhttps://example.com\n"+
		"http://subdomain.example.com\nhttps://subdomain.example.com\n")
}

func main() {
	flag.Usage = printUsage

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	urlArg := flag.Arg(0)

	// If the URL does not start with "http://" or "https://", default to "http://"
	if !strings.HasPrefix(urlArg, "http://") && !strings.HasPrefix(urlArg, "https://") {
		urlArg = "https://" + urlArg
	}

	// Parse the URL and get the hostname
	baseURL, err := url.Parse(urlArg)
	if err != nil || baseURL.Hostname() == "" {
		fmt.Fprintf(os.Stderr, "Invalid URL: %s\n", urlArg)
		os.Exit(1)
	}

	// Ensure the hostname part of the URL is in FQDN format
	fqdnRe := regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !fqdnRe.MatchString(baseURL.Hostname()) {
		printUsage()
		os.Exit(1)
	}

	// Attempt to fetch the provided URL
	resp, err := http.Get(baseURL.String())
	if err != nil {
		// Print an error and exit if the fetch fails
		fmt.Fprintf(os.Stderr, "Failed to fetch %s: %s\n", baseURL.String(), err)
		os.Exit(1)
	}
	// Close the response body after checking
	resp.Body.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		var prevProcessedURLs int
		noChangeCount := 0 // Counter for no change in processed URLs
		for {
			select {
			case <-time.After(3 * time.Second):
				if checker.ProcessedURLs == prevProcessedURLs {
					noChangeCount++
					if noChangeCount >= 35 { // No change in processed URLs for 35 seconds
						fmt.Println("No new URLs found for 35 seconds. Exiting...")
						os.Exit(0)
					}
					fmt.Println("No new URLs have been found. The parsing process will continue...")
				} else {
					fmt.Printf("Processed URLs: %d\n", checker.ProcessedURLs)
					noChangeCount = 0 // Reset the counter if there's a change in processed URLs
				}
				prevProcessedURLs = checker.ProcessedURLs
			}
		}
	}()

	go checker.CheckLinksRecursively(baseURL, baseURL, baseURL.String(), &wg)

	// Wait for all goroutines to finish
	wg.Wait()

	checker.ProcessLinkCheckResults()

	if checker.DidLinkCheckingStart() {
		fmt.Println("Link checking completed.\n" +
			"Results can be found in the 'statuses' folder.\n" +
			"Each file contains the list of URLs returned with the corresponding status code.\n" +
			"The '404.txt' follows the more detailed format:\n" +
			"<URL of the page where broken links were found> -> <the broken link>")
	}
}
