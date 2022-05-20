package acl

import (
	"bytes"
	"fmt"
	"net"
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

func (m *mock) GetDevice(deviceID uint32) (*types.Device, error) {
	return nil, nil
}

func (m *mock) SetAddress(deviceID uint32, address, mask, gateway net.IP) (*types.Result, error) {
	return nil, nil
}

func (m *mock) GetTime(deviceID uint32) (*types.Time, error) {
	return nil, nil
}

func (m *mock) SetTime(deviceID uint32, datetime time.Time) (*types.Time, error) {
	return nil, nil
}

func (m *mock) GetDoorControlState(deviceID uint32, door byte) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *mock) SetDoorControlState(deviceID uint32, door uint8, state types.ControlState, delay uint8) (*types.DoorControlState, error) {
	return nil, nil
}

func (m *mock) GetListener(deviceID uint32) (*types.Listener, error) {
	return nil, nil
}

func (m *mock) SetListener(deviceID uint32, address net.UDPAddr) (*types.Result, error) {
	return nil, nil
}

func (m *mock) GetStatus(deviceID uint32) (*types.Status, error) {
	return nil, nil
}

func (m *mock) GetCards(deviceID uint32) (uint32, error) {
	return m.getCards(deviceID)
}

func (m *mock) GetCardByIndex(deviceID, index uint32) (*types.Card, error) {
	return m.getCardByIndex(deviceID, index)
}

func (m *mock) GetCardByID(deviceID, cardID uint32) (*types.Card, error) {
	return m.getCardByID(deviceID, cardID)
}

func (m *mock) PutCard(deviceID uint32, card types.Card) (bool, error) {
	return m.putCard(deviceID, card)
}

func (m *mock) DeleteCard(deviceID uint32, cardNumber uint32) (bool, error) {
	return m.deleteCard(deviceID, cardNumber)
}

func (m *mock) DeleteCards(deviceID uint32) (bool, error) {
	return m.deleteCards(deviceID)
}

func (m *mock) GetTimeProfile(deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
	return m.getTimeProfile(deviceID, profileID)
}

func (m *mock) SetTimeProfile(deviceID uint32, profile types.TimeProfile) (bool, error) {
	return false, nil
}

func (m *mock) ClearTimeProfiles(deviceID uint32) (bool, error) {
	return false, nil
}

func (m *mock) ClearTaskList(deviceID uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *mock) AddTask(deviceID uint32, task types.Task) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *mock) RefreshTaskList(deviceID uint32) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (m *mock) RecordSpecialEvents(deviceID uint32, enable bool) (bool, error) {
	return false, nil
}

func (m *mock) GetEvent(deviceID, index uint32) (*types.Event, error) {
	return nil, nil
}

func (m *mock) GetEventIndex(deviceID uint32) (*types.EventIndex, error) {
	return nil, nil
}

func (m *mock) SetEventIndex(deviceID, index uint32) (*types.EventIndexResult, error) {
	return nil, nil
}

func (m *mock) Listen(listener uhppote.Listener, q chan os.Signal) error {
	return nil
}

func (m *mock) OpenDoor(deviceID uint32, door uint8) (*types.Result, error) {
	return nil, nil
}

func (m *mock) DeviceList() map[uint32]uhppote.Device {
	return map[uint32]uhppote.Device{}
}

func (m *mock) ListenAddr() *net.UDPAddr {
	return nil
}

var date = func(s string) *types.Date {
	d, _ := time.ParseInLocation("2006-01-02", s, time.Local)
	p := types.Date(d)
	return &p
}

var deviceA = uhppote.Device{
	DeviceID: 12345,
	Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
}

var deviceB = uhppote.Device{
	DeviceID: 54321,
	Doors:    []string{"D1", "D2", "D3", "D4"},
}

var aclA = ACL{
	12345: map[uint32]types.Card{
		65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
		65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
		65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
	},
}

var cardsA = []types.Card{
	types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 0}},
	types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]uint8{1: 1, 2: 0, 3: 0, 4: 1}},
	types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]uint8{1: 0, 2: 0, 3: 0, 4: 0}},
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

	if string(b.Bytes()) != expected {
		t.Errorf("Invalid result from ACL.Print\n   expected:\n%s\n   got:\n%s\n", expected, string(b.Bytes()))
	}
}

func TestClean(t *testing.T) {
	s := string([]byte{70, 114, 111, 110, 116, 9})

	if clean(s) != "front" {
		t.Errorf("Clean did not strip trailing tab from string")
	}
}
