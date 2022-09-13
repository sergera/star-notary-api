package domain

import (
	"errors"
	"regexp"
)

type StarRangeModel struct {
	Start       string
	End         string
	OldestFirst bool
}

func (s StarRangeModel) Validate() error {
	errorMsg := "invalid range"

	pattern := "^[1-9](?:[0-9]+)?$"

	match, err := regexp.MatchString(pattern, s.Start)
	if err != nil {
		return err
	}

	if !match {
		return errors.New(errorMsg)
	}

	match, err = regexp.MatchString(pattern, s.End)
	if err != nil {
		return err
	}

	if !match {
		return errors.New(errorMsg)
	}

	return nil
}
