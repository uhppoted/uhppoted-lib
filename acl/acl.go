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
}

type doormap map[string]struct {
	deviceID uint32
	door     uint8
	name     string
}

type card struct {
	cardnumber uint32
	from       types.Date
	to         types.Date
	doors      []int
}

func (acl *ACL) Print(w io.Writer) {
	if acl != nil {
		devices := []uint32{}
		for k, _ := range *acl {
			devices = append(devices, k)
		}

		sort.SliceStable(devices, func(i, j int) bool { return devices[i] < devices[j] })

		for _, k := range devices {
			v := (*acl)[k]

			cards := []uint32{}
			for c, _ := range v {
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
				return m, fmt.Errorf("Ambiguous reference to door '%s': defined for both devices %v and %v", dd, e.deviceID, d.DeviceID)
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
