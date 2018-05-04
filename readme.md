[![Build Status](https://travis-ci.org/zlepper/encoding-html.svg?branch=master)](https://travis-ci.org/zlepper/encoding-html)

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

At the time of writing, that printed: 
```
{Posts:[{Title:The NetHack dev team is happy to announce the release of NetHack 3.6.1 Link:https://groups.google.com/forum/m/#!topic/rec.games.roguelike.nethack/XhcIrLlNzpA} {Title:Show HN: A fast, hopefully accurate, fuzzy matching library written in Go Link:https://github.com/sahilm/fuzzy} {Title:Larry Harvey, co-founder of Burning Man, has died Link:https://www.nytimes.com/2018/04/28/obituaries/larry-harvey-burning-man-festival-dead-at-70.html} {Title:Ask HN: My startup has basically failed. What now? Link:item?id=16949209} {Title:Kasparov versus the World Link:https://en.wikipedia.org/wiki/Kasparov_versus_the_World} {Title:Show HN: A proof-of-concept FoundationDB based network block device backend Link:https://github.com/dividuum/fdb-nbd} {Title:OpenEMR v5.0.1 Link:http://www.openhealthnews.com/content/openemr-community-releases-monumental-upgrade-their-open-source-ehr-update-ready-download} {Title:It’s Impossible to Prove Your Laptop Hasn’t Been Hacked Link:https://theintercept.com/2018/04/28/computer-malware-tampering/} {Title:HyperTools: A Python toolbox for gaining insights into high-dimensional data Link:http://hypertools.readthedocs.io/en/latest/#} {Title:Nintendo's secretive creative process Link:https://amp.theguardian.com/games/2018/apr/25/nintendo-interview-secret-innovation-lab-ideas-working} {Title:VoiceOps is hiring in SF to build AI for b2b voice data Link:https://voiceops.com/careers.html} {Title:Show HN: Generating fun Stack Exchange questions using Markov chains Link:https://se-simulator.lw1.at/} {Title:The myopia boom (2015) Link:https://www.nature.com/news/the-myopia-boom-1.17120} {Title:Seattle vacates hundreds of marijuana charges going back 30 years Link:https://www.theroot.com/seattle-vacates-hundreds-of-marijuana-possession-charge-1825622917} {Title:In theory, rocks from Oman could store hundreds of years of human CO2 emissions Link:https://www.nytimes.com/interactive/2018/04/26/climate/oman-rocks.html} {Title:The quadratic formula and low-precision arithmetic Link:https://www.johndcook.com/blog/2018/04/28/quadratic-formula/} {Title:Implementing and Understanding Type Classes (2014) Link:http://okmij.org/ftp/Computation/typeclass.html} {Title:Drawing with boids Link:https://miniatureape.github.io/boiddraw/} {Title:Lessons learned from a failing local mall Link:https://www.strongtowns.org/journal/2018/4/23/bon-ton-gone} {Title:French museum discovers half of its collection are fakes Link:https://www.telegraph.co.uk/news/2018/04/28/french-museum-discovers-half-collection-fakes/} {Title:World's oldest spider discovered in Australian outback Link:https://phys.org/news/2018-04-world-oldest-spider-australian-outback.html} {Title:Statement on Nature Machine Intelligence Link:https://openaccess.engineering.oregonstate.edu/home} {Title:The Wren Programming Language Link:https://github.com/munificent/wren} {Title:Facebook Warns Investors to Expect 'Additional Incidents' of User Data Abuse Link:https://www.siliconvalley.com/2018/04/27/facebook-got-an-earnings-boost-but-heres-the-fine-print/} {Title:Open3D: A Modern Library for 3D Data Processing Home Code Docs C++ API Link:http://www.open-3d.org/} {Title:A Layman’s Intro to Western Classical Music Link:https://quariety.com/2018/04/28/a-laymans-intro-to-western-classical-music/} {Title:EU agrees on total ban of bee-harming pesticides Link:https://www.theguardian.com/environment/2018/apr/27/eu-agrees-total-ban-on-bee-harming-pesticides} {Title:What it means to “disagree and commit” and how I do it (2016) Link:http://www.amazonianblog.com/2016/11/what-it-means-to-disagree-and-commit-and-how-i-do-it.html} {Title:Native Clojure with GraalVM Link:https://www.innoq.com/en/blog/native-clojure-and-graalvm/} {Title:Bulldoze the business school Link:https://www.theguardian.com/news/2018/apr/27/bulldoze-the-business-school}]}
```

## Tag options
Everything in encoding-html is specified using tags, the currently available tags are as follows:

### `css`
Specifies the css selector for finding the element. An element
will always be selected from using the parent fields element as root.
This allows for selecting in arrays

If the selector is not specified, then the field will be ignored.
If a selector matches multiple elements, and the field is not an array,
the first element will be used.


### `extract`
Specifies how to get the text to work on.
Valid options are `text` or `attr`. `text` will get all the inner text nodes of the html.
`attr` will get the value of an attribute. What attribute to fetch is specified
using the `attr` tag.

If `extract` is not specified, text will be selected.
If an unknown option is specified, an error will be returning from the decode call.

The extracted values will automatically be parsed into the requested type using the
`strconv.ParseFloat|Int|Bool|UInt()` function in the standard library. If the value
cannot be parsed, and no `default` value has been provided, the entire decode will
return an error.

#### `attr`
Specifies what attribute should be extracted from the matching html element.
If `extract:"attr"` is specified, and this tag is not, an error will be returned.
If the attribute does not exist on the element, the empty string `""` will be considered the value
of the attribute.

### `default`
Specifies a default value that should be set, provided the selected content
was a zero value, or that the actual content could no be converted into the
specified type.

If the default value cannot be converted, the entire parsing will fail and
return an error.