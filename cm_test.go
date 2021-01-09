package hitron

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"deviceId":"74:9B:DE:AD:BE:EF","modelName":"CODA-4680-TPIA",
		"vendorName":"Hitron Technologies","SerialNum":"12345678901",
		"HwVersion":"1A","ApiVersion":"1.11","SoftwareVersion":"7.1.1.2.2b9"
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	v, err := d.CMVersion(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMVersion{
		Error:           NoError,
		DeviceID:        "74:9B:DE:AD:BE:EF",
		ModelName:       "CODA-4680-TPIA",
		VendorName:      "Hitron Technologies",
		SerialNum:       "12345678901",
		HwVersion:       "1A",
		APIVersion:      "1.11",
		SoftwareVersion: "7.1.1.2.2b9",
	}, v)
}

func TestCMDocsisProvision(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"hwInit":"Success","findDownstream":"Success","ranging":"Success",
		"dhcp":"Success","timeOfday":"Success","downloadCfg":"Success",
		"registration":"Success","eaeStatus":"Disable",
		"bpiStatus":"AUTH:start, TEK:start",
		"networkAccess":"Permitted","trafficStatus":"Enable"
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	p, err := d.CMDocsisProvision(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMDocsisProvision{
		Error:          NoError,
		HWInit:         "Success",
		FindDownstream: "Success",
		Ranging:        "Success",
		DHCP:           "Success",
		TimeOfday:      "Success",
		DownloadCfg:    "Success",
		Registration:   "Success",
		EAEStatus:      "Disable",
		BPIStatus:      "AUTH:start, TEK:start",
		NetworkAccess:  "Permitted",
		TrafficStatus:  "Enable",
	}, p)
}

func TestCMDsInfo(t *testing.T) {
	body := `{"errCode":"000","errMsg":"","Freq_List":[
		{"portId":"1","frequency":"615000000","modulation":"QAM256",
		"signalStrength":"3.200","snr":"38.605","channelId":"11",
		"dsoctets":"4829493","correcteds":"0","uncorrect":"0"},
		{"portId":"2","frequency":"603000000","modulation":"QAM256",
		"signalStrength":"2.000","snr":"37.636","channelId":"9",
		"dsoctets":"4960133","correcteds":"0","uncorrect":"0"}
	]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	p, err := d.CMDsInfo(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMDsInfo{
		Error: Error{"000", ""},
		Ports: []PortInfo{
			{
				PortID:         "1",
				Frequency:      615000000,
				Modulation:     "QAM256",
				SignalStrength: 3.2,
				SNR:            38.605,
				ChannelID:      "11",
				DsOctets:       4829493,
			},
			{
				PortID:         "2",
				Frequency:      603000000,
				Modulation:     "QAM256",
				SignalStrength: 2.0,
				SNR:            37.636,
				ChannelID:      "9",
				DsOctets:       4960133,
			},
		},
	}, p)
}

func TestCMUsInfo(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
	"Freq_List":[
		{"portId":"1","frequency":"32300000","modulationType":"64QAM",
		"signalStrength":"45.270","bandwidth":"6400000","channelId":"3"},
		{"portId":"8","frequency":"0","modulationType":"QAM_NONE",
		"signalStrength":"-","bandwidth":"1600000","channelId":"0"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	p, err := d.CMUsInfo(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMUsInfo{
		Error: Error{"000", ""},
		Ports: []PortInfo{
			{
				PortID:         "1",
				Frequency:      32300000,
				Bandwidth:      6400000,
				Modulation:     "64QAM",
				SignalStrength: 45.270,
				ChannelID:      "3",
			},
			{
				PortID:     "8",
				Modulation: "QAM_NONE",
				ChannelID:  "0",
				Bandwidth:  1600000,
			},
		},
	}, p)
}

func TestCMSysInfo(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
	"ntAccess":"Permitted",
	"ip":["7.96.63.138"],
	"subMask":"255.255.255.0",
	"gw":"7.96.63.1",
	"lease":"D: 6 H: 09 M: 25 S: 55","Configname":"bac110000106749be82df7e0",
	"DsDataRate":"1040000000","UsDataRate":"31200000","macAddr":"74:9b:DE:AD:BE:EF"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	p, err := d.CMSysInfo(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMSysInfo{
		Error:         Error{"000", ""},
		NetworkAccess: "Permitted",
		IP:            net.ParseIP("7.96.63.138"),
		SubMask:       net.CIDRMask(24, 32),
		GW:            net.ParseIP("7.96.63.1"),
		Configname:    "bac110000106749be82df7e0",
		DsDataRate:    1040000000,
		UsDataRate:    31200000,
		MacAddr:       mustMAC("74:9b:DE:AD:BE:EF"),
		Lease: (6 * 24 * time.Hour) +
			(9 * time.Hour) +
			(25 * time.Minute) +
			(55 * time.Second),
	}, p)
}

func TestCMDsOFDM(t *testing.T) {
	body := ` {"errCode":"000","errMsg":"",
	"Freq_List":[
		{"receive":0,"ffttype":"NA","Subcarr0freqFreq":"NA",
		"plclock":" NO","ncplock":" NO","mdc1lock":" NO","plcpower":"NA"},
		{"receive":1,"ffttype":"4K","Subcarr0freqFreq":" 275600000",
		"plclock":"YES","ncplock":"YES","mdc1lock":"YES","plcpower":"  0.799999"}
		]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	p, err := d.CMDsOfdm(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMDsOfdm{
		Error: Error{"000", ""},
		Receivers: []OFDMReceiver{
			{},
			{
				ID: 1, FFTType: "4K", SubcarrierFreq: 275600000,
				PLCLocked: true, NCPLocked: true, MDC1Locked: true,
				PLCPower: 0.799999,
			},
		},
	}, p)
}

func TestCMUsOFDM(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
	"Freq_List":[
		{"uschindex":0,"state":"  DISABLED","digAtten":"    0.0000",
		"digAttenBo":"    0.0000","channelBw":"    0.0000",
		"repPower":"    0.0000","repPower1_6":"    0.0000",
		"fftVal":"        2K"},
		{"uschindex":1,"state":"  DISABLED","digAtten":"    0.0000",
		"digAttenBo":"    0.0000","channelBw":"    0.0000",
		"repPower":"    0.0000","repPower1_6":"    0.0000",
		"fftVal":"        2K"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	p, err := d.CMUsOfdm(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMUsOfdm{
		Error: Error{"000", ""},
		Channels: []OFDMAChannel{
			{ID: 0, FFTSize: "2K"},
			{ID: 1, FFTSize: "2K"},
		},
	}, p)
}

func TestCMLog(t *testing.T) {
	msg1 := "a message"
	msg2 := "another message"
	body := fmt.Sprintf(`{"errCode":"000","errMsg":"",
	"Log_List":[
		{"index":1,"time":"11\/15\/2020 03:57:38","type":"68010300","priority":"4",
		"event":"%s"},
		{"index":2,"time":"11\/16\/2020 17:19:06","type":"74010100","priority":"6",
		"event":"%s"}
	]}`, msg1, msg2)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()
	d := &CableModem{credentials{}, mustParse(srv.URL), srv.Client()}

	t1, err := time.Parse(time.RFC3339, "2020-11-15T03:57:38Z")
	assert.NoError(t, err)
	t2, err := time.Parse(time.RFC3339, "2020-11-16T17:19:06Z")
	assert.NoError(t, err)

	p, err := d.CMLog(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, CMLog{
		Error: NoError,
		Logs: []LogEntry{
			{
				ID: 1, Time: t1, Type: "68010300", Severity: "Error",
				Event: msg1,
			},
			{
				ID: 2, Time: t2, Type: "74010100", Severity: "Notice",
				Event: msg2,
			},
		},
	}, p)
}
