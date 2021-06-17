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
	Event          *StatusEvent   `json:"event,omitempty"`
}

type StatusEvent struct {
	Index      uint32          `json:"index"`
	Type       byte            `json:"type"`
	Granted    bool            `json:"access-granted"`
	Door       byte            `json:"door"`
	Direction  uint8           `json:"direction"`
	CardNumber uint32          `json:"card-number"`
	Timestamp  *types.DateTime `json:"timestamp,omitempty"`
	Reason     uint8           `json:"reason"`
}

type GetStatusRequest struct {
	DeviceID DeviceID
}

type GetStatusResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Status   Status   `json:"status"`
}

func (u *UHPPOTED) GetStatus(request GetStatusRequest) (*GetStatusResponse, error) {
	u.debug("get-status", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	status, err := u.UHPPOTE.GetStatus(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving status for %v (%w)", device, err))
	}

	response := GetStatusResponse{
		DeviceID: DeviceID(status.SerialNumber),
		Status: Status{
			DoorState:      status.DoorState,
			DoorButton:     status.DoorButton,
			SystemError:    status.SystemError,
			SystemDateTime: status.SystemDateTime,
			SequenceId:     status.SequenceId,
			SpecialInfo:    status.SpecialInfo,
			RelayState:     status.RelayState,
			InputState:     status.InputState,
		},
	}

	if status.Event != nil {
		response.Status.Event = &StatusEvent{
			Index:      status.Event.Index,
			Type:       status.Event.Type,
			Granted:    status.Event.Granted,
			Door:       status.Event.Door,
			Direction:  status.Event.Direction,
			CardNumber: status.Event.CardNumber,
			Timestamp:  status.Event.Timestamp,
			Reason:     status.Event.Reason,
		}
	}

	u.debug("get-status", fmt.Sprintf("response %+v", response))

	return &response, nil
}
