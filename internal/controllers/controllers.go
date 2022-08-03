package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/sergera/star-notary-backend/internal/conf"
	"github.com/sergera/star-notary-backend/internal/models"
	"github.com/sergera/star-notary-backend/internal/repositories"
)

type StarController struct {
	repo *repositories.StarRepository
}

func NewStarController() *StarController {
	return &StarController{
		repositories.NewStarRepository(conf.DBHost, conf.DBPort, conf.DBName, conf.DBUser, conf.DBPassword, false),
	}
}

func (sc *StarController) CreateStar(w http.ResponseWriter, r *http.Request) {
	defer sc.repo.Close()
	sc.repo.Open()

	var m models.StarModel

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateOwner()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateTokenId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateCoordinates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ValidateName()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = sc.repo.InsertWalletIfAbsent(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sc.repo.Create(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
