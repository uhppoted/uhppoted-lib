package uhppoted

import (
	"github.com/uhppoted/uhppote-core/types"
)

type IUHPPOTED interface {
	GetTimeProfiles(request GetTimeProfilesRequest) (*GetTimeProfilesResponse, error)
	PutTimeProfiles(request PutTimeProfilesRequest) (*PutTimeProfilesResponse, int, error)
	GetTimeProfile(request GetTimeProfileRequest) (*GetTimeProfileResponse, error)
	PutTimeProfile(request PutTimeProfileRequest) (*PutTimeProfileResponse, error)
	ClearTimeProfiles(request ClearTimeProfilesRequest) (*ClearTimeProfilesResponse, error)
	PutTaskList(request PutTaskListRequest) (*PutTaskListResponse, int, error)
}

type GetTimeProfilesRequest struct {
	DeviceID uint32
	From     int
	To       int
}

type GetTimeProfilesResponse struct {
	DeviceID DeviceID            `json:"device-id"`
	Profiles []types.TimeProfile `json:"profiles"`
}

type PutTimeProfilesRequest struct {
	DeviceID uint32
	Profiles []types.TimeProfile `json:"profiles"`
}

type PutTimeProfilesResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Warnings []error  `json:"warnings"`
}

type GetTimeProfileRequest struct {
	DeviceID  uint32
	ProfileID uint8
}

type GetTimeProfileResponse struct {
	DeviceID    DeviceID          `json:"device-id"`
	TimeProfile types.TimeProfile `json:"time-profile"`
}

type PutTimeProfileRequest struct {
	DeviceID    uint32
	TimeProfile types.TimeProfile
}

type PutTimeProfileResponse struct {
	DeviceID    DeviceID          `json:"device-id"`
	TimeProfile types.TimeProfile `json:"time-profile"`
}

type ClearTimeProfilesRequest struct {
	DeviceID uint32
}

type ClearTimeProfilesResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Cleared  bool     `json:"cleared"`
}

type PutTaskListRequest struct {
	DeviceID uint32
	Tasks    []types.Task `json:"tasks"`
}

type PutTaskListResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Warnings []error  `json:"warnings"`
}
