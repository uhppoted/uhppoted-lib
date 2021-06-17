package config

import (
	"net"
	"testing"
)

func TestDeviceMarshall(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.100:60000")

	expected := `# DEVICES
UTO311-L0x.405419896.name = test
UTO311-L0x.405419896.address = 192.168.1.100:60000
UTO311-L0x.405419896.rollover = 10000
UTO311-L0x.405419896.door.1 = Gryffindor
UTO311-L0x.405419896.door.2 = Ravenclaw
UTO311-L0x.405419896.door.3 = Hufflepuff
UTO311-L0x.405419896.door.4 = Slytherin

`
	device := Device{
		Name:     "test",
		Address:  addr,
		Rollover: 10000,
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

func TestDeviceMarshallWithNullAddress(t *testing.T) {
	expected := `# DEVICES
UTO311-L0x.405419896.name = test
UTO311-L0x.405419896.rollover = 10000
UTO311-L0x.405419896.door.1 = Gryffindor
UTO311-L0x.405419896.door.2 = Ravenclaw
UTO311-L0x.405419896.door.3 = Hufflepuff
UTO311-L0x.405419896.door.4 = Slytherin

`
	device := Device{
		Name:     "test",
		Rollover: 10000,
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
