package uhppoted

import (
	"errors"
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

func (u *UHPPOTED) GetEvent(request GetEventRequest) (*GetEventResponse, error) {
	u.debug("get-events", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	eventID := request.EventID

	record, err := u.UHPPOTE.GetEvent(device, eventID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting event for ID %v from %v (%w)", eventID, device, err))
	}

	if record == nil {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("No event record for ID %v for %v", eventID, device))
	}

	if record.Index != eventID {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("No event record for ID %v for %v", eventID, device))
	}

	response := GetEventResponse{
		DeviceID: DeviceID(record.SerialNumber),
		Event: Event{
			DeviceID:   device,
			Index:      record.Index,
			Type:       record.Type,
			Granted:    record.Granted,
			Door:       record.Door,
			Direction:  record.Direction,
			CardNumber: record.CardNumber,
			Timestamp:  record.Timestamp,
			Reason:     record.Reason,
		},
	}

	u.debug("get-event", fmt.Sprintf("response %+v", response))

	return &response, nil
}

// Retrieves up to MAX events starting with the current controller event index. The current controller
// event index is updated on completion of this request.
func (u *UHPPOTED) GetEvents(request GetEventsRequest) (*GetEventsResponse, error) {
	u.debug("get-events", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	events := []Event{}

	first, err := u.UHPPOTE.GetEvent(device, 0)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting first event index from %v (%w)", device, err))
	}

	last, err := u.UHPPOTE.GetEvent(device, 0xffffffff)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting last event index from %v (%w)", device, err))
	}

	if first != nil && first.Index > 0 && last != nil && last.Index > 0 && last.Index != first.Index {
		current, err := u.UHPPOTE.GetEventIndex(device)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting current event index from %v (%w)", device, err))
		} else if current == nil {
			return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting current event index from %v (%w)", device, errors.New("Record not found")))
		}

		max := request.Max
		index := current.Index
		next := index + 1

		if index == 0 || index < first.Index {
			index = first.Index
			next = first.Index
		}

		for len(events) < max && index != last.Index {
			e, err := u.UHPPOTE.GetEvent(device, next)
			if err != nil {
				return nil, err
			} else if e == nil || e.Index != next {
				if last.Index < first.Index {
					next = 1
				} else {
					break
				}
			} else {
				events = append(events, Event{
					DeviceID:   device,
					Index:      e.Index,
					Type:       e.Type,
					Granted:    e.Granted,
					Door:       e.Door,
					Direction:  e.Direction,
					CardNumber: e.CardNumber,
					Timestamp:  e.Timestamp,
					Reason:     e.Reason,
				})

				index = next
				next = index + 1
			}
		}

		if _, err := u.UHPPOTE.SetEventIndex(device, index); err != nil {
			return nil, err
		}
	}

	response := GetEventsResponse{
		DeviceID: DeviceID(device),
		Events:   events,
	}

	u.debug("get-events", fmt.Sprintf("response %+v", response))

	return &response, nil
}

// Unwraps the request and dispatches the corresponding controller command to enable or disable
// door open, door close and door button press events for the controller.
func (u *UHPPOTED) RecordSpecialEvents(request RecordSpecialEventsRequest) (*RecordSpecialEventsResponse, error) {
	u.debug("record-special-events", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	enable := request.Enable

	updated, err := u.UHPPOTE.RecordSpecialEvents(device, enable)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error updating 'record special events' flag for %v (%w)", device, err))
	}

	response := RecordSpecialEventsResponse{
		DeviceID: DeviceID(device),
		Enable:   enable,
		Updated:  updated,
	}

	u.debug("record-special-events", fmt.Sprintf("response %+v", response))

	return &response, nil
}
