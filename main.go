package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	fs := custom404(http.Dir("./static"))

	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.HandleFunc("/list", getAnthologyList)
	mux.HandleFunc("/view", getFromAnthology)
	mux.HandleFunc("/convert", convertArticle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := recoverPanic(logRequest(rateLimiter(logRequest(mux))))

	log.Println("Serving on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
