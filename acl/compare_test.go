package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"reflect"
	"testing"
)

func TestCompareWithoutPIN(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			523321456: types.Card{CardNumber: 523321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 7531},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			523321456: types.Card{CardNumber: 523321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1357},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 523321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 7531},
				types.Card{CardNumber: 923321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
		},
	}

	diff, err := Compare(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}

func TestCompareWithPIN(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			523321456: types.Card{CardNumber: 523321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 7531},
			423321456: types.Card{CardNumber: 423321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 2468},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			523321456: types.Card{CardNumber: 523321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1375},
			423321456: types.Card{CardNumber: 423321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 2468},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 423321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 2468},
				types.Card{CardNumber: 923321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 523321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}, PIN: 1375},
				types.Card{CardNumber: 823321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2023-01-02"), To: date("2023-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
		},
	}

	diff, err := CompareWithPIN(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}

func TestCompareWithTimeProfiles(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 29, 3: 0, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 29, 3: 0, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
		},
	}

	diff, err := Compare(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}

func TestCompareWithMultipleDevices(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
		},
		54321: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 1, 4: 0}},
			},
		},
	}

	diff, err := Compare(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}
