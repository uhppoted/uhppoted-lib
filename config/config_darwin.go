package config

const (
	// DefaultConfig is the default file path for the uhppoted configuration file
	DefaultConfig = "/usr/local/etc/com.github.uhppoted/uhppoted.conf"

	_etc string = "/usr/local/etc/com.github.uhppoted"
	_var string = "/usr/local/var/com.github.uhppoted"

	restUsers  string = "/usr/local/etc/com.github.uhppoted/rest/users"
	restGroups string = "/usr/local/etc/com.github.uhppoted/rest/groups"
	restHOTP   string = "/usr/local/etc/com.github.uhppoted/rest/counters"

	mqttBrokerCertificate string = "/usr/local/etc/com.github.uhppoted/mqtt/broker.cert"
	mqttClientCertificate string = "/usr/local/etc/com.github.uhppoted/mqtt/client.cert"
	mqttClientKey         string = "/usr/local/etc/com.github.uhppoted/mqtt/client.key"
	mqttUsers             string = "/usr/local/etc/com.github.uhppoted/mqtt.permissions.users"
	mqttGroups            string = "/usr/local/etc/com.github.uhppoted/mqtt.permissions.groups"
	mqttCards             string = "/usr/local/etc/com.github.uhppoted/mqtt/cards"
	hotpSecrets           string = "/usr/local/etc/com.github.uhppoted/mqtt.hotp.secrets"
	rsaKeyDir             string = "/usr/local/etc/com.github.uhppoted/mqtt/rsa"

	eventIDs     string = "/usr/local/var/com.github.uhppoted/mqtt.events.retrieved"
	hotpCounters string = "/usr/local/var/com.github.uhppoted/mqtt.hotp.counters"
	nonceServer  string = "/usr/local/var/com.github.uhppoted/mqtt.nonce"
	nonceClients string = "/usr/local/var/com.github.uhppoted/mqtt.nonce.counters"

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
	httpdAuditFile       string = _var + "/httpd/audit/audit.log"
)
