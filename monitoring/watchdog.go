package monitoring

import (
	"fmt"
	"math"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

type Watchdog struct {
	healthcheck *HealthCheck
	state       struct {
		Started     time.Time
		HealthCheck struct {
			Alerted bool
		}
	}
}

func NewWatchdog(h *HealthCheck) Watchdog {
	return Watchdog{
		healthcheck: h,
		state: struct {
			Started     time.Time
			HealthCheck struct {
				Alerted bool
			}
		}{
			Started: time.Now(),
			HealthCheck: struct {
				Alerted bool
			}{
				Alerted: false,
			},
		},
	}
}

func (w *Watchdog) ID() string {
	return "watchdog"
}

func (w *Watchdog) Exec(handler MonitoringHandler) error {
	debugf("watchdog", "exec")

	warnings := uint(0)
	errors := uint(0)
	healthCheckRunning := false

	// Verify health-check

	delay := math.Max(MIN_DELAY, w.healthcheck.interval.Seconds()+PADDING)

	dt := time.Since(w.state.Started).Round(time.Second)
	if w.healthcheck.state.Touched != nil {
		dt = time.Since(*w.healthcheck.state.Touched)
		if math.Abs(dt.Seconds()) < delay {
			healthCheckRunning = true
		}
	}

	if math.Abs(dt.Seconds()) > delay {
		errors += 1
		if !w.state.HealthCheck.Alerted {
			msg := fmt.Sprintf("'health-check' subsystem has not run since %v (%v)", types.DateTime(w.state.Started), dt)

			errorf("watchdog", msg)
			if err := handler.Alert(w, msg); err == nil {
				w.state.HealthCheck.Alerted = true
			}
		}
	} else {
		if w.state.HealthCheck.Alerted {
			infof("watchdog", "'health-check' subsystem is running")
			w.state.HealthCheck.Alerted = false
		}
	}

	// Report on known devices
	if healthCheckRunning {
		warnings += w.healthcheck.state.Warnings
		errors += w.healthcheck.state.Errors
	}

	// 'k, done
	if errors > 0 && warnings > 0 {
		warnf("watchdog", "%v, %v", Errors(errors), Warnings(warnings))
	} else if errors > 0 {
		warnf("watchdog", "%%v", Errors(errors))
	} else if warnings > 0 {
		warnf("watchdog", "%v", Warnings(warnings))
	} else {
		infof("watchdog", "OK")
	}

	if errors > 0 && warnings > 0 {
		handler.Alive(w, fmt.Sprintf("%s, %s", Errors(errors), Warnings(warnings)))
	} else if errors > 0 {
		handler.Alive(w, fmt.Sprintf("%v", Errors(errors)))
	} else if warnings > 0 {
		handler.Alive(w, fmt.Sprintf("%v", Warnings(warnings)))
	} else {
		handler.Alive(w, "OK")
	}

	return nil
}
