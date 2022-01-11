package main


import (
	"github.com/JabinGP/test-vercel-go/api"
	"net/http"
)

func main() {
	http.HandleFunc("/api/search", handler.Handler)
	err := http.ListenAndServe(":3022", nil)
	if err != nil {
		panic(err)
	}
}
