package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <section class="page">
        <h1>{{.Title}}</h1>
        {{range .Story}}
            <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
            <li><a href="/{{.Arc}}">{{.Text}}</a></li>
            {{end}}
        </ul>
    </section>>
	<style>
		h1 {
			text-align: center;
			position: relative;
			color: brown;
			font-family: Georgia;
		}
		body{
			font-family: Trebuchet MS;
		}
        .page {
            width: 80%;
            max-width: 500px;
            margin: auto;
            margin-top: 40px;
            margin-bottom: 40px;
            padding: 80px;
            background: #FFFCF6;
        }
    </style>
</body>
</html>`

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Story map[string]Chapter

func main() {
	jsonFile, err := os.Open("gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	story := parseStory(byteValue)

	h := NewHandler(story, nil)
	fmt.Println("Starting a server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}

func NewHandler(s Story, t *template.Template) http.Handler {
	if t == nil {
		t = tpl
	}
	return handler{s, t}
}

type handler struct {
	s Story
	t *template.Template
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func parseStory(storyBytes []byte) (story Story) {
	err := json.Unmarshal(storyBytes, &story)
	if err != nil {
		log.Fatal(err)
	}
	return
}
