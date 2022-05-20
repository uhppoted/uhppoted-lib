package acl

import (
	"github.com/uhppoted/uhppote-core/uhppote"
)

func GetCard(u uhppote.IUHPPOTE, devices []uhppote.Device, cardID uint32) (map[string]Permission, error) {
	acl := map[string]Permission{}
	lookup, err := mapDeviceDoors(devices)
	if err != nil {
		return acl, err
	}

	for _, device := range devices {
		card, err := u.GetCardByID(device.DeviceID, cardID)
		if err != nil {
			return acl, err
		}

		if card != nil {
			for ix, v := range card.Doors {
				ok := false
				profile := 0

				switch {
				case v == 0:
					ok = false

				case v == 1:
					ok = true

				case v >= 2 && v <= 254:
					ok = true
					profile = int(v)
				}

				if ok {
					for _, v := range lookup {
						if v.deviceID == device.DeviceID && v.door == ix {
							acl[v.name] = Permission{
								From:    *card.From,
								To:      *card.To,
								Profile: profile,
							}
						}
					}
				}
			}
		}
	}

	return acl, nil
}
