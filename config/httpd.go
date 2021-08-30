package config

import (
	"time"
)

type HTTPD struct {
	HttpEnabled              bool   `conf:"http.enabled"`
	HttpPort                 uint16 `conf:"http.port"`
	HttpsEnabled             bool   `conf:"https.enabled"`
	HttpsPort                uint16 `conf:"https.port"`
	CACertificate            string `conf:"tls.ca"`
	TLSCertificate           string `conf:"tls.certificate"`
	TLSKey                   string `conf:"tls.key"`
	RequireClientCertificate bool   `conf:"tls.client.certificates.required"`
	Security                 struct {
		Auth          string        `conf:"auth"`
		AuthDB        string        `conf:"local.db"`
		CookieMaxAge  int           `conf:"cookie.max-age"`
		LoginExpiry   string        `conf:"login.expiry"`
		SessionExpiry string        `conf:"session.expiry"`
		StaleTime     time.Duration `conf:"stale-time"`
	} `conf:"security"`
	RequestTimeout time.Duration `conf:"request.timeout"`
	System         struct {
		Controllers string        `conf:"controllers"`
		Doors       string        `conf:"doors"`
		Groups      string        `conf:"groups"`
		Cards       string        `conf:"cards"`
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
			ACL    string `conf:"acl"`
			System string `conf:"system"`
			Cards  string `conf:"cards"`
			Doors  string `conf:"doors"`
		} `conf:"rules"`
	} `conf:"db"`
	Audit struct {
		File string `conf:"file"`
	} `conf:"audit"`
	Retention time.Duration `conf:"retention"`
}

func NewHTTPD() *HTTPD {
	return &HTTPD{
		HttpEnabled:              false,
		HttpsEnabled:             true,
		CACertificate:            httpdCACertificate,
		TLSCertificate:           httpdTLSCertificate,
		TLSKey:                   httpdTLSKey,
		RequireClientCertificate: false,
		Security: struct {
			Auth          string        `conf:"auth"`
			AuthDB        string        `conf:"local.db"`
			CookieMaxAge  int           `conf:"cookie.max-age"`
			LoginExpiry   string        `conf:"login.expiry"`
			SessionExpiry string        `conf:"session.expiry"`
			StaleTime     time.Duration `conf:"stale-time"`
		}{
			Auth:          "basic",
			AuthDB:        httpdAuthDB,
			CookieMaxAge:  24,
			LoginExpiry:   "5m",
			SessionExpiry: "60m",
			StaleTime:     6 * time.Hour,
		},
		RequestTimeout: 5 * time.Second,
		System: struct {
			Controllers string        `conf:"controllers"`
			Doors       string        `conf:"doors"`
			Groups      string        `conf:"groups"`
			Cards       string        `conf:"cards"`
			Refresh     time.Duration `conf:"refresh"`
			Windows     struct {
				Ok          time.Duration `conf:"ok"`
				Uncertain   time.Duration `conf:"uncertain"`
				Systime     time.Duration `conf:"systime"`
				CacheExpiry time.Duration `conf:"expires"`
			} `conf:"windows"`
		}{
			Controllers: httpdControllersFile,
			Doors:       httpdDoorsFile,
			Groups:      httpdGroupsFile,
			Cards:       httpdCardsFile,
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
				ACL    string `conf:"acl"`
				System string `conf:"system"`
				Cards  string `conf:"cards"`
				Doors  string `conf:"doors"`
			} `conf:"rules"`
		}{
			Rules: struct {
				ACL    string `conf:"acl"`
				System string `conf:"system"`
				Cards  string `conf:"cards"`
				Doors  string `conf:"doors"`
			}{
				ACL:    httpdDBACLRules,
				System: httpdDBSystemRules,
				Cards:  httpdDBCardRules,
				Doors:  httpdDBDoorRules,
			},
		},
		Audit: struct {
			File string `conf:"file"`
		}{
			File: httpdAuditFile,
		},
		Retention: 6 * time.Hour,
	}
}
