package hitron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:govet
func TestWiFiModeString(t *testing.T) {
	testdata := []struct {
		in       WiFiMode
		expected string
	}{
		{WiFiModeB, "802.11b"},
		{WiFiModeG, "802.11g"},
		{WiFiModeN, "802.11n"},
		{WiFiModeA, "802.11a"},
		{WiFiModeAC, "802.11ac"},

		{WiFiModeB | WiFiModeG | WiFiModeN, "802.11b/g/n"},
		{WiFiModeG | WiFiModeN, "802.11g/n"},
		{WiFiModeA | WiFiModeN | WiFiModeAC, "802.11a/n/ac"},
	}

	for _, d := range testdata {
		d := d
		t.Run(d.expected, func(t *testing.T) {
			assert.Equal(t, d.expected, d.in.String())
		})
	}
}
