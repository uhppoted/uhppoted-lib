package uhppoted

import (
	"net"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

type IUHPPOTED interface {
	GetDevices(request GetDevicesRequest) (*GetDevicesResponse, error)
	GetDevice(request GetDeviceRequest) (*GetDeviceResponse, error)
	GetTime(request GetTimeRequest) (*GetTimeResponse, error)
	SetTime(request SetTimeRequest) (*SetTimeResponse, error)
	GetDoorDelay(request GetDoorDelayRequest) (*GetDoorDelayResponse, error)
	SetDoorDelay(request SetDoorDelayRequest) (*SetDoorDelayResponse, error)
	GetDoorControl(request GetDoorControlRequest) (*GetDoorControlResponse, error)
	SetDoorControl(request SetDoorControlRequest) (*SetDoorControlResponse, error)
	RecordSpecialEvents(request RecordSpecialEventsRequest) (*RecordSpecialEventsResponse, error)
	GetStatus(request GetStatusRequest) (*GetStatusResponse, error)
	GetCardRecords(request GetCardRecordsRequest) (*GetCardRecordsResponse, error)
	GetCards(request GetCardsRequest) (*GetCardsResponse, error)
	DeleteCards(request DeleteCardsRequest) (*DeleteCardsResponse, error)
	GetCard(request GetCardRequest) (*GetCardResponse, error)
	PutCard(request PutCardRequest) (*PutCardResponse, error)
	DeleteCard(request DeleteCardRequest) (*DeleteCardResponse, error)
	GetTimeProfiles(request GetTimeProfilesRequest) (*GetTimeProfilesResponse, error)
	PutTimeProfiles(request PutTimeProfilesRequest) (*PutTimeProfilesResponse, int, error)
	GetTimeProfile(request GetTimeProfileRequest) (*GetTimeProfileResponse, error)
	PutTimeProfile(request PutTimeProfileRequest) (*PutTimeProfileResponse, error)
	ClearTimeProfiles(request ClearTimeProfilesRequest) (*ClearTimeProfilesResponse, error)
	PutTaskList(request PutTaskListRequest) (*PutTaskListResponse, int, error)
	GetEventRange(request GetEventRangeRequest) (*GetEventRangeResponse, error)
	GetEvent(request GetEventRequest) (*GetEventResponse, error)
	OpenDoor(request OpenDoorRequest) (*OpenDoorResponse, error)
}

type GetDevicesRequest struct {
}

type GetDevicesResponse struct {
	Devices map[uint32]DeviceSummary `json:"devices"`
}

type GetDeviceRequest struct {
	DeviceID DeviceID
}

type GetDeviceResponse struct {
	DeviceType string           `json:"device-type"`
	DeviceID   DeviceID         `json:"device-id"`
	IpAddress  net.IP           `json:"ip-address"`
	SubnetMask net.IP           `json:"subnet-mask"`
	Gateway    net.IP           `json:"gateway-address"`
	MacAddress types.MacAddress `json:"mac-address"`
	Version    types.Version    `json:"version"`
	Date       types.Date       `json:"date"`
	Address    net.UDPAddr      `json:"address"`
	TimeZone   *time.Location   `json:"timezone,omitempty"`
}

type GetTimeRequest struct {
	DeviceID DeviceID
}

type GetTimeResponse struct {
	DeviceID DeviceID       `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

type SetTimeRequest struct {
	DeviceID DeviceID
	DateTime types.DateTime
}

type SetTimeResponse struct {
	DeviceID DeviceID       `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

type GetDoorDelayRequest struct {
	DeviceID DeviceID
	Door     uint8
}

type GetDoorDelayResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Door     uint8    `json:"door"`
	Delay    uint8    `json:"delay"`
}

type SetDoorDelayRequest struct {
	DeviceID DeviceID
	Door     uint8
	Delay    uint8
}

type SetDoorDelayResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Door     uint8    `json:"door"`
	Delay    uint8    `json:"delay"`
}

type GetDoorControlRequest struct {
	DeviceID DeviceID
	Door     uint8
}

type GetDoorControlResponse struct {
	DeviceID DeviceID           `json:"device-id"`
	Door     uint8              `json:"door"`
	Control  types.ControlState `json:"control"`
}

type SetDoorControlRequest struct {
	DeviceID DeviceID
	Door     uint8
	Control  types.ControlState
}

type SetDoorControlResponse struct {
	DeviceID DeviceID           `json:"device-id"`
	Door     uint8              `json:"door"`
	Control  types.ControlState `json:"control"`
}

type GetStatusRequest struct {
	DeviceID DeviceID
}

type GetStatusResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Status   Status   `json:"status"`
}

type GetCardRecordsRequest struct {
	DeviceID DeviceID
}

type GetCardRecordsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Cards    uint32   `json:"cards"`
}

type GetCardsRequest struct {
	DeviceID DeviceID
}

type GetCardsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Cards    []uint32 `json:"cards"`
}

type DeleteCardsRequest struct {
	DeviceID DeviceID
}

type DeleteCardsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Deleted  bool     `json:"deleted"`
}

type GetCardRequest struct {
	DeviceID   DeviceID
	CardNumber uint32
}

type GetCardResponse struct {
	DeviceID DeviceID   `json:"device-id"`
	Card     types.Card `json:"card"`
}

type PutCardRequest struct {
	DeviceID DeviceID
	Card     types.Card
}

type PutCardResponse struct {
	DeviceID DeviceID   `json:"device-id"`
	Card     types.Card `json:"card"`
}

type DeleteCardRequest struct {
	DeviceID   DeviceID
	CardNumber uint32
}

type DeleteCardResponse struct {
	DeviceID   DeviceID `json:"device-id"`
	CardNumber uint32   `json:"card-number"`
	Deleted    bool     `json:"deleted"`
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

type GetEventRangeRequest struct {
	DeviceID DeviceID
	Start    *types.DateTime
	End      *types.DateTime
}

type GetEventRangeResponse struct {
	DeviceID DeviceID    `json:"device-id,omitempty"`
	Dates    *DateRange  `json:"dates,omitempty"`
	Events   *EventRange `json:"events,omitempty"`
}

type GetEventRequest struct {
	DeviceID DeviceID
	EventID  uint32
}

type GetEventResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Event    Event    `json:"event"`
}

type RecordSpecialEventsRequest struct {
	DeviceID DeviceID
	Enable   bool
}

type RecordSpecialEventsResponse struct {
	DeviceID DeviceID
	Enable   bool
	Updated  bool
}

type OpenDoorRequest struct {
	DeviceID DeviceID
	Door     uint8
}

type OpenDoorResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Door     uint8    `json:"door"`
	Opened   bool     `json:"opened"`
}
