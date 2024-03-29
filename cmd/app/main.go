package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sergera/star-notary-backend/internal/api"
	"github.com/sergera/star-notary-backend/internal/conf"
	"github.com/sergera/star-notary-backend/internal/notifier"
	"github.com/sergera/star-notary-backend/pkg/cors"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	conf, err := conf.ConfSingleton()
	var originPatterns []string
	var port string
	if err == nil {
		originPatterns = strings.Split(conf.CORSAllowedURLs, ",")
		port = conf.Port
	}

	mux := http.NewServeMux()
	starAPI := api.NewStarAPI()
	cors := cors.NewCors(
		originPatterns,
		[]cors.HTTPVerb{cors.Options, cors.Get},
	)

	mux.HandleFunc("/star-range", cors.WrapHandlerFunc(starAPI.GetStarRange))
	mux.HandleFunc("/create", starAPI.CreateStar)
	mux.HandleFunc("/set-name", starAPI.SetName)
	mux.HandleFunc("/set-price", starAPI.SetPrice)
	mux.HandleFunc("/remove-from-sale", starAPI.RemoveFromSale)
	mux.HandleFunc("/purchase", starAPI.Purchase)

	starNotifier := notifier.StarNotifierSingleton()
	mux.HandleFunc("/notify-stars", cors.WrapHandlerFunc(starNotifier.Subscribe))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Starting application on port %s", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
