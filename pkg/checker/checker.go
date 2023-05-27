package checker

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
)

type LinkStatus struct {
	Link   string
	Status string
}

func CheckLink(link string) LinkStatus {
	// send a GET request to the link
	resp, err := http.Get(link)
	if err != nil || resp.StatusCode != 200 {
		return LinkStatus{Link: link, Status: "Broken"}
	} else {
		return LinkStatus{Link: link, Status: "OK"}
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
						links = append(links, a.Val)
					}
				}
			}
		}
	}
}
