package link

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	var l []Link
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					text := extractText(n)
					l = append(l, Link{"Href:" + a.Val, "Text:" + editText(text)})
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
