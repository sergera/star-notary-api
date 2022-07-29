package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sergera/star-notary-backend/internal/conf"
	"github.com/sergera/star-notary-backend/internal/controllers"
)

func main() {

	conf.Setup()

	mux := http.NewServeMux()

	ctl := controllers.NewStarController()

	mux.HandleFunc("/create", ctl.CreateStar)

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
