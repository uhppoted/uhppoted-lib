package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
)

func TestParseTable(t *testing.T) {
	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	table := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N"},
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	list, warnings, err := ParseTable(&table, devices, true)
	if err != nil {
		t.Fatalf("Unexpected error parsing table: %v", err)
	}

	if len(warnings) != 0 {
		t.Errorf("Returned warnings - expected:\n%+v\ngot:\n%+v\n", 0, warnings)
	}

	if list == nil {
		t.Fatalf("ParseTable returned invalid result: %v", list)
	}

	if !reflect.DeepEqual(*list, expected) {
		t.Errorf("Returned incorrect ACL - expected:\n%+v\ngot:\n%+v\n", expected, *list)
	}
}

func TestParseTableWithMultipleDevices(t *testing.T) {
	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 0, 4: 1}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 1, 4: 1}},
		},
	}

	table := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop", "D1", "D2", "D3", "D4"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N", "Y", "Y", "N", "Y"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y", "Y", "N", "Y", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N", "N", "Y", "Y", "Y"},
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

	list, warnings, err := ParseTable(&table, devices, true)
	if err != nil {
		t.Fatalf("Unexpected error parsing table: %v", err)
	}

	if len(warnings) != 0 {
		t.Errorf("Returned warnings - expected:\n%+v\ngot:\n%+v\n", 0, warnings)
	}

	if list == nil {
		t.Fatalf("ParseTable returned invalid result: %v", list)
	}

	if !reflect.DeepEqual(*list, expected) {
		t.Errorf("Returned incorrect ACL - expected:\n%+v\ngot:\n%+v\n", expected, *list)
	}
}

func TestParseTableWithDuplicateCardNumbers(t *testing.T) {
	expected := ACL{
		12345: map[uint32]types.Card{
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	errors := []error{
		&DuplicateCardError{65537},
	}

	table := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N"},
			[]string{"65537", "2020-01-01", "2020-12-31", "Y", "N", "N", "Y"},
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	list, warnings, err := ParseTable(&table, devices, false)
	if err != nil {
		t.Fatalf("Unexpected error parsing table: %v", err)
	}

	if !reflect.DeepEqual(warnings, errors) {
		t.Errorf("Returned unexpected warnings - expected:\n%+v\ngot:\n%+v\n", errors, warnings)
	}

	if list == nil {
		t.Fatalf("ParseTable returned invalid result: %v", list)
	}

	if !reflect.DeepEqual(*list, expected) {
		t.Errorf("Returned incorrect ACL - expected:\n%+v\ngot:\n%+v\n", expected, *list)
	}
}

func TestParseTableWithDuplicateCardNumbersAndStrict(t *testing.T) {
	table := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N"},
			[]string{"65537", "2020-01-01", "2020-12-31", "Y", "N", "N", "Y"},
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	_, _, err := ParseTable(&table, devices, true)
	if err == nil {
		t.Fatalf("Expected error parsing table with duplicate card numbers and 'strict', got %v", err)
	}
}

func TestMakeTable(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N"},
		},
	}

	d := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	devices := []uhppote.Device{d}

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithTimeProfiles(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 29}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "29"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N"},
		},
	}

	d := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	devices := []uhppote.Device{d}

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithMultipleDevices(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 0, 4: 1}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 1, 4: 1}},
		},
	}

	expected := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop", "D1", "D2", "D3", "D4"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N", "Y", "Y", "N", "Y"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y", "Y", "N", "Y", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N", "N", "Y", "Y", "Y"},
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

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithBlankDoors(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	expected := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Workshop"},
		Records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N"},
		},
	}

	devices := []uhppote.Device{
		uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "", "Workshop"},
		},
	}

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithMissingACL(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	_, err := MakeTable(acl, devices)
	if err == nil {
		t.Fatalf("Expected error creating table")
	}
}

func TestMakeRecordsetWithMismatchedDates(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 0, 4: 1}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-01"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-01-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 1, 4: 1}},
		},
	}

	expected := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop", "D1", "D2", "D3", "D4"},
		Records: [][]string{
			[]string{"65537", "2020-01-01", "2020-12-31", "Y", "N", "N", "N", "Y", "Y", "N", "Y"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y", "Y", "N", "Y", "Y"},
			[]string{"65539", "2020-01-03", "2020-12-31", "N", "N", "N", "N", "N", "Y", "Y", "Y"},
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

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithMismatchedCards(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: types.MustParseDate("2020-01-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 0}},
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-02"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-02-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-03-04"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		54321: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: types.MustParseDate("2020-01-01"), To: types.MustParseDate("2020-12-31"), Doors: map[uint8]uint8{1: 1, 2: 1, 3: 0, 4: 1}},
			65538: types.Card{CardNumber: 65538, From: types.MustParseDate("2020-03-01"), To: types.MustParseDate("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 1, 4: 1}},
			65539: types.Card{CardNumber: 65539, From: types.MustParseDate("2020-01-03"), To: types.MustParseDate("2020-11-30"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 1, 4: 1}},
			65540: types.Card{CardNumber: 65540, From: types.MustParseDate("2019-01-01"), To: types.MustParseDate("2021-12-31"), Doors: map[uint8]uint8{1: 0, 2: 1, 3: 0, 4: 1}},
		},
	}

	expected := Table{
		Header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop", "D1", "D2", "D3", "D4"},
		Records: [][]string{
			[]string{"65536", "2020-01-01", "2020-12-31", "Y", "N", "Y", "N", "N", "N", "N", "N"},
			[]string{"65537", "2020-01-01", "2020-12-31", "Y", "N", "N", "N", "Y", "Y", "N", "Y"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y", "Y", "N", "Y", "Y"},
			[]string{"65539", "2020-01-03", "2020-12-31", "N", "N", "N", "N", "N", "Y", "Y", "Y"},
			[]string{"65540", "2019-01-01", "2021-12-31", "N", "N", "N", "N", "N", "Y", "N", "Y"},
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

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}
