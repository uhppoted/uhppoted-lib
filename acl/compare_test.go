package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 0, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 0, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
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

func TestCompareWithTimeProfiles(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 29, 3: 0, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 29, 3: 0, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
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
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 0, 4: 0}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 0, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
		},
		54321: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]int{1: 0, 2: 0, 3: 1, 4: 0}},
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
