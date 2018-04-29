# Encoding-html
A library for decoding html into golang structs. Useful e.g. 
for making crawlers to interact with pages that does not have an actual api. 

## Installation
```
go get github.com/zlepper/encoding-html
```

# Examples
Getting the front page of hackernews:

```go
package main

import (
	"github.com/zlepper/encoding-html"
	"net/http"
	"log"
)

type Post struct {
	Title string `css:".title a"`
	Link string `css:".title a" extract:"attr" attr:"href"`
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

```