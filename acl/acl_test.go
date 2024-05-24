package acl

import (
	"bytes"
	"fmt"
	"net"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

type mock struct {
	getCards       func(uint32) (uint32, error)
	getCardByIndex func(uint32, uint32) (*types.Card, error)
	getCardByID    func(uint32, uint32) (*types.Card, error)
	putCard        func(uint32, types.Card) (bool, error)
	deleteCard     func(uint32, uint32) (bool, error)
	deleteCards    func(uint32) (bool, error)
	getTimeProfile func(uint32, uint8) (*types.TimeProfile, error)
}

func (m *mock) GetDevices() ([]types.Device, error) {
	return nil, nil
}

func (m *mock) GetDevice(controller uint32) (*types.Device, error) {
	return nil, nil
}

func (m *mock) SetAddress(controller uint32, address, mask, gateway net.IP) (*types.Result, error) {
	return nil, nil
}

func (m *mock) GetTime(controller uint32) (*types.Time, error) {
	return nil, nil
}

func (m *mock) SetTime(controller uint32, datetime time.Time) (*types.Time, error) {
	return nil, nil
}

func (m *mock) GetDoorControlState(controller uint32, door byte) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *mock) SetDoorControlState(controller uint32, door uint8, state types.ControlState, delay uint8) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *mock) SetDoorPasscodes(controller uint32, door uint8, passcodes ...uint32) (bool, error) {
	return false, nil
}

func (m *mock) SetInterlock(controller uint32, interlock types.Interlock) (bool, error) {
	return false, nil
}

func (m *mock) ActivateKeypads(controller uint32, keypads map[uint8]bool) (bool, error) {
	return false, nil
}

func (m *mock) GetListener(controller uint32) (netip.AddrPort, error) {
	return netip.AddrPort{}, nil
}

func (m *mock) SetListener(controller uint32, address net.UDPAddr) (*types.Result, error) {
	return nil, nil
}

func (m *mock) GetStatus(controller uint32) (*types.Status, error) {
	return nil, nil
}

func (m *mock) GetCards(controller uint32) (uint32, error) {
	return m.getCards(controller)
}

func (m *mock) GetCardByIndex(controller, index uint32) (*types.Card, error) {
	return m.getCardByIndex(controller, index)
}

func (m *mock) GetCardByID(controller, cardID uint32) (*types.Card, error) {
	return m.getCardByID(controller, cardID)
}

func (m *mock) PutCard(controller uint32, card types.Card, formats ...types.CardFormat) (bool, error) {
	return m.putCard(controller, card)
}

func (m *mock) DeleteCard(controller uint32, cardNumber uint32) (bool, error) {
	return m.deleteCard(controller, cardNumber)
}

func (m *mock) DeleteCards(controller uint32) (bool, error) {
	return m.deleteCards(controller)
}

func (m *mock) GetTimeProfile(controller uint32, profileID uint8) (*types.TimeProfile, error) {
	return m.getTimeProfile(controller, profileID)
}

func (m *mock) SetTimeProfile(controller uint32, profile types.TimeProfile) (bool, error) {
	return false, nil
}

func (m *mock) ClearTimeProfiles(controller uint32) (bool, error) {
	return false, nil
}

func (m *mock) ClearTaskList(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *mock) AddTask(controller uint32, task types.Task) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *mock) RefreshTaskList(controller uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *mock) RecordSpecialEvents(controller uint32, enable bool) (bool, error) {
	return false, nil
}

func (m *mock) GetEvent(controller, index uint32) (*types.Event, error) {
	return nil, nil
}

func (m *mock) GetEventIndex(controller uint32) (*types.EventIndex, error) {
	return nil, nil
}

func (m *mock) SetEventIndex(controller, index uint32) (*types.EventIndexResult, error) {
	return nil, nil
}

func (m *mock) Listen(listener uhppote.Listener, q chan os.Signal) error {
	return nil
}

func (m *mock) OpenDoor(controller uint32, door uint8) (*types.Result, error) {
	return nil, nil
}

func (m *mock) SetPCControl(controller uint32, enable bool) (bool, error) {
	return true, nil
}

func (m *mock) RestoreDefaultParameters(controller uint32) (bool, error) {
	return true, nil
}

func (m *mock) DeviceList() map[uint32]uhppote.Device {
	return map[uint32]uhppote.Device{}
}

func (m *mock) ListenAddrList() []netip.AddrPort {
	return nil
}

var date = func(s string) types.Date {
	if d, err := time.ParseInLocation("2006-01-02", s, time.Local); err != nil {
		return types.Date{}
	} else {
		return types.Date(d)
	}
}

var deviceA = uhppote.Device{
	DeviceID: 12345,
	Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
}

func TestACLPrintf(t *testing.T) {
	expected := `12345
  65531 65531    2020-01-02 2020-10-31 Y N N N
  65532 65532    2020-02-03 2020-11-30 Y N N Y
  65533 65533    2020-03-04 2020-12-31 N N N N
67890
  65531 65531    2020-01-02 2020-10-31 Y N N N
  65532 65532    2020-02-03 2020-11-30 Y N N Y
  65534 65534    2020-03-04 2020-12-31 N N N N
`

	acl := ACL{
		12345: map[uint32]types.Card{
			65531: types.Card{CardNumber: 65531, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65532: types.Card{CardNumber: 65532, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65533: types.Card{CardNumber: 65533, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
		67890: map[uint32]types.Card{
			65531: types.Card{CardNumber: 65531, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
			65532: types.Card{CardNumber: 65532, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
			65534: types.Card{CardNumber: 65534, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
		},
	}

	var b bytes.Buffer

	acl.Print(&b)

	if b.String() != expected {
		t.Errorf("Invalid result from ACL.Print\n   expected:\n%s\n   got:\n%s\n", expected, b.String())
	}
}

func TestClean(t *testing.T) {
	s := string([]byte{70, 114, 111, 110, 116, 9})

	if clean(s) != "front" {
		t.Errorf("Clean did not strip trailing tab from string")
	}
}
