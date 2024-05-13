package uhppoted

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

type stub struct {
	getTimeProfile      func(controller uint32, profileID uint8) (*types.TimeProfile, error)
	setTimeProfile      func(controller uint32, profile types.TimeProfile) (bool, error)
	clearTimeProfiles   func(controller uint32) (bool, error)
	getEventIndex       func(controller uint32) (*types.EventIndex, error)
	setEventIndex       func(controller, index uint32) (*types.EventIndexResult, error)
	getEvent            func(controller, index uint32) (*types.Event, error)
	recordSpecialEvents func(controller uint32, enable bool) (bool, error)
}

func (m *stub) DeviceList() map[uint32]uhppote.Device {
	return nil
}

func (m *stub) ListenAddrList() []netip.AddrPort {
	return nil
}

func (m *stub) GetDevices() ([]types.Device, error) {
	return nil, nil
}

func (m *stub) GetDevice(controller uint32) (*types.Device, error) {
	return nil, nil
}

func (m *stub) SetAddress(controller uint32, address, mask, gateway net.IP) (*types.Result, error) {
	return nil, nil
}

func (m *stub) GetTime(serialNumber uint32) (*types.Time, error) {
	return nil, nil
}

func (m *stub) SetTime(serialNumber uint32, datetime time.Time) (*types.Time, error) {
	return nil, nil
}

func (m *stub) GetListener(controller uint32) (*types.Listener, error) {
	return nil, nil
}

func (m *stub) SetListener(controller uint32, address net.UDPAddr) (*types.Result, error) {
	return nil, nil
}

func (m *stub) GetStatus(serialNumber uint32) (*types.Status, error) {
	return nil, nil
}

func (m *stub) GetCards(controller uint32) (uint32, error) {
	return 0, nil
}

func (m *stub) GetCardByIndex(controller, index uint32) (*types.Card, error) {
	return nil, nil
}

func (m *stub) GetCardByID(controller, cardNumber uint32) (*types.Card, error) {
	return nil, nil
}

func (m *stub) PutCard(controller uint32, card types.Card, formats ...types.CardFormat) (bool, error) {
	return false, nil
}

func (m *stub) DeleteCard(controller uint32, cardNumber uint32) (bool, error) {
	return false, nil
}

func (m *stub) DeleteCards(controller uint32) (bool, error) {
	return false, nil
}

func (m *stub) GetDoorControlState(controller uint32, door byte) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *stub) SetDoorControlState(controller uint32, door uint8, state types.ControlState, delay uint8) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *stub) SetDoorPasscodes(controller uint32, door uint8, passcodes ...uint32) (bool, error) {
	return false, nil
}

func (m *stub) SetInterlock(controller uint32, interlock types.Interlock) (bool, error) {
	return false, nil
}

func (m *stub) ActivateKeypads(controller uint32, keypads map[uint8]bool) (bool, error) {
	return false, nil
}

func (m *stub) GetTimeProfile(controller uint32, profileID uint8) (*types.TimeProfile, error) {
	if m.getTimeProfile != nil {
		return m.getTimeProfile(controller, profileID)
	}

	return nil, fmt.Errorf("Not implemented")
}

func (m *stub) SetTimeProfile(controller uint32, profile types.TimeProfile) (bool, error) {
	if m.setTimeProfile != nil {
		return m.setTimeProfile(controller, profile)
	}

	return false, fmt.Errorf("Not implemented")
}

func (m *stub) ClearTimeProfiles(controller uint32) (bool, error) {
	if m.clearTimeProfiles != nil {
		return m.clearTimeProfiles(controller)
	}

	return false, fmt.Errorf("Not implemented")
}

func (m *stub) ClearTaskList(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *stub) AddTask(controller uint32, task types.Task) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *stub) RefreshTaskList(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *stub) OpenDoor(controller uint32, door uint8) (*types.Result, error) {
	return nil, nil
}

func (m *stub) GetEventIndex(controller uint32) (*types.EventIndex, error) {
	if m.getEventIndex != nil {
		return m.getEventIndex(controller)
	}

	return nil, fmt.Errorf("Not implemented")
}

func (m *stub) SetEventIndex(controller, index uint32) (*types.EventIndexResult, error) {
	if m.setEventIndex != nil {
		return m.setEventIndex(controller, index)
	}

	return nil, fmt.Errorf("Not implemented")
}

func (m *stub) GetEvent(controller, index uint32) (*types.Event, error) {
	if m.getEvent != nil {
		return m.getEvent(controller, index)
	}

	return nil, fmt.Errorf("Not implemented")
}

func (m *stub) RecordSpecialEvents(controller uint32, enable bool) (bool, error) {
	if m.recordSpecialEvents != nil {
		return m.recordSpecialEvents(controller, enable)
	}

	return false, fmt.Errorf("Not implemented")
}

func (m *stub) SetPCControl(controller uint32, enable bool) (bool, error) {
	return true, nil
}

func (m *stub) RestoreDefaultParameters(controller uint32) (bool, error) {
	return true, nil
}

func (m *stub) Listen(listener uhppote.Listener, q chan os.Signal) error {
	return nil
}
