package main

import (
	"fmt"
	"link"
	"log"
	"strings"
)

var exampleHtml = `
<html>
<body>
<a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>
`

func main() {
	ex := strings.NewReader(exampleHtml)
	links, err := link.Parse(ex)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(links)
}
