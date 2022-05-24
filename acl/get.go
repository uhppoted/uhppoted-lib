package acl

import (
	"sync"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func GetACL(u uhppote.IUHPPOTE, devices []uhppote.Device) (ACL, []error) {
	acl := sync.Map{}
	errors := []error{}
	guard := sync.RWMutex{}

	for _, device := range devices {
		acl.Store(device.DeviceID, map[uint32]types.Card{})
	}

	var wg sync.WaitGroup

	for _, d := range devices {
		device := d
		wg.Add(1)
		go func() {
			if cards, err := getACL(u, device.DeviceID); err != nil {
				guard.Lock()
				errors = append(errors, err)
				guard.Unlock()
			} else {
				acl.Store(device.DeviceID, cards)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	a := make(ACL)
	acl.Range(func(k, v interface{}) bool {
		a[k.(uint32)] = v.(map[uint32]types.Card)
		return true
	})

	return a, errors
}

func getACL(u uhppote.IUHPPOTE, deviceID uint32) (map[uint32]types.Card, error) {
	cards := map[uint32]types.Card{}

	N, err := u.GetCards(deviceID)
	if err != nil {
		return cards, err
	}

	var index uint32 = 1
	for count := 0; count < int(N); {
		card, err := u.GetCardByIndex(deviceID, index, nil)
		if err != nil {
			return nil, err
		}

		if card != nil {
			cards[card.CardNumber] = card.Clone()
			count++
		}

		index++
	}

	return cards, nil
}
