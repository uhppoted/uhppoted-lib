package config

import ()

type REST struct {
	HttpEnabled               bool   `conf:"http.enabled"`
	HttpPort                  uint16 `conf:"http.port"`
	HttpsEnabled              bool   `conf:"https.enabled"`
	HttpsPort                 uint16 `conf:"https.port"`
	TLSKeyFile                string `conf:"tls.key"`
	TLSCertificateFile        string `conf:"tls.certificate"`
	CACertificateFile         string `conf:"tls.ca"`
	RequireClientCertificates bool   `conf:"tls.client.certificates"`
	CORSEnabled               bool   `conf:"CORS.enabled"`
	AuthEnabled               bool   `conf:"auth.enabled"`
	Users                     string `conf:"auth.users"`
	Groups                    string `conf:"auth.groups"`
	HOTP                      HOTP   `conf:"auth.hotp"`
	Locale                    string `conf:"translation.locale"`
	Protocol                  string `conf:"protocol.version"`
}

type OpenAPI struct {
	Enabled   bool   `conf:"enabled"`
	Directory string `conf:"directory"`
}

func NewREST() *REST {
	return &REST{
		HttpEnabled:               false,
		HttpPort:                  8080,
		HttpsEnabled:              true,
		HttpsPort:                 8443,
		TLSKeyFile:                "uhppoted.key",
		TLSCertificateFile:        "uhppoted.cert",
		CACertificateFile:         "ca.cert",
		RequireClientCertificates: true,
		CORSEnabled:               false,
		AuthEnabled:               false,
		Users:                     restUsers,
		Groups:                    restGroups,
		HOTP: HOTP{
			Range:    8,
			Counters: restHOTP,
		},
	}
}

func NewOpenAPI() *OpenAPI {
	return &OpenAPI{
		Enabled:   false,
		Directory: "./openapi",
	}
}
