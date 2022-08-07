package main

import (
	"fmt"
	"log"
	"strings"
)

var exampleHtml = `
<a href="/dog">
  <span>Something in a span</span>
  Text not in a span
  <b>Bold text!</b>
</a>
`

func main() {
	ex := strings.NewReader(exampleHtml)
	links, err := link.Parse(ex)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(links)
}
