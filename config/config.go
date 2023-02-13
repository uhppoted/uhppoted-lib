package config

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/encoding/conf"
	"github.com/uhppoted/uhppoted-lib/monitoring"
)

type DeviceMap map[uint32]*Device

type Device struct {
	Name     string
	Address  *net.UDPAddr
	Doors    []string
	TimeZone string
}

type kv struct {
	Key       string
	Value     interface{}
	IsDefault bool
}

const pretty = `# SYSTEM{{range .system}}
{{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}

# REST{{range .rest}}
{{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}

# MQTT{{range .mqtt}}
{{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}

# AWS{{range .aws}}
{{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}

# HTTPD{{range .httpd}}
{{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}

# Wild Apricot{{range .wildapricot}}
{{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}

# OPEN API{{range .openapi}}
{{if .IsDefault}}# {{end}}{{.Key}} = {{.Value}}{{end}}

# DEVICES{{range $id,$device := .devices}}
UT0311-L0x.{{$id}}.name = {{$device.Name}}{{if $device.Address}}
UT0311-L0x.{{$id}}.address = {{$device.Address}}{{end}}
UT0311-L0x.{{$id}}.door.1 = {{index $device.Doors 0}}
UT0311-L0x.{{$id}}.door.2 = {{index $device.Doors 1}}
UT0311-L0x.{{$id}}.door.3 = {{index $device.Doors 2}}
UT0311-L0x.{{$id}}.door.4 = {{index $device.Doors 3}}
UT0311-L0x.{{$id}}.timezone = {{$device.TimeZone}}
{{else}}
# Example configuration for UTO311-L04 with serial number 405419896
# UT0311-L0x.405419896.name = D405419896
# UT0311-L0x.405419896.address = 192.168.1.100:60000
# UT0311-L0x.405419896.door.1 = Front Door
# UT0311-L0x.405419896.door.2 = Side Door
# UT0311-L0x.405419896.door.3 = Garage
# UT0311-L0x.405419896.door.4 = Workshop
# UT0311-L0x.405419896.timezone = UTC+2
{{end}}`

// const dump = `# SYSTEM{{range .system}}
// {{.Key}} = {{.Value}}{{end}}
//
// # REST{{range .rest}}
// {{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}
//
// # MQTT{{range .mqtt}}
// {{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}
//
// # AWS{{range .aws}}
// {{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}
//
// # HTTPD{{range .httpd}}
// {{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}
//
// # WildApricot{{range .wildapricot}}
// {{if .IsDefault}}; {{end}}{{.Key}} = {{.Value}}{{end}}
//
// # OPEN API{{range .openapi}}
// {{if .IsDefault}}# {{end}}{{.Key}} = {{.Value}}{{end}}
//
// # DEVICES{{range $id,$device := .devices}}
// UT0311-L0x.{{$id}}.name = {{$device.Name}}
// UT0311-L0x.{{$id}}.address = {{$device.Address}}
// UT0311-L0x.{{$id}}.door.1 = {{index $device.Doors 0}}
// UT0311-L0x.{{$id}}.door.2 = {{index $device.Doors 1}}
// UT0311-L0x.{{$id}}.door.3 = {{index $device.Doors 2}}
// UT0311-L0x.{{$id}}.door.4 = {{index $device.Doors 3}}
// {{else}}
// # Example configuration for UTO311-L04 with serial number 405419896
// # UT0311-L0x.405419896.name = D405419896
// # UT0311-L0x.405419896.address = 192.168.1.100:60000
// # UT0311-L0x.405419896.door.1 = Front Door
// # UT0311-L0x.405419896.door.2 = Side Door
// # UT0311-L0x.405419896.door.3 = Garage
// # UT0311-L0x.405419896.door.4 = Workshop
// {{end}}`

type Config struct {
	System      `conf:""`
	Devices     DeviceMap `conf:"/^UT0311-L0x\\.([0-9]+)\\.(.*)/"`
	REST        `conf:"rest"`
	MQTT        `conf:"mqtt"`
	AWS         `conf:"aws"`
	HTTPD       `conf:"httpd"`
	WildApricot `conf:"wild-apricot"`
	OpenAPI     `conf:"openapi"`
}

type System struct {
	BindAddress         *types.BindAddr      `conf:"bind.address"`
	BroadcastAddress    *types.BroadcastAddr `conf:"broadcast.address"`
	ListenAddress       *types.ListenAddr    `conf:"listen.address"`
	Timeout             time.Duration        `conf:"timeout"`
	HealthCheckInterval time.Duration        `conf:"monitoring.healthcheck.interval"`
	HealthCheckIdle     time.Duration        `conf:"monitoring.healthcheck.idle"`
	HealthCheckIgnore   time.Duration        `conf:"monitoring.healthcheck.ignore"`
	WatchdogInterval    time.Duration        `conf:"monitoring.watchdog.interval"`
}

type Lockfile struct {
	File   string `conf:"file"`
	Remove bool   `conf:"remove"`
}

func NewConfig() *Config {
	bind, broadcast, listen := DefaultIpAddresses()

	c := Config{
		System: System{
			BindAddress:         &bind,
			BroadcastAddress:    &broadcast,
			ListenAddress:       &listen,
			Timeout:             2500 * time.Millisecond,
			HealthCheckInterval: 15 * time.Second,
			HealthCheckIdle:     monitoring.IDLE,
			HealthCheckIgnore:   monitoring.IGNORE,
			WatchdogInterval:    5 * time.Second,
		},
		REST:        *NewREST(),
		MQTT:        *NewMQTT(),
		AWS:         *NewAWS(),
		HTTPD:       *NewHTTPD(),
		WildApricot: *NewWildApricot(),
		OpenAPI:     *NewOpenAPI(),
		Devices:     make(DeviceMap, 0),
	}

	return &c
}

func (c *Config) Load(path string) error {
	if path == "" {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := c.Read(f); err != nil {
		return err
	}

	// generate random 'temporary' HMAC key just to avoid defaulting to ""
	if c.MQTT.HMAC.Key == "" {
		hmac := make([]byte, 16)
		if _, err := rand.Read(hmac); err != nil {
			return err
		}

		c.MQTT.HMAC.Key = fmt.Sprintf("%032v", hex.EncodeToString(hmac))
	}

	return nil
}

func (c *Config) Validate() error {
	if c != nil {
		// validate bind.address port
		port := c.System.BindAddress.Port

		if port == 60000 {
			return fmt.Errorf("port %v is not a valid port for bind.address", port)
		}

		if port != 0 && port == c.System.BroadcastAddress.Port {
			return fmt.Errorf("bind.address port (%v) must not be the same as the broadcast.address port", port)
		}

		if port != 0 && port == c.System.ListenAddress.Port {
			return fmt.Errorf("bind.address port (%v) must not be the same as the listen.address port", port)
		}

		// validate broadcast.address port
		if c.System.BroadcastAddress.Port == 0 {
			return fmt.Errorf("port %v is not a valid port for broadcast.address", c.System.BroadcastAddress.Port)
		}

		// validate listen.address port
		if c.System.ListenAddress.Port == 0 {
			return fmt.Errorf("port %v is not a valid port for listen.address", c.System.ListenAddress.Port)
		}

		// check for duplicate doors
		doors := make(map[string]bool)
		for _, device := range c.Devices {
			for _, door := range device.Doors {
				d := strings.ReplaceAll(strings.ToLower(door), " ", "")

				if d != "" && doors[d] {
					return fmt.Errorf("door '%s' is defined more than once in configuration", door)
				}

				doors[d] = true
			}
		}
	}

	return nil
}

func (c *Config) Read(r io.Reader) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return conf.Unmarshal(bytes, c)
}

func (c *Config) Write(w io.Writer) error {
	defc := NewConfig()
	defv := map[string][]kv{
		"system":      listify("", &defc.System),
		"rest":        listify("rest.", &defc.REST),
		"mqtt":        listify("mqtt.", &defc.MQTT),
		"aws":         listify("aws.", &defc.AWS),
		"httpd":       listify("httpd.", &defc.HTTPD),
		"wildapricot": listify("wild-apricot.", &defc.WildApricot),
		"openapi":     listify("openapi.", &defc.OpenAPI),
	}

	config := map[string]interface{}{
		"system":      listify("", &c.System),
		"rest":        listify("rest.", &c.REST),
		"mqtt":        listify("mqtt.", &c.MQTT),
		"aws":         listify("aws.", &c.AWS),
		"httpd":       listify("httpd.", &c.HTTPD),
		"wildapricot": listify("wild-apricot.", &c.WildApricot),
		"openapi":     listify("openapi.", &c.OpenAPI),
		"devices":     c.Devices,
	}

	for k, l := range defv {
		list := config[k].([]kv)
		for i, v := range list {
			if v == l[i] {
				list[i].IsDefault = true
			}
		}
	}

	return template.Must(template.New("uhppoted.conf").Parse(pretty)).Execute(w, config)
}

func listify(parent string, s interface{}) []kv {
	list := []kv{}

	g := func(tag string, v interface{}) bool {
		list = append(list, kv{parent + tag, fmt.Sprintf("%v", v), false})
		return true
	}

	conf.Range(s, g)

	return list
}

// Ref. https://stackoverflow.com/questions/23529663/how-to-get-all-addresses-and-masks-from-local-interfaces-in-go
func DefaultIpAddresses() (types.BindAddr, types.BroadcastAddr, types.ListenAddr) {
	bind := types.BindAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 0,
		Zone: "",
	}

	broadcast := types.BroadcastAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60000,
		Zone: "",
	}

	listen := types.ListenAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60001,
		Zone: "",
	}

	copy(bind.IP, net.IPv4zero)
	copy(broadcast.IP, net.IPv4bcast)
	copy(listen.IP, net.IPv4zero)

	if ifaces, err := net.Interfaces(); err == nil {
	loop:
		for _, i := range ifaces {
			if addrs, err := i.Addrs(); err == nil {
				for _, a := range addrs {
					switch v := a.(type) {
					case *net.IPNet:
						if v.IP.To4() != nil && i.Flags&net.FlagLoopback == 0 {
							copy(bind.IP, v.IP.To4())
							copy(listen.IP, v.IP.To4())
							if i.Flags&net.FlagBroadcast != 0 {
								addr := v.IP.To4()
								mask := v.Mask
								binary.BigEndian.PutUint32(broadcast.IP, binary.BigEndian.Uint32(addr)|^binary.BigEndian.Uint32(mask))
							}
							break loop
						}
					}
				}
			}
		}
	}

	return bind, broadcast, listen
}

func (f DeviceMap) MarshalConf(tag string) ([]byte, error) {
	var s strings.Builder

	if len(f) > 0 {
		fmt.Fprintf(&s, "# DEVICES\n")
		for id, device := range f {
			fmt.Fprintf(&s, "UTO311-L0x.%d.name = %s\n", id, device.Name)

			if device.Address != nil {
				fmt.Fprintf(&s, "UTO311-L0x.%d.address = %s\n", id, device.Address)
			}

			for d, door := range device.Doors {
				fmt.Fprintf(&s, "UTO311-L0x.%d.door.%d = %s\n", id, d+1, door)
			}
			fmt.Fprintf(&s, "\n")
		}
	}

	return []byte(s.String()), nil
}

func (f *DeviceMap) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(`^/(.*?)/$`)
	match := re.FindStringSubmatch(tag)
	if len(match) < 2 {
		return f, fmt.Errorf("invalid 'conf' regular expression tag: %s", tag)
	}

	re, err := regexp.Compile(match[1])
	if err != nil {
		return f, err
	}

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) > 1 {
			id, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return f, fmt.Errorf("invalid 'testMap' key %s: %v", key, err)
			}

			d, ok := (*f)[uint32(id)]
			if !ok || d == nil {
				d = &Device{
					Doors: make([]string, 4),
				}

				(*f)[uint32(id)] = d
			}

			switch match[2] {
			case "name":
				d.Name = value

			case "address":
				address, err := net.ResolveUDPAddr("udp", value)
				if err != nil {
					return f, fmt.Errorf("device %v, invalid address '%s': %v", id, value, err)
				} else {
					d.Address = &net.UDPAddr{
						IP:   make(net.IP, net.IPv4len),
						Port: address.Port,
						Zone: "",
					}

					copy(d.Address.IP, address.IP.To4())
				}

			case "door.1":
				d.Doors[0] = value

			case "door.2":
				d.Doors[1] = value

			case "door.3":
				d.Doors[2] = value

			case "door.4":
				d.Doors[3] = value

			case "timezone":
				d.TimeZone = value
			}
		}
	}

	return f, nil
}

// func resolve(v string) (*net.UDPAddr, error) {
// 	address, err := net.ResolveUDPAddr("udp", v)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	addr := net.UDPAddr{
// 		IP:   make(net.IP, net.IPv4len),
// 		Port: address.Port,
// 		Zone: "",
// 	}
//
// 	copy(addr.IP, address.IP.To4())
//
// 	return &addr, nil
// }
