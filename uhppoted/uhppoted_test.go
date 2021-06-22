package uhppoted

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

type stub struct {
	getTimeProfile      func(deviceID uint32, profileID uint8) (*types.TimeProfile, error)
	setTimeProfile      func(deviceID uint32, profile types.TimeProfile) (bool, error)
	clearTimeProfiles   func(deviceID uint32) (bool, error)
	recordSpecialEvents func(deviceID uint32, enable bool) (bool, error)
}

func (m *stub) DeviceList() map[uint32]uhppote.Device {
	return nil
}

func (m *stub) ListenAddr() *net.UDPAddr {
	return nil
}

func (m *stub) GetDevices() ([]types.Device, error) {
	return nil, nil
}

func (m *stub) GetDevice(deviceID uint32) (*types.Device, error) {
	return nil, nil
}

func (m *stub) SetAddress(deviceID uint32, address, mask, gateway net.IP) (*types.Result, error) {
	return nil, nil
}

func (m *stub) GetTime(serialNumber uint32) (*types.Time, error) {
	return nil, nil
}

func (m *stub) SetTime(serialNumber uint32, datetime time.Time) (*types.Time, error) {
	return nil, nil
}

func (m *stub) GetListener(deviceID uint32) (*types.Listener, error) {
	return nil, nil
}

func (m *stub) SetListener(deviceID uint32, address net.UDPAddr) (*types.Result, error) {
	return nil, nil
}

func (m *stub) GetStatus(serialNumber uint32) (*types.Status, error) {
	return nil, nil
}

func (m *stub) GetCards(deviceID uint32) (uint32, error) {
	return 0, nil
}

func (m *stub) GetCardByIndex(deviceID, index uint32) (*types.Card, error) {
	return nil, nil
}

func (m *stub) GetCardByID(deviceID, cardNumber uint32) (*types.Card, error) {
	return nil, nil
}

func (m *stub) PutCard(deviceID uint32, card types.Card) (bool, error) {
	return false, nil
}

func (m *stub) DeleteCard(deviceID uint32, cardNumber uint32) (bool, error) {
	return false, nil
}

func (m *stub) DeleteCards(deviceID uint32) (bool, error) {
	return false, nil
}

func (m *stub) GetDoorControlState(deviceID uint32, door byte) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *stub) SetDoorControlState(deviceID uint32, door uint8, state uint8, delay uint8) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *stub) GetTimeProfile(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
	if m.getTimeProfile != nil {
		return m.getTimeProfile(deviceID, profileID)
	}

	return nil, fmt.Errorf("Not implemented")
}

func (m *stub) SetTimeProfile(deviceID uint32, profile types.TimeProfile) (bool, error) {
	if m.setTimeProfile != nil {
		return m.setTimeProfile(deviceID, profile)
	}

	return false, fmt.Errorf("Not implemented")
}

func (m *stub) ClearTimeProfiles(deviceID uint32) (bool, error) {
	if m.clearTimeProfiles != nil {
		return m.clearTimeProfiles(deviceID)
	}

	return false, fmt.Errorf("Not implemented")
}

func (m *stub) ClearTaskList(deviceID uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *stub) AddTask(deviceID uint32, task types.Task) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *stub) RefreshTaskList(deviceID uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *stub) OpenDoor(deviceID uint32, door uint8) (*types.Result, error) {
	return nil, nil
}

func (m *stub) GetEventIndex(deviceID uint32) (*types.EventIndex, error) {
	return nil, nil
}

func (m *stub) SetEventIndex(deviceID, index uint32) (*types.EventIndexResult, error) {
	return nil, nil
}

func (m *stub) GetEvent(deviceID, index uint32) (*types.Event, error) {
	return nil, nil
}

func (m *stub) RecordSpecialEvents(deviceID uint32, enable bool) (bool, error) {
	if m.recordSpecialEvents != nil {
		return m.recordSpecialEvents(deviceID, enable)
	}

	return false, fmt.Errorf("Not implemented")
}

func (m *stub) Listen(listener uhppote.Listener, q chan os.Signal) error {
	return nil
}
