package acl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uhppoted/uhppote-core/uhppote"
)

func Revoke(u uhppote.IUHPPOTE, devices []uhppote.Device, cardID uint32, doors []string) error {
	m, err := mapDeviceDoors(devices)
	if err != nil {
		return err
	}

	list := []string{}
	if reflect.DeepEqual(doors, []string{"ALL"}) {
		for k := range m {
			list = append(list, k)
		}
	} else {
		list = append(list, doors...)
	}

	for _, dd := range list {
		door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
		if _, ok := m[door]; !ok {
			return fmt.Errorf("door '%v' is not defined in the device configuration", dd)
		}
	}

	for _, d := range devices {
		l := []uint8{}

		for _, dd := range list {
			door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
			if e, ok := m[door]; ok && e.deviceID == d.DeviceID {
				l = append(l, e.door)
			}
		}

		if err := revoke(u, d.DeviceID, cardID, l); err != nil {
			return err
		}
	}

	return nil
}

func revoke(u uhppote.IUHPPOTE, deviceID uint32, cardID uint32, doors []uint8) error {
	if len(doors) == 0 {
		return nil
	}

	card, err := u.GetCardByID(deviceID, cardID)
	if err != nil {
		return err
	} else if card == nil {
		return nil
	}

	for _, d := range doors {
		card.Doors[d] = 0
	}

	if ok, err := u.PutCard(deviceID, *card); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("failed to update access rights for card '%v' on device '%v'", cardID, deviceID)
	}

	return nil
}
