package uhppoted

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

type Event struct {
	DeviceID   uint32         `json:"device-id"`
	Index      uint32         `json:"event-id"`
	Type       uint8          `json:"event-type"`
	Granted    bool           `json:"access-granted"`
	Door       uint8          `json:"door-id"`
	Direction  uint8          `json:"direction"`
	CardNumber uint32         `json:"card-number"`
	Timestamp  types.DateTime `json:"timestamp"`
	Reason     uint8          `json:"event-reason"`
}

func (e Event) IsZero() bool {
	return e.Index == 0
}

func (u *UHPPOTED) GetEventIndices(deviceID uint32) (uint32, uint32, uint32, error) {
	var first uint32 = 0
	var last uint32 = 0
	var current uint32 = 0

	if v, err := u.UHPPOTE.GetEvent(deviceID, 0); err != nil {
		return 0, 0, 0, err
	} else if v != nil {
		first = v.Index
	}

	if v, err := u.UHPPOTE.GetEvent(deviceID, 0xffffffff); err != nil {
		return 0, 0, 0, err
	} else if v != nil {
		last = v.Index
	}

	if v, err := u.UHPPOTE.GetEventIndex(deviceID); err != nil {
		return 0, 0, 0, err
	} else if v != nil {
		current = v.Index
	}

	return first, last, current, nil
}

func (u *UHPPOTED) GetEvent(deviceID uint32, index uint32) (*Event, error) {
	event, err := u.UHPPOTE.GetEvent(deviceID, index)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, err)
	} else if event == nil {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, fmt.Errorf("%v: no event %v", deviceID, index))
	} else if index != 0 && index != 0xffffffff && event.Index != index {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, fmt.Errorf("%v: no event %v", deviceID, index))
	}

	return &Event{
		DeviceID:   uint32(event.SerialNumber),
		Index:      event.Index,
		Type:       event.Type,
		Granted:    event.Granted,
		Door:       event.Door,
		Direction:  event.Direction,
		CardNumber: event.CardNumber,
		Timestamp:  event.Timestamp,
		Reason:     event.Reason,
	}, nil
}

// Retrieves up to N events subsequent to the 'current' event index (or the 'first' event if the current event index
// is less than the first event index). The on-device index is updated to the index of the last retrieved event.
func (u *UHPPOTED) GetEvents(deviceID uint32, N int) ([]Event, error) {
	var first uint32 = 0
	var current uint32 = 0

	if v, err := u.UHPPOTE.GetEvent(deviceID, 0); err != nil {
		return nil, err
	} else if v != nil {
		first = v.Index
	}

	if v, err := u.UHPPOTE.GetEventIndex(deviceID); err != nil {
		return nil, err
	} else if v != nil {
		current = v.Index
	}

	index := current + 1
	if index < first {
		index = first
	}

	events := []Event{}

	for len(events) < N {
		event, err := u.UHPPOTE.GetEvent(deviceID, index)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		if event == nil {
			break
		}

		events = append(events, Event{
			DeviceID:   uint32(event.SerialNumber),
			Index:      event.Index,
			Type:       event.Type,
			Granted:    event.Granted,
			Door:       event.Door,
			Direction:  event.Direction,
			CardNumber: event.CardNumber,
			Timestamp:  event.Timestamp,
			Reason:     event.Reason,
		})

		current = event.Index
		index++
	}

	response, err := u.UHPPOTE.SetEventIndex(deviceID, current)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, err)
	} else if response == nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("no response to set-event-index %v for %v", current, deviceID))
	} else if response.Index != current {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("failed to update %v event index to %v", deviceID, current))
	}

	return events, nil
}

// Unwraps the request and dispatches the corresponding controller command to enable or disable
// door open, door close and door button press events for the controller.
func (u *UHPPOTED) RecordSpecialEvents(deviceID uint32, enable bool) (bool, error) {
	u.debug("record-special-events", fmt.Sprintf("%v enable:%v", deviceID, enable))

	updated, err := u.UHPPOTE.RecordSpecialEvents(deviceID, enable)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error enabling/disabling 'record special events' (%w)", deviceID, err))
	}

	u.debug("record-special-events", fmt.Sprintf("updated %+v", updated))

	return updated, nil
}

func (u *UHPPOTED) FetchEvents(controller uint32, from, to uint32) ([]Event, error) {
	first, err := u.UHPPOTE.GetEvent(controller, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve 'first' event for controller %d (%w)", controller, err)
	} else if first == nil {
		return nil, fmt.Errorf("no 'first' event record returned for controller %d", controller)
	}

	last, err := u.UHPPOTE.GetEvent(controller, 0xffffffff)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve 'last' event for controller %d (%w)", controller, err)
	} else if first == nil {
		return nil, fmt.Errorf("no 'last' event record returned for controller %d", controller)
	}

	if last.Index >= first.Index {
		from = min(max(from, first.Index), last.Index)
		to = min(max(to, first.Index), last.Index)
	}

	events := []Event{}
	for index := from; index <= to; index++ {
		record, err := u.UHPPOTE.GetEvent(controller, index)
		if err != nil {
			u.warn("fetch-events", fmt.Errorf("failed to retrieve event for controller %d, ID %d (%w)", controller, index, err))
		} else if record == nil {
			u.warn("fetch-events", fmt.Errorf("no event record for controller %d, index %d", controller, index))
		} else if record.Index != index {
			u.warn("fetch-events", fmt.Errorf("no event record for controller %d, index %d", controller, index))
		} else {
			events = append(events, Event{
				DeviceID:   uint32(record.SerialNumber),
				Index:      record.Index,
				Type:       record.Type,
				Granted:    record.Granted,
				Door:       record.Door,
				Direction:  record.Direction,
				CardNumber: record.CardNumber,
				Timestamp:  record.Timestamp,
				Reason:     record.Reason,
			})
		}
	}

	return events, nil
}
