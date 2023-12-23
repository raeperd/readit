package main

import (
	_ "embed"
	"log"
	"strings"

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
	for _, element := range query.Select(node) {
		for _, attr := range element.Attr {
			if attr.Key == "href" {
				log.Println(attr.Val)
			}
		}
	}
}
