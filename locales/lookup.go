package locales

import (
	"github.com/uhppoted/uhppoted-lib/locales/en"
)

func Lookup(key string) (string, bool) {
	if v, ok := en.Dictionary[key]; ok {
		return v, true
	}

	return "", false
}
