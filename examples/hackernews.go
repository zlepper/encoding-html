package main

import (
	"github.com/zlepper/encoding-html"
	"log"
	"net/http"
)

type Post struct {
	Title string `css:".title a"`
	Link  string `css:".title a" extract:"attr" attr:"href"`
}
type HN struct {
	Posts []Post `css:".itemlist .athing"`
}

func main() {
	resp, err := http.Get("https://news.ycombinator.com/")
	if err != nil {
		log.Fatal(err)
	}

	var hn HN
	err = html.NewDecoder(resp.Body).Decode(&hn)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", hn)
}
