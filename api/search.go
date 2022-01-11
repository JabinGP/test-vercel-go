package handler

import (
	"embed"
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
	"github.com/blevesearch/bleve/v2"
)

//go:embed site.json
var f embed.FS

var posts = []Post{}
var mapping = bleve.NewIndexMapping()
var index bleve.Index

type Post struct {
        Slug    string
        Title   string
        Content string
        Excerpt string
}

func init() {
	err := initPosts()
	if err != nil {
		panic(err)
	}
	err = initIndex()
	if err != nil {
		panic(err)
	}
}


func initIndex() error {
	// "" for inmemory
        idx, err := bleve.New("", mapping)
        if err != nil {
                return err
        }
        for i, post := range posts {
                idx.Index(
                        fmt.Sprintf("%d", i),
                        post,
                )
                idx.SetInternal([]byte(strconv.Itoa(i)), []byte(post.Excerpt))
        }
	index = idx
	return nil
}


func initPosts() error {
        jsonBts, err := f.ReadFile("site.json")
        if err != nil {
                return err
        }
        err = json.Unmarshal(jsonBts, &posts)
        if err != nil {
                return err
        }
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	keywords := r.FormValue("keywords")
	query := bleve.NewMatchQuery(keywords)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"Slug", "Title", "Excerpt"}
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
        type Result = map[string]interface{}
        var res []Result
        for _, hit := range searchResults.Hits {
                res = append(res, Result{
                        "slug":    hit.Fields["Slug"],
                        "title":   hit.Fields["Title"],
                        "excerpt": hit.Fields["Excerpt"],
                })
        }
	resBts, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, string(resBts))
}
