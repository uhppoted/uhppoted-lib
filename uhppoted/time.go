package uhppoted

import (
	"fmt"
	"time"
)

func (u *UHPPOTED) GetTime(request GetTimeRequest) (*GetTimeResponse, error) {
	u.debug("get-time", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	result, err := u.UHPPOTE.GetTime(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error getting time for %v (%w)", device, err))
	}

	response := GetTimeResponse{
		DeviceID: DeviceID(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug("get-time", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) SetTime(request SetTimeRequest) (*SetTimeResponse, error) {
	u.debug("set-time", fmt.Sprintf("request  %v", request))

	device := uint32(request.DeviceID)
	result, err := u.UHPPOTE.SetTime(device, time.Time(request.DateTime))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error setting time for %v (%w)", device, err))
	}

	response := SetTimeResponse{
		DeviceID: DeviceID(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug("set-time", fmt.Sprintf("response %+v", response))

	return &response, nil
}
