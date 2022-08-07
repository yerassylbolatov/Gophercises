package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	urlFlag := flag.String("flag", "http://calhoun.io/", "Point the URL for a sitemap generation")
	flag.Parse()

	resp, err := http.Get(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	links, _ := link.Parse(resp.Body)
	fmt.Println(links)
}
