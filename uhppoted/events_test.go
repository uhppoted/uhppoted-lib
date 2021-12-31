package uhppoted

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

func TestGetEvents(t *testing.T) {
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-02-10 07:12:01", time.Local)
	index := uint32(17)

	events := []Event{
		Event{
			DeviceID:   405419896,
			Index:      18,
			Type:       2,
			Granted:    true,
			Door:       2,
			Direction:  1,
			CardNumber: 6154413,
			Timestamp:  types.DateTime(timestamp),
			Reason:     6,
		},
		Event{
			DeviceID:   405419896,
			Index:      19,
			Type:       2,
			Granted:    false,
			Door:       2,
			Direction:  2,
			CardNumber: 6154414,
			Timestamp:  types.DateTime(timestamp),
			Reason:     15,
		},
	}

	request := GetEventsRequest{
		DeviceID: 405419896,
		Max:      2,
	}

	expected := GetEventsResponse{
		DeviceID: 405419896,
		Events:   events,
	}

	mock := stub{
		getEventIndex: func(deviceID uint32) (*types.EventIndex, error) {
			if deviceID == 405419896 {
				return &types.EventIndex{
					SerialNumber: 405419896,
					Index:        index,
				}, nil
			}

			return nil, fmt.Errorf("Invalid arguments")
		},

		setEventIndex: func(deviceID, ix uint32) (*types.EventIndexResult, error) {
			if deviceID == 405419896 {
				index = ix
				return &types.EventIndexResult{
					SerialNumber: types.SerialNumber(deviceID),
					Index:        index,
					Changed:      true,
				}, nil
			}

			return nil, fmt.Errorf("Invalid arguments")
		},

		getEvent: func(deviceID, index uint32) (*types.Event, error) {
			switch {
			case deviceID == 405419896 && index == 0:
				return &types.Event{
					SerialNumber: 405419896,
					Index:        18,
					Type:         2,
					Granted:      true,
					Door:         2,
					Direction:    1,
					CardNumber:   6154413,
					Timestamp:    types.DateTime(timestamp),
					Reason:       6,
				}, nil

			case deviceID == 405419896 && index == 0xffffffff:
				return &types.Event{
					SerialNumber: 405419896,
					Index:        19,
					Type:         2,
					Granted:      false,
					Door:         2,
					Direction:    2,
					CardNumber:   6154414,
					Timestamp:    types.DateTime(timestamp),
					Reason:       15,
				}, nil

			case deviceID == 405419896 && index == 18:
				return &types.Event{
					SerialNumber: 405419896,
					Index:        18,
					Type:         2,
					Granted:      true,
					Door:         2,
					Direction:    1,
					CardNumber:   6154413,
					Timestamp:    types.DateTime(timestamp),
					Reason:       6,
				}, nil
			case deviceID == 405419896 && index == 19:
				return &types.Event{
					SerialNumber: 405419896,
					Index:        19,
					Type:         2,
					Granted:      false,
					Door:         2,
					Direction:    2,
					CardNumber:   6154414,
					Timestamp:    types.DateTime(timestamp),
					Reason:       15,
				}, nil
			}

			return nil, nil
		},
	}

	u := UHPPOTED{
		UHPPOTE:         &mock,
		ListenBatchSize: 0,
		Log:             nil,
	}

	response, err := u.GetEvents(request)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if response == nil {
		t.Fatalf("Invalid response (%v)", response)
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("Incorrect response:\n   expected: %+v\n   got:      %+v\n", expected, *response)
	}

	if index != uint32(19) {
		t.Errorf("Failed to update controller event index - expected: %v\n   got:      %v\n", 19, index)
	}
}

func TestRecordSpecialEvents(t *testing.T) {
	request := RecordSpecialEventsRequest{
		DeviceID: 405419896,
		Enable:   true,
	}

	expected := RecordSpecialEventsResponse{
		DeviceID: 405419896,
		Enable:   true,
		Updated:  true,
	}

	mock := stub{
		recordSpecialEvents: func(deviceID uint32, enable bool) (bool, error) {
			if deviceID == 405419896 && enable == true {
				return true, nil
			}

			return false, fmt.Errorf("Invalid arguments")
		},
	}

	u := UHPPOTED{
		UHPPOTE:         &mock,
		ListenBatchSize: 0,
		Log:             nil,
	}

	response, err := u.RecordSpecialEvents(request)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if response == nil {
		t.Fatalf("Invalid response (%v)", response)
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("Incorrect response:\n   expected: %+v\n   got:      %+v\n", expected, *response)
	}
}

func TestRecordSpecialEventsWithFail(t *testing.T) {
	request := RecordSpecialEventsRequest{
		DeviceID: 405419896,
		Enable:   true,
	}

	expected := RecordSpecialEventsResponse{
		DeviceID: 405419896,
		Enable:   true,
		Updated:  false,
	}

	mock := stub{
		recordSpecialEvents: func(deviceID uint32, enable bool) (bool, error) {
			return false, nil
		},
	}

	u := UHPPOTED{
		UHPPOTE:         &mock,
		ListenBatchSize: 0,
		Log:             nil,
	}

	response, err := u.RecordSpecialEvents(request)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if response == nil {
		t.Fatalf("Invalid response (%v)", response)
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("Incorrect response:\n   expected: %+v\n   got:      %+v\n", expected, *response)
	}
}
