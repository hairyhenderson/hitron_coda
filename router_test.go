package hitron

import (
	"context"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCableModem(srv *httptest.Server) *CableModem {
	return &CableModem{credentials: credentials{}, base: mustParse(srv.URL), hc: srv.Client()}
}

//nolint:funlen
func TestRouterSysInfo(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"deviceId":"74:9B:DE:AD:BE:EF",
		"modelName":"CODA-4680-TPIA",
		"vendorName":"Hitron Technologies",
		"SerialNum":"123456789012",
		"HwVersion":"1A",
		"ApiVersion":"1.11",
		"SoftwareVersion":"7.1.1.2.2b9",
		
		"sysTime":"2020-11-17 02:12:33",
		"tz":"13_2_1",
		"lanName":"brlan0",
		"privLanIP":"192.168.0.1\/24",
		"lanRx":"19601748772",
		"lanTx":"141585555187",
		"wanName":"erouter0",
		"wanIP":["23.233.27.226","2607:f2c0:f200:a03:59e0:7e1e:f96b:923b"],
		"wanRx":"139788502458",
		"wanRxPkts":"175946286",
		"wanTx":"18787516468",
		"wanTxPkts":"52845543",
		"dns":["127.0.0.1","2607:f2c0::2"],
		"rfMac":"74:9B:DE:AD:BE:EF",
		"secDNS":"",
		"systemLanUptime": "468117",
		"systemWanUptime":"468083",
		"routerMode":"Dualstack"
		}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterSysInfo(context.Background())
	assert.NoError(t, err)

	loc, err := time.LoadLocation("Etc/UTC")
	assert.NoError(t, err)

	systime, err := time.ParseInLocation("2006-01-02 15:04:05", "2020-11-17 02:12:33", loc)
	assert.NoError(t, err)

	//nolint:dupl
	assert.EqualValues(t, RouterSysInfo{
		Error: NoError,
		CMVersion: CMVersion{
			DeviceID:        "74:9B:DE:AD:BE:EF",
			ModelName:       "CODA-4680-TPIA",
			VendorName:      "Hitron Technologies",
			SerialNum:       "123456789012",
			HwVersion:       "1A",
			APIVersion:      "1.11",
			SoftwareVersion: "7.1.1.2.2b9",
		},
		SystemTime: systime,
		LANName:    "brlan0",
		WanName:    "erouter0",
		RouterMode: "Dualstack",
		PrivLanIP:  net.ParseIP("192.168.0.1"),
		PrivLanNet: &net.IPNet{IP: net.ParseIP("192.168.0.0").To4(), Mask: net.CIDRMask(24, 32)},

		LanRx:           19601748772,
		LanTx:           141585555187,
		WanIP:           []net.IP{net.ParseIP("23.233.27.226"), net.ParseIP("2607:f2c0:f200:a03:59e0:7e1e:f96b:923b")},
		WanRx:           139788502458,
		WanRxPkts:       175946286,
		WanTx:           18787516468,
		WanTxPkts:       52845543,
		DNS:             []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("2607:f2c0::2")},
		RFMac:           net.HardwareAddr{0x74, 0x9b, 0xde, 0xad, 0xbe, 0xef},
		SystemLanUptime: 130*time.Hour + 1*time.Minute + 57*time.Second,
		SystemWanUptime: 130*time.Hour + 1*time.Minute + 23*time.Second,
	}, p)

	// different body from CODA-4582 in bridge mode
	body = `{
		"errCode": "000",
		"errMsg": "",
		"deviceId": "",
		"modelName": "CODA-4582",
		"vendorName": "RES",
		"SerialNum": "",
		"HwVersion": "2A",
		"ApiVersion": "1.12.1",
		"SoftwareVersion": "7.1.1.2.8b4",
		"deploymentName": "",
		"VendorLogo": "0",
		"sysTime": "2022-07-29 23:26:32",
		"tz": "0",
		"lanName": null,
		"privLanIP": "192.168.0.1\/24",
		"lanRx": "1.05G Bytes",
		"lanTx": "18.15G Bytes",
		"wanName": null,
		"wanIP": [ "0.0.0.0\n", "" ],
		"wanRx": "18.23G Bytes",
		"wanRxPkts": null,
		"wanTx": "989.1M Bytes",
		"wanTxPkts": null,
		"dns": [ "", "" ],
		"rfMac": "",
		"secDNS": "none",
		"systemLanUptime": "000 days 04h:43m:05s",
		"systemWanUptime": "000 days 00h:00m:00s",
		"routerMode": "Bridge",
		"lanAsWanMode": "0",
		"privLanIp": "192.168.0.1\/24"
	}`
	srv = staticResponseServer(t, body)

	d = testCableModem(srv)

	p, err = d.RouterSysInfo(context.Background())
	require.NoError(t, err)

	systime, err = time.ParseInLocation("2006-01-02 15:04:05", "2022-07-29 23:26:32", loc)
	assert.NoError(t, err)

	assert.EqualValues(t, RouterSysInfo{
		Error: NoError,
		CMVersion: CMVersion{
			ModelName:       "CODA-4582",
			VendorName:      "RES",
			HwVersion:       "2A",
			APIVersion:      "1.12.1",
			SoftwareVersion: "7.1.1.2.8b4",
		},
		SystemTime: systime,
		RouterMode: "Bridge",
		PrivLanIP:  net.ParseIP("192.168.0.1"),
		PrivLanNet: &net.IPNet{IP: net.ParseIP("192.168.0.0").To4(), Mask: net.CIDRMask(24, 32)},

		LanRx:           1127428915,
		LanTx:           19488414105,
		WanIP:           []net.IP{net.ParseIP("0.0.0.0")},
		WanRx:           19574313451,
		WanTx:           1037146521,
		DNS:             []net.IP{},
		SystemLanUptime: 4*time.Hour + 43*time.Minute + 5*time.Second,
	}, p)
}

func TestParseUptime(t *testing.T) {
	d, err := parseUptime("")
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), d)

	d, err = parseUptime("468117")
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(468117)*time.Second, d)

	d, err = parseUptime("000 days 00h:00m:00s")
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), d)

	d, err = parseUptime("000 days 04h:43m:05s")
	assert.NoError(t, err)
	assert.Equal(t, 4*time.Hour+43*time.Minute+5*time.Second, d)

	d, err = parseUptime("015 days 22h:13m:57s")
	assert.NoError(t, err)
	assert.Equal(t,
		(15*24*time.Hour)+
			22*time.Hour+
			13*time.Minute+
			57*time.Second, d)

	d, err = parseUptime("400 days 22h:13m:57s")
	assert.NoError(t, err)
	assert.Equal(t,
		(400*24*time.Hour)+
			22*time.Hour+
			13*time.Minute+
			57*time.Second, d)
}

func TestRouterSysInfo_String(t *testing.T) {
	loc, err := time.LoadLocation("Etc/UTC")
	assert.NoError(t, err)

	systime, err := time.ParseInLocation("2006-01-02 15:04:05", "2020-11-17 02:12:33", loc)
	assert.NoError(t, err)

	//nolint:dupl
	s := RouterSysInfo{
		Error: NoError,
		CMVersion: CMVersion{
			DeviceID:        "74:9B:DE:AD:BE:EF",
			ModelName:       "CODA-4680-TPIA",
			VendorName:      "Hitron Technologies",
			SerialNum:       "123456789012",
			HwVersion:       "1A",
			APIVersion:      "1.11",
			SoftwareVersion: "7.1.1.2.2b9",
		},
		SystemTime: systime,
		LANName:    "brlan0",
		WanName:    "erouter0",
		RouterMode: "Dualstack",
		PrivLanIP:  net.ParseIP("192.168.0.1"),
		PrivLanNet: &net.IPNet{IP: net.ParseIP("192.168.0.0").To4(), Mask: net.CIDRMask(24, 32)},

		LanRx:           19601748772,
		LanTx:           141585555187,
		WanIP:           []net.IP{net.ParseIP("23.233.27.226"), net.ParseIP("2607:f2c0:f200:a03:59e0:7e1e:f96b:923b")},
		WanRx:           139788502458,
		WanRxPkts:       175946286,
		WanTx:           18787516468,
		WanTxPkts:       52845543,
		DNS:             []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("2607:f2c0::2")},
		RFMac:           net.HardwareAddr{0x74, 0x9b, 0xde, 0xad, 0xbe, 0xef},
		SystemLanUptime: 130*time.Hour + 1*time.Minute + 57*time.Second,
		SystemWanUptime: 130*time.Hour + 1*time.Minute + 23*time.Second,
	}

	assert.Equal(t, `CMVersion:
	DeviceID: 74:9B:DE:AD:BE:EF
	ModelName: CODA-4680-TPIA
	VendorName: Hitron Technologies
	SerialNum: 123456789012
	HwVersion: 1A
	APIVersion: 1.11
	SoftwareVersion: 7.1.1.2.2b9
SystemTime: 2020-11-17 02:12:33 +0000 UTC
LAN: brlan0 (IP 192.168.0.1) (Net 192.168.0.0/24)
	Rx/Tx: 18.3G/131.9G
WAN: erouter0 (23.233.27.226, 2607:f2c0:f200:a03:59e0:7e1e:f96b:923b)
	Rx/Tx: 130.2G/17.5G
	Rx/Tx Packets: 175,946,286/52,845,543
DNS: 127.0.0.1, 2607:f2c0::2
RFMac: 74:9b:de:ad:be:ef
System Uptime: LAN 130h1m57s, WAN 130h1m23s
RouterMode: Dualstack
`, s.String())
}

func TestRouterCapability(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"gatewayOnOff":"ON","routerMode":"Dualstack","uPnpOnOff":"ON",
		"HnapOnOff":"OFF","UsbOnOff":"OFF","sipAlgOnOff":"OFF"
	}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterCapability(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterCapability{
		Error: NoError,

		RouterMode: "Dualstack",
		Gateway:    true,
		UPnP:       true,
	}, p)
}

func TestRouterLocation(t *testing.T) {
	body := `{"errCode":"000","errMsg":"","locationText":"Basement"}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterLocation(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterLocation{
		Error:        NoError,
		LocationText: "Basement",
	}, p)
}

func TestRouterDMZ(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"enable":"OFF","host":"0.0.0.0",
		"privateLan":"192.168.0.1","subMask":"255.255.255.0"
	}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterDMZ(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterDMZ{
		Error:      NoError,
		Enable:     false,
		Host:       net.ParseIP("0.0.0.0"),
		PrivateLan: net.ParseIP("192.168.0.1"),
		Mask:       net.ParseIP("255.255.255.0"),
	}, p)
}

func TestRouterPortForwardStatus(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"allRulesOnOff":"ON",
		"privateLan":"192.168.0.1","subMask":"255.255.255.0"
	}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterPortForwardStatus(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterPortForwardStatus{
		Error:      NoError,
		Enable:     true,
		PrivateLan: net.ParseIP("192.168.0.1"),
		Mask:       net.ParseIP("255.255.255.0"),
	}, p)
}

func TestRouterPortForwardall(t *testing.T) {
	body := `{"errCode":"000","errMsg":"","total":2,
		"Rules_List":[
			{"appName":"custom",
			"pubStart":"1024","pubEnd":"2048","priStart":"1024","priEnd":"2048",
			"protocol":"UDP","localIpAddr":"192.168.0.2",
			"remoteIpStar":"0.0.0.0","remoteIpEnd":"255.255.255.255",
			"ruleOnOff":"ON","origin":"1","id":"1"},
			{"appName":"SSH",
			"pubStart":"22","pubEnd":"22","priStart":"2222","priEnd":"2222",
			"protocol":"TCP","localIpAddr":"192.168.0.16",
			"remoteIpStar":"10.0.0.1","remoteIpEnd":"11.4.3.2",
			"ruleOnOff":"OFF","origin":"1","id":"2"}
		]}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterPortForwardall(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterPortForwardall{
		Error: NoError,
		Total: 2,
		Rules: []PortForwardRule{
			{
				Enable: true, ID: 1, Origin: 1,
				AppName:      "custom",
				Protocol:     "UDP",
				PublicPorts:  PortRange{1024, 2048},
				PrivatePorts: PortRange{1024, 2048},
				LocalIP:      net.ParseIP("192.168.0.2"),
				RemoteIPs:    IPRange{net.IPv4zero, net.IPv4bcast},
			},
			{
				Enable: false, ID: 2, Origin: 1,
				AppName:      "SSH",
				Protocol:     "TCP",
				PublicPorts:  PortRange{22, 22},
				PrivatePorts: PortRange{2222, 2222},
				LocalIP:      net.ParseIP("192.168.0.16"),
				RemoteIPs:    IPRange{net.ParseIP("10.0.0.1"), net.ParseIP("11.4.3.2")},
			},
		},
	}, p)
}

func TestRouterPortTriggerStatus(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"allRulesOnOff":"ON"
	}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterPortTriggerStatus(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterPortTriggerStatus{
		Error:  NoError,
		Enable: true,
	}, p)
}

func TestRouterPortTriggerall(t *testing.T) {
	body := `{"errCode":"000","errMsg":"","total":1,
		"Rules_List":[
			{"ruleOnOff":"ON","appName":"foo","protocol":"BOTH",
			"pubStart":"80","pubEnd":"88","priStart":"8080","priEnd":"8088",
			"timeout":"50","twowayOnOff":"ON","id":"1"}
		]}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterPortTriggerall(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterPortTriggerall{
		Error: NoError,
		Total: 1,
		Rules: []PortTriggerRule{
			{
				Enable: true, ID: 1,
				AppName:      "foo",
				Protocol:     "BOTH",
				TriggerPorts: PortRange{80, 88},
				TargetPorts:  PortRange{8080, 8088},
				TwoWay:       true,
				Timeout:      50 * time.Millisecond,
			},
		},
	}, p)
}

func TestRouterTR069(t *testing.T) {
	body := `{"errCode":"000","errMsg":"","tr069url":"http://example.com"}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.RouterTR069(context.Background())
	assert.NoError(t, err)

	assert.EqualValues(t, RouterTR069{
		Error:    NoError,
		TR069URL: "http://example.com",
	}, p)
}
