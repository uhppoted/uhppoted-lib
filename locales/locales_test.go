package locales

import (
	"embed"
	"testing"
)

//go:embed test.json

var testfs embed.FS

func TestLoadLocale(t *testing.T) {
	if err := Load(testfs, "test.json"); err != nil {
		t.Fatalf("%v", err)
	}

	expected := map[string]string{
		"event.type.0": "test.none",
		"event.type.1": "test.swipe",
		"event.type.2": "door",
		"event.type.3": "test.alarm",
	}

	for k, v := range expected {
		if w, ok := dictionary[k]; !ok || w != v {
			t.Errorf("Incorrect dictionary entry for '%v' - expected:'%v', got:'%v'", k, v, w)
		}
	}
}
