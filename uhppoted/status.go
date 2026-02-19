package uhppoted

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

type Status struct {
	DoorState      map[uint8]bool `json:"door-states"`
	DoorButton     map[uint8]bool `json:"door-buttons"`
	SystemError    uint8          `json:"system-error"`
	SystemDateTime types.DateTime `json:"system-datetime"`
	SequenceId     uint32         `json:"sequence-id"`
	SpecialInfo    uint8          `json:"special-info"`
	RelayState     uint8          `json:"relay-state"`
	InputState     uint8          `json:"input-state"`
	Event          Event          `json:"event"`
}

func (u *UHPPOTED) GetStatus(deviceID uint32) (*Status, error) {
	status, err := u.UHPPOTE.GetStatus(deviceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error retrieving status for %v (%w)", deviceID, err))
	}

	sysdatetime := func() types.DateTime {
		if status.SystemDateTime.IsZero() {
			return types.DateTime{}
		} else {
			return status.SystemDateTime
		}
	}

	response := Status{
		DoorState:      status.DoorState,
		DoorButton:     status.DoorButton,
		SystemError:    status.SystemError,
		SystemDateTime: sysdatetime(),
		SequenceId:     status.SequenceId,
		SpecialInfo:    status.SpecialInfo,
		RelayState:     status.RelayState,
		InputState:     status.InputState,
	}

	if !status.Event.IsZero() {
		response.Event = Event{
			Index:      status.Event.Index,
			Type:       status.Event.Type,
			Granted:    status.Event.Granted,
			Door:       status.Event.Door,
			Direction:  status.Event.Direction,
			CardNumber: status.Event.CardNumber,
			Reason:     status.Event.Reason,
		}

		if !status.Event.Timestamp.IsZero() {
			response.Event.Timestamp = status.Event.Timestamp
		}
	}

	return &response, nil
}
