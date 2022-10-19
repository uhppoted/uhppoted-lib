package config

const (
	// DefaultConfig is the default file path for the uhppoted configuration file
	DefaultConfig = "/etc/uhppoted/uhppoted.conf"

	_etc string = "/etc/uhppoted"
	_var string = "/var/uhppoted"

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

	httpdAuthDB         string = _etc + "/httpd/auth.json"
	httpdCACertificate  string = _etc + "/httpd/ca.cert"
	httpdTLSCertificate string = _etc + "/httpd/uhppoted.cert"
	httpdTLSKey         string = _etc + "/httpd/uhppoted.key"

	httpdRulesACL         string = _etc + "/httpd/acl.grl"
	httpdRulesInterfaces  string = ""
	httpdRulesControllers string = ""
	httpdRulesCards       string = ""
	httpdRulesDoors       string = ""
	httpdRulesGroups      string = ""
	httpdRulesEvents      string = ""
	httpdRulesLogs        string = ""
	httpdRulesUsers       string = ""

	httpdInterfacesFile  string = _var + "/httpd/system/interfaces.json"
	httpdControllersFile string = _var + "/httpd/system/controllers.json"
	httpdDoorsFile       string = _var + "/httpd/system/doors.json"
	httpdGroupsFile      string = _var + "/httpd/system/groups.json"
	httpdCardsFile       string = _var + "/httpd/system/cards.json"
	httpdEventsFile      string = _var + "/httpd/system/events.json"
	httpdLogsFile        string = _var + "/httpd/system/logs.json"
	httpdUsersFile       string = _var + "/httpd/system/users.json"
	httpdHistoryFile     string = _var + "/httpd/system/history.json"
	httpdAuditFile       string = _var + "/httpd/audit/audit.log"

	SheetsCredentials string = _etc + "/sheets/credentials.json"
)
