package hitron

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"sntpOnOff":"ON","sntpTimeZone":"13_1_0","sntpSrvName":"pool.ntp.org",
		"daylightOnOff":"ON","daylightTime":"0"
	}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	p, err := d.Time(context.Background())
	assert.NoError(t, err)

	tz, _ := time.LoadLocation("Africa/Monrovia")

	assert.EqualValues(t, Time{
		Error:        NoError,
		Enable:       true,
		Daylight:     true,
		DaylightTime: 0,
		TZ:           tz,
		SNTPServer:   "pool.ntp.org",
	}, p)
}

func TestDNS(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"lanDnsOnOff":"ON",
		"landns1":"192.168.0.1","landns2":"",
		"dnsProxyOnOff":"ON","domainSuffix":"ht.home",
		"proxyName1":"foo","proxyName2":"bar"
	}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	ctx := context.Background()

	p, err := d.DNS(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, DNS{
		Error:        NoError,
		AutoEnable:   true,
		ProxyEnable:  true,
		LanDNS1:      net.ParseIP("192.168.0.1"),
		LanDNS2:      nil,
		DomainSuffix: "ht.home",
		ProxyName1:   "foo",
		ProxyName2:   "bar",
	}, p)
}

func TestDDNS(t *testing.T) {
	body := `{"errCode":"000","errMsg":"",
		"ddnsOnOff":"ON",
		"ddnsSrvProvider":10,"ddnsUsername":"foo","ddnsPassword":"bar",
		"ddnsHostnames":"foo.example.com","ddnsUpdateInterval":"604800"
	}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	ctx := context.Background()

	p, err := d.DDNS(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, DDNS{
		Error:          NoError,
		Enable:         true,
		Provider:       "ipv6tb@he.net",
		Username:       "foo",
		Password:       "bar",
		Hostname:       "foo.example.com",
		UpdateInterval: 7 * 24 * time.Hour,
	}, p)
}

func TestHosts(t *testing.T) {
	body := `{"errCode":"000","errMsg":"","HostNumberOfEntries":2,
		"Hosts_List":[
			{"hostName":"Unknown","macAddr":"de:ad:be:ef:ca:fe",
			"ip":"192.168.0.15","addressSource":"DHCP-IP",
			"connectType":"Ethernet","connectTo":"CA:FE:DE:AD:FA:CE",
			"comnum":1,"appEnable":"TRUE","action":"Resume"},
			{"hostName":"Unknown","macAddr":"13:37:be:ef:ca:fe",
			"ip":"192.168.0.16","addressSource":"DHCP-IP",
			"connectType":"Ethernet","connectTo":"CA:FE:DE:AD:FA:CE",
			"comnum":1,"appEnable":"TRUE","action":"Resume"}
		]}`

	srv := staticResponseServer(t, body)
	d := testCableModem(srv)

	ctx := context.Background()

	p, err := d.Hosts(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, Hosts{
		Error: NoError,
		Hosts: []Host{
			{
				Name:          "Unknown",
				AddressSource: "DHCP-IP",
				IP:            net.ParseIP("192.168.0.15"),
				MacAddr:       net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe},
				ConnectTo:     net.HardwareAddr{0xCA, 0xFE, 0xDE, 0xAD, 0xFA, 0xCE},
				Comnum:        1,
				AppEnable:     true,
				Action:        "Resume",
				ConnectType:   "Ethernet",
			},
			{
				Name:          "Unknown",
				AddressSource: "DHCP-IP",
				IP:            net.ParseIP("192.168.0.16"),
				MacAddr:       net.HardwareAddr{0x13, 0x37, 0xbe, 0xef, 0xca, 0xfe},
				ConnectTo:     net.HardwareAddr{0xCA, 0xFE, 0xDE, 0xAD, 0xFA, 0xCE},
				Comnum:        1,
				AppEnable:     true,
				Action:        "Resume",
				ConnectType:   "Ethernet",
			},
		},
	}, p)
}

func TestUsersCSRF(t *testing.T) {
	body := `{"errCode":"000","errMsg":"","CSRF":"abcdefgh1234.4321abcdefgh"}`

	srv := staticResponseServer(t, body)

	d := testCableModem(srv)
	ctx := context.Background()

	p, err := d.UsersCSRF(ctx)
	assert.NoError(t, err)

	assert.EqualValues(t, UsersCSRF{
		Error: NoError,
		CSRF:  "abcdefgh1234.4321abcdefgh",
	}, p)
}
