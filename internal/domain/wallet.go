package domain

import (
	"errors"
	"regexp"
)

type WalletModel struct {
	Address string `json:"address"`
	Id      string `json:"id"`
}

func (wallet WalletModel) Validate() error {
	pattern := "^0x[0-9a-zA-Z]{40}$"

	match, err := regexp.MatchString(pattern, wallet.Address)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New("invalid wallet address")
}
