package acl

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"codeberg.org/uhppoted/uhppoted-core/types"
	"codeberg.org/uhppoted/uhppoted-core/uhppote"
)

func ParseTSV(f io.Reader, controllers []uhppote.Device, strict bool) (ACL, []error, error) {
	acl := make(ACL)
	for _, device := range controllers {
		acl[device.DeviceID] = make(map[uint32]types.Card)
	}

	r := csv.NewReader(f)
	r.Comma = '\t'

	header, err := r.Read()
	if err != nil {
		return nil, nil, err
	}

	index, err := parseHeader(header, controllers)
	if err != nil {
		return nil, nil, err
	} else if index == nil {
		return nil, nil, fmt.Errorf("invalid TSV header")
	}

	line := 0
	list := []map[uint32]types.Card{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}

		line += 1
		cards, err := parseRecord(record, *index)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing TSV - line %d: %w", line, err)
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
						warning := fmt.Errorf("duplicate card number (%v)", card.CardNumber)
						for _, w := range warnings {
							if reflect.DeepEqual(w, warning) {
								continue loop
							}
						}

						warnings = append(warnings, fmt.Errorf("duplicate card number (%v)", card.CardNumber))
					}

					continue
				}

				acl[id][card.CardNumber] = card
			}
		}
	}

	return acl, warnings, nil
}

func MakeTSV(acl ACL, controllers []uhppote.Device, withPIN bool, withFirstCard bool, f io.Writer) error {
	t, err := MakeTable(acl, controllers, withPIN, withFirstCard)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.Comma = '\t'

	if err := w.Write(t.Header); err != nil {
		return err
	}

	for _, r := range t.Records {
		if err := w.Write(r); err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}

func MakeTSVWithPIN(acl ACL, controllers []uhppote.Device, f io.Writer) error {
	t, err := MakeTable(acl, controllers, true, false)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.Comma = '\t'

	if err := w.Write(t.Header); err != nil {
		return err
	}

	for _, r := range t.Records {
		if err := w.Write(r); err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}
