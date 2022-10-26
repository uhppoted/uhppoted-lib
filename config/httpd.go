package config

import (
	"time"
)

type HTTPD struct {
	HTML                     string `conf:"html"`
	HttpEnabled              bool   `conf:"http.enabled"`
	HttpPort                 uint16 `conf:"http.port"`
	HttpsEnabled             bool   `conf:"https.enabled"`
	HttpsPort                uint16 `conf:"https.port"`
	CACertificate            string `conf:"tls.ca"`
	TLSCertificate           string `conf:"tls.certificate"`
	TLSKey                   string `conf:"tls.key"`
	RequireClientCertificate bool   `conf:"tls.client.certificates.required"`
	Security                 struct {
		Auth          string `conf:"auth"`
		AuthDB        string `conf:"local.db"`
		CookieMaxAge  int    `conf:"cookie.max-age"`
		LoginExpiry   string `conf:"login.expiry"`
		SessionExpiry string `conf:"session.expiry"`
		OTP           struct {
			Issuer string `conf:"issuer"`
		} `conf:"otp"`
	} `conf:"security"`
	RequestTimeout time.Duration `conf:"request.timeout"`
	System         struct {
		Interfaces  string        `conf:"interfaces"`
		Controllers string        `conf:"controllers"`
		Doors       string        `conf:"doors"`
		Groups      string        `conf:"groups"`
		Cards       string        `conf:"cards"`
		Events      string        `conf:"events"`
		Logs        string        `conf:"logs"`
		Users       string        `conf:"users"`
		History     string        `conf:"history"`
		Refresh     time.Duration `conf:"refresh"`
		Windows     struct {
			Ok          time.Duration `conf:"ok"`
			Uncertain   time.Duration `conf:"uncertain"`
			Systime     time.Duration `conf:"systime"`
			CacheExpiry time.Duration `conf:"expires"`
		} `conf:"windows"`
	} `conf:"system"`
	DB struct {
		Rules struct {
			ACL         string `conf:"acl"`
			Interfaces  string `conf:"interfaces"`
			Controllers string `conf:"controllers"`
			Cards       string `conf:"cards"`
			Doors       string `conf:"doors"`
			Groups      string `conf:"groups"`
			Events      string `conf:"events"`
			Logs        string `conf:"logs"`
			Users       string `conf:"users"`
		} `conf:"rules"`
	} `conf:"db"`
	Audit struct {
		File string `conf:"file"`
	} `conf:"audit"`
	Retention time.Duration `conf:"retention"`
	Timezones string        `conf:"timezones"`
}

func NewHTTPD() *HTTPD {
	return &HTTPD{
		HTML:                     "",
		HttpEnabled:              false,
		HttpPort:                 8080,
		HttpsEnabled:             true,
		HttpsPort:                8443,
		CACertificate:            httpdCACertificate,
		TLSCertificate:           httpdTLSCertificate,
		TLSKey:                   httpdTLSKey,
		RequireClientCertificate: false,
		Security: struct {
			Auth          string `conf:"auth"`
			AuthDB        string `conf:"local.db"`
			CookieMaxAge  int    `conf:"cookie.max-age"`
			LoginExpiry   string `conf:"login.expiry"`
			SessionExpiry string `conf:"session.expiry"`
			OTP           struct {
				Issuer string `conf:"issuer"`
			} `conf:"otp"`
		}{
			Auth:          "basic",
			AuthDB:        httpdAuthDB,
			CookieMaxAge:  24,
			LoginExpiry:   "1m",
			SessionExpiry: "60m",
			OTP: struct {
				Issuer string `conf:"issuer"`
			}{
				Issuer: "uhppoted-httpd",
			},
		},
		RequestTimeout: 5 * time.Second,
		System: struct {
			Interfaces  string        `conf:"interfaces"`
			Controllers string        `conf:"controllers"`
			Doors       string        `conf:"doors"`
			Groups      string        `conf:"groups"`
			Cards       string        `conf:"cards"`
			Events      string        `conf:"events"`
			Logs        string        `conf:"logs"`
			Users       string        `conf:"users"`
			History     string        `conf:"history"`
			Refresh     time.Duration `conf:"refresh"`
			Windows     struct {
				Ok          time.Duration `conf:"ok"`
				Uncertain   time.Duration `conf:"uncertain"`
				Systime     time.Duration `conf:"systime"`
				CacheExpiry time.Duration `conf:"expires"`
			} `conf:"windows"`
		}{
			Interfaces:  httpdInterfacesFile,
			Controllers: httpdControllersFile,
			Doors:       httpdDoorsFile,
			Groups:      httpdGroupsFile,
			Cards:       httpdCardsFile,
			Events:      httpdEventsFile,
			Logs:        httpdLogsFile,
			Users:       httpdUsersFile,
			History:     httpdHistoryFile,
			Refresh:     30 * time.Second,
			Windows: struct {
				Ok          time.Duration `conf:"ok"`
				Uncertain   time.Duration `conf:"uncertain"`
				Systime     time.Duration `conf:"systime"`
				CacheExpiry time.Duration `conf:"expires"`
			}{
				Ok:          60 * time.Second,
				Uncertain:   300 * time.Second,
				Systime:     300 * time.Second,
				CacheExpiry: 120 * time.Second,
			},
		},
		DB: struct {
			Rules struct {
				ACL         string `conf:"acl"`
				Interfaces  string `conf:"interfaces"`
				Controllers string `conf:"controllers"`
				Cards       string `conf:"cards"`
				Doors       string `conf:"doors"`
				Groups      string `conf:"groups"`
				Events      string `conf:"events"`
				Logs        string `conf:"logs"`
				Users       string `conf:"users"`
			} `conf:"rules"`
		}{
			Rules: struct {
				ACL         string `conf:"acl"`
				Interfaces  string `conf:"interfaces"`
				Controllers string `conf:"controllers"`
				Cards       string `conf:"cards"`
				Doors       string `conf:"doors"`
				Groups      string `conf:"groups"`
				Events      string `conf:"events"`
				Logs        string `conf:"logs"`
				Users       string `conf:"users"`
			}{
				ACL:         httpdRulesACL,
				Interfaces:  httpdRulesInterfaces,
				Controllers: httpdRulesControllers,
				Cards:       httpdRulesCards,
				Doors:       httpdRulesDoors,
				Groups:      httpdRulesGroups,
				Events:      httpdRulesEvents,
				Logs:        httpdRulesLogs,
				Users:       httpdRulesUsers,
			},
		},
		Audit: struct {
			File string `conf:"file"`
		}{
			File: httpdAuditFile,
		},
		Retention: 6 * time.Hour,
		Timezones: "",
	}
}
