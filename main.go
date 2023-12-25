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
	htmls := make([]string, len(urls))
	for i, url := range urls {
		var body string
		err := requests.URL(url).Client(&client).ToString(&body).Fetch(ctx)
		if err != nil {
			log.Fatal(err)
		}
		htmls[i] = body
	}

	articleQuery := ArticleQuery{
		Title:    *css.MustParse("div.post-header > h1"),
		Contents: *css.MustParse("div.post-content-inner"),
	}
	articles := make([]Article, len(htmls))
	for i, body := range htmls {
		node, err := html.Parse(strings.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}
		article := Article{}
		for _, element := range articleQuery.Title.Select(node) {
			article.Title = element.FirstChild.Data
		}
		for _, element := range articleQuery.Contents.Select(node) {
			var contents strings.Builder
			for c := element.FirstChild; c != nil; c = c.NextSibling {
				contents.WriteString(c.Data)
			}
			article.Contents = contents.String()
		}
		articles[i] = article
	}

	for _, article := range articles {
		log.Printf("title: %v contents: %v", article.Title, article.Contents)
	}
}

type ArticleQuery struct {
	Title    css.Selector
	Contents css.Selector
}

type Article struct {
	Title    string
	Contents string
}
