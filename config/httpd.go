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
		Controllers string `conf:"controllers"`
		Doors       string `conf:"doors"`
	} `conf:"system"`
	DB struct {
		File  string `conf:"file"`
		Rules struct {
			ACL    string `conf:"acl"`
			System string `conf:"system"`
			Cards  string `conf:"cards"`
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
			Controllers string `conf:"controllers"`
			Doors       string `conf:"doors"`
		}{
			Controllers: httpdControllersFile,
			Doors:       httpdDoorsFile,
		},
		DB: struct {
			File  string `conf:"file"`
			Rules struct {
				ACL    string `conf:"acl"`
				System string `conf:"system"`
				Cards  string `conf:"cards"`
			} `conf:"rules"`
		}{
			File: httpdDBFile,
			Rules: struct {
				ACL    string `conf:"acl"`
				System string `conf:"system"`
				Cards  string `conf:"cards"`
			}{
				ACL:    httpdDBACLRules,
				System: httpdDBSystemRules,
				Cards:  httpdDBCardRules,
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
