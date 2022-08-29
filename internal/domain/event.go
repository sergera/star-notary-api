package domain

import "time"

type Event struct {
	TokenId     string    `json:"token_id"`
	Owner       string    `json:"owner"`
	Coordinates string    `json:"coordinates"`
	Name        string    `json:"name"`
	Price       string    `json:"price"`
	IsForSale   bool      `json:"is_for_sale"`
	Date        time.Time `json:"date"`
}
