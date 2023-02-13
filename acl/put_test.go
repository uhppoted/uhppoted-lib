package acl

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

func TestPutACL(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1234},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1234},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
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
		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == cardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithTimeProfiles(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 29, 4: 0}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 29, 4: 0}, PIN: 4321},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
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
		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == cardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
		getTimeProfile: func(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
			if profileID == 29 {
				return &types.TimeProfile{}, nil
			}

			return nil, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithInvalidTimeProfile(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 55, 4: 0}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{65538},
			Errors:    []error{fmt.Errorf("time profile 55 is not defined for 12345")},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
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
		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == cardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
		getTimeProfile: func(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
			if profileID == 29 {
				return &types.TimeProfile{}, nil
			}

			return nil, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLDryRun(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
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
		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == cardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, true)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithMultipleDevices(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}, PIN: 4321},
			types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 4322},
		},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},

		54321: Report{
			Unchanged: []uint32{65536},
			Updated:   []uint32{65537, 65538},
			Added:     []uint32{},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}, PIN: 4321},
			types.Card{CardNumber: 65539, From: date("2020-01-01"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards[deviceID])), nil
		},
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards[deviceID] {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards[deviceID]) {
				return nil, nil
			}
			return &cards[deviceID][index-1], nil
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

		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards[deviceID] {
				if c.CardNumber == cardNumber {
					cards[deviceID] = append(cards[deviceID][:ix], cards[deviceID][ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		if len(cards) != len(expected) {
			t.Errorf("Internal card lists not updated correctly - expected:%v devices, got:%v devices", len(expected), len(cards))
		} else {
			for k, p := range expected {
				q := cards[k]
				if !reflect.DeepEqual(p, q) {
					t.Errorf("Device %v: internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", k, p, q)
				}
			}
		}
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}
func TestPutACLWithConcurrency(t *testing.T) {
	delays := map[uint32]time.Duration{
		12345: 500 * time.Millisecond,
		54321: 1500 * time.Millisecond,
	}

	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}, PIN: 4321},
			types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 4322},
		},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},

		54321: Report{
			Unchanged: []uint32{65536},
			Updated:   []uint32{65537, 65538},
			Added:     []uint32{},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := struct {
		cards map[uint32][]types.Card
		sync.RWMutex
	}{
		cards: map[uint32][]types.Card{
			12345: []types.Card{
				types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
				types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}, PIN: 4321},
				types.Card{CardNumber: 65539, From: date("2020-01-01"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}},
			},

			54321: []types.Card{
				types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
				types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1222},
				types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 4322},
				types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			cards.RLock()
			defer cards.RUnlock()

			return uint32(len(cards.cards[deviceID])), nil
		},

		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards.cards[deviceID] {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},

		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			cards.RLock()
			defer cards.RUnlock()

			if int(index) < 0 || int(index) > len(cards.cards[deviceID]) {
				return nil, nil
			}
			return &cards.cards[deviceID][index-1], nil
		},

		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			cards.RLock()
			for ix, c := range cards.cards[deviceID] {
				if c.CardNumber == card.CardNumber {
					cards.cards[deviceID][ix] = card
					cards.RUnlock()
					return true, nil
				}
			}
			cards.RUnlock()

			cards.Lock()
			cards.cards[deviceID] = append(cards.cards[deviceID], card)
			cards.Unlock()

			time.Sleep(delays[deviceID])

			return true, nil
		},

		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			cards.Lock()
			defer cards.Unlock()

			for ix, c := range cards.cards[deviceID] {
				if c.CardNumber == cardNumber {
					cards.cards[deviceID] = append(cards.cards[deviceID][:ix], cards.cards[deviceID][ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	cards.RLock()
	defer cards.RUnlock()
	if !reflect.DeepEqual(cards.cards, expected) {
		if len(cards.cards) != len(expected) {
			t.Errorf("Internal card lists not updated correctly - expected:%v devices, got:%v devices", len(expected), len(cards.cards))
		} else {
			for k, p := range expected {
				q := cards.cards[k]
				if !reflect.DeepEqual(p, q) {
					t.Errorf("Device %v: internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", k, p, q)
				}
			}
		}
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithFailures(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{65538},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			if card.CardNumber == 65538 {
				return false, nil
			}

			for ix, c := range cards {
				if c.CardNumber == card.CardNumber {
					cards[ix] = card
					return true, nil
				}
			}

			cards = append(cards, card)

			return true, nil
		},
		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == cardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithConcurrentErrors(t *testing.T) {
	errors := []error{fmt.Errorf("RANDOM")}

	delays := map[uint32]time.Duration{
		12345: 500 * time.Millisecond,
		54321: 1500 * time.Millisecond,
	}

	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 1}, PIN: 4321},
			types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},

		54321: Report{
			Unchanged: []uint32{},
			Updated:   []uint32{},
			Added:     []uint32{},
			Deleted:   []uint32{},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}, PIN: 4321},
			types.Card{CardNumber: 65539, From: date("2020-01-01"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 1, 4: 1}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1222},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 0}, PIN: 4322},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			time.Sleep(delays[deviceID])

			if deviceID == 54321 {
				return uint32(len(cards[deviceID])), errors[0]
			} else {
				return uint32(len(cards[deviceID])), nil
			}
		},
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards[deviceID] {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards[deviceID]) {
				return nil, nil
			}
			return &cards[deviceID][index-1], nil
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

		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards[deviceID] {
				if c.CardNumber == cardNumber {
					cards[deviceID] = append(cards[deviceID][:ix], cards[deviceID][ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if !reflect.DeepEqual(err, errors) {
		t.Errorf("Expected errors putting ACL - expected:%v, got:%v", errors, err)
	}

	if !reflect.DeepEqual(cards, expected) {
		if len(cards) != len(expected) {
			t.Errorf("Internal card lists not updated correctly - expected:%v devices, got:%v devices", len(expected), len(cards))
		} else {
			for k, p := range expected {
				q := cards[k]
				if !reflect.DeepEqual(p, q) {
					t.Errorf("Device %v: internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", k, p, q)
				}
			}
		}
	}

	if !reflect.DeepEqual(rpt, report) {
		if len(rpt) != len(report) {
			t.Errorf("Returned report does not match expected - expected:%v devices, got:%v", len(report), len(rpt))
		} else {
			for k, p := range report {
				q := rpt[k]
				if !reflect.DeepEqual(p, q) {
					t.Errorf("Device %v report does not match expected:\n    expected:%+v\n    got:     %+v", k, p, q)
				}
			}
		}
	}
}

func TestPutACLWithErrors(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{65538},
			Errors:    []error{fmt.Errorf("Mysterious error updating card %v", 65538)},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}, PIN: 1221},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}, PIN: 4321},
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
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			if card.CardNumber == 65538 {
				return false, fmt.Errorf("Mysterious error updating card %v", card.CardNumber)
			}

			for ix, c := range cards {
				if c.CardNumber == card.CardNumber {
					cards[ix] = card
					return true, nil
				}
			}

			cards = append(cards, card)

			return true, nil
		},
		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == cardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithNoCurrentPermissions(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1221},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{},
			Updated:   []uint32{65537},
			Added:     []uint32{},
			Deleted:   []uint32{},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}, PIN: 1221},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByID: func(deviceID, cardID uint32) (*types.Card, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
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
		deleteCard: func(deviceID uint32, cardNumber uint32) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == cardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if len(err) > 0 {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}
