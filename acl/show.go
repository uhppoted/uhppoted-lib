package acl

import (
	"codeberg.org/uhppoted/uhppoted-core/uhppote"
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
							permission := Permission{
								From:      card.From,
								To:        card.To,
								Profile:   profile,
								FirstCard: false,
							}

							switch {
							case ix == 1 && card.FirstCard.Door1:
								permission.FirstCard = true

							case ix == 2 && card.FirstCard.Door2:
								permission.FirstCard = true

							case ix == 3 && card.FirstCard.Door3:
								permission.FirstCard = true

							case ix == 4 && card.FirstCard.Door4:
								permission.FirstCard = true
							}

							acl[v.name] = permission
						}
					}
				}
			}
		}
	}

	return acl, nil
}
