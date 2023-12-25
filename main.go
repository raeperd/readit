package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
)

//go:embed testdata/techblog.woowahan.com.html
var seed string

func main() {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(seed))
	if err != nil {
		log.Fatal(err)
	}
	urls := make([]string, 0)
	document.Find("div.post-item > a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			urls = append(urls, href)
		}
	})

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
		Title:    "div.post-header > h1",
		Contents: "div.post-content-inner",
	}
	articles := make([]Article, len(htmls))
	for i, body := range htmls {
		document, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}
		article := Article{
			Title:    document.Find(articleQuery.Title).Text(),
			Contents: document.Find(articleQuery.Contents).Text(),
		}
		articles[i] = article
	}

	for _, article := range articles {
		log.Printf("title: %v contents: %v", article.Title, article.Contents)
	}
}

type ArticleQuery struct {
	Title    string
	Contents string
}

type Article struct {
	Title    string
	Contents string
}
