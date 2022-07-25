package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sergera/star-notary-backend/internal/conf"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Printf("got / request\n")
	io.WriteString(w, "Hello World!\n")
}

func main() {

	conf.Setup()

	mux := http.NewServeMux()

	mux.HandleFunc("/", getRoot)

	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mux,
	}

	fmt.Printf("Staring application on port %s", conf.Port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
