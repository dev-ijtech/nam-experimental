package southbound

import (
	"context"
	"time"

	"github.com/dev-ijtech/nam-experimental"
)

type JuniperDevice struct {
	Device         *nam.Device
	SouthboundImpl *SouthboundImpl
}

type JuniperSystemInformation struct {
	HardwareModel string `xml:"system-information>hardware-model"`
	OsName        string `xml:"system-information>os-name"`
	OsVersion     string `xml:"system-information>os-version"`
	SerialNumber  string `xml:"system-information>serial-number"`
	HostName      string `xml:"system-information>host-name"`
}

type JuniperSoftwareInformation struct {
	HardwareModel string `xml:"software-information>product-model"`
	OsVersion     string `xml:"software-information>junos-version"`
	HostName      string `xml:"software-information>host-name"`
}

func (j JuniperDevice) GetDeviceDetails(ctx context.Context) (*nam.Device, error) {
	swInfo, err := execRpc[JuniperSoftwareInformation]("<get-software-information/>", j.Device.ManagementIPv4, j.SouthboundImpl)
	if err != nil {
		return &nam.Device{}, err
	}

	j.Device.Name = swInfo.HostName
	j.Device.Version = swInfo.OsVersion
	j.Device.UpdatedAt = time.Now()

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				j.SouthboundImpl.Logger.Printf("GetDeviceDetails err: %s\n", err)
			}
			j.SouthboundImpl.Logger.Printf("GetDeviceDetails cancelled!\n")
			return &nam.Device{}, nil
		default:
			return j.Device, nil
		}
	}

}
