package config

import (
	"time"
)

type WildApricot struct {
	HTTP struct {
		ClientTimeout time.Duration `conf:"client-timeout"`
		Retries       int           `conf:"retries"`
		RetryDelay    time.Duration `conf:"retry-delay"`
		PageSize      uint32        `conf:"page-size"`
		MaxPages      uint32        `conf:"max-pages"`
		PageDelay     time.Duration `conf:"page-delay"`
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
			PageSize      uint32        `conf:"page-size"`
			MaxPages      uint32        `conf:"max-pages"`
			PageDelay     time.Duration `conf:"page-delay"`
		}{
			ClientTimeout: 10 * time.Second,
			Retries:       3,
			RetryDelay:    5 * time.Second,
			PageSize:      100,
			MaxPages:      10,
			PageDelay:     100 * time.Millisecond,
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
