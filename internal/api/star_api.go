package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sergera/star-notary-backend/internal/conf"
	"github.com/sergera/star-notary-backend/internal/domain"
	"github.com/sergera/star-notary-backend/internal/repositories"
)

type StarAPI struct {
	conn       *repositories.DBConnection
	starRepo   *repositories.StarRepository
	walletRepo *repositories.WalletRepository
}

func NewStarAPI() *StarAPI {
	/* TODO: resolve database session lifecycle */
	conf := conf.GetConf()
	conn := repositories.NewDBConnection(conf.DBHost, conf.DBPort, conf.DBName, conf.DBUser, conf.DBPassword, false)
	conn.Open()
	starRepo := repositories.NewStarRepository(conn)
	walletRepo := repositories.NewWalletRepository(conn)
	return &StarAPI{conn, starRepo, walletRepo}
}

func (sc *StarAPI) CreateStar(w http.ResponseWriter, r *http.Request) {
	var e domain.Event

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var m = domain.StarModel{
		TokenId:     e.TokenId,
		Coordinates: e.Coordinates,
		Name:        e.Name,
		Price:       "0",
		IsForSale:   false,
		Date:        e.Date,
		Wallet: &domain.WalletModel{
			Address: e.Owner,
		},
	}

	if err := m.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := sc.walletRepo.InsertWalletIfAbsent(m.Wallet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sc.starRepo.CreateStar(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (sc *StarAPI) GetStarRange(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	oldestFirst := r.URL.Query().Get("oldest-first")

	if oldestFirst == "" {
		oldestFirst = "false"
	}

	oldestFirstBool, err := strconv.ParseBool(oldestFirst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := domain.StarRangeModel{
		Start:       start,
		End:         end,
		OldestFirst: oldestFirstBool,
	}

	if err := m.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stars, err := sc.starRepo.GetStarRange(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	starsInBytes, err := json.Marshal(stars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(starsInBytes)
}

func CorsHandler(h http.HandlerFunc) http.HandlerFunc {
	conf := conf.GetConf()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", conf.CORSAllowedURLs)
		if r.Method == "OPTIONS" {
			//handle preflight in here
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept")
			w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
		} else {
			h.ServeHTTP(w, r)
		}
	}
}
