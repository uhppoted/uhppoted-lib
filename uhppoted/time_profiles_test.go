package uhppoted

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

func pdate(s string) *types.Date {
	d, _ := types.DateFromString(s)

	return &d
}

func hhmm(s string) types.HHmm {
	t, _ := types.HHmmFromString(s)

	return *t
}

func TestGetTimeProfiles(t *testing.T) {
	profile := types.TimeProfile{
		ID:              29,
		LinkedProfileID: 19,
		From:            pdate("2021-04-01"),
		To:              pdate("2021-12-29"),
		Weekdays: types.Weekdays{
			time.Monday:    true,
			time.Tuesday:   true,
			time.Wednesday: false,
			time.Thursday:  true,
			time.Friday:    false,
			time.Saturday:  true,
			time.Sunday:    true,
		},
		Segments: types.Segments{
			1: types.Segment{Start: hhmm("08:30"), End: hhmm("09:45")},
			2: types.Segment{Start: hhmm("11:35"), End: hhmm("13:15")},
			3: types.Segment{Start: hhmm("14:01"), End: hhmm("17:59")},
		},
	}

	request := GetTimeProfilesRequest{
		DeviceID: 405419896,
	}

	expected := GetTimeProfilesResponse{
		DeviceID: 405419896,
		Profiles: []types.TimeProfile{profile},
	}

	mock := stub{
		getTimeProfile: func(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
			if deviceID == 405419896 {
				if profileID == profile.ID {
					return &profile, nil
				}
				return nil, nil
			}

			return nil, fmt.Errorf("Invalid arguments")
		},
	}

	u := UHPPOTED{
		UHPPOTE:         &mock,
		ListenBatchSize: 0,
	}

	response, err := u.GetTimeProfiles(request)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if response == nil {
		t.Fatalf("Invalid response (%v)", response)
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("Incorrected response:\n   expected: %+v\n   got:      %+v\n", expected, *response)
	}
}

func TestGetTimeProfile(t *testing.T) {
	profile := types.TimeProfile{
		ID:              29,
		LinkedProfileID: 19,
		From:            pdate("2021-04-01"),
		To:              pdate("2021-12-29"),
		Weekdays: types.Weekdays{
			time.Monday:    true,
			time.Tuesday:   true,
			time.Wednesday: false,
			time.Thursday:  true,
			time.Friday:    false,
			time.Saturday:  true,
			time.Sunday:    true,
		},
		Segments: types.Segments{
			1: types.Segment{Start: hhmm("08:30"), End: hhmm("09:45")},
			2: types.Segment{Start: hhmm("11:35"), End: hhmm("13:15")},
			3: types.Segment{Start: hhmm("14:01"), End: hhmm("17:59")},
		},
	}

	request := GetTimeProfileRequest{
		DeviceID:  405419896,
		ProfileID: 29,
	}

	expected := GetTimeProfileResponse{
		DeviceID:    405419896,
		TimeProfile: profile,
	}

	mock := stub{
		getTimeProfile: func(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
			if deviceID == 405419896 && profileID == profile.ID {
				return &profile, nil
			}

			return nil, fmt.Errorf("Invalid arguments")
		},
	}

	u := UHPPOTED{
		UHPPOTE:         &mock,
		ListenBatchSize: 0,
	}

	response, err := u.GetTimeProfile(request)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if response == nil {
		t.Fatalf("Invalid response (%v)", response)
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("Incorrected response:\n   expected: %+v\n   got:      %+v\n", expected, *response)
	}
}

func TestSetTimeProfile(t *testing.T) {
	linked := types.TimeProfile{
		ID:   19,
		From: pdate("2021-04-01"),
		To:   pdate("2021-12-29"),
		Weekdays: types.Weekdays{
			time.Monday:    true,
			time.Tuesday:   true,
			time.Wednesday: false,
			time.Thursday:  true,
			time.Friday:    false,
			time.Saturday:  true,
			time.Sunday:    true,
		},
		Segments: types.Segments{
			1: types.Segment{Start: hhmm("08:30"), End: hhmm("09:45")},
			2: types.Segment{Start: hhmm("11:35"), End: hhmm("13:15")},
			3: types.Segment{Start: hhmm("14:01"), End: hhmm("17:59")},
		},
	}

	profile := types.TimeProfile{
		ID:              29,
		LinkedProfileID: 19,
		From:            pdate("2021-04-01"),
		To:              pdate("2021-12-29"),
		Weekdays: types.Weekdays{
			time.Monday:    true,
			time.Tuesday:   true,
			time.Wednesday: false,
			time.Thursday:  true,
			time.Friday:    false,
			time.Saturday:  true,
			time.Sunday:    true,
		},
		Segments: types.Segments{
			1: types.Segment{Start: hhmm("08:30"), End: hhmm("09:45")},
			2: types.Segment{Start: hhmm("11:35"), End: hhmm("13:15")},
			3: types.Segment{Start: hhmm("14:01"), End: hhmm("17:59")},
		},
	}

	request := PutTimeProfileRequest{
		DeviceID:    405419896,
		TimeProfile: profile,
	}

	expected := PutTimeProfileResponse{
		DeviceID:    405419896,
		TimeProfile: profile,
	}

	mock := stub{
		getTimeProfile: func(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
			if deviceID == 405419896 && profileID == linked.ID {
				return &linked, nil
			}

			return nil, fmt.Errorf("Invalid arguments")
		},

		setTimeProfile: func(deviceID uint32, profile types.TimeProfile) (bool, error) {
			if deviceID == 405419896 && profile.ID == 29 {
				return true, nil
			}

			return false, fmt.Errorf("Invalid arguments")
		},
	}

	u := UHPPOTED{
		UHPPOTE:         &mock,
		ListenBatchSize: 0,
	}

	response, err := u.PutTimeProfile(request)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if response == nil {
		t.Fatalf("Invalid response (%v)", response)
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("Incorrected response:\n   expected: %+v\n   got:      %+v\n", expected, *response)
	}
}

func TestClearTimeProfiles(t *testing.T) {
	request := ClearTimeProfilesRequest{
		DeviceID: 405419896,
	}

	expected := ClearTimeProfilesResponse{
		DeviceID: 405419896,
		Cleared:  true,
	}

	mock := stub{
		clearTimeProfiles: func(deviceID uint32) (bool, error) {
			if deviceID == 405419896 {
				return true, nil
			}

			return false, fmt.Errorf("Invalid arguments")
		},
	}

	u := UHPPOTED{
		UHPPOTE:         &mock,
		ListenBatchSize: 0,
	}

	response, err := u.ClearTimeProfiles(request)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if response == nil {
		t.Fatalf("Invalid response (%v)", response)
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("Incorrected response:\n   expected: %+v\n   got:      %+v\n", expected, *response)
	}
}
