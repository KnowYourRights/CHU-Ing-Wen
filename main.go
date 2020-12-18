package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))

	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.HandleFunc("/list", getAnthologyList)
	mux.HandleFunc("/view", getFromAnthology)
	mux.HandleFunc("/convert", convertArticle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Serving on port", port)
	log.Fatal(http.ListenAndServe(":"+port, logRequest(rateLimiter(mux))))
}
