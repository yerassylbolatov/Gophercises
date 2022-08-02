package urlshort

import (
	"gopkg.in/yaml.v3"
	"net/http"
)

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if target, ok := pathToUrls[path]; ok {
			http.Redirect(w, r, target, http.StatusFound)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathToUrlMap := buildMap(parsedYaml)
	return MapHandler(pathToUrlMap, fallback), nil
}

type yamlMap struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(yml []byte) ([]yamlMap, error) {
	var pathUrl []yamlMap
	err := yaml.Unmarshal(yml, &pathUrl)
	if err != nil {
		return nil, err
	}
	return pathUrl, nil
}

func buildMap(parsedYaml []yamlMap) map[string]string {
	pathToUrls := make(map[string]string)
	for _, yamlEntry := range parsedYaml {
		pathToUrls[yamlEntry.Path] = yamlEntry.Url
	}
	return pathToUrls
}
