package uhppoted

import (
	"testing"
)

func TestName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"authorise", "authorise"},
		{"authorise|authorize", "authorise"},
	}

	for _, test := range tests {
		if n := name(test.name); n != test.expected {
			t.Errorf("incorrect name() - expected:%v, got:%v", test.expected, n)
		}
	}
}

func TestAltName(t *testing.T) {
	tests := []struct {
		arg      string
		name     string
		expected bool
	}{
		{"authorise", "authorise", true},
		{"authorise", "authorize", false},
		{"authorise", "authorise|authorize", true},
		{"authorize", "authorise|authorize", true},
	}

	for _, test := range tests {
		if ok := alt(test.name, test.arg); ok != test.expected {
			t.Errorf("incorrect alt(%v) - expected:%v, got:%v", test.name, test.expected, ok)
		}
	}
}
