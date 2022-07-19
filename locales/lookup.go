package locales

import (
	"encoding/json"
	"io/fs"

	"github.com/uhppoted/uhppoted-lib/locales/en"
)

var dictionary = en.Dictionary

func Load(f fs.FS, file string) error {
	d := struct {
		Dictionary map[string]string `json:"dictionary"`
	}{}

	if bytes, err := fs.ReadFile(f, file); err != nil {
		return err
	} else if err := json.Unmarshal(bytes, &d); err != nil {
		return err
	}

	for k, v := range d.Dictionary {
		dictionary[k] = v
	}

	return nil
}

func Lookup(key string) (string, bool) {
	if v, ok := dictionary[key]; ok {
		return v, true
	}

	return "", false
}
