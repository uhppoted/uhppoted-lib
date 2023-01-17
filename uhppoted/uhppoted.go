package uhppoted

import (
	"errors"
	"log"
	"net/http"

	"github.com/uhppoted/uhppote-core/uhppote"
)

const (
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusNotFound            = http.StatusNotFound
	StatusUnauthorized        = http.StatusUnauthorized
	StatusInternalServerError = http.StatusInternalServerError
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("not authorized")
	ErrInternalServerError = errors.New("internal server error")
)

type UHPPOTED struct {
	UHPPOTE         uhppote.IUHPPOTE
	ListenBatchSize int
	Log             *log.Logger
}

func (u *UHPPOTED) debug(tag string, msg interface{}) {
	if u != nil && u.Log != nil {
		u.Log.Printf("DEBUG  %-12s %v", tag, msg)
	}
}

func (u *UHPPOTED) info(tag string, msg interface{}) {
	if u != nil && u.Log != nil {
		u.Log.Printf("INFO   %-12s %v", tag, msg)
	}
}

func (u *UHPPOTED) warn(tag string, err error) {
	if u != nil && u.Log != nil {
		u.Log.Printf("WARN   %-12s %v", tag, err)
	}
}
