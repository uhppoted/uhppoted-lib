package acl

import (
	"fmt"
	"sync"

	"codeberg.org/uhppoted/uhppoted-core/types"
	"codeberg.org/uhppoted/uhppoted-core/uhppote"
)

func PutACL(u uhppote.IUHPPOTE, acl ACL, dryrun bool, withPIN bool, withFirstCard bool, formats ...types.CardFormat) (map[uint32]Report, []error) {
	f := func(u uhppote.IUHPPOTE, deviceID uint32, c types.Card) (bool, error) {
		return putCard(u, deviceID, c, withPIN, withFirstCard, formats...)
	}

	eq := func(p, q types.Card) bool {
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

		if withPIN && p.PIN != q.PIN {
			return false
		}

		if withFirstCard && !p.FirstCard.Equal(q.FirstCard) {
			return false
		}

		return true
	}

	return putACLImpl(u, acl, dryrun, f, eq)
}

func putACLImpl(u uhppote.IUHPPOTE, acl ACL, dryrun bool, write put, eq equivalent) (map[uint32]Report, []error) {
	report := sync.Map{}
	errors := []error{}
	guard := sync.RWMutex{}

	for id := range acl {
		report.Store(id, Report{
			Unchanged: []uint32{},
			Updated:   []uint32{},
			Added:     []uint32{},
			Deleted:   []uint32{},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		})
	}

	var wg sync.WaitGroup

	for k, v := range acl {
		id := k
		cards := v

		wg.Go(func() {
			var rpt *Report
			var err error

			if dryrun {
				rpt, err = fakePutACL(u, id, cards, eq)
			} else {
				rpt, err = putACL(u, id, cards, write, eq)
			}

			if rpt != nil {
				report.Store(id, *rpt)
			}

			if err != nil {
				guard.Lock()
				errors = append(errors, err)
				guard.Unlock()
			}

		})
	}

	wg.Wait()

	r := map[uint32]Report{}
	report.Range(func(k, v any) bool {
		r[k.(uint32)] = v.(Report)
		return true
	})

	return r, errors
}

func putACL(u uhppote.IUHPPOTE, deviceID uint32, cards map[uint32]types.Card, write put, eq equivalent) (*Report, error) {
	current, err := getACL(u, deviceID)
	if err != nil {
		return nil, err
	}

	diff := compare(deviceID, current, cards, eq)

	report := Report{
		Unchanged: []uint32{},
		Updated:   []uint32{},
		Added:     []uint32{},
		Deleted:   []uint32{},
		Failed:    []uint32{},
		Errored:   []uint32{},
		Errors:    []error{},
	}

	for _, card := range diff.Unchanged {
		report.Unchanged = append(report.Unchanged, card.CardNumber)
	}

	for _, card := range diff.Updated {
		if err := validate(u, deviceID, card); err != nil {
			report.Errored = append(report.Errored, card.CardNumber)
			report.Errors = append(report.Errors, err)
		} else {
			if ok, err := write(u, deviceID, card); err != nil {
				report.Errored = append(report.Errored, card.CardNumber)
				report.Errors = append(report.Errors, err)
			} else if !ok {
				report.Failed = append(report.Failed, card.CardNumber)
			} else {
				report.Updated = append(report.Updated, card.CardNumber)
			}
		}
	}

	for _, card := range diff.Added {
		if err := validate(u, deviceID, card); err != nil {
			report.Errored = append(report.Errored, card.CardNumber)
			report.Errors = append(report.Errors, err)
		} else {
			if ok, err := write(u, deviceID, card); err != nil {
				report.Errored = append(report.Errored, card.CardNumber)
				report.Errors = append(report.Errors, err)
			} else if !ok {
				report.Failed = append(report.Failed, card.CardNumber)
			} else {
				report.Added = append(report.Added, card.CardNumber)
			}
		}
	}

	for _, card := range diff.Deleted {
		if ok, err := u.DeleteCard(deviceID, card.CardNumber); err != nil {
			report.Errored = append(report.Errored, card.CardNumber)
			report.Errors = append(report.Errors, err)
		} else if !ok {
			report.Failed = append(report.Failed, card.CardNumber)
		} else {
			report.Deleted = append(report.Deleted, card.CardNumber)
		}
	}

	return &report, nil
}

func fakePutACL(u uhppote.IUHPPOTE, deviceID uint32, cards map[uint32]types.Card, eq equivalent) (*Report, error) {
	current, err := getACL(u, deviceID)
	if err != nil {
		return nil, err
	}

	diff := compare(deviceID, current, cards, eq)

	report := Report{
		Unchanged: []uint32{},
		Updated:   []uint32{},
		Added:     []uint32{},
		Deleted:   []uint32{},
		Failed:    []uint32{},
		Errored:   []uint32{},
		Errors:    []error{},
	}

	for _, card := range diff.Unchanged {
		report.Unchanged = append(report.Unchanged, card.CardNumber)
	}

	for _, card := range diff.Updated {
		report.Updated = append(report.Updated, card.CardNumber)
	}

	for _, card := range diff.Added {
		report.Added = append(report.Added, card.CardNumber)
	}

	for _, card := range diff.Deleted {
		report.Deleted = append(report.Deleted, card.CardNumber)
	}

	return &report, nil
}

func validate(u uhppote.IUHPPOTE, deviceID uint32, card types.Card) error {
	for _, door := range []uint8{1, 2, 3, 4} {
		if v, ok := card.Doors[door]; ok && v >= 2 && v <= 254 {
			if profile, err := u.GetTimeProfile(deviceID, uint8(v)); err != nil {
				return err
			} else if profile == nil {
				return fmt.Errorf("time profile %v is not defined for %v", v, deviceID)
			}
		}
	}

	return nil
}
