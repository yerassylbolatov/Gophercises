package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"strings"
)

type link struct {
	Href string
	Text string
}

func main() {
	filename := flag.String("file", "ex4.html", "Specify the name of the html document for parsing")
	flag.Parse()
	read, err := readHtml(*filename)
	if err != nil {
		log.Fatal(err)
	}
	parsedHtml, err := parse(read)
	fmt.Println(parsedHtml)
}

func readHtml(filename string) (string, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(f), nil
}

func parse(s string) ([]link, error) {
	r := strings.NewReader(s)
	var l []link
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					text := extractText(n)
					l = append(l, link{"Href:" + a.Val, "Text:" + editText(text)})
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return l, nil
}

func editText(text string) string {
	ret := strings.Join(strings.Fields(text), " ")
	return ret
}

func extractText(n *html.Node) string {
	text := ""
	if n.Type == html.TextNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return text
}
