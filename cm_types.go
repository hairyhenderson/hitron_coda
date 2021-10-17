package hitron

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// CMVersion contains version information for the cable modem
type CMVersion struct {
	Error
	DeviceID        string `json:"deviceId"` // generally the same as the external MAC Address
	ModelName       string `json:"modelName"`
	VendorName      string `json:"vendorName"` // Usually always "Hitron Technologies"
	SerialNum       string `json:"SerialNum"`
	HwVersion       string `json:"HwVersion"`
	APIVersion      string `json:"ApiVersion"`
	SoftwareVersion string `json:"SoftwareVersion"`
}

func (v CMVersion) String() string {
	if v.Error != NoError && v.Error.Message != "" {
		return v.Error.String()
	}

	sb := strings.Builder{}
	sb.WriteString("DeviceID: ")
	sb.WriteString(v.DeviceID)
	sb.WriteString("\n")
	sb.WriteString("ModelName: ")
	sb.WriteString(v.ModelName)
	sb.WriteString("\n")
	sb.WriteString("VendorName: ")
	sb.WriteString(v.VendorName)
	sb.WriteString("\n")
	sb.WriteString("SerialNum: ")
	sb.WriteString(v.SerialNum)
	sb.WriteString("\n")
	sb.WriteString("HwVersion: ")
	sb.WriteString(v.HwVersion)
	sb.WriteString("\n")
	sb.WriteString("APIVersion: ")
	sb.WriteString(v.APIVersion)
	sb.WriteString("\n")
	sb.WriteString("SoftwareVersion: ")
	sb.WriteString(v.SoftwareVersion)
	sb.WriteString("\n")

	return sb.String()
}

// CMDocsisProvision contains DOCSIS provisioning state
// For each step:
// - "Process" indicates the CODA-4x8x is attempting to complete a connection step.
// - "Success" indicates the CODA-4x8x has completed a connection step.
// - "Disable" indicates the relevant feature has been turned off.
type CMDocsisProvision struct {
	Error
	HWInit         string `json:"hwInit"`         // "Success"
	FindDownstream string `json:"findDownstream"` // "Success"
	Ranging        string `json:"ranging"`        // "Success"
	DHCP           string `json:"dhcp"`           // "Success"
	TimeOfday      string `json:"timeOfday"`      // "Success"
	DownloadCfg    string `json:"downloadCfg"`    // "Success"
	Registration   string `json:"registration"`   // "Success"
	EAEStatus      string `json:"eaeStatus"`      // "Disable" - EARLY AUTHENTICATION AND ENCRYPTION
	BPIStatus      string `json:"bpiStatus"`      // "AUTH:start, TEK:start" - Baseline Privacy Interface
	NetworkAccess  string `json:"networkAccess"`  // "Permitted"
	TrafficStatus  string `json:"trafficStatus"`  // "Enable"
}

// BPIStatus - TODO
// type BPIStatus struct {
// 	AUTH string // Authorization finite state machine
// 	TEK  string // Traffic encryption keys FSM
// }

// CMDsInfo - Downstream Port Info
type CMDsInfo struct {
	Error
	Ports PortInfos `json:"Freq_List"`
}

// CMUsInfo - Upstream Port Info
type CMUsInfo struct {
	Error
	Ports PortInfos `json:"Freq_List"`
}

// PortInfos -
type PortInfos []PortInfo

// UnmarshalJSON - implements json.Unmarshaler for the PortInfos type. If one of
// the PortInfo entries can't be unmarshaled, we just skip it in this case.
func (p *PortInfos) UnmarshalJSON(b []byte) error {
	raw := []json.RawMessage{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	for _, port := range raw {
		pt := &PortInfo{}

		if pt.UnmarshalJSON(port) != nil {
			// if we can't handle it, just skip it
			continue
		}

		*p = append(*p, *pt)
	}

	return err
}

// PortInfo -
type PortInfo struct {
	PortID         string  `json:"portId"`                // "1"
	Frequency      int64   `json:"frequency,string"`      // in Hz
	Modulation     string  `json:"modulationType"`        // "QAM256"
	SignalStrength float64 `json:"signalStrength,string"` // signal strength of the downstream channel, in dBmV (dB above/below 1 mV)
	ChannelID      string  `json:"channelId"`             // "11"
	// upstream-only
	Bandwidth int64 `json:"bandwidth,string,omitempty"` // maximum available upstream bandwidth (in bits/sec, maybe?)
	// downstream-only
	SNR        float64 `json:"snr,string"`        // signal-to-noise ratio of the downstream data channel, in dB
	DsOctets   int64   `json:"dsoctets,string"`   // number of octets/bytes received
	Correcteds int64   `json:"correcteds,string"` // number of blocks that were corrupt, and were corrected.
	Uncorrect  int64   `json:"uncorrect,string"`  // number of blocks that were corrupt, but were unable to be corrected.
}

// UnmarshalJSON - implements json.Unmarshaler
func (p *PortInfo) UnmarshalJSON(b []byte) error {
	raw := struct {
		PortID         string `json:"portId"`
		Frequency      int64  `json:"frequency,string"`
		Modulation     string
		ModulationType string
		Bandwidth      int64 `json:"bandwidth,string"`
		SignalStrength string
		SNR            float64 `json:"snr,string"`
		ChannelID      string
		DSOctets       int64 `json:"dsoctets,string"`
		Correcteds     int64 `json:"correcteds,string"`
		Uncorrect      int64 `json:"uncorrect,string"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal PortInfo %q: %w", string(b), err)
	}

	p.PortID = raw.PortID
	p.ChannelID = raw.ChannelID
	p.Frequency = raw.Frequency
	p.Bandwidth = raw.Bandwidth
	p.DsOctets = raw.DSOctets
	p.Correcteds = raw.Correcteds
	p.Uncorrect = raw.Uncorrect
	p.SNR = raw.SNR

	p.Modulation = raw.Modulation
	if raw.Modulation == "" && raw.ModulationType != "" {
		p.Modulation = raw.ModulationType
	}

	p.SignalStrength = atof64(raw.SignalStrength)

	return nil
}

// CMSysInfo -
//nolint:lll
type CMSysInfo struct {
	Error
	NetworkAccess string           // Permitted/Denied - whether or not your service provider allows you to access the Internet over the CABLE connection.
	IP            net.IP           // WAN IP address negotiated by DHCP
	SubMask       net.IPMask       // WAN Subnet Mask
	GW            net.IP           // WAN Gateway IP
	Lease         time.Duration    // WAN DHCP "D: 6 H: 11 M: 10 S: 20"
	Configname    string           // (??)
	DsDataRate    int64            // WAN downstream data rate (bits/sec)
	UsDataRate    int64            // WAN upstream data rate (bits/sec)
	MacAddr       net.HardwareAddr // WAN MAC address
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *CMSysInfo) UnmarshalJSON(b []byte) error {
	raw := struct {
		// unmarshals fine
		Error
		NetworkAccess string   `json:"ntAccess"`   // "Permitted" (??)
		IPs           []net.IP `json:"ip"`         // should only be single element
		GW            net.IP   `json:"gw"`         // "7.96.63.1"
		Configname    string   `json:"Configname"` // "bac110000106749be82df7e0"
		// will need custom unmarshaling
		SubMask    string `json:"subMask"`    // "255.255.255.0"
		Lease      string `json:"lease"`      // "D: 6 H: 11 M: 10 S: 20"
		DsDataRate string `json:"DsDataRate"` // "1040000000"
		UsDataRate string `json:"UsDataRate"` // "31200000"
		MacAddr    string `json:"macAddr"`    // "74:9b:e8:2d:f7:e0"
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal SysInfo %q: %w", string(b), err)
	}

	s.Error = raw.Error
	s.NetworkAccess = raw.NetworkAccess
	s.IP = raw.IPs[0]
	s.GW = raw.GW
	s.Configname = raw.Configname

	maskIP := net.ParseIP(raw.SubMask)
	if maskIP.To4() != nil {
		maskIP = maskIP.To4()
	}

	s.SubMask = net.IPMask(maskIP)

	s.MacAddr, err = net.ParseMAC(raw.MacAddr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal MacAddr %q: %w", raw.MacAddr, err)
	}

	s.DsDataRate = atoi64(raw.DsDataRate)
	s.UsDataRate = atoi64(raw.UsDataRate)
	s.Lease = parseDHCPLeaseDuration(raw.Lease)

	return nil
}

// parseDHCPLeaseDuration - parses the 'lease' field from the CM/SysInfo
// response, which is a string of the form "D: n H: nn M: nn: S: nn", where n
// are 0-9. This implementation assumes that only the characters D, H, M, S, :,
// ' ' (space), and 0-9 will be present. Any other characters will be silently
// ignored.
func parseDHCPLeaseDuration(s string) (dur time.Duration) {
	if s == "" {
		return 0
	}

	//nolint:gomnd
	multMap := map[byte]time.Duration{
		'D': 24 * time.Hour,
		'H': time.Hour,
		'M': time.Minute,
		'S': time.Second,
	}

	// multiplier - gets set based on the D/H/M/S prefix
	mul := time.Duration(1)

	for s != "" {
		switch s[0] {
		case 'D', 'H', 'M', 'S':
			mul = multMap[s[0]]
			s = s[1:]

			continue
		case ':', ' ':
			s = s[1:]

			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// time to consume all of the 0-9 characters in a row...
			i, t := 0, 0
			for ; i < len(s); i++ {
				c := s[i]
				if c < '0' || c > '9' {
					break
				}

				t = t*10 + int(c) - '0'
			}

			s = s[i:]
			dur += time.Duration(t) * mul
		default:
			s = s[1:]

			continue
		}
	}

	return dur
}

// CMDsOfdm -
type CMDsOfdm struct {
	Error
	Receivers []OFDMReceiver `json:"Freq_List"`
}

// OFDMReceiver - OFDM Downstream Receiver information
//nolint:lll
type OFDMReceiver struct {
	ID             int     // OFDM Receiver index
	FFTType        string  // Type of FFT in use (NA/4K/etc...)
	SubcarrierFreq int64   // Frequency in Hz of the first OFDM subcarrier
	PLCLocked      bool    // whether or not this OFDM connection's Physical Link Channel data is locked. The PLC tells the CODA-4x8x how to decode the OFDM signal, and what power level to use. Once the CODA4x8x receives a PLC without uncorrectable errors, the PLC is locked and subsequent communication can continue.
	NCPLocked      bool    // whether or not this OFDM connection’s next codeword pointer (NCP) data is locked. The NCP tells the CODA-4x8x which codewords are to be used for OFDM communication, and which profile to use for each codeword. Once the CODA-4x8x receives an NCP without uncorrectable errors, the NCP is locked and subsequent communication can continue.
	MDC1Locked     bool    // whether or not this OFDM connection’s Multipath Delay Commutator (MDC) data is locked. This provides information about the method of Fast Fourier Transform (FFT) to be used on the OFDM connection. Once the CODA-4x8x receives an MDC1 without errors, the MDC1 is locked and subsequent communication can continue.
	PLCPower       float64 // power level the CODA-4x8x has been instructed to use on this OFDM connection by the physical link channel (PLC) data, in dBmV (decibels above/below 1 millivolt).
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *OFDMReceiver) UnmarshalJSON(b []byte) error {
	raw := struct {
		ID             int    `json:"receive"`
		FFTType        string `json:"ffttype"`
		SubcarrierFreq string `json:"Subcarr0freqFreq"`
		PLCLocked      string `json:"plclock"`
		NCPLocked      string `json:"ncplock"`
		MDC1Locked     string `json:"mdc1lock"`
		PLCPower       string `json:"plcpower"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal OFDMReceiver %q: %w", string(b), err)
	}

	s.ID = raw.ID
	if raw.FFTType != "NA" {
		s.FFTType = raw.FFTType
	}

	s.SubcarrierFreq = atoi64(raw.SubcarrierFreq)
	s.PLCPower = atof64(raw.PLCPower)

	s.PLCLocked = raw.PLCLocked == yes
	s.NCPLocked = raw.NCPLocked == yes
	s.MDC1Locked = raw.MDC1Locked == yes

	return nil
}

// CMUsOfdm - OFDM/OFDMA Upstream Channel Info
type CMUsOfdm struct {
	Error
	Channels []OFDMAChannel `json:"Freq_List"`
}

// OFDMAChannel - OFDM/OFDMA Channel
//nolint:lll
type OFDMAChannel struct {
	ID          int     // Channel index
	Enable      bool    //
	DigAtten    float64 // the digital attenuation, or signal loss, of the transmission medium on which the channel's signal is carried, in decibels (dB).
	DigAttenBo  float64 // the measured digital attenuation of the channel's signal, in decibels (dB). Digital attenuation is affected by the frequency of the signal; a higher-frequency signal will suffer more attenuation than a lower-frequency signal.
	ChannelBw   float64 // the bandwidth of this channel, expressed as the number of subchannels multiplied by the channel's Fast Fourier Transform size, in megahertz (MHz).
	RepPower    float64 // the reported power of this channel, in quarter-decibels above/below 1 millivolt (quarter-dBmV).
	RepPower1_6 float64 // the target power (P1.6r_n, or power spectral density in 1.6MHz) of this channel, in quarter-decibels above/below 1 millivolt (quarter- dBmV).
	FFTSize     string  // the type of Fast Fourier Transform in use on the relevant channel.
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *OFDMAChannel) UnmarshalJSON(b []byte) error {
	raw := struct {
		ID          int    `json:"uschindex"`
		State       string `json:"state"`       // :"  DISABLED"
		DigAtten    string `json:"digAtten"`    // :"    0.0000"
		DigAttenBo  string `json:"digAttenBo"`  // :"    0.0000"
		ChannelBw   string `json:"channelBw"`   // :"    0.0000"
		RepPower    string `json:"repPower"`    // :"    0.0000"
		RepPower1_6 string `json:"repPower1_6"` // :"    0.0000"
		FFTVal      string `json:"fftVal"`      // :"        2K"
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal OFDMAChannel %q: %w", string(b), err)
	}

	s.ID = raw.ID
	s.Enable = strings.TrimSpace(raw.State) == "ENABLED"
	s.FFTSize = strings.TrimSpace(raw.FFTVal)

	// there may be spaces, and may be "NA" or empty - unparseable should just be interpreted as 0
	s.DigAtten = atof64(raw.DigAtten)
	s.DigAttenBo = atof64(raw.DigAttenBo)
	s.ChannelBw = atof64(raw.ChannelBw)
	s.RepPower = atof64(raw.RepPower)
	s.RepPower1_6 = atof64(raw.RepPower1_6)

	return nil
}

// CMLog -
type CMLog struct {
	Error
	Logs []LogEntry `json:"Log_List"`
}

// LogEntry -
type LogEntry struct {
	ID       int
	Time     time.Time
	Type     string
	Severity string // syslog-style severity string & mapping
	Event    string
}

// UnmarshalJSON - implements json.Unmarshaler
func (s *LogEntry) UnmarshalJSON(b []byte) error {
	raw := struct {
		ID       int `json:"index"`
		Time     string
		Type     string
		Priority string
		Event    string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal LogEntry %q: %w", string(b), err)
	}

	s.ID = raw.ID
	s.Type = raw.Type
	s.Event = raw.Event

	sevMap := map[string]string{
		"1": "Emergency",
		"2": "Alert",
		"3": "Critical",
		"4": "Error",
		"5": "Warning",
		"6": "Notice",
		"7": "Information",
		"8": "Debug",
	}
	s.Severity = sevMap[raw.Priority]

	// Date format is MM/DD/YYYY HH:MM:SS - no timezone
	t, err := time.Parse("01/02/2006 15:04:05", raw.Time)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp %q: %w", raw.Time, err)
	}

	s.Time = t

	return nil
}
