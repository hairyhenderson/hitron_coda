package hitron

import (
	"encoding/json"
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
		t.Run(d.in, func(t *testing.T) {
			out := parseDHCPLeaseDuration(d.in)
			assert.Equal(t, d.dur, out)
		})
	}
}

func TestPortInfoUnmarshalJSON(t *testing.T) {
	in := `{
	"portId": "1",
	"frequency": "615000000",
	"modulation": "QAM256",
	"signalStrength": "2.500",
	"snr": "37.636",
	"channelId": "11",
	"dsoctets": "20883398",
	"correcteds": "1",
	"uncorrect": "0"
}`
	expected := &PortInfo{
		PortID:         "1",
		Frequency:      615000000,
		Modulation:     "QAM256",
		SignalStrength: 2.5,
		ChannelID:      "11",
		SNR:            37.636,
		DsOctets:       20883398,
		Correcteds:     1,
		Uncorrect:      0,
	}
	out := &PortInfo{}
	err := json.Unmarshal([]byte(in), out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)

	in = `{
	"portId": "8",
	"frequency": "",
	"modulation": "QAM256",
	"signalStrength": false,
	"snr": false,
	"channelId": "0",
	"dsoctets": "0",
	"correcteds": "0",
	"uncorrect": "0"
}`
	out = &PortInfo{}
	err = json.Unmarshal([]byte(in), out)
	assert.Error(t, err)
}

func TestCMDsInfoUnmarshalJSON(t *testing.T) {
	in := `{
	"errCode": "000",
	"errMsg": "",
	"Freq_List": [
		{
		"portId": "1",
		"frequency": "615000000",
		"modulation": "QAM256",
		"signalStrength": "2.500",
		"snr": "37.636",
		"channelId": "11",
		"dsoctets": "20883398",
		"correcteds": "1",
		"uncorrect": "0"
		},
		{
		"portId": "8",
		"frequency": "",
		"modulation": "QAM256",
		"signalStrength": false,
		"snr": false,
		"channelId": "0",
		"dsoctets": "0",
		"correcteds": "0",
		"uncorrect": "0"
		}
	]
}`
	expected := &CMDsInfo{
		Error: NoError,
		Ports: []PortInfo{
			{
				PortID:         "1",
				Frequency:      615000000,
				Modulation:     "QAM256",
				SignalStrength: 2.5,
				ChannelID:      "11",
				SNR:            37.636,
				DsOctets:       20883398,
				Correcteds:     1,
				Uncorrect:      0,
			},
		},
	}

	out := &CMDsInfo{}
	err := json.Unmarshal([]byte(in), out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}
