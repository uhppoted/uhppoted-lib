package acl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func parseHeader(header []string, devices []uhppote.Device) (*index, error) {
	columns := make(map[string]struct {
		door  string
		index int
	})

	index := index{
		cardnumber: 0,
		from:       0,
		to:         0,
		doors:      make(map[uint32][]int),
		PIN:        0,
	}

	for _, d := range devices {
		index.doors[d.DeviceID] = make([]int, 4)
	}

	for c, field := range header {
		key := clean(field)
		ix := c + 1

		if columns[key].index != 0 {
			return nil, fmt.Errorf("duplicate column name '%s'", field)
		}

		columns[key] = struct {
			door  string
			index int
		}{
			door:  field,
			index: ix,
		}
	}

loop:
	for c, v := range columns {
		if c != "cardnumber" && c != "from" && c != "to" && c != "pin" {
			for _, device := range devices {
				for _, door := range device.Doors {
					if d := clean(door); d == c {
						continue loop
					}
				}
			}

			return nil, fmt.Errorf("no configured door matches '%s'", v.door)
		}
	}

	if c, ok := columns["cardnumber"]; ok {
		index.cardnumber = c.index
	}

	if c, ok := columns["from"]; ok {
		index.from = c.index
	}

	if c, ok := columns["to"]; ok {
		index.to = c.index
	}

	if c, ok := columns["pin"]; ok {
		index.PIN = c.index
	}

	for _, device := range devices {
		for i, door := range device.Doors {
			if d := clean(door); d != "" {
				if c, ok := columns[d]; ok {
					index.doors[device.DeviceID][i] = c.index
				}
			}
		}
	}

	if index.cardnumber == 0 {
		return nil, fmt.Errorf("missing 'Card Number' column")
	}

	if index.from == 0 {
		return nil, fmt.Errorf("missing 'From' column")
	}

	if index.to == 0 {
		return nil, fmt.Errorf("missing 'To' column")
	}

	//	for _, device := range devices {
	//		for i, door := range device.Doors {
	//			if d := clean(door); d != "" {
	//				if index.doors[device.DeviceID][i] == 0 {
	//					return nil, fmt.Errorf("missing column for door '%s'", door)
	//				}
	//			}
	//		}
	//	}

	return &index, nil
}

func parseRecord(record []string, index index) (map[uint32]types.Card, error) {
	cards := make(map[uint32]types.Card, 0)

	for k, v := range index.doors {
		cardno, err := getCardNumber(record, index)
		if err != nil {
			return nil, err
		}

		pin, err := getPIN(record, index)
		if err != nil {
			return nil, err
		}

		from, err := getFromDate(record, index)
		if err != nil {
			return nil, err
		}

		to, err := getToDate(record, index)
		if err != nil {
			return nil, err
		}

		doors, err := getDoors(record, v)
		if err != nil {
			return nil, err
		}

		cards[k] = types.Card{
			CardNumber: cardno,
			From:       from,
			To:         to,
			Doors:      doors,
			PIN:        types.PIN(pin),
		}
	}

	return cards, nil
}

func getCardNumber(record []string, index index) (uint32, error) {
	f := field(record, index.cardnumber)
	cardnumber, err := strconv.ParseUint(f, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid card number '%s' (%w)", f, err)
	}

	return uint32(cardnumber), nil
}

func getPIN(record []string, index index) (uint32, error) {
	if index.PIN > 0 {
		f := field(record, index.PIN)

		if f == "" {
			return 0, nil
		} else if pin, err := strconv.ParseUint(f, 10, 32); err != nil {
			return 0, fmt.Errorf("invalid card PIN '%s' (%w)", f, err)
		} else if pin > 999999 {
			return 0, fmt.Errorf("invalid card PIN '%s' (%v)", f, pin)
		} else {
			return uint32(pin), nil
		}
	}

	return 0, nil
}

func getFromDate(record []string, index index) (types.Date, error) {
	f := field(record, index.from)

	if date, err := time.ParseInLocation("2006-01-02", f, time.Local); err != nil {
		return types.Date{}, fmt.Errorf("invalid 'from' date '%s' (%w)", f, err)
	} else {
		return types.Date(date), nil
	}
}

func getToDate(record []string, index index) (types.Date, error) {
	f := field(record, index.to)

	if date, err := time.ParseInLocation("2006-01-02", f, time.Local); err != nil {
		return types.Date{}, fmt.Errorf("invalid 'to' date '%s' (%w)", f, err)
	} else {
		return types.Date(date), nil
	}
}

func getDoors(record []string, v []int) (map[uint8]uint8, error) {
	doors := map[uint8]uint8{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
	}

	re := regexp.MustCompile("[0-9]+")
	for i, d := range v {
		if d == 0 {
			continue
		}

		v := field(record, d)
		if v == "N" {
			doors[uint8(i+1)] = 0
		} else if v == "Y" {
			doors[uint8(i+1)] = 1
		} else if matched := re.MatchString(v); matched {
			if profile, _ := strconv.Atoi(v); profile < 2 || profile > 254 {
				return doors, fmt.Errorf("invalid time profile (%v) for door %v (valid profiles are in the interval [2..254])", v, record[d])
			} else {
				doors[uint8(i+1)] = uint8(profile)
			}
		} else {
			return doors, fmt.Errorf("expected 'Y/N/<profile ID>' for door: '%s'", record[d])
		}
	}

	return doors, nil
}

func field(record []string, ix int) string {
	return strings.TrimSpace(record[ix-1])
}
