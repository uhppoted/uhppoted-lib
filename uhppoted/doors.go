package uhppoted

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

func (u *UHPPOTED) GetDoorDelay(request GetDoorDelayRequest) (*GetDoorDelayResponse, error) {
	u.debug("get-door-delay", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	door := request.Door
	result, err := u.UHPPOTE.GetDoorControlState(device, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting door %v delay for %v (%w)", door, device, err))
	}

	response := GetDoorDelayResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug("get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) SetDoorDelay(deviceID uint32, door uint8, delay uint8) error {
	u.debug("set-door-delay", fmt.Sprintf("%v door:%v delay:%v", deviceID, door, delay))

	state, err := u.UHPPOTE.GetDoorControlState(deviceID, door)
	if err != nil {
		return fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("%v  error getting door %v delay (%w)", deviceID, door, err))
	}

	response, err := u.UHPPOTE.SetDoorControlState(deviceID, door, state.ControlState, delay)
	if err != nil {
		return fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("%v  error setting door %v delay (%ws)", deviceID, door, err))
	}

	u.debug("set-door-delay", fmt.Sprintf("response %+v", response))

	return nil
}

func (u *UHPPOTED) GetDoorControl(request GetDoorControlRequest) (*GetDoorControlResponse, error) {
	u.debug("get-door-control", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	door := request.Door
	result, err := u.UHPPOTE.GetDoorControlState(device, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting door %v control for %v (%w)", door, device, err))
	}

	response := GetDoorControlResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Control:  result.ControlState,
	}

	u.debug("get-door-control", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) SetDoorControl(deviceID uint32, door uint8, mode types.ControlState) error {
	u.debug("set-door-control", fmt.Sprintf("%v door:%v mode:%v", deviceID, door, mode))

	state, err := u.UHPPOTE.GetDoorControlState(deviceID, door)
	if err != nil {
		return fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("%v  error getting door %v control mode (%w)", deviceID, door, err))
	}

	response, err := u.UHPPOTE.SetDoorControlState(deviceID, door, mode, state.Delay)
	if err != nil {
		return fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("%v  error setting door %v control mode %v (%w)", deviceID, door, mode, err))
	}

	u.debug("set-door-control", fmt.Sprintf("response %+v", response))

	return nil
}

func (u *UHPPOTED) OpenDoor(request OpenDoorRequest) (*OpenDoorResponse, error) {
	u.debug("open-door", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	door := request.Door
	result, err := u.UHPPOTE.OpenDoor(device, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error opening door %v on %v (%w)", door, device, err))
	}

	response := OpenDoorResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     door,
		Opened:   result.Succeeded,
	}

	u.debug("open-door", fmt.Sprintf("response %+v", response))

	return &response, nil
}
