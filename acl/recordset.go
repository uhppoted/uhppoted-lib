package acl

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"strings"

	"codeberg.org/uhppoted/uhppoted-core/types"
	"codeberg.org/uhppoted/uhppoted-core/uhppote"
)

type DuplicateCardError struct {
	CardNumber uint32
}

func (e *DuplicateCardError) Error() string {
	return fmt.Sprintf("%-10v Duplicate card number", e.CardNumber)
}

func ParseTable(table *Table, controllers []uhppote.Device, strict bool, withPINs bool, withFirstCard bool) (*ACL, []error, error) {
	acl := make(ACL)
	for _, device := range controllers {
		acl[device.DeviceID] = make(map[uint32]types.Card)
	}

	index, err := parseHeader(table.Header, controllers, withPINs, withFirstCard)
	if err != nil {
		return nil, nil, err
	} else if index == nil {
		return nil, nil, fmt.Errorf("invalid table header")
	}

	list := []map[uint32]types.Card{}
	for row, record := range table.Records {
		cards, err := parseRecord(record, *index, withPINs, withFirstCard)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing table - row %d: %w", row+1, err)
		}

		list = append(list, cards)
	}

	duplicates := map[uint32]int{}
	for _, cards := range list {
		for _, card := range cards {
			count := duplicates[card.CardNumber]
			duplicates[card.CardNumber] = count + 1
			break
		}
	}

	warnings := []error{}
	for _, cards := range list {
	loop:
		for id, card := range cards {
			if acl[id] != nil {
				if count := duplicates[card.CardNumber]; count > 1 {
					if strict {
						return nil, nil, fmt.Errorf("duplicate card number (%v)", card.CardNumber)
					} else {
						warning := &DuplicateCardError{card.CardNumber}
						for i := range warnings {
							if reflect.DeepEqual(warnings[i], warning) {
								continue loop
							}
						}

						warnings = append(warnings, warning)
					}

					continue
				}

				acl[id][card.CardNumber] = card
			}
		}
	}

	return &acl, warnings, nil
}

func MakeTable(acl ACL, controllers []uhppote.Device, withPIN bool, withFirstCard bool) (*Table, error) {
	if header, err := makeHeader(controllers, withPIN, withFirstCard); err != nil {
		return nil, err
	} else {
		// .. build column index from header
		index := struct {
			PIN       int
			firstcard int
			doors     map[string]int
		}{
			PIN:       -1,
			firstcard: -1,
			doors:     map[string]int{},
		}

		cdr := header[3:]
		if withPIN {
			cdr = header[4:]
		}

		for i, h := range cdr {
			col := clean(h)

			if withFirstCard && col == "firstcard" {
				index.firstcard = i + 1
			} else {
				index.doors[col] = i + 1
			}
		}

		// ... map card fields to record fields
		cards := map[uint32]card{}
		for _, d := range controllers {
			v, ok := acl[d.DeviceID]
			if !ok {
				return nil, fmt.Errorf("ACL missing for device %v", d.DeviceID)
			}

			jndex := []int{0, 0, 0, 0}
			for i, door := range d.Doors {
				if clean(door) != "" {
					jndex[i] = index.doors[clean(door)]
				}
			}

			for cardno, c := range v {
				record, ok := cards[cardno]
				if !ok {
					record = card{
						cardnumber: c.CardNumber,
						PIN:        uint32(c.PIN),
						from:       c.From,
						to:         c.To,
						doors:      make([]int, len(index.doors)),
						firstcard:  false,
					}
				}

				if c.From.Before(record.from) {
					record.from = c.From
				}

				if c.To.After(record.to) {
					record.to = c.To
				}

				if uint32(c.PIN) != record.PIN {
					record.PIN = math.MaxUint32
				}

				for _, i := range []uint8{1, 2, 3, 4} {
					ix := jndex[i-1]

					if ix == 0 && clean(d.Doors[i-1]) != "" {
						return nil, fmt.Errorf("missing door ID for device %v, door:%v", d.DeviceID, i)
					}

					if ix != 0 {
						record.doors[ix-1] = int(c.Doors[i])
					}
				}

				// // ... set first-card if any first-card privilege is set and any door has access
				if withFirstCard {
					if c.Doors[1] > 0 && c.FirstCard.Door1 {
						record.firstcard = true
					}

					if c.Doors[2] > 0 && c.FirstCard.Door2 {
						record.firstcard = true
					}

					if c.Doors[3] > 0 && c.FirstCard.Door3 {
						record.firstcard = true
					}

					if c.Doors[4] > 0 && c.FirstCard.Door4 {
						record.firstcard = true
					}
				}

				cards[cardno] = record
			}
		}

		keys := []uint32{}
		for k := range cards {
			keys = append(keys, k)
		}

		slices.Sort(keys)

		records := [][]string{}
		for _, k := range keys {
			c := cards[k]
			record := []string{
				fmt.Sprintf("%v", c.cardnumber),
				fmt.Sprintf("%v", c.from),
				fmt.Sprintf("%v", c.to),
			}

			if withPIN {
				pin := fmt.Sprintf("%v", c.PIN)
				if c.PIN == math.MaxUint32 {
					pin = "****"
				} else if c.PIN == 0 {
					pin = ""
				}

				record = []string{
					fmt.Sprintf("%v", c.cardnumber),
					fmt.Sprintf("%v", pin),
					fmt.Sprintf("%v", c.from),
					fmt.Sprintf("%v", c.to),
				}
			}

			for _, v := range c.doors {
				switch {
				case v == 0:
					record = append(record, "N")

				case v == 1:
					record = append(record, "Y")

				case v > 1 && v < 255:
					record = append(record, fmt.Sprintf("%v", v))
				default:
					record = append(record, "N")
				}
			}

			if withFirstCard {
				if c.firstcard {
					record = append(record, "Y")
				} else {
					record = append(record, "N")
				}
			}

			records = append(records, record)
		}

		rs := Table{
			Header:  header,
			Records: records,
		}

		return &rs, nil
	}
}

// func MakeTableWithPIN(acl ACL, controllers []uhppote.Device) (*Table, error) {
// 	if header, err := makeHeader(controllers, true, false); err != nil {
// 		return nil, err
// 	} else {
// 		index := map[string]int{}
// 		for i, h := range header {
// 			if i > 3 {
// 				index[clean(h)] = i - 3
// 			}
// 		}
//
// 		cards := map[uint32]card{}
// 		for _, d := range controllers {
// 			v, ok := acl[d.DeviceID]
// 			if !ok {
// 				return nil, fmt.Errorf("ACL missing for device %v", d.DeviceID)
// 			}
//
// 			jndex := []int{0, 0, 0, 0}
// 			for i, door := range d.Doors {
// 				if clean(door) != "" {
// 					jndex[i] = index[clean(door)]
// 				}
// 			}
//
// 			for cardno, c := range v {
// 				record, ok := cards[cardno]
// 				if !ok {
// 					record = card{
// 						cardnumber: c.CardNumber,
// 						PIN:        uint32(c.PIN),
// 						from:       c.From,
// 						to:         c.To,
// 						doors:      make([]int, len(index)),
// 					}
// 				}
//
// 				if c.From.Before(record.from) {
// 					record.from = c.From
// 				}
//
// 				if c.To.After(record.to) {
// 					record.to = c.To
// 				}
//
// 				if uint32(c.PIN) != record.PIN {
// 					record.PIN = math.MaxUint32
// 				}
//
// 				for i := uint8(1); i <= 4; i++ {
// 					ix := jndex[i-1]
//
// 					if ix == 0 && clean(d.Doors[i-1]) != "" {
// 						return nil, fmt.Errorf("missing door ID for device %v, door:%v", d.DeviceID, i)
// 					}
//
// 					if ix != 0 {
// 						record.doors[ix-1] = int(c.Doors[i])
// 					}
// 				}
//
// 				cards[cardno] = record
// 			}
// 		}
//
// 		keys := []uint32{}
// 		for k := range cards {
// 			keys = append(keys, k)
// 		}
//
// 		slices.Sort(keys)
//
// 		records := [][]string{}
// 		for _, k := range keys {
// 			c := cards[k]
//
// 			pin := fmt.Sprintf("%v", c.PIN)
// 			if c.PIN == math.MaxUint32 {
// 				pin = "****"
// 			} else if c.PIN == 0 {
// 				pin = ""
// 			}
//
// 			record := []string{
// 				fmt.Sprintf("%v", c.cardnumber),
// 				fmt.Sprintf("%v", pin),
// 				fmt.Sprintf("%v", c.from),
// 				fmt.Sprintf("%v", c.to),
// 			}
//
// 			for _, v := range c.doors {
// 				switch {
// 				case v == 0:
// 					record = append(record, "N")
//
// 				case v == 1:
// 					record = append(record, "Y")
//
// 				case v > 1 && v < 255:
// 					record = append(record, fmt.Sprintf("%v", v))
// 				default:
// 					record = append(record, "N")
// 				}
// 			}
//
// 			records = append(records, record)
// 		}
//
// 		rs := Table{
// 			Header:  header,
// 			Records: records,
// 		}
//
// 		return &rs, nil
// 	}
// }

func makeHeader(controllers []uhppote.Device, withPIN bool, withFirstCard bool) ([]string, error) {
	keys := []uint32{}
	for _, d := range controllers {
		keys = append(keys, d.DeviceID)
	}

	slices.Sort(keys)

	header := []string{
		"Card Number",
		"From",
		"To",
	}

	if withPIN {
		header = []string{
			"Card Number",
			"PIN",
			"From",
			"To",
		}
	}

	for _, id := range keys {
		for _, d := range controllers {
			if d.DeviceID == id {
				for _, door := range d.Doors {
					if clean(door) != "" {
						header = append(header, strings.TrimSpace(door))
					}
				}
			}
		}
	}

	if withFirstCard {
		header = append(header, "FirstCard")
	}

	return header, nil
}
