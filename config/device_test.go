package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func TestDeviceMarshal(t *testing.T) {
	expected := `# DEVICES
UT0311-L0x.405419896.name = test
UT0311-L0x.405419896.address = 192.168.1.100
UT0311-L0x.405419896.door.1 = Gryffindor
UT0311-L0x.405419896.door.2 = Ravenclaw
UT0311-L0x.405419896.door.3 = Hufflepuff
UT0311-L0x.405419896.door.4 = Slytherin

`
	device := Device{
		Name:     "test",
		Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
		Doors:    []string{"Gryffindor", "Ravenclaw", "Hufflepuff", "Slytherin"},
		TimeZone: "UTC",
	}

	devices := DeviceMap{
		405419896: &device,
	}

	bytes, err := devices.MarshalConf("devices")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(bytes) != expected {
		t.Errorf("Incorrectly marshalled device list\n   expected:%v\n   got:     %v", expected, string(bytes))
	}
}

func TestDeviceMarshalWithNullAddress(t *testing.T) {
	expected := `# DEVICES
UT0311-L0x.405419896.name = test
UT0311-L0x.405419896.door.1 = Gryffindor
UT0311-L0x.405419896.door.2 = Ravenclaw
UT0311-L0x.405419896.door.3 = Hufflepuff
UT0311-L0x.405419896.door.4 = Slytherin

`
	device := Device{
		Name:     "test",
		Doors:    []string{"Gryffindor", "Ravenclaw", "Hufflepuff", "Slytherin"},
		TimeZone: "UTC",
	}

	devices := DeviceMap{
		405419896: &device,
	}

	bytes, err := devices.MarshalConf("devices")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(bytes) != expected {
		t.Errorf("Incorrectly marshalled device list\n   expected:%v\n   got:     %v", expected, string(bytes))
	}
}

func TestDeviceMarshalTCP(t *testing.T) {
	expected := `# DEVICES
UT0311-L0x.405419896.name = test
UT0311-L0x.405419896.address = tcp:192.168.1.100
UT0311-L0x.405419896.door.1 = Gryffindor
UT0311-L0x.405419896.door.2 = Ravenclaw
UT0311-L0x.405419896.door.3 = Hufflepuff
UT0311-L0x.405419896.door.4 = Slytherin

`
	device := Device{
		Name:     "test",
		Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
		Doors:    []string{"Gryffindor", "Ravenclaw", "Hufflepuff", "Slytherin"},
		TimeZone: "UTC",
		Protocol: "tcp",
	}

	devices := DeviceMap{
		405419896: &device,
	}

	bytes, err := devices.MarshalConf("devices")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(bytes) != expected {
		t.Errorf("Incorrectly marshalled device list\n   expected:%v\n   got:     %v", expected, string(bytes))
	}
}

func TestDeviceUnmarshal(t *testing.T) {
	tag := `/^UT0311-L0x\.([0-9]+)\.(.*)/`
	values := map[string]string{
		"UT0311-L0x.405419896.name":     "Alpha",
		"UT0311-L0x.405419896.address":  "192.168.1.100:60000",
		"UT0311-L0x.405419896.door.1":   "Gryffindor",
		"UT0311-L0x.405419896.door.2":   "Hufflepuff",
		"UT0311-L0x.405419896.door.3":   "Ravenclaw",
		"UT0311-L0x.405419896.door.4":   "Slytherin",
		"UT0311-L0x.405419896.timezone": "CEST",
	}

	expected := DeviceMap{
		405419896: &Device{
			Name:     "Alpha",
			Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
			Doors:    []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"},
			TimeZone: "CEST",
			Protocol: "udp",
		},
	}

	devices := DeviceMap{}

	unmarshalled, err := devices.UnmarshalConf(tag, values)
	if err != nil {
		t.Fatalf("error unmarshalling 'conf' controllers section (%v)", err)
	}

	if !reflect.DeepEqual(*unmarshalled.(*DeviceMap), expected) {
		t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %+v\ngot:      %+v", expected, *unmarshalled.(*DeviceMap))

		m := unmarshalled.(*DeviceMap)

		if len(expected) != len(*m) {
			t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %v controllers\ngot:     %v controllers", len(expected), len(*m))
		} else {
			for k, v := range expected {
				u := (*m)[k]
				if !reflect.DeepEqual(u, v) {
					t.Errorf("incorrectly unmarshalled 'conf' controller %v section\nexpected: %+v\ngot:      %+v", k, *v, *u)
				}
			}
		}
	}

	if !reflect.DeepEqual(devices, expected) {
		t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %+v\ngot:      %+v", expected, devices)
	}
}

func TestDeviceUnmarshalUDP(t *testing.T) {
	tag := `/^UT0311-L0x\.([0-9]+)\.(.*)/`
	values := map[string]string{
		"UT0311-L0x.405419896.name":     "Alpha",
		"UT0311-L0x.405419896.address":  "udp:192.168.1.100:60000",
		"UT0311-L0x.405419896.door.1":   "Gryffindor",
		"UT0311-L0x.405419896.door.2":   "Hufflepuff",
		"UT0311-L0x.405419896.door.3":   "Ravenclaw",
		"UT0311-L0x.405419896.door.4":   "Slytherin",
		"UT0311-L0x.405419896.timezone": "CEST",
	}

	expected := DeviceMap{
		405419896: &Device{
			Name:     "Alpha",
			Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
			Doors:    []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"},
			TimeZone: "CEST",
			Protocol: "udp",
		},
	}

	devices := DeviceMap{}

	unmarshalled, err := devices.UnmarshalConf(tag, values)
	if err != nil {
		t.Fatalf("error unmarshalling 'conf' controllers section (%v)", err)
	}

	if !reflect.DeepEqual(*unmarshalled.(*DeviceMap), expected) {
		t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %+v\ngot:      %+v", expected, *unmarshalled.(*DeviceMap))

		m := unmarshalled.(*DeviceMap)

		if len(expected) != len(*m) {
			t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %v controllers\ngot:     %v controllers", len(expected), len(*m))
		} else {
			for k, v := range expected {
				u := (*m)[k]
				if !reflect.DeepEqual(u, v) {
					t.Errorf("incorrectly unmarshalled 'conf' controller %v section\nexpected: %+v\ngot:      %+v", k, *v, *u)
				}
			}
		}
	}

	if !reflect.DeepEqual(devices, expected) {
		t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %+v\ngot:      %+v", expected, devices)
	}
}

func TestDeviceUnmarshalTCP(t *testing.T) {
	tag := `/^UT0311-L0x\.([0-9]+)\.(.*)/`
	values := map[string]string{
		"UT0311-L0x.405419896.name":     "Alpha",
		"UT0311-L0x.405419896.address":  "tcp:192.168.1.100:60000",
		"UT0311-L0x.405419896.door.1":   "Gryffindor",
		"UT0311-L0x.405419896.door.2":   "Hufflepuff",
		"UT0311-L0x.405419896.door.3":   "Ravenclaw",
		"UT0311-L0x.405419896.door.4":   "Slytherin",
		"UT0311-L0x.405419896.timezone": "CEST",
	}

	expected := DeviceMap{
		405419896: &Device{
			Name:     "Alpha",
			Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
			Doors:    []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"},
			TimeZone: "CEST",
			Protocol: "tcp",
		},
	}

	devices := DeviceMap{}

	unmarshalled, err := devices.UnmarshalConf(tag, values)
	if err != nil {
		t.Fatalf("error unmarshalling 'conf' controllers section (%v)", err)
	}

	if !reflect.DeepEqual(*unmarshalled.(*DeviceMap), expected) {
		t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %+v\ngot:      %+v", expected, *unmarshalled.(*DeviceMap))

		m := unmarshalled.(*DeviceMap)

		if len(expected) != len(*m) {
			t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %v controllers\ngot:     %v controllers", len(expected), len(*m))
		} else {
			for k, v := range expected {
				u := (*m)[k]
				if !reflect.DeepEqual(u, v) {
					t.Errorf("incorrectly unmarshalled 'conf' controller %v section\nexpected: %+v\ngot:      %+v", k, *v, *u)
				}
			}
		}
	}

	if !reflect.DeepEqual(devices, expected) {
		t.Errorf("incorrectly unmarshalled 'conf' controllers section\nexpected: %+v\ngot:      %+v", expected, devices)
	}
}

func TestDeviceMapToControllers(t *testing.T) {
	LA, _ := time.LoadLocation("America/Los_Angeles")

	devices := DeviceMap{
		405419896: &Device{
			Name:     "Alpha",
			Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
			Doors:    []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"},
			TimeZone: "America/Los_Angeles",
			Protocol: "tcp",
		},
		303986753: &Device{
			Name:    "Beta",
			Address: types.MustParseControllerAddr("192.168.1.100:60000"),
			Doors:   []string{"Great Hall", "Kitchen", "Dungeon", "Hogsmeade"},
		},
	}

	expected := []uhppote.Device{
		uhppote.Device{
			Name:     "Beta",
			DeviceID: 303986753,
			Doors:    []string{"Great Hall", "Kitchen", "Dungeon", "Hogsmeade"},
			Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
			TimeZone: time.Local,
			Protocol: "udp",
		},
		uhppote.Device{
			Name:     "Alpha",
			DeviceID: 405419896,
			Doors:    []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"},
			Address:  types.MustParseControllerAddr("192.168.1.100:60000"),
			TimeZone: LA,
			Protocol: "tcp",
		},
	}

	controllers := devices.ToControllers()

	if !reflect.DeepEqual(controllers, expected) {
		t.Errorf("invalid controllers list\n   expected: %v\n   got:      %v", expected, controllers)
	}
}
