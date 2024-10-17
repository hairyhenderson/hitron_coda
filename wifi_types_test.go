package hitron

import (
	"encoding/json"
	"net"
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
		t.Run(d.expected, func(t *testing.T) {
			assert.Equal(t, d.expected, d.in.String())
		})
	}
}

func mustParseMac(in string) net.HardwareAddr {
	mac, err := net.ParseMAC(in)
	if err != nil {
		panic(err)
	}

	return mac
}

func TestWiFiClientEntry_UnmarshalJSON(t *testing.T) {
	testdata := []struct {
		in       []byte
		expected WiFiClientEntry
	}{
		{
			in: []byte(`{
				"index": 1, "band": "2.4G", "ssid": "XXX",
				"hostname": "with-int-aid", "mac":"AA:6F:99:8D:CB:A4",
				"aid": 1, "rssi": "-70", "br": "51M",
				"pm": "11NG_HT20", "ch": "11", "bw": "20MHz"
			}`),
			expected: WiFiClientEntry{
				Index:     1,
				Band:      "2.4G",
				SSID:      "XXX",
				Hostname:  "with-int-aid",
				MACAddr:   mustParseMac("AA:6F:99:8D:CB:A4"),
				AID:       1,
				RSSI:      -70,
				DataRate:  51 * 1024 * 1024,
				PhyMode:   "11NG_HT20",
				Channel:   11,
				Bandwidth: 20_000_000,
			},
		},
		{
			in: []byte(`{
				"index": 1, "band": "2.4G", "ssid": "XXX",
				"hostname": "with-string-aid", "mac":"AA:6F:99:8D:CB:A4",
				"aid": "42", "rssi": "-70", "br": "51M",
				"pm": "11NG_HT20", "ch": "11", "bw": "20MHz"
			}`),
			expected: WiFiClientEntry{
				Index:     1,
				Band:      "2.4G",
				SSID:      "XXX",
				Hostname:  "with-string-aid",
				MACAddr:   mustParseMac("AA:6F:99:8D:CB:A4"),
				AID:       42,
				RSSI:      -70,
				DataRate:  51 * 1024 * 1024,
				PhyMode:   "11NG_HT20",
				Channel:   11,
				Bandwidth: 20_000_000,
			},
		},
	}

	for _, d := range testdata {
		entry := WiFiClientEntry{}
		err := json.Unmarshal(d.in, &entry)
		assert.NoError(t, err)
		assert.Equal(t, d.expected, entry)
	}
}
