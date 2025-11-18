package config

import (
	"time"
)

type WildApricot struct {
	HTTP struct {
		ClientTimeout time.Duration `conf:"client-timeout"`
		Retries       int           `conf:"retries"`
		RetryDelay    time.Duration `conf:"retry-delay"`
		PageSize      int           `conf:"page-size"`
		MaxPages      int           `conf:"max-pages"`
	} `conf:"http"`

	Fields struct {
		CardNumber string `conf:"card-number"`
		PIN        string `conf:"PIN"`
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
			PageSize      int           `conf:"page-size"`
			MaxPages      int           `conf:"max-pages"`
		}{
			ClientTimeout: 10 * time.Second,
			Retries:       3,
			RetryDelay:    5 * time.Second,
			PageSize:      100,
			MaxPages:      10,
		},

		Fields: struct {
			CardNumber string `conf:"card-number"`
			PIN        string `conf:"PIN"`
		}{
			CardNumber: "Card Number",
			PIN:        "PIN",
		},

		FacilityCode: "",
	}
}
