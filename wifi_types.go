package hitron

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// WiFiAccessControl -
type WiFiAccessControl struct {
	Error
	BlockType string
	RulesList []WiFiAccessControlRule
}

// WiFiAccessControlRule -
type WiFiAccessControlRule struct {
	ID       int
	Hostname string
	MACAddr  net.HardwareAddr
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *WiFiAccessControlRule) UnmarshalJSON(b []byte) error {
	raw := struct {
		ID       int
		Hostname string `json:"hostName"`
		MACAddr  string `json:"macAddr"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiAccessControlRule %q: %w", string(b), err)
	}

	s.ID = raw.ID
	s.Hostname = raw.Hostname
	s.MACAddr, _ = net.ParseMAC(raw.MACAddr)

	return nil
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *WiFiAccessControl) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error

		BlockType           string
		RuleNumberOfEntries int
		RulesList           []WiFiAccessControlRule `json:"Rules_List"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiAccessControl %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.BlockType = raw.BlockType
	s.RulesList = raw.RulesList

	return nil
}

// WiFiAccessControlStatus -
type WiFiAccessControlStatus struct {
	Error
	BlockType string
}

// WiFiGuestSSID -
type WiFiGuestSSID struct {
	Error
	Enable   bool
	SSID     string
	SSID5G   string
	Password string
	MaxUsers int
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *WiFiGuestSSID) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error
		SSIDName              string
		SSIDName5G            string
		Enable                string
		Pswd                  string
		AdminGuestAccProvider string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiGuestSSID %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.SSID = raw.SSIDName
	s.SSID5G = raw.SSIDName5G
	s.Password = raw.Pswd
	s.Enable = raw.Enable == on
	s.MaxUsers, _ = strconv.Atoi(raw.AdminGuestAccProvider)

	return nil
}

// WiFiRadios -
type WiFiRadios struct {
	Error
	Radios []WiFiRadio `json:"Raidos_List"` // Yup, there's a typo in the Cable Modem
}

// WiFiRadio -
type WiFiRadio struct {
	Error
	Vendor         string   // "1" (??)
	Band           string   // 2.4G/5G
	ChanBandwidth  string   // 20MHz,20/40MHz,40MHz,80MHz
	RadioURI       string   // Path to the radio: /1/Device/WiFi/Radios/<n>
	Channel        int      // only applicable if AutoChannel == false
	CurrentChannel int      //
	Mode           WiFiMode //
	Enable         bool     //
	EnableDCS      bool     // Dynamic Channel Selection
	EnableDFS      bool     // Dynamic Frequency Selection - use frequencies reserved for radars
	EnableWPS      bool     // WiFi Protected Setup
	IGMPSnoop      bool     //
	AutoChannel    bool     //
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *WiFiRadio) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error
		Vendor            string
		Band              string
		WlsOnOff          string
		WlsDcsOnOff       string
		WlsMode           string
		NBandwidth        string `json:"n_bandwidth"`
		WlsChannel        interface{}
		AutoChannel       string
		WlsDfsOnOff       string
		WlsCurrentChannel string
		WlswpsOnOff       string
		IgmpSnoop         string
		RadioURI          string `json:"Radio_URI"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiRadio %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.Vendor = raw.Vendor
	s.Band = raw.Band
	s.ChanBandwidth = raw.NBandwidth
	s.RadioURI = raw.RadioURI

	s.Enable = raw.WlsOnOff == on
	s.EnableDCS = raw.WlsDcsOnOff == on
	s.EnableDFS = raw.WlsDfsOnOff == on
	s.AutoChannel = raw.AutoChannel == on
	s.EnableWPS = raw.WlswpsOnOff == on
	s.IGMPSnoop = raw.IgmpSnoop == on

	s.CurrentChannel, _ = strconv.Atoi(raw.WlsCurrentChannel)

	switch c := raw.WlsChannel.(type) {
	case int:
		s.Channel = c
	case string:
		s.Channel, _ = strconv.Atoi(c)
	}

	// map according to the CODA's UI
	wirelessMode := map[string]WiFiMode{
		"0": WiFiModeB,
		"1": WiFiModeG,
		"2": WiFiModeN,
		"3": WiFiModeB | WiFiModeG | WiFiModeN,
		"4": WiFiModeG | WiFiModeN,
		"5": WiFiModeB | WiFiModeG,
		"6": WiFiModeAC,
		"7": WiFiModeA,
		"8": WiFiModeA | WiFiModeN,
		"9": WiFiModeA | WiFiModeN | WiFiModeAC,
	}
	s.Mode = wirelessMode[raw.WlsMode]

	return nil
}

// WiFiMode enumerates the available WiFi modes (802.11a/b/g/n/ac)
type WiFiMode uint8

// Useful WiFi constants for indicating protocol support. For multi-protocol
// support, join with binary ORs (|).
const (
	WiFiModeA WiFiMode = 1 << iota
	WiFiModeB
	WiFiModeG
	WiFiModeN
	WiFiModeAC
)

func (m WiFiMode) String() string {
	suffixes := []string{}
	if m&WiFiModeA != 0 {
		suffixes = append(suffixes, "a")
	}

	if m&WiFiModeB != 0 {
		suffixes = append(suffixes, "b")
	}

	if m&WiFiModeG != 0 {
		suffixes = append(suffixes, "g")
	}

	if m&WiFiModeN != 0 {
		suffixes = append(suffixes, "n")
	}

	if m&WiFiModeAC != 0 {
		suffixes = append(suffixes, "ac")
	}

	suffix := strings.Join(suffixes, "/")

	return "802.11" + suffix
}

// WiFiRadiosAdvanced -
type WiFiRadiosAdvanced struct {
	Error
	Radios []WiFiRadioAdvanced `json:"Advanced_List"`
}

// WiFiRadioAdvanced -
type WiFiRadioAdvanced struct {
	Error
	WiFiRadio
	SSID           string
	BGMode         string
	NOperatingMode string
	NGuardInterval string
	TxStream       string
	RxStream       string
	NMCS           int
	NCoexistence   bool
	NRDG           bool
	Namsdu         bool
	Nautoba        bool
	Nbadecline     bool
	BandSteering   bool
	ShowMSO        bool
}

// UnmarshalJSON - implements json.Unmarshaler
//nolint:funlen
func (s *WiFiRadioAdvanced) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error

		Vendor            string
		Band              string
		WlsOnOff          string
		WlsDcsOnOff       string
		WlsMode           string
		NBandwidth        string `json:"n_bandwidth"`
		WlsChannel        interface{}
		AutoChannel       string
		WlsDfsOnOff       string
		WlsCurrentChannel string
		WlswpsOnOff       string
		IgmpSnoop         string
		RadioURI          string `json:"Radio_URI"`

		BGMode         string `json:"bgMode"`
		SSIDName       string `json:"ssidName"`
		NCoexistence   string `json:"n_coexistence"`
		BandSteering   string `json:"bandsteering"`
		NOperatingMode string `json:"n_OperatingMode"`
		NGuardInterval string `json:"n_GuardInterval"`
		NMCS           string `json:"n_mcs"`
		NRDG           string `json:"n_rdg"`
		Namsdu         string `json:"n_amsdu"`
		Nautoba        string `json:"n_autoba"`
		Nbadecline     string `json:"n_badecline"`
		TxStream       string `json:"tx_stream"`
		RxStream       string `json:"rx_stream"`
		ShowMSO        string `json:"showMSO"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiRadioAdvanced %q: %w", string(b), err)
	}

	s.Error = raw.Error

	s.SSID = raw.SSIDName
	s.NOperatingMode = raw.NOperatingMode
	s.NGuardInterval = raw.NGuardInterval
	s.TxStream = raw.TxStream
	s.RxStream = raw.RxStream

	s.BandSteering = raw.BandSteering == on
	s.ShowMSO = raw.ShowMSO == "true"

	s.NCoexistence = raw.NCoexistence == enabled
	s.NRDG = raw.NRDG == enabled
	s.Namsdu = raw.Namsdu == enabled
	s.Nautoba = raw.Nautoba == enabled
	s.Nbadecline = raw.Nbadecline == enabled

	s.NMCS, _ = strconv.Atoi(raw.NMCS)

	s.Vendor = raw.Vendor
	s.Band = raw.Band
	s.ChanBandwidth = raw.NBandwidth
	s.RadioURI = raw.RadioURI

	s.Enable = raw.WlsOnOff == on
	s.EnableDCS = raw.WlsDcsOnOff == on
	s.EnableDFS = raw.WlsDfsOnOff == on
	s.AutoChannel = raw.AutoChannel == on
	s.EnableWPS = raw.WlswpsOnOff == on
	s.IGMPSnoop = raw.IgmpSnoop == on

	s.CurrentChannel, _ = strconv.Atoi(raw.WlsCurrentChannel)

	switch c := raw.WlsChannel.(type) {
	case int:
		s.Channel = c
	case string:
		s.Channel, _ = strconv.Atoi(c)
	}

	// map according to the CODA's UI
	wirelessMode := map[string]WiFiMode{
		"0": WiFiModeB,
		"1": WiFiModeG,
		"2": WiFiModeN,
		"3": WiFiModeB | WiFiModeG | WiFiModeN,
		"4": WiFiModeG | WiFiModeN,
		"5": WiFiModeB | WiFiModeG,
		"6": WiFiModeAC,
		"7": WiFiModeA,
		"8": WiFiModeA | WiFiModeN,
		"9": WiFiModeA | WiFiModeN | WiFiModeAC,
	}

	s.Mode = wirelessMode[raw.WlsMode]

	return nil
}

// WiFiRadiosSurvey -
type WiFiRadiosSurvey struct {
	Error
	APs []WiFiAP `json:"APs_List"`
}

// WiFiAP - information about a WiFi Access Point/Network
type WiFiAP struct {
	Band     string
	Channel  int
	SSID     string
	BSSID    net.HardwareAddr
	Signal   int    // percentage
	WMode    string // W-MODE
	Security string
	WPS      bool
	ExtCh    string
	NT       string
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *WiFiAP) UnmarshalJSON(b []byte) error {
	raw := struct {
		Band       string
		SSIDName   string `json:"ssidName"`
		WlsChannel string `json:"wlsChannel"`
		BSSID      string `json:"bssid"`
		Signal     string `json:"signal"`
		WMode      string `json:"wmode"`
		Security   string `json:"security"`
		WPS        string `json:"wps"`
		ExtCh      string `json:"extch"`
		NT         string `json:"nt"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiAP %q: %w", string(b), err)
	}

	s.SSID = raw.SSIDName
	s.Band = raw.Band
	s.WMode = raw.WMode
	s.ExtCh = raw.ExtCh
	s.NT = raw.NT

	s.WPS = raw.WPS == yes

	s.BSSID, _ = net.ParseMAC(strings.TrimSpace(raw.BSSID))
	s.Channel, _ = strconv.Atoi(raw.WlsChannel)

	return nil
}

// WiFiSSIDs -
type WiFiSSIDs struct {
	Error
	SSIDs  []SSID      `json:"SSIDs_List"`
	Guests []GuestSSID `json:"Guests_List"`
}

// SSID -
type SSID struct {
	BSSID        net.HardwareAddr
	Name         string `json:"ssidName"`
	Band         string
	IfName       string
	AuthMode     string
	SecurityMode string
	EncryptType  string
	Passphrase   string
	URI          string
	DefaultKey   string
	ID           int
	Radio        int
	Enable       bool
	EnableWPS    bool
	Visible      bool
	EnableWMM    bool
	EnableWLS    bool
	BandSteering bool
	Primary      bool
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *SSID) UnmarshalJSON(b []byte) error {
	raw := struct {
		ID          string `json:"id"`
		SSIDName    string `json:"ssidName"`
		Band        string `json:"band"`
		Enable      string `json:"enable"`
		WlswpsOnOff string `json:"wlswpsOnOff"`
		IfName      string `json:"ifName"`
		BSSID       string `json:"bssid"`
		Radio       string `json:"radio"`
		Visible     string `json:"visible"`
		WMM         string `json:"wmm"`
		AuthMode    string `json:"authMode"`
		SecuMode    string `json:"SecuMode"`
		EncryptType string `json:"encryptType"`
		Passphrase  string `json:"passPhrase"`
		WlsEnable   string `json:"wlsEnable"`
		URI         string `json:"SSID_URI"`
		DefaultKey  string `json:"defaultKey"`
		Bandsteer   string `json:"bandsteer"`
		Primary     string `json:"primary"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal SSID %q: %w", string(b), err)
	}

	s.Name = raw.SSIDName
	s.Band = raw.Band
	s.IfName = raw.IfName
	s.AuthMode = raw.AuthMode
	s.SecurityMode = raw.SecuMode
	s.EncryptType = raw.EncryptType
	s.Passphrase = raw.Passphrase
	s.URI = raw.URI
	s.DefaultKey = raw.DefaultKey

	s.Radio, _ = strconv.Atoi(raw.Radio)

	s.Enable = raw.Enable == on
	s.EnableWPS = raw.WlswpsOnOff == on
	s.EnableWLS = raw.WlsEnable == on
	s.Visible = raw.Visible == on
	s.EnableWMM = raw.WMM == on
	s.BandSteering = raw.Bandsteer == on
	s.Primary = raw.Primary == yes

	s.BSSID, _ = net.ParseMAC(strings.TrimSpace(raw.BSSID))

	s.ID, _ = strconv.Atoi(raw.ID)

	return nil
}

// GuestSSID -
type GuestSSID struct {
	Enable bool
	IfName string
	Relate string
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *GuestSSID) UnmarshalJSON(b []byte) error {
	raw := struct {
		Enable string
		IfName string
		Relate string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal GuestSSID %q: %w", string(b), err)
	}

	s.Enable = raw.Enable == on
	s.IfName = raw.IfName
	s.Relate = raw.Relate

	return nil
}

// WiFiWPS -
type WiFiWPS struct {
	Error
	Method      string
	ClientPin   string
	Status      string
	TimeElapsed time.Duration
	Enable      bool
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *WiFiWPS) UnmarshalJSON(b []byte) error {
	raw := struct {
		Error
		WlswpsOnOff       string
		WlsWpsMethod      string
		WlsWpsClientPin   string
		WlsWpsStatus      string
		WlsWpsTimeElapsed string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiWPS %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.Method = raw.WlsWpsMethod
	s.ClientPin = raw.WlsWpsClientPin
	s.Status = raw.WlsWpsStatus

	s.Enable = raw.WlswpsOnOff == on

	t, _ := strconv.Atoi(raw.WlsWpsTimeElapsed)
	s.TimeElapsed = time.Duration(t) * time.Second

	return nil
}

// WiFiClient -
type WiFiClient struct {
	Error
	Clients []WiFiClientEntry `json:"Client_List"`
}

// WiFiClientEntry -
type WiFiClientEntry struct {
	Index     int
	Band      string
	SSID      string
	Hostname  string
	MACAddr   net.HardwareAddr
	AID       int
	RSSI      int // Received Signal Strength Indicator. Estimated measure of power level that a client is receiving from AP.
	DataRate  int64
	PhyMode   string
	Channel   int
	Bandwidth int64
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *WiFiClientEntry) UnmarshalJSON(b []byte) error {
	raw := struct {
		AID       string `json:"aid"`
		Index     int    `json:"index"`
		Band      string `json:"band"`
		SSID      string `json:"ssid"`
		Hostname  string `json:"hostname"`
		MACAddr   string `json:"mac"`
		RSSI      string `json:"rssi"`
		DataRate  string `json:"br"`
		PhyMode   string `json:"pm"`
		Channel   string `json:"ch"`
		Bandwidth string `json:"bw"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WiFiClientEntry %q: %w", string(b), err)
	}

	s.AID, _ = strconv.Atoi(raw.AID)
	s.Index = raw.Index
	s.Band = raw.Band
	s.SSID = raw.SSID
	s.Hostname = raw.Hostname
	s.MACAddr, _ = net.ParseMAC(raw.MACAddr)
	s.RSSI, _ = strconv.Atoi(raw.RSSI)
	s.PhyMode = raw.PhyMode
	s.Channel, _ = strconv.Atoi(raw.Channel)

	if strings.HasSuffix(raw.Bandwidth, "MHz") {
		s.Bandwidth, _ = strconv.ParseInt(raw.Bandwidth[0:len(raw.Bandwidth)-3], 10, 64)
		s.Bandwidth *= 1_000_000
	}

	if strings.HasSuffix(raw.DataRate, "M") {
		s.DataRate, _ = strconv.ParseInt(raw.DataRate[0:len(raw.DataRate)-1], 10, 64)
		// track as bits per second (value was in mebibits/sec)
		//nolint:gomnd
		s.DataRate *= 1024 * 1024
	}

	return nil
}
