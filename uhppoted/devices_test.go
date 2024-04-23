package uhppoted

import (
	"encoding/json"
	"net"
	"testing"
)

func TestDeviceSummaryToJSON(t *testing.T) {
	expected := `{"device-type":"UTO311-L04","ip-address":"192.168.1.100","port":60000}`

	device := DeviceSummary{
		DeviceType: "UTO311-L04",
		Address:    net.ParseIP("192.168.1.100"),
		Port:       60000,
	}

	bytes, err := json.Marshal(device)

	if err != nil {
		t.Fatalf("error marshalling DeviceSummary to JSON (%v)", err)
	}

	if string(bytes) != expected {
		t.Fatalf("incorrectly marshalled DeviceSummary\nexpected: %v\ngot:     %v", expected, string(bytes))
	}
}
