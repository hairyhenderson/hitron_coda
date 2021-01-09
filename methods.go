// File generated with 'go generate'. Do not edit!

package hitron

import (
	"context"
)


// CMDocsisProvision - /CM/DocsisProvision
func (c *CableModem) CMDocsisProvision(ctx context.Context) (out CMDocsisProvision, err error) {
	err = c.getJSON(ctx, "/CM/DocsisProvision", &out)

	return out, err
}

// CMDsInfo - /CM/DsInfo
func (c *CableModem) CMDsInfo(ctx context.Context) (out CMDsInfo, err error) {
	err = c.getJSON(ctx, "/CM/DsInfo", &out)

	return out, err
}

// CMDsOfdm - /CM/DsOfdm
func (c *CableModem) CMDsOfdm(ctx context.Context) (out CMDsOfdm, err error) {
	err = c.getJSON(ctx, "/CM/DsOfdm", &out)

	return out, err
}

// CMLog - /CM/Log
func (c *CableModem) CMLog(ctx context.Context) (out CMLog, err error) {
	err = c.getJSON(ctx, "/CM/Log", &out)

	return out, err
}

// CMSysInfo - /CM/SysInfo
func (c *CableModem) CMSysInfo(ctx context.Context) (out CMSysInfo, err error) {
	err = c.getJSON(ctx, "/CM/SysInfo", &out)

	return out, err
}

// CMUsInfo - /CM/UsInfo
func (c *CableModem) CMUsInfo(ctx context.Context) (out CMUsInfo, err error) {
	err = c.getJSON(ctx, "/CM/UsInfo", &out)

	return out, err
}

// CMUsOfdm - /CM/UsOfdm
func (c *CableModem) CMUsOfdm(ctx context.Context) (out CMUsOfdm, err error) {
	err = c.getJSON(ctx, "/CM/UsOfdm", &out)

	return out, err
}

// CMVersion - /CM/Version
func (c *CableModem) CMVersion(ctx context.Context) (out CMVersion, err error) {
	err = c.getJSON(ctx, "/CM/Version", &out)

	return out, err
}

// DDNS - /DDNS
func (c *CableModem) DDNS(ctx context.Context) (out DDNS, err error) {
	err = c.getJSON(ctx, "/DDNS", &out)

	return out, err
}

// DNS - /DNS
func (c *CableModem) DNS(ctx context.Context) (out DNS, err error) {
	err = c.getJSON(ctx, "/DNS", &out)

	return out, err
}

// Hosts - /Hosts
func (c *CableModem) Hosts(ctx context.Context) (out Hosts, err error) {
	err = c.getJSON(ctx, "/Hosts", &out)

	return out, err
}

// RouterCapability - /Router/Capability
func (c *CableModem) RouterCapability(ctx context.Context) (out RouterCapability, err error) {
	err = c.getJSON(ctx, "/Router/Capability", &out)

	return out, err
}

// RouterDMZ - /Router/DMZ
func (c *CableModem) RouterDMZ(ctx context.Context) (out RouterDMZ, err error) {
	err = c.getJSON(ctx, "/Router/DMZ", &out)

	return out, err
}

// RouterLocation - /Router/Location
func (c *CableModem) RouterLocation(ctx context.Context) (out RouterLocation, err error) {
	err = c.getJSON(ctx, "/Router/Location", &out)

	return out, err
}

// RouterPortForwardStatus - /Router/PortForward/Status
func (c *CableModem) RouterPortForwardStatus(ctx context.Context) (out RouterPortForwardStatus, err error) {
	err = c.getJSON(ctx, "/Router/PortForward/Status", &out)

	return out, err
}

// RouterPortForwardall - /Router/PortForward/all
func (c *CableModem) RouterPortForwardall(ctx context.Context) (out RouterPortForwardall, err error) {
	err = c.getJSON(ctx, "/Router/PortForward/all", &out)

	return out, err
}

// RouterPortTriggerStatus - /Router/PortTrigger/Status
func (c *CableModem) RouterPortTriggerStatus(ctx context.Context) (out RouterPortTriggerStatus, err error) {
	err = c.getJSON(ctx, "/Router/PortTrigger/Status", &out)

	return out, err
}

// RouterPortTriggerall - /Router/PortTrigger/all
func (c *CableModem) RouterPortTriggerall(ctx context.Context) (out RouterPortTriggerall, err error) {
	err = c.getJSON(ctx, "/Router/PortTrigger/all", &out)

	return out, err
}

// RouterSysInfo - /Router/SysInfo
func (c *CableModem) RouterSysInfo(ctx context.Context) (out RouterSysInfo, err error) {
	err = c.getJSON(ctx, "/Router/SysInfo", &out)

	return out, err
}

// RouterTR069 - /Router/TR069
func (c *CableModem) RouterTR069(ctx context.Context) (out RouterTR069, err error) {
	err = c.getJSON(ctx, "/Router/TR069", &out)

	return out, err
}

// Time - /Time
func (c *CableModem) Time(ctx context.Context) (out Time, err error) {
	err = c.getJSON(ctx, "/Time", &out)

	return out, err
}

// WiFiAccessControl - /WiFi/AccessControl
func (c *CableModem) WiFiAccessControl(ctx context.Context) (out WiFiAccessControl, err error) {
	err = c.getJSON(ctx, "/WiFi/AccessControl", &out)

	return out, err
}

// WiFiAccessControlStatus - /WiFi/AccessControl/Status
func (c *CableModem) WiFiAccessControlStatus(ctx context.Context) (out WiFiAccessControlStatus, err error) {
	err = c.getJSON(ctx, "/WiFi/AccessControl/Status", &out)

	return out, err
}

// WiFiGuestSSID - /WiFi/GuestSSID
func (c *CableModem) WiFiGuestSSID(ctx context.Context) (out WiFiGuestSSID, err error) {
	err = c.getJSON(ctx, "/WiFi/GuestSSID", &out)

	return out, err
}

// WiFiRadios - /WiFi/Radios
func (c *CableModem) WiFiRadios(ctx context.Context) (out WiFiRadios, err error) {
	err = c.getJSON(ctx, "/WiFi/Radios", &out)

	return out, err
}

// WiFiRadiosAdvanced - /WiFi/Radios/Advanced
func (c *CableModem) WiFiRadiosAdvanced(ctx context.Context) (out WiFiRadiosAdvanced, err error) {
	err = c.getJSON(ctx, "/WiFi/Radios/Advanced", &out)

	return out, err
}

// WiFiRadiosSurvey - /WiFi/Radios/Survey
func (c *CableModem) WiFiRadiosSurvey(ctx context.Context) (out WiFiRadiosSurvey, err error) {
	err = c.getJSON(ctx, "/WiFi/Radios/Survey", &out)

	return out, err
}

// WiFiSSIDs - /WiFi/SSIDs
func (c *CableModem) WiFiSSIDs(ctx context.Context) (out WiFiSSIDs, err error) {
	err = c.getJSON(ctx, "/WiFi/SSIDs", &out)

	return out, err
}

// WiFiWPS - /WiFi/WPS
func (c *CableModem) WiFiWPS(ctx context.Context) (out WiFiWPS, err error) {
	err = c.getJSON(ctx, "/WiFi/WPS", &out)

	return out, err
}

