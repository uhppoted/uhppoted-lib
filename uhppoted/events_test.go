package uhppoted

import (
	"fmt"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

func TestGetEventIndices(t *testing.T) {
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-02-10 07:12:01", time.Local)
	index := uint32(17)

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
					Index:        39,
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
					Index:        107,
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

	first, last, current, err := u.GetEventIndices(405419896)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if first != 39 {
		t.Errorf("Incorrect 'FIRST' event - expected:%v, got:%v", 39, first)
	}

	if last != 107 {
		t.Errorf("Incorrect 'LAST' event - expected:%v, got:%v", 107, last)
	}

	if current != 17 {
		t.Errorf("Incorrect 'CURRENT' event - expected:%v, got:%v", 17, current)
	}
}

func TestRecordSpecialEvents(t *testing.T) {
	deviceID := uint32(405419896)
	enabled := true

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

	updated, err := u.RecordSpecialEvents(deviceID, enabled)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if !updated {
		t.Errorf("Incorrect response - expected: %+v, got: %+v", true, updated)
	}
}

func TestRecordSpecialEventsWithFail(t *testing.T) {
	deviceID := uint32(405419896)
	enabled := true

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

	updated, err := u.RecordSpecialEvents(deviceID, enabled)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if updated {
		t.Errorf("Incorrect response: - expected:%+v, got:%+v", false, updated)
	}
}
