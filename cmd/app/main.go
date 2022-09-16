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

	conf := conf.ConfSingleton()

	mux := http.NewServeMux()

	starAPI := api.NewStarAPI()

	mux.HandleFunc("/star-range", api.CorsHandler(starAPI.GetStarRange))
	mux.HandleFunc("/create", starAPI.CreateStar)
	mux.HandleFunc("/set-name", starAPI.SetName)
	mux.HandleFunc("/set-price", starAPI.SetPrice)
	mux.HandleFunc("/remove-from-sale", starAPI.RemoveFromSale)
	mux.HandleFunc("/purchase", starAPI.Purchase)

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
