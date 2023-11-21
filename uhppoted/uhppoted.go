package uhppoted

import (
	"errors"
	"net/http"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/log"
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
}

func (u *UHPPOTED) debug(tag string, msg interface{}) {
	log.Debugf("%-12s %v", tag, msg)
}

// func (u *UHPPOTED) info(tag string, msg interface{}) {
// 	log.Infof("%-12s %v", tag, msg)
// }

func (u *UHPPOTED) warn(tag string, err error) {
	log.Warnf("%-12s %v", tag, err)
}
