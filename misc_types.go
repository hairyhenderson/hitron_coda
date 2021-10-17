package hitron

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	on      = "ON"
	yes     = "YES"
	enabled = "Enabled"
)

// Error contains common error code information.
type Error struct {
	Code    string `json:"errCode"`
	Message string `json:"errMsg"`
}

func (e Error) String() string {
	if e.Code == "000" {
		return ""
	}

	return fmt.Sprintf("Error %s: %s", e.Code, e.Message)
}

// NoError represents the successful state
//nolint:gochecknoglobals
var NoError = Error{Code: "000", Message: ""}

// Time -
type Time struct {
	Error
	Enable       bool
	Daylight     bool
	DaylightTime int
	TZ           *time.Location
	SNTPServer   string
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *Time) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error
		SntpOnOff     string
		SntpTimeZone  string
		SntpSrvName   string
		DaylightOnOff string
		DaylightTime  string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Time %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.Enable = raw.SntpOnOff == on
	s.Daylight = raw.DaylightOnOff == on
	s.SNTPServer = raw.SntpSrvName

	s.TZ, _ = tzToLocation(raw.SntpTimeZone)

	return nil
}

// DNS -
type DNS struct {
	Error
	AutoEnable   bool // whether to set DNS servers automatically, or manually
	ProxyEnable  bool // enable DNS proxying to DHCP clients on the LAN
	LanDNS1      net.IP
	LanDNS2      net.IP
	DomainSuffix string
	ProxyName1   string
	ProxyName2   string
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *DNS) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error
		LanDNSOnOff   string
		DNSProxyOnOff string
		LanDNS1       net.IP
		LanDNS2       net.IP
		DomainSuffix  string
		ProxyName1    string
		ProxyName2    string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal DNS %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.LanDNS1 = raw.LanDNS1
	s.LanDNS2 = raw.LanDNS2
	s.DomainSuffix = raw.DomainSuffix
	s.ProxyName1 = raw.ProxyName1
	s.ProxyName2 = raw.ProxyName2
	s.AutoEnable = raw.LanDNSOnOff == on
	s.ProxyEnable = raw.DNSProxyOnOff == on

	return nil
}

// DDNS -
type DDNS struct {
	Error
	Enable         bool
	Provider       string
	Username       string
	Password       string
	Hostname       string
	UpdateInterval time.Duration
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *DDNS) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error
		DDNSOnOff          string
		DDNSSrvProvider    int
		DDNSUsername       string
		DDNSPassword       string
		DDNSHostnames      string
		DDNSUpdateInterval string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal DDNS %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.Username = raw.DDNSUsername
	s.Password = raw.DDNSPassword
	s.Hostname = raw.DDNSHostnames

	s.Enable = raw.DDNSOnOff == on

	// Just default to 0 if it can't be parsed
	d, _ := strconv.Atoi(raw.DDNSUpdateInterval)
	s.UpdateInterval = time.Duration(d) * time.Second

	providers := map[int]string{
		1:  "dyndns@dyndns.org",
		2:  "default@freedns.afraid.org",
		3:  "default@zoneedit.com",
		4:  "default@no-ip.com",
		5:  "default@easydns.com",
		6:  "default@tzo.com",
		7:  "dyndns@3322.org",
		8:  "default@sitelutions.com",
		9:  "default@dnsomatic.com",
		10: "ipv6tb@he.net",
		11: "default@dynsip.org",
	}
	s.Provider = providers[raw.DDNSSrvProvider]

	return nil
}

// Hosts -
type Hosts struct {
	Error
	Hosts []Host `json:"Hosts_List"`
}

// Host -
type Host struct {
	Name          string
	AddressSource string
	MacAddr       net.HardwareAddr
	IP            net.IP
	ConnectType   string
	ConnectTo     net.HardwareAddr
	Comnum        int
	AppEnable     bool
	Action        string
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *Host) UnmarshalJSON(b []byte) error {
	raw := struct {
		HostName      string
		AddressSource string
		MacAddr       string
		IP            net.IP
		ConnectType   string
		ConnectTo     string
		Comnum        int
		AppEnable     string
		Action        string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Host %q: %w", string(b), err)
	}

	s.Name = raw.HostName
	s.AddressSource = raw.AddressSource
	s.IP = raw.IP
	s.ConnectType = raw.ConnectType
	s.Comnum = raw.Comnum
	s.Action = raw.Action

	s.AppEnable = raw.AppEnable == "TRUE"

	s.MacAddr, _ = net.ParseMAC(raw.MacAddr)
	s.ConnectTo, _ = net.ParseMAC(raw.ConnectTo)

	return nil
}

type UsersCSRF struct {
	Error
	CSRF string
}
