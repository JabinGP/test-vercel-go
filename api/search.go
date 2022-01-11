package handler

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
)

//go:embed site.json
var f embed.FS

var posts = []map[string]interface{}{}

func init() {
	data, _ := f.ReadFile("site.json")
	err := json.Unmarshal(data, &posts)
	if err != nil {
		panic(err)
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {

	for _, post := range posts {
		fmt.Fprintf(w, "<h1>%s</h1>", post["title"])
	}
}
