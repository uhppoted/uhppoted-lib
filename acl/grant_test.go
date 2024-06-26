package acl

import (
	"reflect"
	"testing"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func TestGrant(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-01-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 1}, PIN: 5432},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	err := Grant(&u, devices, 65538, types.MustParseDate("2023-01-01"), types.MustParseDate("2023-12-31"), 0, []string{"Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'grant': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestGrantWithTimeProfile(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-01-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 29, 4: 1}, PIN: 5432},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

		getTimeProfile: func(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
			if profileID == 29 {
				return &types.TimeProfile{}, nil
			}

			return nil, nil
		},
	}

	err := Grant(&u, devices, 65538, types.MustParseDate("2023-01-01"), types.MustParseDate("2023-12-31"), 29, []string{"Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'grant': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestGrantWithUndefinedTimeProfile(t *testing.T) {
	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

		getTimeProfile: func(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
			if profileID == 29 {
				return &types.TimeProfile{}, nil
			}

			return nil, nil
		},
	}

	if err := Grant(&u, devices, 65538, types.MustParseDate("2023-01-01"), types.MustParseDate("2023-12-31"), 55, []string{"Garage"}); err == nil {
		t.Fatalf("Expected error invoking 'grant' with undefined time profile, got: %v", err)
	}
}

func TestGrantWithAmbiguousDoors(t *testing.T) {
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

	err := Grant(&u, devices, 65538, types.MustParseDate("2023-01-01"), types.MustParseDate("2023-12-31"), 0, []string{"Garage"})
	if err == nil {
		t.Fatalf("Expected error invoking 'grant', got '%v'", err)
	}
}

func TestGrantWithNewCard(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 5432},
		types.Card{CardNumber: 65536, From: types.MustParseDate("2023-01-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 1, 4: 0}},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 5432},
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

	err := Grant(&u, devices, 65536, types.MustParseDate("2023-01-01"), types.MustParseDate("2023-12-31"), 0, []string{"Side Door", "Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'grant': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Card not added to device internal card list:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestGrantWithNarrowerDateRange(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 1}, PIN: 5432},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	d := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	devices := []uhppote.Device{d}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
		types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	err := Grant(&u, devices, 65538, types.MustParseDate("2023-04-01"), types.MustParseDate("2023-10-31"), 0, []string{"Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'grant': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestGrantAcrossMultipleDevices(t *testing.T) {
	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-01-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 1}, PIN: 5432},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-02-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-01-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-04-03"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-02-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-03-02"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-04-03"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	err := Grant(&u, devices, 65538, types.MustParseDate("2023-01-01"), types.MustParseDate("2023-12-31"), 0, []string{"Garage", "D2"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'grant': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestGrantAll(t *testing.T) {
	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-03-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}, PIN: 5432},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-02-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-03-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-04-03"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-01-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-02-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-03-02"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-04-03"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	err := Grant(&u, devices, 65538, types.MustParseDate("2023-03-02"), types.MustParseDate("2023-10-31"), 0, []string{"ALL"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'grant ALL': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestGrantWithInvalidDoor(t *testing.T) {
	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-02-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-03-02"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-04-03"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-02"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-02-03"), To: types.MustParseDate("2023-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 5432},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-03-04"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: []types.Card{
			types.Card{CardNumber: 65537, From: types.MustParseDate("2023-02-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: types.MustParseDate("2023-03-02"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
			types.Card{CardNumber: 65539, From: types.MustParseDate("2023-04-03"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	err := Grant(&u, devices, 65538, types.MustParseDate("2023-01-01"), types.MustParseDate("2023-12-31"), 0, []string{"Garage", "D2X"})
	if err == nil {
		t.Errorf("Expected error invoking 'grant' with invalid door")
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestGrantWithNoCurrentPermissions(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-04-01"), To: types.MustParseDate("2023-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1221},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: types.MustParseDate("2023-01-01"), To: types.MustParseDate("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
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

	err := Grant(&u, devices, 65537, types.MustParseDate("2023-04-01"), types.MustParseDate("2023-10-31"), 0, []string{"Garage"})
	if err != nil {
		t.Fatalf("Unexpected error invoking 'grant': %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}
