package config

import (
	"golang.org/x/sys/windows"
	"path/filepath"
)

var DefaultConfig = filepath.Join(workdir(), "uhppoted.conf")

var _etc string = workdir()
var _var string = workdir()

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

var httpdAuthDB string = filepath.Join(_etc, "httpd", "auth.json")
var httpdCACertificate string = filepath.Join(_etc, "httpd", "ca.cert")
var httpdTLSCertificate string = filepath.Join(_etc, "httpd", "uhppoted.cert")
var httpdTLSKey string = filepath.Join(_etc, "httpd", "uhppoted.key")

var httpdRulesACL string = filepath.Join(_etc, "httpd", "acl.grl")
var httpdRulesInterfaces string = ""
var httpdRulesControllers string = ""
var httpdRulesCards string = ""
var httpdRulesDoors string = ""
var httpdRulesGroups string = ""
var httpdRulesEvents string = ""
var httpdRulesLogs string = ""
var httpdRulesUsers string = ""

var httpdInterfacesFile string = filepath.Join(_var, "httpd", "system", "interfaces.json")
var httpdControllersFile string = filepath.Join(_var, "httpd", "system", "controllers.json")
var httpdDoorsFile string = filepath.Join(_var, "httpd", "system", "doors.json")
var httpdGroupsFile string = filepath.Join(_var, "httpd", "system", "groups.json")
var httpdCardsFile string = filepath.Join(_var, "httpd", "system", "cards.json")
var httpdEventsFile string = filepath.Join(_var, "httpd", "system", "events.json")
var httpdLogsFile string = filepath.Join(_var, "httpd", "system", "logs.json")
var httpdUsersFile string = filepath.Join(_var, "httpd", "system", "users.json")
var httpdAuditFile string = filepath.Join(_var, "httpd", "audit", "audit.log")

func workdir() string {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return `C:\uhppoted`
	}

	return filepath.Join(programData, "uhppoted")
}
