package config

import (
	"golang.org/x/sys/windows"
	"path/filepath"
)

// DefaultConfig is the default file path for the uhppoted configuration file
var DefaultConfig = filepath.Join(workdir(), "uhppoted.conf")

var restUsers string = filepath.Join(workdir(), "rest", "users")
var restGroups string = filepath.Join(workdir(), "rest", "groups")
var restHOTP string = filepath.Join(workdir(), "rest", "counters")

var mqttBrokerCertificate string = filepath.Join(workdir(), "mqtt", "broker.cert")
var mqttClientCertificate string = filepath.Join(workdir(), "mqtt", "client.cert")
var mqttClientKey string = filepath.Join(workdir(), "mqtt", "client.key")
var mqttUsers string = filepath.Join(workdir(), "mqtt.permissions.users")
var mqttGroups string = filepath.Join(workdir(), "mqtt.permissions.groups")
var mqttCards string = filepath.Join(workdir(), "mqtt", "cards")
var hotpSecrets string = filepath.Join(workdir(), "mqtt.hotp.secrets")
var rsaKeyDir string = filepath.Join(workdir(), "mqtt", "rsa")

var eventIDs string = filepath.Join(workdir(), "mqtt.events.retrieved")
var hotpCounters string = filepath.Join(workdir(), "mqtt.hotp.counters")
var nonceServer string = filepath.Join(workdir(), "mqtt.nonce")
var nonceClients string = filepath.Join(workdir(), "mqtt.nonce.counters")

var httpdAuthDB string = filepath.Join(workdir(), "httpd", "auth.json")
var httpdCACertificate string = filepath.Join(workdir(), "httpd", "ca.cert")
var httpdTLSCertificate string = filepath.Join(workdir(), "httpd", "uhppoted.cert")
var httpdTLSKey string = filepath.Join(workdir(), "httpd", "uhppoted.key")
var httpdControllersFile string = filepath.Join(workdir(), "httpd", "system", "controllers.json")
var httpdDoorsFile string = filepath.Join(workdir(), "httpd", "system", "doors.json")
var httpdGroupsFile string = filepath.Join(workdir(), "httpd", "system", "groups.json")
var httpdCardsFile string = filepath.Join(workdir(), "httpd", "system", "cards.json")
var httpdEventsFile string = filepath.Join(workdir(), "httpd", "system", "events.json")
var httpdRulesACL string = filepath.Join(workdir(), "httpd", "acl.grl")
var httpdRulesSystem string = filepath.Join(workdir(), "httpd", "system.grl")
var httpdRulesCards string = filepath.Join(workdir(), "httpd", "cards.grl")
var httpdRulesDoors string = filepath.Join(workdir(), "httpd", "doors.grl")
var httpdRulesGroups string = filepath.Join(workdir(), "httpd", "groups.grl")
var httpdAuditFile string = filepath.Join(workdir(), "httpd", "audit", "audit.log")

func workdir() string {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return `C:\uhppoted`
	}

	return filepath.Join(programData, "uhppoted")
}
