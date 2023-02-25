package hitron

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWiFiAccessControl(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"blockType":"Block Listed","RuleNumberOfEntries":1,
		"Rules_List":[{"id":0,"hostName":"foo","macAddr":"AA:BB:CC:DD:EE:FF"}]
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiAccessControl(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiAccessControl{
		Error:     NoError,
		BlockType: "Block Listed",
		RulesList: []WiFiAccessControlRule{
			{ID: 0, Hostname: "foo", MACAddr: net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}},
		},
	}, p)
}

func TestWiFiAccessControlStatus(t *testing.T) {
	body := `{"errCode":"000","errMsg":"", "blockType":"Allow Listed"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiAccessControlStatus(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiAccessControlStatus{
		Error:     NoError,
		BlockType: "Allow Listed",
	}, p)
}

func TestWiFiGuestSSID(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"ssidName":"CODA-Guest","ssidName5G":"CODA-Guest-5G",
		"enable":"ON","pswd":"GuestPassword",
		"adminGuestAccProvider":"10"
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiGuestSSID(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiGuestSSID{
		Error:    NoError,
		Enable:   true,
		SSID:     "CODA-Guest",
		SSID5G:   "CODA-Guest-5G",
		Password: "GuestPassword",
		MaxUsers: 10,
	}, p)
}

func TestWiFiRadios(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"RadioNumberOfEntries":2,
		"Raidos_List":[
			{"vendor":"1","band":"2.4G","wlsOnOff":"ON","wlsDcsOnOff":"ON",
			"wlsMode":"4","n_bandwidth":"20/40MHZ","wlsChannel":"3",
			"autoChannel":"OFF","wlsDfsOnOff":"OFF","wlsCurrentChannel":"3",
			"wlswpsOnOff":"ON","igmpSnoop":"ON",
			"Radio_URI":"\/1\/Device\/WiFi\/Radios\/1","ssid_list":null
			},
			{"vendor":"1","band":"5G","wlsOnOff":"ON","wlsDcsOnOff":"OFF",
			"wlsMode":"9","n_bandwidth":"80MHZ","wlsChannel":0,
			"autoChannel":"ON","wlsDfsOnOff":"OFF","wlsCurrentChannel":"36",
			"wlswpsOnOff":"ON","igmpSnoop":"ON",
			"Radio_URI":"\/1\/Device\/WiFi\/Radios\/2","ssid_list":null}
		]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiRadios(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiRadios{
		Error: NoError,
		Radios: []WiFiRadio{
			{
				Vendor: "1", Band: "2.4G", Enable: true,
				AutoChannel:    false,
				Channel:        3,
				CurrentChannel: 3,
				EnableDCS:      true,
				Mode:           WiFiModeG | WiFiModeN,
				ChanBandwidth:  "20/40MHZ",
				EnableWPS:      true,
				IGMPSnoop:      true,
				RadioURI:       "/1/Device/WiFi/Radios/1",
			},
			{
				Vendor: "1", Band: "5G", Enable: true,
				AutoChannel:    true,
				Channel:        0,
				CurrentChannel: 36,
				EnableDCS:      false,
				Mode:           WiFiModeA | WiFiModeN | WiFiModeAC,
				ChanBandwidth:  "80MHZ",
				EnableWPS:      true,
				IGMPSnoop:      true,
				RadioURI:       "/1/Device/WiFi/Radios/2",
			},
		},
	}, p)
}

func TestWiFiRadioDetails(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"vendor":"1","band":"5G","wlsOnOff":"ON","wlsDcsOnOff":"OFF",
		"wlsMode":"9","n_bandwidth":"80MHZ","wlsChannel":0,"autoChannel":"ON",
		"wlsDfsOnOff":"OFF","wlsCurrentChannel":"36","wlswpsOnOff":"ON",
		"igmpSnoop":"ON","Radio_URI":"\/1\/Device\/WiFi\/Radios\/2",
		"ssid_list":""}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiRadioDetails(ctx, 2)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiRadio{
		Error:  NoError,
		Vendor: "1", Band: "5G", Enable: true,
		AutoChannel:    true,
		Channel:        0,
		CurrentChannel: 36,
		EnableDCS:      false,
		Mode:           WiFiModeA | WiFiModeN | WiFiModeAC,
		ChanBandwidth:  "80MHZ",
		EnableWPS:      true,
		IGMPSnoop:      true,
		RadioURI:       "/1/Device/WiFi/Radios/2",
	}, p)
}

//nolint:funlen
func TestWiFiRadiosAdvanced(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"AdvancedNumberOfEntries":2, "Advanced_List":[
			{"vendor":"1","band":"2.4G","wlsOnOff":"ON","wlsDcsOnOff":"ON",
			"wlsMode":"4","n_bandwidth":"20\/40MHZ","wlsChannel":"3",
			"autoChannel":"OFF","wlsDfsOnOff":"OFF","wlsCurrentChannel":"3",
			"wlswpsOnOff":"ON","igmpSnoop":"ON",
			"Radio_URI":"\/1\/Device\/WiFi\/Radios\/1","ssid_list":"",
			"bgMode":"","n_coexistence":"Enabled",
			"n_OperatingMode":"Mixed Mode","n_GuardInterval":"Long",
			"n_mcs":"0","n_rdg":"Disabled","n_amsdu":"Enabled",
			"n_autoba":"Enabled","n_badecline":"Disabled","tx_stream":"",
			"rx_stream":"","bandsteering":"ON","ssidName":"CODA",
			"showMSO":"false"},
			{"vendor":"1","band":"5G","wlsOnOff":"ON","wlsDcsOnOff":"OFF",
			"wlsMode":"9","n_bandwidth":"80MHZ","wlsChannel":0,
			"autoChannel":"ON","wlsDfsOnOff":"OFF","wlsCurrentChannel":"36",
			"wlswpsOnOff":"ON","igmpSnoop":"ON",
			"Radio_URI":"\/1\/Device\/WiFi\/Radios\/2","ssid_list":"",
			"bgMode":"","n_coexistence":"Enabled",
			"n_OperatingMode":"Mixed Mode","n_GuardInterval":"Long",
			"n_mcs":"0","n_rdg":"Disabled","n_amsdu":"Disabled",
			"n_autoba":"Enabled","n_badecline":"Disabled","tx_stream":"",
			"rx_stream":"","bandsteering":"ON","ssidName":"CODA",
			"showMSO":"false"}
	]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiRadiosAdvanced(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiRadiosAdvanced{
		Error: NoError,
		Radios: []WiFiRadioAdvanced{
			{
				WiFiRadio: WiFiRadio{
					Vendor: "1", Band: "2.4G", Enable: true,
					AutoChannel:    false,
					Channel:        3,
					CurrentChannel: 3,
					EnableDCS:      true,
					Mode:           WiFiModeG | WiFiModeN,
					ChanBandwidth:  "20/40MHZ",
					EnableWPS:      true,
					IGMPSnoop:      true,
					RadioURI:       "/1/Device/WiFi/Radios/1",
				},
				SSID:           "CODA",
				BGMode:         "",
				NCoexistence:   true,
				NOperatingMode: "Mixed Mode",
				NGuardInterval: "Long",
				NMCS:           0,
				NRDG:           false,
				Namsdu:         true,
				Nautoba:        true,
				Nbadecline:     false,
				TxStream:       "",
				RxStream:       "",
				BandSteering:   true,
				ShowMSO:        false,
			},
			{
				WiFiRadio: WiFiRadio{
					Vendor: "1", Band: "5G", Enable: true,
					AutoChannel:    true,
					Channel:        0,
					CurrentChannel: 36,
					EnableDCS:      false,
					Mode:           WiFiModeA | WiFiModeN | WiFiModeAC,
					ChanBandwidth:  "80MHZ",
					EnableWPS:      true,
					IGMPSnoop:      true,
					RadioURI:       "/1/Device/WiFi/Radios/2",
				},
				SSID:           "CODA",
				BGMode:         "",
				NCoexistence:   true,
				NOperatingMode: "Mixed Mode",
				NGuardInterval: "Long",
				NMCS:           0,
				NRDG:           false,
				Namsdu:         false,
				Nautoba:        true,
				Nbadecline:     false,
				TxStream:       "",
				RxStream:       "",
				BandSteering:   true,
				ShowMSO:        false,
			},
		},
	}, p)
}

func TestWiFiRadiosSurvey(t *testing.T) {
	body := `{"errCode":"000","errMsg":"", "APNumberOfEntries":2,
		"APs_List":[
			{"band":"2.4G","wlsChannel":"11",
			"ssidName":"","bssid":"ca:fe:de:ad:be:ef\n",
			"security":"WPA2","signal":"-35",
			"wmode":"IEEE80211_MODE_11NG_HT20","extch":"NONE","nt":"N\/A",
			"wps":"NO"},
			{"band":"5G","wlsChannel":"149",
			"ssidName":"CODA","bssid":"ca:fe:de:ad:fa:ce\n",
			"security":"WPA2","signal":"-80",
			"wmode":"IEEE80211_MODE_11AC_VHT80","extch":"NONE","nt":"N\/A",
			"wps":"NO"}
	]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiRadiosSurvey(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiRadiosSurvey{
		Error: NoError,
		APs: []WiFiAP{
			{
				Band: "2.4G", Channel: 11,
				SSID:  "",
				BSSID: net.HardwareAddr{0xca, 0xfe, 0xde, 0xad, 0xbe, 0xef},
				WMode: "IEEE80211_MODE_11NG_HT20",
				ExtCh: "NONE",
				NT:    "N/A",
			},
			{
				Band: "5G", Channel: 149,
				SSID:  "CODA",
				BSSID: net.HardwareAddr{0xca, 0xfe, 0xde, 0xad, 0xfa, 0xce},
				WMode: "IEEE80211_MODE_11AC_VHT80",
				ExtCh: "NONE",
				NT:    "N/A",
			},
		},
	}, p)
}

func TestWiFiSSIDs(t *testing.T) {
	body := ` {"errCode":"000","errMsg":"",
		"SSIDNumberOfEntries":2, "SSIDs_List":[
			{"id":"1","ssidName":"CODA","band":"2.4G","enable":"ON",
			"wlswpsOnOff":"ON","ifName":"ath0","bssid":"CA:FE:DE:AD:BE:EF",
			"radio":"1","visible":"ON","wmm":"ON","authMode":"4","SecuMode":"2",
			"encryptType":"3","passPhrase":"supersecret",
			"wlsEnable":"ON", "SSID_URI":"\/1\/Device\/WiFi\/SSIDs\/1",
			"defaultKey":"1234567890","bandsteer":"ON","primary":"YES"},
			{"id":"2","ssidName":"CODA","band":"5G","enable":"ON",
			"wlswpsOnOff":"ON","ifName":"ath1","bssid":"CA:FE:DE:AD:FA:CE",
			"radio":"2","visible":"ON","wmm":"ON","authMode":"4","SecuMode":"2",
			"encryptType":"3","passPhrase":"supersecret",
			"wlsEnable":"ON","SSID_URI":"\/1\/Device\/WiFi\/SSIDs\/2",
			"defaultKey":"1234567890","bandsteer":"ON","primary":"YES"}
		],
		"Guests_List":[{"enable":"OFF","ifName":"ath6","relate":"ath0"}]
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiSSIDs(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiSSIDs{
		Error: NoError,
		SSIDs: []SSID{
			{
				ID: 1, IfName: "ath0", Name: "CODA", Band: "2.4G",
				BSSID:  net.HardwareAddr{0xca, 0xfe, 0xde, 0xad, 0xbe, 0xef},
				Enable: true, EnableWPS: true, EnableWLS: true, EnableWMM: true,
				Visible: true, BandSteering: true, Primary: true,
				Radio: 1, AuthMode: "4", SecurityMode: "2", EncryptType: "3",
				Passphrase: "supersecret",
				URI:        "/1/Device/WiFi/SSIDs/1",
				DefaultKey: "1234567890",
			},
			{
				ID: 2, IfName: "ath1", Name: "CODA", Band: "5G",
				BSSID:  net.HardwareAddr{0xca, 0xfe, 0xde, 0xad, 0xfa, 0xce},
				Enable: true, EnableWPS: true, EnableWLS: true, EnableWMM: true,
				Visible: true, BandSteering: true, Primary: true,
				Radio: 2, AuthMode: "4", SecurityMode: "2", EncryptType: "3",
				Passphrase: "supersecret",
				URI:        "/1/Device/WiFi/SSIDs/2",
				DefaultKey: "1234567890",
			},
		},
		Guests: []GuestSSID{{Enable: false, IfName: "ath6", Relate: "ath0"}},
	}, p)
}

func TestWiFiWPS(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"wlswpsOnOff":"ON","wlsWpsMethod":"PushButton","wlsWpsClientPin":"",
		"wlsWpsStatus":"In Progress","wlsWpsTimeElapsed":"42"
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiWPS(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiWPS{
		Error:       NoError,
		Enable:      true,
		Method:      "PushButton",
		ClientPin:   "",
		Status:      "In Progress",
		TimeElapsed: 42 * time.Second,
	}, p)
}

func TestWiFiClient(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"ClientNumberOfEntries":2,
		"Client_List":[
			{"index":1,"band":"2.4G","ssid":"CODA","hostname":"foo",
			"mac":"CA:FE:DE:AD:BE:EF","aid":"1","rssi":"-84","br":"10M",
			"pm":"IEEE80211_MODE_11NG_HT20","ch":"3","bw":"20MHz"},
			{"index":2,"band":"5G","ssid":"CODA","hostname":"bar",
			"mac":"BE:EF:C0:FF:EE:99","aid":"2","rssi":"-28","br":"866M",
			"pm":"IEEE80211_MODE_11AC_VHT80","ch":"40","bw":"80MHz"}
		]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := testCableModem(srv)

	ctx := ContextWithDebugLogger(context.Background(), t)

	p, err := d.WiFiClient(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, WiFiClient{
		Error: NoError,
		Clients: []WiFiClientEntry{
			{
				Index:     1,
				AID:       1,
				Band:      "2.4G",
				SSID:      "CODA",
				Hostname:  "foo",
				MACAddr:   net.HardwareAddr{0xca, 0xfe, 0xde, 0xad, 0xbe, 0xef},
				RSSI:      -84,
				DataRate:  10 * mib,
				PhyMode:   "IEEE80211_MODE_11NG_HT20",
				Channel:   3,
				Bandwidth: 20_000_000,
			},
			{
				Index:     2,
				AID:       2,
				Band:      "5G",
				SSID:      "CODA",
				Hostname:  "bar",
				MACAddr:   net.HardwareAddr{0xbe, 0xef, 0xc0, 0xff, 0xee, 0x99},
				RSSI:      -28,
				DataRate:  866 * mib,
				PhyMode:   "IEEE80211_MODE_11AC_VHT80",
				Channel:   40,
				Bandwidth: 80_000_000,
			},
		},
	}, p)
}
