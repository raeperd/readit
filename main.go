package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/ericchiang/css"
	"golang.org/x/net/html"
)

//go:embed testdata/techblog.woowahan.com.html
var seed string

func main() {
	query, err := css.Parse("div.post-item > a")
	if err != nil {
		log.Fatal(err)
	}
	node, err := html.Parse(strings.NewReader(seed))
	if err != nil {
		log.Fatal(err)
	}
	urls := make([]string, 0)
	for _, element := range query.Select(node) {
		for _, attr := range element.Attr {
			if attr.Key == "href" {
				urls = append(urls, attr.Val)
			}
		}
	}
	client := http.Client{}
	ctx := context.Background()
	for _, url := range urls {
		var body string
		err := requests.URL(url).Client(&client).ToString(&body).Fetch(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(body)
	}
}
