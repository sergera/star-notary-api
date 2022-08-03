package models

import (
	"errors"
	"regexp"
	"strings"
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

func (s StarModel) ValidateTokenId() error {
	pattern := "^[1-9](?:[0-9]+)?$"

	match, err := regexp.MatchString(pattern, s.TokenId)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New("invalid token id")
}

func (s StarModel) ValidateOwner() error {
	pattern := "^0x[0-9a-zA-Z]{40}$"

	match, err := regexp.MatchString(pattern, s.Owner)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New("invalid owner address")
}

func (s StarModel) ValidateCoordinates() error {
	pattern := "^(?:[0-1][0-9]|2[0-4])(?:[0-5][0-9]|60){2}\\.(?:[0-9][0-9])[+-](?:[0-8][0-9]|90)(?:[0-5][0-9]|60){2}\\.(?:[0-9][0-9])$"

	match, err := regexp.MatchString(pattern, s.Coordinates)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New("invalid coordinates")
}

func (s StarModel) ValidateName() error {
	errorMsg := "invalid name"

	nameLength := len(s.Name)
	if nameLength > 32 || nameLength < 4 {
		return errors.New(errorMsg)
	}

	pattern := "^(?:[a-z]+)(?: [a-z]+)*$"

	match, err := regexp.MatchString(pattern, s.Name)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New(errorMsg)
}

func (s StarModel) ValidatePrice() error {
	errorMsg := "invalid price"

	integerPattern := "^[1-9][0-9]*$"
	fractionPattern := "^[0-9]*[1-9]$"

	priceSlice := strings.Split(s.Price, ".")

	switch len(priceSlice) {
	case 1:
		if len(priceSlice[0]) > 12 {
			return errors.New(errorMsg)
		}

		match, err := regexp.MatchString(integerPattern, priceSlice[0])
		if err != nil {
			return err
		}

		if match {
			return nil
		}

		return errors.New(errorMsg)
	case 2:
		if len(priceSlice[0]) > 12 {
			return errors.New(errorMsg)
		}

		if len(priceSlice[1]) > 18 {
			return errors.New(errorMsg)
		}

		match, err := regexp.MatchString(integerPattern, priceSlice[0])
		if err != nil {
			return err
		}

		if !match {
			return errors.New(errorMsg)
		}

		match, err = regexp.MatchString(fractionPattern, priceSlice[1])
		if err != nil {
			return err
		}

		if !match {
			return errors.New(errorMsg)
		}

		return nil
	default:
		return errors.New(errorMsg)
	}
}
