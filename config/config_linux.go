package config

const (
	// DefaultConfig is the default file path for the uhppoted configuration file
	DefaultConfig = "/etc/uhppoted/uhppoted.conf"

	restUsers  string = "/etc/uhppoted/rest/users"
	restGroups string = "/etc/uhppoted/rest/groups"
	restHOTP   string = "/etc/uhppoted/rest/counters"

	mqttBrokerCertificate string = "/etc/uhppoted/mqtt/broker.cert"
	mqttClientCertificate string = "/etc/uhppoted/mqtt/client.cert"
	mqttClientKey         string = "/etc/uhppoted/mqtt/client.key"
	mqttUsers             string = "/etc/uhppoted/mqtt.permissions.users"
	mqttGroups            string = "/etc/uhppoted/mqtt.permissions.groups"
	mqttCards             string = "/etc/uhppoted/mqtt/cards"
	hotpSecrets           string = "/etc/uhppoted/mqtt.hotp.secrets"
	rsaKeyDir             string = "/etc/uhppoted/mqtt/rsa"

	eventIDs     string = "/var/uhppoted/mqtt.events.retrieved"
	hotpCounters string = "/var/uhppoted/mqtt.hotp.counters"
	nonceServer  string = "/var/uhppoted/mqtt.nonce"
	nonceClients string = "/var/uhppoted/mqtt.nonce.counters"

	httpdAuthDB          string = "/etc/uhppoted/httpd/auth.json"
	httpdCACertificate   string = "/etc/uhppoted/httpd/ca.cert"
	httpdTLSCertificate  string = "/etc/uhppoted/httpd/uhppoted.cert"
	httpdTLSKey          string = "/etc/uhppoted/httpd/uhppoted.key"
	httpdControllersFile string = "/var/uhppoted/httpd/system/controlers.json"
	httpdDoorsFile       string = "/var/uhppoted/httpd/system/doors.json"
	httpdGroupsFile      string = "/var/uhppoted/httpd/system/groups.json"
	httpdCardsFile       string = "/var/uhppoted/httpd/system/cards.json"
	httpdEventsFile      string = "/var/uhppoted/httpd/system/events.json"
	httpdRulesACL        string = "/var/uhppoted/httpd/acl.grl"
	httpdRulesSystem     string = "/var/uhppoted/httpd/system.grl"
	httpdRulesCards      string = "/var/uhppoted/httpd/cards.grl"
	httpdRulesDoors      string = "/var/uhppoted/httpd/doors.grl"
	httpdRulesGroups     string = "/var/uhppoted/httpd/groups.grl"
	httpdAuditFile       string = "/var/uhppoted/httpd/audit/audit.log"
)
