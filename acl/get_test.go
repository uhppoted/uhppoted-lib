package acl

import (
	"fmt"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
)

func TestGetACL(t *testing.T) {
	errors := []error{}

	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	acl, err := GetACL(&u, devices)
	if !reflect.DeepEqual(err, errors) {
		t.Fatalf("Unexpected error getting ACL: %v", err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	if l, ok := acl[deviceA.DeviceID]; !ok {
		t.Errorf("Missing access list for device ID %v", deviceA.DeviceID)
	} else {
		e := expected[deviceA.DeviceID]
		if len(l) != len(e) {
			t.Errorf("device %v: record counts do not match - expected %d, got %d", deviceA.DeviceID, len(e), len(l))
		}

		for _, card := range e {
			if c, ok := l[card.CardNumber]; !ok {
				t.Errorf("device %v: missing record for card %v", deviceA.DeviceID, card.CardNumber)
			} else if !reflect.DeepEqual(c, card) {
				t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", deviceA.DeviceID, card.CardNumber, card, c)
			}
		}
	}
}

func TestGetACLWithPIN(t *testing.T) {
	errors := []error{}

	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 7531},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 1357},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 0},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 7531},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 1357},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	acl, err := GetACL(&u, devices)
	if !reflect.DeepEqual(err, errors) {
		t.Fatalf("Unexpected error getting ACL: %v", err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	if l, ok := acl[deviceA.DeviceID]; !ok {
		t.Errorf("Missing access list for device ID %v", deviceA.DeviceID)
	} else {
		e := expected[deviceA.DeviceID]
		if len(l) != len(e) {
			t.Errorf("device %v: record counts do not match - expected %d, got %d", deviceA.DeviceID, len(e), len(l))
		}

		for _, card := range e {
			if c, ok := l[card.CardNumber]; !ok {
				t.Errorf("device %v: missing record for card %v", deviceA.DeviceID, card.CardNumber)
			} else if !reflect.DeepEqual(c, card) {
				t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", deviceA.DeviceID, card.CardNumber, card, c)
			}
		}
	}
}

func TestGetACLWithMultipleDevices(t *testing.T) {
	errors := []error{}
	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			switch deviceID {
			case 12345:
				time.Sleep(500 * time.Millisecond)
			case 54321:
				time.Sleep(1500 * time.Millisecond)
			}

			list, ok := cards[deviceID]
			if !ok {
				return 0, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			return uint32(len(list)), nil
		},

		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			list, ok := cards[deviceID]
			if !ok {
				return nil, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			if int(index) < 0 || int(index) > len(list) {
				return nil, nil
			}
			return &list[index-1], nil
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	acl, err := GetACL(&u, devices)
	if !reflect.DeepEqual(err, errors) {
		t.Fatalf("Unexpected error getting ACL: %v", err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	for _, d := range devices {
		if l, ok := acl[d.DeviceID]; !ok {
			t.Errorf("Missing access list for device ID %v", d.DeviceID)
		} else {
			e := expected[d.DeviceID]
			if len(l) != len(e) {
				t.Errorf("device %v: record counts do not match - expected %d, got %d", d.DeviceID, len(e), len(l))
			}

			for _, card := range e {
				if c, ok := l[card.CardNumber]; !ok {
					t.Errorf("device %v: missing record for card %v", d.DeviceID, card.CardNumber)
				} else if !reflect.DeepEqual(c, card) {
					t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", d.DeviceID, card.CardNumber, card, c)
				}
			}
		}
	}
}

func TestGetACLWithMultipleDevicesAndPINs(t *testing.T) {
	errors := []error{}
	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 7531},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 1357},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 0},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}, PIN: 8642},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 2468},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 0},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 0},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 7531},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 1357},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 0},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}, PIN: 8642},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 2468},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 0},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 0},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			switch deviceID {
			case 12345:
				time.Sleep(500 * time.Millisecond)
			case 54321:
				time.Sleep(1500 * time.Millisecond)
			}

			list, ok := cards[deviceID]
			if !ok {
				return 0, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			return uint32(len(list)), nil
		},

		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			list, ok := cards[deviceID]
			if !ok {
				return nil, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			if int(index) < 0 || int(index) > len(list) {
				return nil, nil
			}
			return &list[index-1], nil
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	acl, err := GetACL(&u, devices)
	if !reflect.DeepEqual(err, errors) {
		t.Fatalf("Unexpected error getting ACL: %v", err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	for _, d := range devices {
		if l, ok := acl[d.DeviceID]; !ok {
			t.Errorf("Missing access list for device ID %v", d.DeviceID)
		} else {
			e := expected[d.DeviceID]
			if len(l) != len(e) {
				t.Errorf("device %v: record counts do not match - expected %d, got %d", d.DeviceID, len(e), len(l))
			}

			for _, card := range e {
				if c, ok := l[card.CardNumber]; !ok {
					t.Errorf("device %v: missing record for card %v", d.DeviceID, card.CardNumber)
				} else if !reflect.DeepEqual(c, card) {
					t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", d.DeviceID, card.CardNumber, card, c)
				}
			}
		}
	}
}

func TestGetACLWithDeviceError(t *testing.T) {
	errors := []error{fmt.Errorf("TEST ERROR")}

	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			if deviceID == 54321 {
				return 0, errors[0]
			}

			list, ok := cards[deviceID]
			if !ok {
				return 0, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			return uint32(len(list)), nil
		},

		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			list, ok := cards[deviceID]
			if !ok {
				return nil, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			if int(index) < 0 || int(index) > len(list) {
				return nil, nil
			}
			return &list[index-1], nil
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	acl, err := GetACL(&u, devices)
	if !reflect.DeepEqual(err, errors) {
		t.Fatalf("Expected error getting ACL - expected:%v, got:%v", errors, err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	for _, d := range devices {
		if l, ok := acl[d.DeviceID]; !ok {
			t.Errorf("Missing access list for device ID %v", d.DeviceID)
		} else {
			e := expected[d.DeviceID]
			if len(l) != len(e) {
				t.Errorf("device %v: record counts do not match - expected %d, got %d", d.DeviceID, len(e), len(l))
			}

			for _, card := range e {
				if c, ok := l[card.CardNumber]; !ok {
					t.Errorf("device %v: missing record for card %v", d.DeviceID, card.CardNumber)
				} else if !reflect.DeepEqual(c, card) {
					t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", d.DeviceID, card.CardNumber, card, c)
				}
			}
		}
	}
}
