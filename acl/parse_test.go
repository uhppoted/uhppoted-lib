package acl

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func TestParseHeader(t *testing.T) {
	expected := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
		},
		PIN: 0,
	}

	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage"}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	ix, err := parseHeader(header, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing header: %v", err)
	} else if ix == nil {
		t.Fatalf("parseHeader returned 'nil'")
	}

	if !reflect.DeepEqual(*ix, expected) {
		t.Errorf("Invalid index\n   expected: %+v\n   got:      %+v", expected, *ix)
	}
}

func TestParseHeaderWithPIN(t *testing.T) {
	expected := index{
		cardnumber: 1,
		from:       3,
		to:         4,
		doors: map[uint32][]int{
			12345: []int{7, 6, 8, 5},
		},
		PIN: 2,
	}

	header := []string{"Card Number", "PIN", "From", "To", "Workshop", "Side Door", "Front Door", "Garage"}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	ix, err := parseHeader(header, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing header: %v", err)
	} else if ix == nil {
		t.Fatalf("parseHeader returned 'nil'")
	}

	if !reflect.DeepEqual(*ix, expected) {
		t.Errorf("Invalid index\n   expected: %+v\n   got:      %+v", expected, *ix)
	}
}

func TestParseHeaderWithMultipleDevices(t *testing.T) {
	expected := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
			54321: []int{8, 9, 10, 11},
		},
		PIN: 0,
	}

	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage", "D1", "D2", "D3", "D4"}

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

	ix, err := parseHeader(header, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing header: %v", err)
	} else if ix == nil {
		t.Fatalf("parseHeader returned 'nil'")
	}

	if !reflect.DeepEqual(*ix, expected) {
		t.Errorf("Invalid index\n   expected: %+v\n   got:      %+v", expected, *ix)
	}
}

func TestParseHeaderWithMissingColumn(t *testing.T) {
	expected := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
			54321: []int{8, 9, 0, 10},
		},
		PIN: 0,
	}

	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage", "D1", "D2", "D4"}

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

	ix, err := parseHeader(header, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing header: %v", err)
	} else if ix == nil {
		t.Fatalf("parseHeader returned 'nil'")
	}

	if !reflect.DeepEqual(*ix, expected) {
		t.Errorf("Invalid index\n   expected: %+v\n   got:      %+v", expected, *ix)
	}
}

func TestParseHeaderWithInvalidColumn(t *testing.T) {
	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage", "D1", "D2", "D3X", "D4"}

	expected := fmt.Errorf("no configured door matches 'D3X'")

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

	ix, err := parseHeader(header, devices)
	if err == nil {
		t.Fatalf("Expected error parsing header with invalid column: %+v", *ix)
	} else if err.Error() != expected.Error() {
		t.Errorf("Incorrect error message\n   expected: %v\n   got:      %v", expected, err)
	}
}

func TestParseRecord(t *testing.T) {
	ix := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
		},
	}

	record := []string{"8165535", "2021-01-01", "2021-12-31", "Y", "Y", "N", "29", "N", "N", "Y", "Y"}

	expected := map[uint32]types.Card{
		12345: types.Card{
			CardNumber: 8165535,
			From:       types.MustParseDate("2021-01-01"),
			To:         types.MustParseDate("2021-12-31"),
			Doors: map[uint8]uint8{
				1: 0,
				2: 1,
				3: 29,
				4: 1,
			},
		},
	}

	cards, err := parseRecord(record, ix)
	if err != nil {
		t.Fatalf("Unexpected error parsing valid record - %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Incorrect cards list\n   expected: %v\n   got:      %v", expected, cards)
	}
}

func TestParseRecordWithPIN(t *testing.T) {
	ix := index{
		cardnumber: 1,
		from:       3,
		to:         4,
		doors: map[uint32][]int{
			12345: []int{7, 6, 8, 5},
		},
		PIN: 2,
	}

	record := []string{"8165535", "7531", "2021-01-01", "2021-12-31", "Y", "Y", "N", "29", "N", "N", "Y", "Y"}

	expected := map[uint32]types.Card{
		12345: types.Card{
			CardNumber: 8165535,
			From:       types.MustParseDate("2021-01-01"),
			To:         types.MustParseDate("2021-12-31"),
			Doors: map[uint8]uint8{
				1: 0,
				2: 1,
				3: 29,
				4: 1,
			},
			PIN: 7531,
		},
	}

	cards, err := parseRecord(record, ix)
	if err != nil {
		t.Fatalf("Unexpected error parsing valid record - %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Incorrect cards list\n   expected: %v\n   got:      %v", expected, cards)
	}
}

func TestParseRecordWithBlankPIN(t *testing.T) {
	ix := index{
		cardnumber: 1,
		PIN:        2,
		from:       3,
		to:         4,
		doors: map[uint32][]int{
			12345: []int{7, 6, 8, 5},
		},
	}

	record := []string{"8165535", "", "2021-01-01", "2021-12-31", "Y", "Y", "N", "29", "N", "N", "Y", "Y"}

	expected := map[uint32]types.Card{
		12345: types.Card{
			CardNumber: 8165535,
			From:       types.MustParseDate("2021-01-01"),
			To:         types.MustParseDate("2021-12-31"),
			Doors: map[uint8]uint8{
				1: 0,
				2: 1,
				3: 29,
				4: 1,
			},
			PIN: 0,
		},
	}

	cards, err := parseRecord(record, ix)
	if err != nil {
		t.Fatalf("Unexpected error parsing record with blank PIN - %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Incorrect cards list\n   expected: %v\n   got:      %v", expected, cards)
	}
}

func TestParseRecordWithInvalidPermission(t *testing.T) {
	ix := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
		},
	}

	record := []string{"8165535", "2021-01-01", "2021-12-31", "Y", "Y", "X", "29", "N", "N", "Y", "Y"}

	_, err := parseRecord(record, ix)
	if err == nil {
		t.Fatalf("Expected error parsing invalid record, got:%v", err)
	}
}

func TestParseRecordWithInvalidTimeProfile(t *testing.T) {
	ix := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
		},
	}

	record := []string{"8165535", "2021-01-01", "2021-12-31", "Y", "Y", "N", "1", "N", "N", "Y", "Y"}

	_, err := parseRecord(record, ix)
	if err == nil {
		t.Fatalf("Expected error parsing invalid record, got:%v", err)
	}
}
