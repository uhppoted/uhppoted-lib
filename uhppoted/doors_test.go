package uhppoted

import (
	"testing"
)

func TestControlStateToString(t *testing.T) {
	if s := NormallyOpen.String(); s != "normally open" {
		t.Errorf("Invalid string for NormallyOpen - expected:%s, got:%s", "normally open", s)
	}

	if s := NormallyClosed.String(); s != "normally closed" {
		t.Errorf("Invalid string for NormallyClosed - expected:%s, got:%s", "normally closed", s)
	}

	if s := Controlled.String(); s != "controlled" {
		t.Errorf("Invalid string for Controlled - expected:%s, got:%s", "controlled", s)
	}
}
