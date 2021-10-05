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
		err = fmt.Errorf("Invalid DeviceID: %v", v)
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

type EventRange struct {
	First *uint32 `json:"first,omitempty"`
	Last  *uint32 `json:"last,omitempty"`
}

func (e EventRange) String() string {
	first := "-"
	last := "-"

	if e.First != nil && *e.First != 0 {
		first = fmt.Sprintf("%v", *e.First)
	}

	if e.Last != nil && *e.Last != 0 {
		last = fmt.Sprintf("%v", *e.Last)
	}

	return fmt.Sprintf("{ First:%v, Last:%v }", first, last)
}

type EventIndex uint32

func (index EventIndex) increment(rollover uint32) EventIndex {
	ix := uint32(index)

	if ix < 1 {
		ix = 1
	} else if ix >= rollover {
		ix = 1
	} else {
		ix += 1
	}

	return EventIndex(ix)
}

func (index EventIndex) decrement(rollover uint32) EventIndex {
	ix := uint32(index)

	if ix <= 1 {
		ix = rollover
	} else if ix > rollover {
		ix = rollover
	} else {
		ix -= 1
	}

	return EventIndex(ix)
}
