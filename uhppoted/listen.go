package uhppoted

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/log"
	lib "github.com/uhppoted/uhppoted-lib/os"
)

type EventMap struct {
	file      string
	retrieved map[uint32]uint32
}

type EventHandler func(Event) bool

type listener struct {
	onConnected func()
	onEvent     func(*types.Status)
	onError     func(error) bool
}

func (l *listener) OnConnected() {
	go func() {
		l.onConnected()
	}()
}

func (l *listener) OnEvent(event *types.Status) {
	go func() {
		l.onEvent(event)
	}()
}

func (l *listener) OnError(err error) bool {
	return l.onError(err)
}

const BATCHSIZE = 32

func (u *UHPPOTED) Listen(handler EventHandler, received *EventMap, q chan os.Signal) {
	var wg sync.WaitGroup

	devices := u.UHPPOTE.DeviceList()

	for _, d := range devices {
		deviceID := d.ID()
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.retrieve(deviceID, received, handler)
		}()
	}

	wg.Wait()

	u.listen(handler, received, q)
}

func (u *UHPPOTED) retrieve(deviceID uint32, received *EventMap, handler EventHandler) {
	if index, ok := received.retrieved[deviceID]; ok {
		u.info("listen", fmt.Sprintf("Fetching unretrieved events for device ID %v", deviceID))

		event, err := u.UHPPOTE.GetEvent(deviceID, 0xffffffff)
		if err != nil {
			u.warn("listen", fmt.Errorf("unable to retrieve events for device ID %v (%w)", deviceID, err))
			return
		}

		if event.Index == uint32(index) {
			u.info("listen", fmt.Sprintf("No unretrieved events for device ID %v", deviceID))
			return
		}

		from := index
		to := event.Index

		if retrieved := u.fetch(deviceID, from+1, to, handler); retrieved != 0 {
			received.retrieved[deviceID] = retrieved
			if err := received.store(); err != nil {
				u.warn("listen", err)
			}
		}
	}
}

func (u *UHPPOTED) listen(handler EventHandler, received *EventMap, q chan os.Signal) {
	u.info("listen", "Initialising event listener")

	backoffs := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}

	ix := 0
	l := listener{
		onConnected: func() {
			u.info("listen", "Listening")
			ix = 0
		},

		onEvent: func(e *types.Status) {
			u.onEvent(e, received, handler)
		},

		onError: func(err error) bool {
			u.warn("listen", err)
			return true
		},
	}

	// NTS: use 'for {..}' because 'for err := u.UHPPOTE.Listen; ..' only ever executes the
	//      'Listen' once - on loop initialization
	for {
		if err := u.UHPPOTE.Listen(&l, q); err != nil {
			u.warn("listen", err)

			delay := 60 * time.Second
			if ix < len(backoffs) {
				delay = backoffs[ix]
				ix++
			}

			u.info("listen", fmt.Sprintf("Retrying in %v", delay))
			time.Sleep(delay)
			continue
		}

		break
	}
}

func (u *UHPPOTED) onEvent(e *types.Status, received *EventMap, handler EventHandler) {
	u.info("event", fmt.Sprintf("%+v", e))

	deviceID := uint32(e.SerialNumber)
	last := e.Event.Index
	first := e.Event.Index

	retrieved, ok := received.retrieved[deviceID]
	if ok && retrieved != uint32(last) {
		first = retrieved
	}

	if eventID := u.fetch(deviceID, first, last, handler); eventID != 0 {
		received.retrieved[deviceID] = eventID
		if err := received.store(); err != nil {
			u.warn("listen", err)
		}
	}
}

func (u *UHPPOTED) fetch(deviceID uint32, from, to uint32, handler EventHandler) (retrieved uint32) {
	batchSize := BATCHSIZE
	if u.ListenBatchSize > 0 {
		batchSize = u.ListenBatchSize
	}

	first, err := u.UHPPOTE.GetEvent(deviceID, 0)
	if err != nil {
		u.warn("listen", fmt.Errorf("failed to retrieve 'first' event for device %d (%w)", deviceID, err))
		return
	} else if first == nil {
		u.warn("listen", fmt.Errorf("no 'first' event record returned for device %d", deviceID))
		return
	}

	last, err := u.UHPPOTE.GetEvent(deviceID, 0xffffffff)
	if err != nil {
		u.warn("listen", fmt.Errorf("failed to retrieve 'last' event for device %d (%w)", deviceID, err))
		return
	} else if first == nil {
		u.warn("listen", fmt.Errorf("no 'last' event record returned for device %d", deviceID))
		return
	}

	if last.Index >= first.Index {
		if uint32(from) < first.Index || uint32(from) > last.Index {
			from = first.Index
		}

		if uint32(to) < first.Index || uint32(to) > last.Index {
			to = last.Index
		}
	} else {
		if uint32(from) < first.Index && uint32(from) > last.Index {
			from = first.Index
		}

		if uint32(to) < first.Index && uint32(to) > last.Index {
			to = last.Index
		}
	}

	count := 0
	index := from
	for {
		count += 1
		if count > batchSize {
			return
		}

		record, err := u.UHPPOTE.GetEvent(deviceID, uint32(index))
		if err != nil {
			u.warn("listen", fmt.Errorf("failed to retrieve event for device %d, ID %d (%w)", deviceID, index, err))
		} else if record == nil {
			u.warn("listen", fmt.Errorf("no event record for device %d, ID %d", deviceID, index))
		} else if record.Index != uint32(index) {
			u.warn("listen", fmt.Errorf("no event record for device %d, ID %d", deviceID, index))
		} else {
			event := Event{
				DeviceID:   uint32(record.SerialNumber),
				Index:      record.Index,
				Type:       record.Type,
				Granted:    record.Granted,
				Door:       record.Door,
				Direction:  record.Direction,
				CardNumber: record.CardNumber,
				Timestamp:  record.Timestamp,
				Reason:     record.Reason,
			}

			if !handler(event) {
				break
			}

			retrieved = record.Index
		}

		if index == to {
			break
		}

		index++
	}

	return
}

func NewEventMap(file string) *EventMap {
	return &EventMap{
		file:      file,
		retrieved: map[uint32]uint32{},
	}
}

func (m *EventMap) Load() error {
	if m.file == "" {
		return nil
	}

	f, err := os.Open(m.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer f.Close()

	re := regexp.MustCompile(`^\s*(.*?)(?::\s*|\s*=\s*|\s+)(\S.*)\s*`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		match := re.FindStringSubmatch(s.Text())
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])

			if device, err := strconv.ParseUint(key, 10, 32); err != nil {
				log.Warnf("error parsing event map entry '%s': %v", s.Text(), err)
			} else if eventID, err := strconv.ParseUint(value, 10, 32); err != nil {
				log.Warnf("error parsing event map entry '%s': %v", s.Text(), err)
			} else {
				m.retrieved[uint32(device)] = uint32(eventID)
			}
		}
	}

	return s.Err()
}

func (m *EventMap) store() error {
	if m.file == "" || IsDevNull(m.file) {
		return nil
	}

	f, err := os.CreateTemp(os.TempDir(), "uhppoted*.tmp")
	if err != nil {
		return err
	}

	defer os.Remove(f.Name())

	for key, value := range m.retrieved {
		if _, err := fmt.Fprintf(f, "%-16d %v\n", key, value); err != nil {
			f.Close()
			return err
		}
	}

	f.Close()

	return lib.Rename(f.Name(), m.file)
}
