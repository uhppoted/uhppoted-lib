package config

const (
	// DefaultConfig is the default file path for the uhppoted configuration file
	DefaultConfig = "/usr/local/etc/com.github.uhppoted/uhppoted.conf"

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

	httpdAuthDB          string = "/usr/local/etc/com.github.uhppoted/httpd/auth.json"
	httpdCACertificate   string = "/usr/local/etc/com.github.uhppoted/httpd/ca.cert"
	httpdTLSCertificate  string = "/usr/local/etc/com.github.uhppoted/httpd/uhppoted.cert"
	httpdTLSKey          string = "/usr/local/etc/com.github.uhppoted/httpd/uhppoted.key"
	httpdControllersFile string = "/usr/local/var/com.github.uhppoted/httpd/system/controllers.json"
	httpdDoorsFile       string = "/usr/local/var/com.github.uhppoted/httpd/system/doors.json"
	httpdDBFile          string = "/usr/local/var/com.github.uhppoted/httpd/memdb/db.json"
	httpdDBACLRules      string = "/usr/local/etc/com.github.uhppoted/httpd/acl.grl"
	httpdDBSystemRules   string = "/usr/local/etc/com.github.uhppoted/httpd/system.grl"
	httpdDBCardRules     string = "/usr/local/etc/com.github.uhppoted/httpd/cards.grl"
	httpdDBDoorRules     string = "/usr/local/etc/com.github.uhppoted/httpd/doors.grl"
	httpdAuditFile       string = "/usr/local/var/com.github.uhppoted/httpd/audit/audit.log"
)
