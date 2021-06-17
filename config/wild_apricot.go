package config

import (
	"time"
)

type WildApricot struct {
	HTTP struct {
		ClientTimeout time.Duration `conf:"client-timeout"`
		Retries       int           `conf:"retries"`
		RetryDelay    time.Duration `conf:"retry-delay"`
	} `conf:"http"`

	Fields struct {
		CardNumber string `conf:"card-number"`
	} `conf:"fields"`

	DisplayOrder struct {
		Groups string `conf:"groups"`
		Doors  string `conf:"doors"`
	} `conf:"display-order"`

	FacilityCode string `conf:"facility-code"`
}

func NewWildApricot() *WildApricot {
	return &WildApricot{
		HTTP: struct {
			ClientTimeout time.Duration `conf:"client-timeout"`
			Retries       int           `conf:"retries"`
			RetryDelay    time.Duration `conf:"retry-delay"`
		}{
			ClientTimeout: 10 * time.Second,
			Retries:       3,
			RetryDelay:    5 * time.Second,
		},

		Fields: struct {
			CardNumber string `conf:"card-number"`
		}{
			CardNumber: "Card Number",
		},

		FacilityCode: "",
	}
}
