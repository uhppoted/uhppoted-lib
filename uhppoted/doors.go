package uhppoted

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

func (u *UHPPOTED) GetDoorDelay(request GetDoorDelayRequest) (*GetDoorDelayResponse, error) {
	u.debug("get-door-delay", fmt.Sprintf("request  %+v", request))

	controller := uint32(request.DeviceID)
	door := request.Door
	result, err := u.UHPPOTE.GetDoorControlState(controller, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error getting door %v delay for %v (%w)", door, controller, err))
	}

	response := GetDoorDelayResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug("get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) SetDoorDelay(controller uint32, door uint8, delay uint8) error {
	u.debug("set-door-delay", fmt.Sprintf("%v door:%v delay:%v", controller, door, delay))

	state, err := u.UHPPOTE.GetDoorControlState(controller, door)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error getting door %v delay (%w)", controller, door, err))
	}

	response, err := u.UHPPOTE.SetDoorControlState(controller, door, state.ControlState, delay)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error setting door %v delay (%ws)", controller, door, err))
	}

	u.debug("set-door-delay", fmt.Sprintf("response %+v", response))

	return nil
}

func (u *UHPPOTED) GetDoorControl(request GetDoorControlRequest) (*GetDoorControlResponse, error) {
	u.debug("get-door-control", fmt.Sprintf("request  %+v", request))

	controller := uint32(request.DeviceID)
	door := request.Door
	result, err := u.UHPPOTE.GetDoorControlState(controller, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error getting door %v control for %v (%w)", door, controller, err))
	}

	response := GetDoorControlResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Control:  result.ControlState,
	}

	u.debug("get-door-control", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) SetDoorControl(controller uint32, door uint8, mode types.ControlState) error {
	u.debug("set-door-control", fmt.Sprintf("%v door:%v mode:%v", controller, door, mode))

	state, err := u.UHPPOTE.GetDoorControlState(controller, door)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error getting door %v control mode (%w)", controller, door, err))
	}

	response, err := u.UHPPOTE.SetDoorControlState(controller, door, mode, state.Delay)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error setting door %v control mode %v (%w)", controller, door, mode, err))
	}

	u.debug("set-door-control", fmt.Sprintf("response %+v", response))

	return nil
}

func (u *UHPPOTED) SetDoorPasscodes(controller uint32, door uint8, passcodes ...uint32) error {
	u.debug("set-door-passcodes", fmt.Sprintf("%v door:%v", controller, door))

	response, err := u.UHPPOTE.SetDoorPasscodes(controller, door, passcodes...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error setting door %v passcodes (%w)", controller, door, err))
	}

	u.debug("set-door-passcodes", fmt.Sprintf("response %+v", response))

	return nil
}

func (u *UHPPOTED) OpenDoor(request OpenDoorRequest) (*OpenDoorResponse, error) {
	u.debug("open-door", fmt.Sprintf("request  %+v", request))

	controller := uint32(request.DeviceID)
	door := request.Door
	result, err := u.UHPPOTE.OpenDoor(controller, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error opening door %v on %v (%w)", door, controller, err))
	}

	response := OpenDoorResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     door,
		Opened:   result.Succeeded,
	}

	u.debug("open-door", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) SetInterlock(controller uint32, interlock types.Interlock) error {
	u.debug("set-interlock", fmt.Sprintf("%v  mode:%v", controller, interlock))

	response, err := u.UHPPOTE.SetInterlock(controller, interlock)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error setting door interlock %v (%w)", controller, interlock, err))
	}

	u.debug("set-interlock", fmt.Sprintf("%v  response:%+v", controller, response))

	return nil
}

func (u *UHPPOTED) ActivateKeypads(controller uint32, keypads map[uint8]bool) error {
	u.debug("activate-keypads", fmt.Sprintf("%v  mode:%v", controller, keypads))

	response, err := u.UHPPOTE.ActivateKeypads(controller, keypads)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error activating controller access keypads (%w)", controller, err))
	} else if !response {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  failed activate controller access keypads", controller))
	}

	u.debug("activate-keypads", fmt.Sprintf("%v  response:%+v", controller, response))

	return nil
}
