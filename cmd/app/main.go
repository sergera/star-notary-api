package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sergera/star-notary-backend/internal/api"
	"github.com/sergera/star-notary-backend/internal/conf"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	conf := conf.GetConf()

	mux := http.NewServeMux()

	starAPI := api.NewStarAPI()

	mux.HandleFunc("/create", starAPI.CreateStar)
	mux.HandleFunc("/star-range", api.CorsHandler(starAPI.GetStarRange))

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
