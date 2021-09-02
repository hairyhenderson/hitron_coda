package hitron

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

// RouterSysInfo -
type RouterSysInfo struct {
	Error
	CMVersion

	SystemTime      time.Time        // current time
	LANName         string           // :"brlan0",
	PrivLanIP       net.IP           // :"192.168.0.1\/24",
	PrivLanNet      *net.IPNet       //
	LanRx           int64            // :"19601748772",
	LanTx           int64            // :"141585555187",
	WanName         string           // :"erouter0",
	WanIP           []net.IP         // :["23.233.27.226","2607:f2c0:f200:a03:59e0:7e1e:f96b:923b"],
	WanRx           int64            // :"139788502458",
	WanRxPkts       int64            // :"175946286",
	WanTx           int64            // :"18787516468",
	WanTxPkts       int64            // :"52845543",
	DNS             []net.IP         // :["127.0.0.1","2607:f2c0::2"],
	RFMac           net.HardwareAddr // :"74:9B:DE:AD:BE:EF",
	SecDNS          net.IP           // :"",
	SystemLanUptime time.Duration    // : "468117",
	SystemWanUptime time.Duration    // :"468083",
	RouterMode      string           // :"Dualstack"
}

// UnmarshalJSON - implements json.Unmarshaler
//nolint:funlen
func (s *RouterSysInfo) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error
		CMVersion

		SysTime         string   `json:"sysTime"`         // :"2020-11-17 02:12:33",
		TZ              string   `json:"tz"`              // :"13_2_1",
		LANName         string   `json:"lanName"`         // :"brlan0",
		PrivLanIP       string   `json:"privLanIP"`       // :"192.168.0.1\/24",
		LanRx           string   `json:"lanRx"`           // :"19601748772",
		LanTx           string   `json:"lanTx"`           // :"141585555187",
		WanName         string   `json:"wanName"`         // :"erouter0",
		WanIP           []net.IP `json:"wanIP"`           // :["23.233.27.226","2607:f2c0:f200:a03:59e0:7e1e:f96b:923b"],
		WanRx           string   `json:"wanRx"`           // :"139788502458",
		WanRxPkts       string   `json:"wanRxPkts"`       // :"175946286",
		WanTx           string   `json:"wanTx"`           // :"18787516468",
		WanTxPkts       string   `json:"wanTxPkts"`       // :"52845543",
		DNS             []net.IP `json:"dns"`             // :["127.0.0.1","2607:f2c0::2"],
		RFMac           string   `json:"rfMac"`           // :"74:9B:DE:AD:BE:EF",
		SecDNS          net.IP   `json:"secDNS"`          // :"",
		SystemLanUptime string   `json:"systemLanUptime"` // : "468117",
		SystemWanUptime string   `json:"systemWanUptime"` // :"468083",
		RouterMode      string   `json:"routerMode"`      // :"Dualstack"
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RouterSysInfo %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.CMVersion = raw.CMVersion

	s.LANName = raw.LANName
	s.WanName = raw.WanName
	s.RouterMode = raw.RouterMode
	s.WanIP = raw.WanIP
	s.DNS = raw.DNS
	s.SecDNS = raw.SecDNS

	s.RFMac, err = net.ParseMAC(raw.RFMac)
	if err != nil {
		return fmt.Errorf("failed to parse RF MAC Address %q: %w", raw.RFMac, err)
	}

	s.PrivLanIP, s.PrivLanNet, err = net.ParseCIDR(raw.PrivLanIP)
	if err != nil {
		return fmt.Errorf("failed to parse CIDR %q: %w", raw.PrivLanIP, err)
	}

	s.LanRx = atoi64(raw.LanRx)
	s.LanTx = atoi64(raw.LanTx)
	s.WanRx = atoi64(raw.WanRx)
	s.WanRxPkts = atoi64(raw.WanRxPkts)
	s.WanTx = atoi64(raw.WanTx)
	s.WanTxPkts = atoi64(raw.WanTxPkts)

	lanUp := atoi64(raw.SystemLanUptime)
	s.SystemLanUptime = time.Duration(lanUp) * time.Second

	wanUp := atoi64(raw.SystemWanUptime)
	s.SystemWanUptime = time.Duration(wanUp) * time.Second

	l, err := tzToLocation(raw.TZ)
	if err != nil {
		return err
	}

	// Date format is MM/DD/YYYY HH:MM:SS - no timezone
	t, err := time.ParseInLocation("2006-01-02 15:04:05", raw.SysTime, l)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp %q/%s: %w", raw.SysTime, raw.TZ, err)
	}

	s.SystemTime = t

	return nil
}

func tzToLocation(tz string) (*time.Location, error) {
	tzmap := map[string]string{
		"0_1_0":  "Pacific/Kwajalein",
		"1_1_0":  "Pacific/Pago_Pago",
		"2_1_0":  "Pacific/Honolulu",
		"3_1_1":  "America/Anchorage",
		"4_1_1":  "America/Los_Angeles",
		"5_1_0":  "America/Phoenix",
		"5_2_1":  "America/Denver",
		"6_1_1":  "America/Mexico_City",
		"6_2_1":  "America/Chicago",
		"7_1_0":  "America/Indiana/Indianapolis",
		"7_2_1":  "America/New_York",
		"8_1_0":  "America/Caracas",
		"8_2_1":  "America/Halifax",
		"9_1_1":  "America/St_Johns",
		"10_1_1": "America/Sao_Paulo",
		"11_1_1": "Atlantic/South_Georgia",
		"12_1_1": "Atlantic/Azores",
		"13_1_0": "Africa/Monrovia",
		"13_2_1": "Etc/UTC",
		"14_1_0": "Africa/Tunis",
		"14_2_1": "Europe/Rome",
		"15_1_0": "Africa/Johannesburg",
		"16_1_0": "Europe/Athens",
		"17_1_0": "Europe/Samara",
		"18_1_0": "Asia/Yekaterinburg",
		"19_1_0": "Asia/Kolkata",
		"20_1_0": "Asia/Omsk",
		"21_1_0": "Asia/Bangkok",
		"22_1_0": "Asia/Shanghai",
		"22_2_0": "Asia/Taipei",
		"23_1_0": "Asia/Tokyo",
		"24_1_0": "Pacific/Guam",
		"24_2_1": "Australia/Sydney",
		"25_1_0": "Pacific/Bougainville",
		"26_1_1": "Pacific/Auckland",
	}

	l, err := time.LoadLocation(tzmap[tz])
	if err != nil {
		return time.UTC, fmt.Errorf("failed to load location %q: %w", tzmap[tz], err)
	}

	return l, nil
}

// RouterCapability -
type RouterCapability struct {
	Error

	RouterMode string
	Gateway    bool
	UPnP       bool
	HNAP       bool
	USB        bool
	SIPAlg     bool
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *RouterCapability) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error

		RouterMode   string
		GatewayOnOff string
		UPnpOnOff    string
		HnapOnOff    string
		UsbOnOff     string
		SIPAlgOnOff  string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RouterCapability %q: %w", string(b), err)
	}

	s.Error = raw.Error

	s.RouterMode = raw.RouterMode

	s.Gateway = raw.GatewayOnOff == on
	s.UPnP = raw.UPnpOnOff == on
	s.HNAP = raw.HnapOnOff == on
	s.USB = raw.UsbOnOff == on
	s.SIPAlg = raw.SIPAlgOnOff == on

	return nil
}

// RouterLocation -
type RouterLocation struct {
	Error
	LocationText string
}

// RouterDMZ -
type RouterDMZ struct {
	Error
	Enable     bool
	Host       net.IP
	PrivateLan net.IP
	Mask       net.IP
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *RouterDMZ) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error

		Enable     string
		Host       net.IP
		PrivateLan net.IP
		SubMask    net.IP
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RouterDMZ %q: %w", string(b), err)
	}

	s.Error = raw.Error

	s.Enable = raw.Enable == on
	s.Host = raw.Host
	s.PrivateLan = raw.PrivateLan
	s.Mask = raw.SubMask

	return nil
}

// RouterPortForwardStatus -
type RouterPortForwardStatus struct {
	Error
	Enable     bool
	PrivateLan net.IP
	Mask       net.IP
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *RouterPortForwardStatus) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error

		AllRulesOnOff string
		PrivateLan    net.IP
		SubMask       net.IP
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RouterPortForwardStatus %q: %w", string(b), err)
	}

	s.Error = raw.Error

	s.Enable = raw.AllRulesOnOff == on
	s.PrivateLan = raw.PrivateLan
	s.Mask = raw.SubMask

	return nil
}

// RouterPortForwardall -
type RouterPortForwardall struct {
	Error
	Total int
	Rules []PortForwardRule `json:"Rules_List"`
}

// PortForwardRule -
type PortForwardRule struct {
	Enable       bool
	ID           int
	Origin       int
	AppName      string
	Protocol     string
	PublicPorts  PortRange
	PrivatePorts PortRange
	LocalIP      net.IP
	RemoteIPs    IPRange
}

// PortRange -
type PortRange struct {
	Start, End int
}

func parsePort(in string) int {
	p, _ := strconv.Atoi(in)

	return p
}

// IPRange -
type IPRange struct {
	Start, End net.IP
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *PortForwardRule) UnmarshalJSON(b []byte) error {
	raw := struct {
		ID               string
		Origin           string
		AppName          string
		PubStart, PubEnd string
		PriStart, PriEnd string
		RuleOnOff        string
		Protocol         string
		LocalIPAddr      net.IP
		RemoteIPStar     net.IP
		RemoteIPEnd      net.IP
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal PortForwardRule %q: %w", string(b), err)
	}

	s.ID, _ = strconv.Atoi(raw.ID)
	s.Origin, _ = strconv.Atoi(raw.Origin)
	s.AppName = raw.AppName
	s.Protocol = raw.Protocol

	s.Enable = raw.RuleOnOff == on

	s.PublicPorts = PortRange{parsePort(raw.PubStart), parsePort(raw.PubEnd)}
	s.PrivatePorts = PortRange{parsePort(raw.PriStart), parsePort(raw.PriEnd)}

	s.LocalIP = raw.LocalIPAddr

	s.RemoteIPs = IPRange{raw.RemoteIPStar, raw.RemoteIPEnd}

	return nil
}

// RouterPortTriggerStatus -
// Port triggering is a means of automating port forwarding. The CODA-4x8x scans
// outgoing traffic (from the LAN to the WAN) to see if any of the traffic's
// destination ports match those specified in the port triggering rules you
// configure. If any of the ports match, the CODA-4x8x automatically opens the
// incoming ports specified in the rule, in anticipation of incoming traffic.
type RouterPortTriggerStatus struct {
	Error
	Enable bool
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *RouterPortTriggerStatus) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error

		AllRulesOnOff string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RouterPortTriggerStatus %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.Enable = raw.AllRulesOnOff == on

	return nil
}

// RouterPortTriggerall -
type RouterPortTriggerall struct {
	Error
	Total int
	Rules []PortTriggerRule `json:"Rules_List"`
}

// PortTriggerRule -
type PortTriggerRule struct {
	AppName      string
	Protocol     string    // TCP, UDP, BOTH
	TriggerPorts PortRange // outgoing/public
	TargetPorts  PortRange // incoming/private
	ID           int
	Timeout      time.Duration
	Enable       bool
	TwoWay       bool
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *PortTriggerRule) UnmarshalJSON(b []byte) error {
	raw := struct {
		RuleOnOff        string
		ID               string
		AppName          string
		PubStart, PubEnd string
		PriStart, PriEnd string
		Protocol         string
		Timeout          string
		TwoWayOnOff      string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal PortTriggerRule %q: %w", string(b), err)
	}

	s.ID, _ = strconv.Atoi(raw.ID)
	s.Enable = raw.RuleOnOff == on
	s.TwoWay = raw.TwoWayOnOff == on

	s.AppName = raw.AppName
	s.Protocol = raw.Protocol

	s.TriggerPorts = PortRange{parsePort(raw.PubStart), parsePort(raw.PubEnd)}
	s.TargetPorts = PortRange{parsePort(raw.PriStart), parsePort(raw.PriEnd)}

	to, _ := strconv.Atoi(raw.Timeout)
	s.Timeout = time.Duration(to) * time.Millisecond

	return nil
}

// RouterTR069 -
type RouterTR069 struct {
	Error
	TR069URL string
}
