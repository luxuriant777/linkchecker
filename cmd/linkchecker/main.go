package main

import (
	"flag"
	"fmt"
	"linkchecker/pkg/checker"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
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
		urlArg = "http://" + urlArg
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

	var wg sync.WaitGroup
	wg.Add(1)

	go checker.CheckLinksRecursively(baseURL, baseURL, &wg)

	// Wait for all goroutines to finish
	wg.Wait()
}
