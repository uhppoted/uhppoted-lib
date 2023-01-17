package uhppoted

import (
	"encoding/json"
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

type DeviceID uint32

func (id *DeviceID) UnmarshalJSON(bytes []byte) (err error) {
	v := uint32(0)

	if err = json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	if v == 0 {
		err = fmt.Errorf("invalid device ID: %v", v)
		return
	}

	*id = DeviceID(v)

	return
}

type DateRange struct {
	Start *types.DateTime `json:"start,omitempty"`
	End   *types.DateTime `json:"end,omitempty"`
}

func (d *DateRange) String() string {
	if d.Start != nil && d.End != nil {
		return fmt.Sprintf("{ Start:%v, End:%v }", d.Start, d.End)
	}

	if d.Start != nil {
		return fmt.Sprintf("{ Start:%v }", d.Start)
	}

	if d.End != nil {
		return fmt.Sprintf("{ End:%v }", d.End)
	}

	return "{}"
}
