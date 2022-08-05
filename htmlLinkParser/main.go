package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

type link struct {
	Href string
	Text string
}

func main() {
	file, err := os.Open("ex1.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	a := parse(file)
	fmt.Println(a)
}

func parse(file *os.File) []link {
	var l []link
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					text := extractText(n)
					l = append(l, link{a.Val, text})
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return l
}

func extractText(n *html.Node) string {
	text := ""
	if n.Type != html.ElementNode && n.Data != "a" && n.Type != html.CommentNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return text
}
