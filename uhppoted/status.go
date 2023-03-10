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
	Event          *Event         `json:"event,omitempty"`
}

// type StatusEvent struct {
// 	Index      uint32          `json:"index"`
// 	Type       byte            `json:"type"`
// 	Granted    bool            `json:"access-granted"`
// 	Door       byte            `json:"door"`
// 	Direction  uint8           `json:"direction"`
// 	CardNumber uint32          `json:"card-number"`
// 	Timestamp  *types.DateTime `json:"timestamp,omitempty"`
// 	Reason     uint8           `json:"reason"`
// }

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

	if status.Event != nil {
		response.Event = &Event{
			Index:      status.Event.Index,
			Type:       status.Event.Type,
			Granted:    status.Event.Granted,
			Door:       status.Event.Door,
			Direction:  status.Event.Direction,
			CardNumber: status.Event.CardNumber,
			Reason:     status.Event.Reason,
		}

		if status.Event.Timestamp != nil {
			response.Event.Timestamp = *status.Event.Timestamp
		}
	}

	return &response, nil
}
