package config

import (
	"strings"
	"time"
)

type MQTT struct {
	ServerID        string      `conf:"server.ID"`
	Connection      Connection  `conf:"connection"`
	Topics          Topics      `conf:"topic"`
	Locale          string      `conf:"translation.locale"`
	Protocol        string      `conf:"protocol.version"`
	Alerts          Alerts      `conf:"alerts"`
	EventsKeyID     string      `conf:"events.key"`
	SystemKeyID     string      `conf:"system.key"`
	EventIDs        string      `conf:"events.index.filepath"`
	Permissions     Permissions `conf:"permissions"`
	Cards           string      `conf:"cards"`
	HMAC            HMAC        `conf:"security.HMAC"`
	Authentication  string      `conf:"security.authentication"`
	HOTP            HOTP        `conf:"security.hotp"`
	RSA             RSA         `conf:"security.rsa"`
	Nonce           Nonce       `conf:"security.nonce"`
	SignOutgoing    bool        `conf:"security.outgoing.sign"`
	EncryptOutgoing bool        `conf:"security.outgoing.encrypt"`
	LockfileRemove  bool        `conf:"lockfile.remove"`
	Disconnects     Disconnects `conf:"disconnects"`
	ACL             ACL         `conf:"acl"`
}

type Connection struct {
	Broker            string `conf:"broker"`
	ClientID          string `conf:"client.ID"`
	Username          string `conf:"username"`
	Password          string `conf:"password"`
	BrokerCertificate string `conf:"broker.certificate"`
	ClientCertificate string `conf:"client.certificate"`
	ClientKey         string `conf:"client.key"`
	Verify            string `conf:"verify"`
}

type Topics struct {
	Root       string `conf:"root"`
	Requests   string `conf:"requests"`
	Replies    string `conf:"replies"`
	EventsFeed string `conf:"events"`
	LiveEvents string `conf:"events.live"`
	System     string `conf:"system"`
}

type Alerts struct {
	QOS      byte `conf:"qos"`
	Retained bool `conf:"retained"`
}

type HMAC struct {
	Required bool   `conf:"required"`
	Key      string `conf:"key"`
}

type HOTP struct {
	Range    uint64 `conf:"range"`
	Secrets  string `conf:"secrets"`
	Counters string `conf:"counters"`
}

type RSA struct {
	KeyDir string `conf:"keys"`
}

type Nonce struct {
	Required bool   `conf:"required"`
	Server   string `conf:"server"`
	Clients  string `conf:"clients"`
}

type Permissions struct {
	Enabled bool   `conf:"enabled"`
	Users   string `conf:"users"`
	Groups  string `conf:"groups"`
}

type Disconnects struct {
	Enabled  bool          `conf:"enabled"`
	Interval time.Duration `conf:"interval"`
	Max      uint32        `conf:"max"`
}

type ACL struct {
	Verify string `conf:"verify"`
}

func (t *Topics) Resolve(subtopic string) string {
	if strings.HasPrefix(subtopic, "/") {
		return strings.ReplaceAll(strings.TrimPrefix(subtopic, "/"), " ", "")
	}

	if strings.HasPrefix(subtopic, "./") {
		return strings.ReplaceAll(t.Root+"/"+strings.TrimPrefix(subtopic, "./"), " ", "")
	}

	return strings.ReplaceAll(t.Root+"/"+subtopic, " ", "")
}

func NewMQTT() *MQTT {
	return &MQTT{
		ServerID: "uhppoted",
		Connection: Connection{
			Broker:            "tcp://127.0.0.1:1883",
			ClientID:          "uhppoted-mqttd",
			BrokerCertificate: mqttBrokerCertificate,
			ClientCertificate: mqttClientCertificate,
			ClientKey:         mqttClientKey,
		},
		Topics: Topics{
			Root:       "uhppoted/gateway",
			Requests:   "./requests",
			Replies:    "./replies",
			EventsFeed: "./events",
			LiveEvents: "./events/live",
			System:     "./system",
		},
		Alerts: Alerts{
			QOS:      1,
			Retained: true,
		},
		EventsKeyID:     "events",
		SystemKeyID:     "system",
		SignOutgoing:    true,
		EncryptOutgoing: true,
		HMAC: HMAC{
			Required: false,
			Key:      "",
		},
		Authentication: "HOTP, RSA",
		HOTP: HOTP{
			Range:    8,
			Secrets:  hotpSecrets,
			Counters: hotpCounters,
		},
		RSA: RSA{
			KeyDir: rsaKeyDir,
		},
		Nonce: Nonce{
			Required: true,
			Server:   nonceServer,
			Clients:  nonceClients,
		},
		Permissions: Permissions{
			Enabled: false,
			Users:   mqttUsers,
			Groups:  mqttGroups,
		},
		Cards:    mqttCards,
		EventIDs: eventIDs,
		Disconnects: Disconnects{
			Enabled:  true,
			Interval: 5 * time.Minute,
			Max:      10,
		},
		ACL: ACL{
			Verify: "RSA",
		},
	}
}
