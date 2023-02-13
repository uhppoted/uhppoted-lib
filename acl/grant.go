package acl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func Grant(u uhppote.IUHPPOTE, devices []uhppote.Device, cardID uint32, from, to types.Date, profile int, doors []string) error {
	m, err := mapDeviceDoors(devices)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(doors, []string{"ALL"}) {
		for _, d := range devices {
			if err := grantAll(u, d.DeviceID, cardID, from, to); err != nil {
				return err
			}
		}

		return nil
	}

	for _, dd := range doors {
		door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
		if _, ok := m[door]; !ok {
			return fmt.Errorf("door '%v' is not defined in the device configuration", dd)
		}
	}

	for _, d := range devices {
		l := []uint8{}

		for _, dd := range doors {
			door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
			if e, ok := m[door]; ok && e.deviceID == d.DeviceID {
				l = append(l, e.door)
			}
		}

		if err := grant(u, d.DeviceID, cardID, from, to, profile, l); err != nil {
			return err
		}
	}

	return nil
}

func grant(u uhppote.IUHPPOTE, deviceID uint32, cardID uint32, from, to types.Date, profileID int, doors []uint8) error {
	if len(doors) == 0 {
		return nil
	}

	if profileID >= 2 && profileID <= 254 {
		if profile, err := u.GetTimeProfile(deviceID, uint8(profileID)); err != nil {
			return err
		} else if profile == nil {
			return fmt.Errorf("time profile %v is not defined for %v", profileID, deviceID)
		}
	}

	card, err := u.GetCardByID(deviceID, cardID)
	if err != nil {
		return err
	} else if card == nil {
		card = &types.Card{
			CardNumber: cardID,
			From:       &from,
			To:         &to,
			Doors: map[uint8]uint8{
				1: 0,
				2: 0,
				3: 0,
				4: 0,
			},
		}
	}

	revoked := true
	for _, v := range card.Doors {
		switch {
		case v == 1:
			revoked = false
		case v >= 2 && v <= 254:
			revoked = false
		}
	}

	if revoked || card.From.After(from) {
		card.From = &from
	}

	if revoked || card.To.Before(to) {
		card.To = &to
	}

	for _, d := range doors {
		if profileID >= 2 && profileID <= 254 {
			card.Doors[d] = uint8(profileID)
		} else {
			card.Doors[d] = 1
		}
	}

	if ok, err := putCard(u, deviceID, *card); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("failed to update access rights for card '%v' on device '%v'", cardID, deviceID)
	}

	return nil
}

func grantAll(u uhppote.IUHPPOTE, deviceID uint32, cardID uint32, from, to types.Date) error {
	card := types.Card{
		CardNumber: cardID,
		From:       &from,
		To:         &to,
		Doors: map[uint8]uint8{
			1: 1,
			2: 1,
			3: 1,
			4: 1,
		},
	}

	if ok, err := putCard(u, deviceID, card); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("failed to update access rights for card '%v' on device '%v'", cardID, deviceID)
	}

	return nil
}
