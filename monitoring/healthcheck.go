package monitoring

import (
	"fmt"
	"math"
	"net/netip"
	"sync"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

type HealthCheck struct {
	uhppote    uhppote.IUHPPOTE
	interval   time.Duration
	idleTime   time.Duration
	ignoreTime time.Duration
	state      struct {
		Started time.Time
		Touched *time.Time
		Devices struct {
			Status   sync.Map
			Listener sync.Map
			Errors   sync.Map
		}
		Warnings uint
		Errors   uint
	}
}

type status struct {
	Touched time.Time
	Status  types.Status
}

type listener struct {
	Touched time.Time
	Address netip.AddrPort
}

type alerts struct {
	missing      bool
	unexpected   bool
	touched      bool
	synchronized bool
	nolistener   bool
	listener     bool
}

var cache = struct {
	addresses []netip.AddrPort
	touched   time.Time
	sync.RWMutex
}{}

func NewHealthCheck(u uhppote.IUHPPOTE, interval, idleTime, ignoreTime time.Duration) HealthCheck {
	return HealthCheck{
		uhppote:    u,
		interval:   interval,
		idleTime:   idleTime,
		ignoreTime: ignoreTime,
		state: struct {
			Started time.Time
			Touched *time.Time
			Devices struct {
				Status   sync.Map
				Listener sync.Map
				Errors   sync.Map
			}
			Warnings uint
			Errors   uint
		}{
			Started: time.Now(),
			Touched: nil,
			Devices: struct {
				Status   sync.Map
				Listener sync.Map
				Errors   sync.Map
			}{
				Status:   sync.Map{},
				Listener: sync.Map{},
				Errors:   sync.Map{},
			},
			Warnings: 0,
			Errors:   0,
		},
	}
}

func (h *HealthCheck) ID() string {
	return "health-check"
}

func (h *HealthCheck) Exec(handler MonitoringHandler) {
	debugf("health-check", "exec")

	now := time.Now()
	errors := uint(0)
	warnings := uint(0)

	h.update(now)

	e, w := h.known(now, handler)
	errors += e
	warnings += w

	e, w = h.unexpected(now, handler)
	errors += e
	warnings += w

	h.state.Warnings = warnings
	h.state.Errors = errors

	// 'k, done
	if errors > 0 && warnings > 0 {
		warnf("health-check", "%v, %v", Errors(errors), Warnings(warnings))
	} else if errors > 0 {
		warnf("health-check", "%v", Errors(errors))
	} else if warnings > 0 {
		warnf("health-check", "%v", Warnings(warnings))
	} else {
		infof("health-check", "OK")
	}

	if errors > 0 && warnings > 0 {
		handler.Alive(h, fmt.Sprintf("%v, %v", Errors(errors), Warnings(warnings)))
	} else if errors > 0 {
		handler.Alive(h, fmt.Sprintf("%v", Errors(errors)))
	} else if warnings > 0 {
		handler.Alive(h, fmt.Sprintf("%v", Warnings(warnings)))
	} else {
		handler.Alive(h, "OK")
	}
}

func (h *HealthCheck) update(now time.Time) {
	if dt := time.Since(cache.touched); dt > 60*time.Second {
		go h.resolve()
	}

	devices := make(map[uint32]bool)

	if found, err := h.uhppote.GetDevices(); err != nil {
		warnf("health-check", "'keep-alive' error: %v", err)
	} else {
		for _, id := range found {
			devices[uint32(id.SerialNumber)] = true
		}
	}

	for id := range h.uhppote.DeviceList() {
		devices[id] = true
	}

	for id := range devices {
		s, err := h.uhppote.GetStatus(id)
		if err == nil {
			h.state.Devices.Status.Store(id, status{
				Status:  *s,
				Touched: now,
			})
		}

		l, _, err := h.uhppote.GetListener(id)
		if err == nil && l.IsValid() {
			h.state.Devices.Listener.Store(id, listener{
				Address: l,
				Touched: now,
			})
		}
	}

	h.state.Touched = &now
}

// Check known/identified devices
func (h *HealthCheck) known(now time.Time, handler MonitoringHandler) (uint, uint) {
	errors := uint(0)
	warnings := uint(0)

	for id := range h.uhppote.DeviceList() {
		alerted := alerts{
			missing:      false,
			unexpected:   false,
			touched:      false,
			synchronized: false,
			nolistener:   false,
			listener:     false,
		}

		if v, found := h.state.Devices.Errors.Load(id); found {
			alerted.missing = v.(alerts).missing
			alerted.unexpected = v.(alerts).unexpected
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
			alerted.nolistener = v.(alerts).nolistener
			alerted.listener = v.(alerts).listener
		}

		if _, found := h.state.Devices.Status.Load(id); !found {
			errors += 1
			if !alerted.missing {
				if alert(h, handler, id, "device not found") {
					alerted.missing = true
				}
			}
		} else {
			if alerted.missing {
				if info(h, handler, id, "device present") {
					alerted.missing = false
				}
			}
		}

		e, w := h.checkStatus(id, now, &alerted, handler, true)
		errors += e
		warnings += w

		e, w = h.checkListener(id, now, &alerted, handler, true)
		errors += e
		warnings += w

		h.state.Devices.Errors.Store(id, alerted)
	}

	return errors, warnings
}

// Identify and check any unexpected devices
func (h *HealthCheck) unexpected(now time.Time, handler MonitoringHandler) (uint, uint) {
	errors := uint(0)
	warnings := uint(0)

	f := func(key, value interface{}) bool {
		alerted := alerts{
			missing:      false,
			unexpected:   false,
			touched:      false,
			synchronized: false,
		}

		if v, found := h.state.Devices.Errors.Load(key); found {
			alerted.missing = v.(alerts).missing
			alerted.unexpected = v.(alerts).unexpected
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
			alerted.nolistener = v.(alerts).nolistener
			alerted.listener = v.(alerts).listener
		}

		for id := range h.uhppote.DeviceList() {
			if id == key {
				if alerted.unexpected {
					if alert(h, handler, key.(uint32), "added to configuration") {
						alerted.unexpected = false
						h.state.Devices.Errors.Store(id, alerted)
					}
				}

				return true
			}
		}

		touched := value.(status).Touched
		if now.After(touched.Add(h.ignoreTime)) {
			h.state.Devices.Status.Delete(key)
			h.state.Devices.Errors.Delete(key)

			if alerted.unexpected {
				warn(h, handler, key.(uint32), "disappeared")
			}
		} else {
			warnings += 1
			if !alerted.unexpected {
				if warn(h, handler, key.(uint32), "unexpected device") {
					alerted.unexpected = true
				}
			}

			e, w := h.checkStatus(key.(uint32), now, &alerted, handler, false)
			errors += e
			warnings += w

			e, w = h.checkListener(key.(uint32), now, &alerted, handler, false)
			errors += e
			warnings += w

			h.state.Devices.Errors.Store(key, alerted)
		}

		return true
	}

	h.state.Devices.Status.Range(f)

	return errors, warnings
}

func (h *HealthCheck) checkStatus(id uint32, now time.Time, alerted *alerts, handler MonitoringHandler, known bool) (uint, uint) {
	sysdatetime := func(v status) time.Time {
		if t := v.Status.SystemDateTime; t.IsZero() {
			return time.Time{}
		} else {
			return time.Time(t)
		}
	}

	errors := uint(0)
	warnings := uint(0)

	if v, found := h.state.Devices.Status.Load(id); found {
		touched := v.(status).Touched
		t := sysdatetime(v.(status))
		dt := time.Since(t).Round(time.Second)
		dtt := int64(math.Abs(time.Since(touched).Seconds()))

		if now.After(touched.Add(h.idleTime)) {
			if known {
				errors += 1
			} else {
				warnings += 1
			}

			if !alerted.touched {
				msg := fmt.Sprintf("no response for %s", time.Since(touched).Round(time.Second))
				if alert(h, handler, id, msg) {
					alerted.touched = true
					alerted.synchronized = false
				}
			}
		} else {
			if alerted.touched {
				if info(h, handler, id, "connected") {
					alerted.touched = false
				}
			}
		}

		if dtt < DELTA/2 {
			if int64(math.Abs(dt.Seconds())) > DELTA {
				if known {
					errors += 1
				} else {
					warnings += 1
				}

				if !alerted.synchronized {
					msg := fmt.Sprintf("system time not synchronized:%v (%v)", types.DateTime(t), dt)
					if alert(h, handler, id, msg) {
						alerted.synchronized = true
					}
				}
			} else {
				if alerted.synchronized {
					msg := fmt.Sprintf("system time synchronized:%v (%v)", types.DateTime(t), dt)
					if info(h, handler, id, msg) {
						alerted.synchronized = false
					}
				}
			}
		}
	}

	return errors, warnings
}

func (h *HealthCheck) checkListener(id uint32, now time.Time, alerted *alerts, handler MonitoringHandler, known bool) (uint, uint) {
	errors := uint(0)
	warnings := uint(0)

	if v, found := h.state.Devices.Listener.Load(id); found {
		address := v.(listener).Address
		touched := v.(listener).Touched

		if now.After(touched.Add(h.idleTime)) {
			if known {
				errors += 1
			} else {
				warnings += 1
			}

			if !alerted.nolistener {
				msg := fmt.Sprintf("no reply to 'get-listener' for %s", time.Since(touched).Round(time.Second))
				if warn(h, handler, id, msg) {
					alerted.nolistener = true
				}
			}
		} else {
			if alerted.nolistener {
				if info(h, handler, id, "listener identified") {
					alerted.nolistener = false
				}
			}
		}

		cache.RLock()
		defer cache.RUnlock()

		if len(cache.addresses) == 0 {
			return errors, warnings
		}

		for _, expected := range cache.addresses {
			if expected == address {
				if alerted.listener {
					if info(h, handler, id, "listener address/port correct") {
						alerted.listener = false
					}
				}

				return errors, warnings
			}
		}

		if known {
			errors += 1
		} else {
			warnings += 1
		}

		if !alerted.listener {
			msg := fmt.Sprintf("incorrect listener address/port: %s", &address)
			if warn(h, handler, id, msg) {
				alerted.listener = true
			}
		}
	}

	return errors, warnings
}

func (h *HealthCheck) resolve() {
	infof("health-check", "refreshing interface IP address list")

	list := h.uhppote.ListenAddrList()

	cache.Lock()
	defer cache.Unlock()

	cache.addresses = list
	cache.touched = time.Now()
}

func info(h *HealthCheck, handler MonitoringHandler, deviceID uint32, message string) bool {
	msg := fmt.Sprintf("UTC0311-L0x %s %s", types.SerialNumber(deviceID), message)

	infof("health-check", msg)
	if err := handler.Alert(h, msg); err != nil {
		return false
	}

	return true
}

func warn(h *HealthCheck, handler MonitoringHandler, deviceID uint32, message string) bool {
	msg := fmt.Sprintf("UTC0311-L0x %s %s", types.SerialNumber(deviceID), message)

	warnf("health-check", msg)
	if err := handler.Alert(h, msg); err != nil {
		return false
	}

	return true
}

func alert(h *HealthCheck, handler MonitoringHandler, deviceID uint32, message string) bool {
	msg := fmt.Sprintf("UTC0311-L0x %s %s", types.SerialNumber(deviceID), message)
	known := false

	for id := range h.uhppote.DeviceList() {
		if deviceID == id {
			known = true
		}
	}

	if known {
		errorf("health-check", msg)
	} else {
		warnf("health-check", msg)
	}

	if err := handler.Alert(h, msg); err != nil {
		return false
	}

	return true
}
