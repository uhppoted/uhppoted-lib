package uhppoted

import (
	"github.com/uhppoted/uhppote-core/types"
)

type IUHPPOTED interface {
	PutTaskList(request PutTaskListRequest) (*PutTaskListResponse, int, error)
}

type PutTaskListRequest struct {
	DeviceID uint32
	Tasks    []types.Task `json:"tasks"`
}

type PutTaskListResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Warnings []error  `json:"warnings"`
}
