package acl

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

type ACL map[uint32]map[uint32]types.Card

type Permission struct {
	From    types.Date
	To      types.Date
	Profile int
}

type index struct {
	cardnumber int
	from       int
	to         int
	doors      map[uint32][]int
	PIN        int
}

type doormap map[string]struct {
	deviceID uint32
	door     uint8
	name     string
}

type card struct {
	cardnumber uint32
	PIN        uint32
	from       types.Date
	to         types.Date
	doors      []int
}

type equivalent = func(types.Card, types.Card) bool
type put = func(u uhppote.IUHPPOTE, deviceID uint32, c types.Card) (bool, error)

func (acl *ACL) Print(w io.Writer) {
	if acl != nil {
		devices := []uint32{}
		for k := range *acl {
			devices = append(devices, k)
		}

		sort.SliceStable(devices, func(i, j int) bool { return devices[i] < devices[j] })

		for _, k := range devices {
			v := (*acl)[k]

			cards := []uint32{}
			for c := range v {
				cards = append(cards, c)
			}

			sort.SliceStable(cards, func(i, j int) bool { return cards[i] < cards[j] })

			fmt.Fprintf(w, "%v\n", k)
			for _, c := range cards {
				card := v[c]
				fmt.Fprintf(w, "  %v %v\n", c, card)
			}
		}
	}
}

func clean(s string) string {
	return regexp.MustCompile(`[\s\t]+`).ReplaceAllString(strings.ToLower(s), "")
}

func mapDeviceDoors(devices []uhppote.Device) (doormap, error) {
	m := doormap{}

	for _, d := range devices {
		for i, dd := range d.Doors {
			door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
			if e, ok := m[door]; ok {
				return m, fmt.Errorf("ambiguous reference to door '%s': defined for both devices %v and %v", dd, e.deviceID, d.DeviceID)
			}

			m[door] = struct {
				deviceID uint32
				door     uint8
				name     string
			}{
				deviceID: d.DeviceID,
				door:     uint8(i + 1),
				name:     strings.TrimSpace(dd),
			}
		}
	}

	return m, nil
}

func putCard(u uhppote.IUHPPOTE, deviceID uint32, c types.Card, formats ...types.CardFormat) (bool, error) {
	card, err := u.GetCardByID(deviceID, c.CardNumber)
	if err != nil {
		return false, err
	} else if card == nil {
		card = &types.Card{CardNumber: c.CardNumber}
	}

	card.From = c.From
	card.To = c.To
	card.Doors = c.Doors

	return u.PutCard(deviceID, *card, formats...)
}

func putCardWithPIN(u uhppote.IUHPPOTE, deviceID uint32, c types.Card, formats ...types.CardFormat) (bool, error) {
	card, err := u.GetCardByID(deviceID, c.CardNumber)
	if err != nil {
		return false, err
	} else if card == nil {
		card = &types.Card{CardNumber: c.CardNumber}
	}

	card.From = c.From
	card.To = c.To
	card.Doors = c.Doors
	card.PIN = c.PIN

	return u.PutCard(deviceID, *card, formats...)
}

/*
 * Compares two cards, ignoring PIN
 */
func equals(p, q types.Card) bool {
	if p.CardNumber != q.CardNumber {
		return false
	}

	if !p.From.Equals(q.From) {
		return false
	}

	if !p.To.Equals(q.To) {
		return false
	}

	for _, i := range []uint8{1, 2, 3, 4} {
		if p.Doors[i] != q.Doors[i] {
			return false
		}
	}

	return true
}

/*
 * Compares two cards, including PIN
 */
func equalsWithPIN(p, q types.Card) bool {
	if !equals(p, q) {
		return false
	} else if p.PIN != q.PIN {
		return false
	}

	return true
}
