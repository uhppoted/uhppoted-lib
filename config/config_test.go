package config

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/encoding/conf"
)

var configuration = []byte(`# SYSTEM
bind.address = 192.168.1.100:54321
broadcast.address = 192.168.1.255:30000
listen.address = 192.168.1.100:12345
timeout = 3.75s
card.format = Wiegand-26

monitoring.healthcheck.interval = 31s
monitoring.healthcheck.idle = 67s
monitoring.healthcheck.ignore = 97s
monitoring.watchdog.interval = 23s

# MQTT
mqtt.connection.broker = tls://127.0.0.63:8887
mqtt.connection.client.ID = muppet
mqtt.connection.username = me
mqtt.connection.password = pickme
mqtt.connection.broker.certificate = mqtt-broker.cert
mqtt.connection.client.certificate = mqtt-client.cert
mqtt.connection.client.key = mqtt-client.key
mqtt.topic.root = uhppoted-qwerty
mqtt.topic.replies = /uiop
mqtt.topic.events = ./asdf
mqtt.topic.events.real-time = ./eekvents
mqtt.topic.system = sys

# AWS
aws.region = us-west-2

# Wild Apricot
wild-apricot.http.client-timeout = 12s
wild-apricot.http.retries = 6
wild-apricot.http.retry-delay = 9s
wild-apricot.fields.card-number = Ye Olde Cardde Nymber
wild-apricot.fields.PIN = LePIN
wild-apricot.facility-code = 89

# HTTPD
httpd.retention = 4h45m

# DEVICES
UT0311-L0x.405419896.name = Q405419896
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = Front Door
UT0311-L0x.405419896.door.2 = Side Door
UT0311-L0x.405419896.door.3 = Garage
UT0311-L0x.405419896.door.4 = Workshop
UT0311-L0x.405419896.timezone = France/Paris
`)

func TestDefaultConfig(t *testing.T) {
	bind, broadcast, listen := DefaultIpAddresses()

	expected := Config{
		System: System{
			BindAddress:         &bind,
			BroadcastAddress:    &broadcast,
			ListenAddress:       &listen,
			Timeout:             2500 * time.Millisecond,
			HealthCheckInterval: 15 * time.Second,
			HealthCheckIdle:     60 * time.Second,
			HealthCheckIgnore:   5 * time.Minute,
			WatchdogInterval:    5 * time.Second,
			CardFormat:          types.WiegandAny,
		},

		MQTT: MQTT{
			Connection: Connection{
				ClientID: "uhppoted-mqttd",
			},
		},

		AWS: AWS{
			Credentials: "",
			Profile:     "default",
			Region:      "us-east-1",
		},

		WildApricot: WildApricot{
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
				PIN        string `conf:"PIN"`
			}{
				CardNumber: "Card Number",
				PIN:        "PIN",
			},
		},

		HTTPD: HTTPD{
			Retention: 6 * time.Hour,
		},
	}

	config := NewConfig()

	if !reflect.DeepEqual(config.System, expected.System) {
		t.Errorf("Incorrect system default configuration:\nexpected:%+v,\ngot:     %+v", expected.System, config.System)
	}

	if config.MQTT.Connection.ClientID != expected.MQTT.Connection.ClientID {
		t.Errorf("Expected mqtt.connection.client.ID: '%v', got: '%v'", expected.MQTT.Connection.ClientID, config.MQTT.Connection.ClientID)
	}

	if !reflect.DeepEqual(config.AWS, expected.AWS) {
		t.Errorf("Incorrect AWS default configuration:\nexpected:%+v,\ngot:     %+v", expected.AWS, config.AWS)
	}

	if !reflect.DeepEqual(config.WildApricot, expected.WildApricot) {
		t.Errorf("Incorrect WildApricot default configuration:\nexpected:%+v,\ngot:     %+v", expected.WildApricot, config.WildApricot)
	}

	if config.HTTPD.Retention != expected.HTTPD.Retention {
		t.Errorf("Expected http.retention:'%v', got: '%v'", expected.HTTPD.Retention, config.HTTPD.Retention)
	}
}

func TestConfigUnmarshal(t *testing.T) {
	expected := Config{
		System: System{
			BindAddress:         &types.BindAddr{IP: []byte{192, 168, 1, 100}, Port: 54321},
			BroadcastAddress:    &types.BroadcastAddr{IP: []byte{192, 168, 1, 255}, Port: 30000},
			ListenAddress:       &types.ListenAddr{IP: []byte{192, 168, 1, 100}, Port: 12345},
			Timeout:             3750 * time.Millisecond,
			HealthCheckInterval: 31 * time.Second,
			HealthCheckIdle:     67 * time.Second,
			HealthCheckIgnore:   97 * time.Second,
			WatchdogInterval:    23 * time.Second,
			CardFormat:          types.Wiegand26,
		},

		MQTT: MQTT{
			Connection: Connection{
				Broker:            "tls://127.0.0.63:8887",
				ClientID:          "muppet",
				Username:          "me",
				Password:          "pickme",
				BrokerCertificate: "mqtt-broker.cert",
				ClientCertificate: "mqtt-client.cert",
				ClientKey:         "mqtt-client.key",
			},

			Topics: Topics{
				Root:           "uhppoted-qwerty",
				Requests:       "./requests",
				Replies:        "/uiop",
				EventsFeed:     "./asdf",
				RealTimeEvents: "./eekvents",
				System:         "sys",
			},
		},

		AWS: AWS{
			Credentials: "",
			Profile:     "default",
			Region:      "us-west-2",
		},

		WildApricot: WildApricot{
			HTTP: struct {
				ClientTimeout time.Duration `conf:"client-timeout"`
				Retries       int           `conf:"retries"`
				RetryDelay    time.Duration `conf:"retry-delay"`
			}{
				ClientTimeout: 12 * time.Second,
				Retries:       6,
				RetryDelay:    9 * time.Second,
			},

			Fields: struct {
				CardNumber string `conf:"card-number"`
				PIN        string `conf:"PIN"`
			}{
				CardNumber: "Ye Olde Cardde Nymber",
				PIN:        "LePIN",
			},

			FacilityCode: "89",
		},

		HTTPD: HTTPD{
			Retention: 285 * time.Minute,
		},
	}

	config := NewConfig()
	err := conf.Unmarshal(configuration, config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(config.System, expected.System) {
		t.Errorf("Incorrect system configuration:\nexpected:%+v,\ngot:     %+v", expected.System, config.System)
	}

	if !reflect.DeepEqual(config.Connection, expected.Connection) {
		t.Errorf("Incorrect 'mqtt.connection' configuration:\nexpected:%+v,\ngot:     %+v", expected.Connection, config.Connection)
	}

	if !reflect.DeepEqual(config.Topics, expected.Topics) {
		t.Errorf("Incorrect 'mqtt.topics' configuration:\nexpected:%+v,\ngot:     %+v", expected.Topics, config.Topics)
	}

	if config.Topics.Resolve(config.Topics.Requests) != "uhppoted-qwerty/requests" {
		t.Errorf("Expected 'mqtt::topic.requests' %v, got:%v", "wystd-qwerty/requests", config.Topics.Resolve(config.Topics.Requests))
	}

	if config.Topics.Resolve(config.Topics.Replies) != "uiop" {
		t.Errorf("Expected 'mqtt::topic.replies' %v, got:%v", "uiop", config.Topics.Resolve(config.Topics.Replies))
	}

	if config.Topics.Resolve(config.Topics.EventsFeed) != "uhppoted-qwerty/asdf" {
		t.Errorf("Expected 'mqtt::topic.events' %v, got:%v", "uhppoted-qwerty/asdf", config.Topics.Resolve(config.Topics.EventsFeed))
	}

	if config.Topics.Resolve(config.Topics.RealTimeEvents) != "uhppoted-qwerty/eekvents" {
		t.Errorf("Expected 'mqtt::topic.events.real-time' %v, got:%v", "uhppoted-qwerty/eekvents", config.Topics.Resolve(config.Topics.RealTimeEvents))
	}

	if config.Topics.Resolve(config.Topics.System) != "uhppoted-qwerty/sys" {
		t.Errorf("Expected 'mqtt::topic.system' %v, got:%v", "uhppoted-qwerty/sys", config.Topics.Resolve(config.Topics.System))
	}

	if !reflect.DeepEqual(config.AWS, expected.AWS) {
		t.Errorf("Incorrect AWS configuration:\nexpected:%+v,\ngot:     %+v", expected.AWS, config.AWS)
	}

	if !reflect.DeepEqual(config.WildApricot, expected.WildApricot) {
		t.Errorf("Incorrect Wild Apricot configuration:\nexpected:%+v,\ngot:     %+v", expected.WildApricot, config.WildApricot)
	}

	if config.HTTPD.Retention != expected.HTTPD.Retention {
		t.Errorf("Incorrect httpd retention:\nexpected:%+v,\ngot:     %+v", expected.HTTPD.Retention, config.HTTPD.Retention)
	}

	if d := config.Devices[405419896]; d == nil {
		t.Errorf("Expected 'device' for ID '%v', got:'%v'", 405419896, d)
	} else {
		if d.Name != "Q405419896" {
			t.Errorf("Expected 'device.name' %s for ID '%v', got:'%v'", "Q405419896", 405419896, d.Name)
		}

		address := net.UDPAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: 60000,
			Zone: "",
		}

		if !reflect.DeepEqual(d.Address, &address) {
			t.Errorf("Expected 'device.address' %s for ID '%v', got:'%v'", &address, 405419896, d.Address)
		}

		if len(d.Doors) != 4 {
			t.Errorf("Expected 4 entries for 'device.door' %s for ID '%v', got:%v", &address, 405419896, len(d.Doors))
		} else {
			if d.Doors[0] != "Front Door" {
				t.Errorf("Expected 'device.door[0]' %s for ID '%v', got:'%s'", "Front Door", 405419896, d.Doors[0])
			}

			if d.Doors[1] != "Side Door" {
				t.Errorf("Expected 'device.door[1]' %s for ID '%v', got:'%s'", "Side Door", 405419896, d.Doors[1])
			}

			if d.Doors[2] != "Garage" {
				t.Errorf("Expected 'device.door[2]' %s for ID '%v', got:'%s'", "Garage", 405419896, d.Doors[2])
			}

			if d.Doors[3] != "Workshop" {
				t.Errorf("Expected 'device.door[3]' %s for ID '%v', got:'%s'", "Workshop", 405419896, d.Doors[3])
			}
		}

		if d.TimeZone != "France/Paris" {
			t.Errorf("Expected 'device.timezone' %s for ID '%v', got:'%v'", "France/Paris", 405419896, d.TimeZone)
		}

	}
}

func TestDefaultConfigWrite(t *testing.T) {
	bind, broadcast, listen := DefaultIpAddresses()

	expected := fmt.Sprintf(`# SYSTEM
; bind.address = %[1]s
; broadcast.address = %[2]s
; listen.address = %[3]s
; timeout = 2.5s
; monitoring.healthcheck.interval = 15s
; monitoring.healthcheck.idle = 1m0s
; monitoring.healthcheck.ignore = 5m0s
; monitoring.watchdog.interval = 5s
; card.format = any

# REST
; rest.http.enabled = false
; rest.http.port = 8080
; rest.https.enabled = true
; rest.https.port = 8443
; rest.tls.key = uhppoted.key
; rest.tls.certificate = uhppoted.cert
; rest.tls.ca = ca.cert
; rest.tls.client.certificates = true
; rest.CORS.enabled = false
; rest.auth.enabled = false
; rest.auth.users = %[4]s
; rest.auth.groups = %[5]s
; rest.auth.hotp.range = 8
; rest.auth.hotp.secrets = 
; rest.auth.hotp.counters = %[6]s
; rest.translation.locale = 
; rest.protocol.version = 

# MQTT
; mqtt.server.ID = uhppoted
; mqtt.connection.broker = tcp://127.0.0.1:1883
; mqtt.connection.client.ID = uhppoted-mqttd
; mqtt.connection.username = 
; mqtt.connection.password = 
; mqtt.connection.broker.certificate = %[7]s
; mqtt.connection.client.certificate = %[8]s
; mqtt.connection.client.key = %[9]s
; mqtt.connection.verify = 
; mqtt.topic.root = uhppoted/gateway
; mqtt.topic.requests = ./requests
; mqtt.topic.replies = ./replies
; mqtt.topic.events = ./events
; mqtt.topic.events.real-time = ./events/live
; mqtt.topic.system = ./system
; mqtt.translation.locale = 
; mqtt.protocol.version = 
; mqtt.alerts.qos = 1
; mqtt.alerts.retained = true
; mqtt.events.key = events
; mqtt.system.key = system
; mqtt.events.index.filepath = %[10]s
; mqtt.permissions.enabled = false
; mqtt.permissions.users = %[11]s
; mqtt.permissions.groups = %[12]s
; mqtt.cards = %[13]s
; mqtt.security.HMAC.required = false
; mqtt.security.HMAC.key = 
; mqtt.security.authentication = HOTP, RSA
; mqtt.security.hotp.range = 8
; mqtt.security.hotp.secrets = %[14]s
; mqtt.security.hotp.counters = %[15]s
; mqtt.security.rsa.keys = %[16]s
; mqtt.security.nonce.required = true
; mqtt.security.nonce.server = %[17]s
; mqtt.security.nonce.clients = %[18]s
; mqtt.security.outgoing.sign = true
; mqtt.security.outgoing.encrypt = true
; mqtt.lockfile.remove = false
; mqtt.disconnects.enabled = true
; mqtt.disconnects.interval = 5m0s
; mqtt.disconnects.max = 10
; mqtt.acl.verify = RSA

# AWS
; aws.credentials = 
; aws.profile = default
; aws.region = us-east-1

# HTTPD
; httpd.html = 
; httpd.http.enabled = false
; httpd.http.port = 8080
; httpd.https.enabled = true
; httpd.https.port = 8443
; httpd.tls.ca = %[20]s
; httpd.tls.certificate = %[21]s
; httpd.tls.key = %[22]s
; httpd.tls.client.certificates.required = false
; httpd.security.auth = basic
; httpd.security.local.db = %[19]s
; httpd.security.cookie.max-age = 24
; httpd.security.login.expiry = 1m
; httpd.security.session.expiry = 60m
; httpd.security.otp.issuer = uhppoted-httpd
; httpd.security.otp.login = allow
; httpd.request.timeout = 5s
; httpd.system.interfaces = %[23]s
; httpd.system.controllers = %[24]s
; httpd.system.doors = %[25]s
; httpd.system.groups = %[26]s
; httpd.system.cards = %[27]s
; httpd.system.events = %[28]s
; httpd.system.logs = %[29]s
; httpd.system.users = %[30]s
; httpd.system.history = %[31]s
; httpd.system.refresh = 30s
; httpd.system.windows.ok = 1m0s
; httpd.system.windows.uncertain = 5m0s
; httpd.system.windows.systime = 5m0s
; httpd.system.windows.expires = 2m0s
; httpd.db.rules.acl = %[32]s
; httpd.db.rules.interfaces = %[33]s
; httpd.db.rules.controllers = %[34]s
; httpd.db.rules.cards = %[35]s
; httpd.db.rules.doors = %[36]s
; httpd.db.rules.groups = %[37]s
; httpd.db.rules.events = %[38]s
; httpd.db.rules.logs = %[39]s
; httpd.db.rules.users = %[40]s
; httpd.audit.file = %[41]s
; httpd.retention = 6h0m0s
; httpd.timezones = 
; httpd.PIN.enabled = false

# Wild Apricot
; wild-apricot.http.client-timeout = 10s
; wild-apricot.http.retries = 3
; wild-apricot.http.retry-delay = 5s
; wild-apricot.fields.card-number = Card Number
; wild-apricot.fields.PIN = PIN
; wild-apricot.display-order.groups = 
; wild-apricot.display-order.doors = 
; wild-apricot.facility-code = 

# OPEN API
# openapi.enabled = false
# openapi.directory = ./openapi

# DEVICES
# Example configuration for UTO311-L04 with serial number 405419896
# UT0311-L0x.405419896.name = D405419896
# UT0311-L0x.405419896.address = 192.168.1.100:60000
# UT0311-L0x.405419896.door.1 = Front Door
# UT0311-L0x.405419896.door.2 = Side Door
# UT0311-L0x.405419896.door.3 = Garage
# UT0311-L0x.405419896.door.4 = Workshop
# UT0311-L0x.405419896.timezone = UTC+2
`, bind.String(), broadcast.String(), listen.String(),
		restUsers, restGroups, restHOTP,
		mqttBrokerCertificate, mqttClientCertificate, mqttClientKey, eventIDs, mqttUsers, mqttGroups, mqttCards, hotpSecrets, hotpCounters, rsaKeyDir,
		nonceServer, nonceClients,
		httpdAuthDB, httpdCACertificate, httpdTLSCertificate, httpdTLSKey,
		httpdInterfacesFile, httpdControllersFile, httpdDoorsFile, httpdGroupsFile, httpdCardsFile, httpdEventsFile, httpdLogsFile, httpdUsersFile, httpdHistoryFile,
		httpdRulesACL, httpdRulesInterfaces, httpdRulesControllers, httpdRulesCards, httpdRulesDoors, httpdRulesGroups, httpdRulesEvents, httpdRulesLogs, httpdRulesUsers, httpdAuditFile)

	config := NewConfig()

	var b bytes.Buffer

	if err := config.Write(&b); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if s := b.String(); s != expected {
		re := regexp.MustCompile(`(\r)?\n`)
		p := re.Split(s, -1)
		q := re.Split(expected, -1)
		N := len(q)

		if N > len(p) {
			N = len(p)
		}
		i := 0
		for i < N {
			if p[i] != q[i] {
				t.Fatalf("Line %d: output from Config.Writer does not match\n   expected:\n%s\n   got:     \n%s\n", i, q[i], p[i])
			}
			i++
		}

		if i < len(p) {
			t.Fatalf("Line %d: unexpected output from Config.Writer\n   got:\n%s\n", i, p[i])
		}

		if i < len(q) {
			t.Fatalf("Line %d: missing from Config.Writer\n   expected:\n%s\n", i, q[i])
		}
	}
}

func TestConfigWrite(t *testing.T) {
	bind, broadcast, listen := DefaultIpAddresses()

	expected := fmt.Sprintf(`# SYSTEM
; bind.address = %[1]s
; broadcast.address = %[2]s
; listen.address = %[3]s
timeout = %[4]v
; monitoring.healthcheck.interval = 15s
; monitoring.healthcheck.idle = 1m0s
; monitoring.healthcheck.ignore = 5m0s
; monitoring.watchdog.interval = 5s
; card.format = any

# REST
; rest.http.enabled = false
; rest.http.port = 8080
; rest.https.enabled = true
; rest.https.port = 8443
; rest.tls.key = uhppoted.key
; rest.tls.certificate = uhppoted.cert
; rest.tls.ca = ca.cert
; rest.tls.client.certificates = true
; rest.CORS.enabled = false
; rest.auth.enabled = false
; rest.auth.users = %[5]s
; rest.auth.groups = %[6]s
; rest.auth.hotp.range = 8
; rest.auth.hotp.secrets = 
; rest.auth.hotp.counters = %[7]s
; rest.translation.locale = 
; rest.protocol.version = 

# MQTT
; mqtt.server.ID = uhppoted
; mqtt.connection.broker = tcp://127.0.0.1:1883
; mqtt.connection.client.ID = uhppoted-mqttd
; mqtt.connection.username = 
; mqtt.connection.password = 
; mqtt.connection.broker.certificate = %[8]s
; mqtt.connection.client.certificate = %[9]s
; mqtt.connection.client.key = %[10]s
; mqtt.connection.verify = 
; mqtt.topic.root = uhppoted/gateway
; mqtt.topic.requests = ./requests
; mqtt.topic.replies = ./replies
; mqtt.topic.events = ./events
; mqtt.topic.events.real-time = ./events/live
; mqtt.topic.system = ./system
; mqtt.translation.locale = 
; mqtt.protocol.version = 
; mqtt.alerts.qos = 1
; mqtt.alerts.retained = true
; mqtt.events.key = events
; mqtt.system.key = system
; mqtt.events.index.filepath = %[11]s
; mqtt.permissions.enabled = false
; mqtt.permissions.users = %[12]s
; mqtt.permissions.groups = %[13]s
; mqtt.cards = %[14]s
; mqtt.security.HMAC.required = false
; mqtt.security.HMAC.key = 
; mqtt.security.authentication = HOTP, RSA
; mqtt.security.hotp.range = 8
; mqtt.security.hotp.secrets = %[15]s
; mqtt.security.hotp.counters = %[16]s
; mqtt.security.rsa.keys = %[17]s
; mqtt.security.nonce.required = true
; mqtt.security.nonce.server = %[18]s
; mqtt.security.nonce.clients = %[19]s
; mqtt.security.outgoing.sign = true
; mqtt.security.outgoing.encrypt = true
; mqtt.lockfile.remove = false
; mqtt.disconnects.enabled = true
; mqtt.disconnects.interval = 5m0s
; mqtt.disconnects.max = 10
; mqtt.acl.verify = RSA

# AWS
; aws.credentials = 
; aws.profile = default
; aws.region = us-east-1

# HTTPD
; httpd.html = 
; httpd.http.enabled = false
; httpd.http.port = 8080
; httpd.https.enabled = true
; httpd.https.port = 8443
; httpd.tls.ca = %[21]s
; httpd.tls.certificate = %[22]s
; httpd.tls.key = %[23]s
; httpd.tls.client.certificates.required = false
; httpd.security.auth = basic
; httpd.security.local.db = %[20]s
; httpd.security.cookie.max-age = 24
; httpd.security.login.expiry = 1m
; httpd.security.session.expiry = 60m
; httpd.security.otp.issuer = uhppoted-httpd
; httpd.security.otp.login = allow
; httpd.request.timeout = 5s
; httpd.system.interfaces = %[24]s
; httpd.system.controllers = %[25]s
; httpd.system.doors = %[26]s
; httpd.system.groups = %[27]s
; httpd.system.cards = %[28]s
; httpd.system.events = %[29]s
; httpd.system.logs = %[30]s
; httpd.system.users = %[31]s
; httpd.system.history = %[32]s
; httpd.system.refresh = 30s
; httpd.system.windows.ok = 1m0s
; httpd.system.windows.uncertain = 5m0s
; httpd.system.windows.systime = 5m0s
; httpd.system.windows.expires = 2m0s
; httpd.db.rules.acl = %[33]s
; httpd.db.rules.interfaces = %[34]s
; httpd.db.rules.controllers = %[35]s
; httpd.db.rules.cards = %[36]s
; httpd.db.rules.doors = %[37]s
; httpd.db.rules.groups = %[38]s
; httpd.db.rules.events = %[39]s
; httpd.db.rules.logs = %[40]s
; httpd.db.rules.users = %[41]s
; httpd.audit.file = %[42]s
httpd.retention = 5h30m0s
; httpd.timezones = 
; httpd.PIN.enabled = false

# Wild Apricot
; wild-apricot.http.client-timeout = 10s
; wild-apricot.http.retries = 3
; wild-apricot.http.retry-delay = 5s
; wild-apricot.fields.card-number = Card Number
; wild-apricot.fields.PIN = PIN
; wild-apricot.display-order.groups = 
; wild-apricot.display-order.doors = 
; wild-apricot.facility-code = 

# OPEN API
# openapi.enabled = false
# openapi.directory = ./openapi

# DEVICES
UT0311-L0x.303986753.name = Q303986753
UT0311-L0x.303986753.door.1 = A
UT0311-L0x.303986753.door.2 = B
UT0311-L0x.303986753.door.3 = C
UT0311-L0x.303986753.door.4 = D
UT0311-L0x.303986753.timezone = UTC

UT0311-L0x.405419896.name = Z405419896
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = D1
UT0311-L0x.405419896.door.2 = D2
UT0311-L0x.405419896.door.3 = D3
UT0311-L0x.405419896.door.4 = D4
UT0311-L0x.405419896.timezone = France/Paris
`, bind.String(), broadcast.String(), listen.String(), 4500*time.Millisecond,
		restUsers, restGroups, restHOTP,
		mqttBrokerCertificate, mqttClientCertificate, mqttClientKey, eventIDs, mqttUsers, mqttGroups, mqttCards, hotpSecrets, hotpCounters, rsaKeyDir,
		nonceServer, nonceClients,
		httpdAuthDB, httpdCACertificate, httpdTLSCertificate, httpdTLSKey,
		httpdInterfacesFile, httpdControllersFile, httpdDoorsFile, httpdGroupsFile, httpdCardsFile, httpdEventsFile, httpdLogsFile, httpdUsersFile, httpdHistoryFile,
		httpdRulesACL, httpdRulesInterfaces, httpdRulesControllers, httpdRulesCards, httpdRulesDoors, httpdRulesGroups, httpdRulesEvents, httpdRulesLogs, httpdRulesUsers, httpdAuditFile)

	config := NewConfig()

	config.Timeout = 4500 * time.Millisecond

	config.Devices = DeviceMap{
		405419896: &Device{
			Name: "Z405419896",
			Address: &net.UDPAddr{
				IP:   []byte{192, 168, 1, 100},
				Port: 60000,
				Zone: "",
			},
			Doors:    []string{"D1", "D2", "D3", "D4"},
			TimeZone: "France/Paris",
		},

		303986753: &Device{
			Name:     "Q303986753",
			Doors:    []string{"A", "B", "C", "D"},
			TimeZone: "UTC",
		},
	}

	config.HTTPD.Retention = 330 * time.Minute

	var b bytes.Buffer

	if err := config.Write(&b); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if s := b.String(); s != expected {
		re := regexp.MustCompile(`(\r)?\n`)
		p := re.Split(s, -1)
		q := re.Split(expected, -1)
		N := len(q)

		if N > len(p) {
			N = len(p)
		}
		i := 0
		for i < N {
			if p[i] != q[i] {
				t.Fatalf("Line %d: output from Config.Writer does not match\n   expected:\n%s\n   got:     \n%s\n", i, q[i], p[i])
			}
			i++
		}

		if i < len(p) {
			t.Fatalf("Line %d: unexpected output from Config.Writer\n   got:\n%s\n", i, p[i])
		}

		if i < len(q) {
			t.Fatalf("Line %d: missing from Config.Writer\n   expected:\n%s\n", i, q[i])
		}
	}
}

func TestConfigValidateWithBindPort60000(t *testing.T) {
	configuration := []byte(`# SYSTEM
bind.address = 192.168.1.100:60000
broadcast.address = 192.168.1.255:60000
listen.address = 192.168.1.100:60001
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err == nil {
		t.Fatalf("Expected error, got: %v", err)
	}
}

func TestConfigValidateWithTheSameBindAndBroadcastPorts(t *testing.T) {
	configuration := []byte(`# SYSTEM
bind.address = 192.168.1.100:12345
broadcast.address = 192.168.1.255:12345
listen.address = 192.168.1.100:60001
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := fmt.Errorf("bind.address port (12345) must not be the same as the broadcast.address port")

	err := config.Validate()
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("Expected error:%v, got:%v", expected, err)
	}
}

func TestConfigValidateWithTheSameBindAndListenPorts(t *testing.T) {
	configuration := []byte(`# SYSTEM
bind.address = 192.168.1.100:60001
listen.address = 192.168.1.100:60001
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := fmt.Errorf("bind.address port (60001) must not be the same as the listen.address port")

	err := config.Validate()
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("Expected error:%v, got:%v", expected, err)
	}
}

func TestConfigValidateWithInvalidBroadcastPort(t *testing.T) {
	configuration := []byte(`# SYSTEM
bind.address = 192.168.1.100:0
broadcast.address = 192.168.1.255:0
listen.address = 192.168.1.100:60001
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err == nil {
		t.Fatalf("Expected 'invalid broadcast address port' error, got %v", err)
	}
}

func TestConfigValidateWithInvalidListenPort(t *testing.T) {
	configuration := []byte(`# SYSTEM
bind.address = 192.168.1.100:0
listen.address = 192.168.1.100:0
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err == nil {
		t.Fatalf("Expected 'invalid listen port' error, got %v", err)
	}
}

func TestConfigValidateWithValidDevice(t *testing.T) {
	configuration := []byte(`# DEVICES
UT0311-L0x.405419896.name = Q405419896
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = Front Door
UT0311-L0x.405419896.door.2 = Side Door
UT0311-L0x.405419896.door.3 = Garage
UT0311-L0x.405419896.door.4 = Workshop
UT0311-L0x.405419896.timezone = France/Paris
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if err := config.Validate(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestConfigValidateWithInvalidDevice(t *testing.T) {
	configuration := []byte(`# DEVICES
UT0311-L0x.405419896.name = Q405419896
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = Front Door
UT0311-L0x.405419896.door.2 = Side Door
UT0311-L0x.405419896.door.3 = Garage
UT0311-L0x.405419896.door.4 = Front Door
UT0311-L0x.405419896.timezone = France/Paris
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := fmt.Errorf("door 'Front Door' is defined more than once in configuration")

	err := config.Validate()
	if err == nil || err.Error() != expected.Error() {
		t.Errorf("Expected error:%v, got:%v", expected, err)
	}
}

func TestConfigValidateWithBlankDoors(t *testing.T) {
	configuration := []byte(`# DEVICES
UT0311-L0x.405419896.name = Q405419896
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = 
UT0311-L0x.405419896.door.2 = 
UT0311-L0x.405419896.door.3 = 
UT0311-L0x.405419896.door.4 = 
UT0311-L0x.405419896.timezone = France/Paris
`)

	config := NewConfig()
	if err := conf.Unmarshal(configuration, config); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if err := config.Validate(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
