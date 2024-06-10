package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
)

func TestRevoke(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == card.CardNumber {
					cards[ix] = card
					return true, nil
				}
			}

			cards = append(cards, card)

			return true, nil
		},
	}

	err := Revoke(&u, devices, 65538, []string{"Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'revoke': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestRevokeWithAmbiguousDoors(t *testing.T) {
	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"Garage", "D2", "D3", "D4"},
		},
	}

	u := mock{}

	err := Revoke(&u, devices, 65538, []string{"Garage"})
	if err == nil {
		t.Fatalf("Expected error invoking 'revoke', got '%v'", err)
	}
}

func TestRevokeWithNewCard(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == card.CardNumber {
					cards[ix] = card
					return true, nil
				}
			}

			cards = append(cards, card)

			return true, nil
		},
	}

	err := Revoke(&u, devices, 65536, []string{"Side Door", "Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'revoke': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Card not added to device internal card list:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestRevokeWithNarrowerDateRange(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	d := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	devices := []uhppote.Device{d}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == card.CardNumber {
					cards[ix] = card
					return true, nil
				}
			}

			cards = append(cards, card)

			return true, nil
		},
	}

	err := Revoke(&u, devices, 65538, []string{"Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'revoke': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestRevokeAcrossMultipleDevices(t *testing.T) {
	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-02-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-02"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-04-03"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-02-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-02"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-04-03"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards[deviceID] {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards[deviceID] {
				if c.CardNumber == card.CardNumber {
					cards[deviceID][ix] = card
					return true, nil
				}
			}

			cards[deviceID] = append(cards[deviceID], card)

			return true, nil
		},
	}

	err := Revoke(&u, devices, 65538, []string{"Garage", "D2"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'revoke': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestRevokeAll(t *testing.T) {
	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-02-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-02"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-04-03"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-02-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-02"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-04-03"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards[deviceID] {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards[deviceID] {
				if c.CardNumber == card.CardNumber {
					cards[deviceID][ix] = card
					return true, nil
				}
			}

			cards[deviceID] = append(cards[deviceID], card)

			return true, nil
		},
	}

	err := Revoke(&u, devices, 65538, []string{"ALL"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'revoke': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestRevokeWithInvalidDoor(t *testing.T) {
	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-02-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-02"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-04-03"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2020-02-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-02"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2020-04-03"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards[deviceID] {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards[deviceID] {
				if c.CardNumber == card.CardNumber {
					cards[deviceID][ix] = card
					return true, nil
				}
			}

			cards[deviceID] = append(cards[deviceID], card)

			return true, nil
		},
	}

	err := Revoke(&u, devices, 65538, []string{"GarageX", "D2"})
	if err == nil {
		t.Errorf("Expected error invoking 'revoke' with invalid door name")
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}
