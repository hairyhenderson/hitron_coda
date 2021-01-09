package hitron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDHCPLeaseTime(t *testing.T) {
	testdata := []struct {
		in  string
		dur time.Duration
	}{
		{"", 0},
		{"S: 30", 30 * time.Second},
		{"D: 0 H: 1 M: 2 S: 3", time.Hour + 2*time.Minute + 3*time.Second},
		{"D: 0 H: 00 M: 00 S: 00", 0},
		{"D: 6 H: 09 M: 25 S: 55", 6*24*time.Hour + 9*time.Hour + 25*time.Minute + 55*time.Second},
		{"D: 6 weird {corrupt{ ][entry", 6 * 24 * time.Hour},
	}

	for _, d := range testdata {
		d := d
		t.Run(d.in, func(t *testing.T) {
			out := parseDHCPLeaseDuration(d.in)
			assert.Equal(t, d.dur, out)
		})
	}
}
