package models

import (
	"time"
)

type StarModel struct {
	TokenId     string
	Owner       string
	Coordinates string
	Name        string
	Price       string
	Date        time.Time
}
