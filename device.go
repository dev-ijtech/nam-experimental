package nam

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	MaxDeviceNameLen           int = 255
	MaxDeviceManagementIPv4Len int = 15
)

const (
	Cisco   = "cisco"
	Juniper = "juniper"
)

type Device struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	ManagementIPv4 string    `json:"managementIpv4"`
	Vendor         string    `json:"vendor"`
	Version        string    `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// TODO: Change the Valid function to implement the Validator interface properly

// func (d *Device) Valid() map[string]string {
// 	problems := make(map[string]string)

// 	if d.Name == "" {
// 		problems["name"] = "device name is empty."
// 	} else if utf8.RuneCountInString(d.Name) > MaxDeviceNameLen {
// 		problems["name"] = fmt.Sprintf("length of device name is longer than %d characters.", MaxDeviceNameLen)
// 	}

// 	if d.ManagementIPv4 == "" {
// 		problems["managementIpv4"] = "device management IPv4 does not exist."
// 	} else if utf8.RuneCountInString(d.ManagementIPv4) > MaxDeviceNameLen {
// 		problems["managementIpv4"] = "bad IPv4 address given."
// 	}

// 	if d.Vendor == "" {
// 		problems["vendor"] = "device vendor is empty."
// 	}

// 	switch strings.ToLower(d.Vendor) {
// 	case Cisco:
// 	case Juniper:

// 	default:
// 		{
// 			problems["vendor"] = fmt.Sprintf("unrecognised device vendor %s.", d.Vendor)
// 		}
// 	}

// 	if d.Version == "" {
// 		problems["vendor"] = "device version is empty."
// 	}

// 	return problems
// }

func (d *Device) Valid() error {
	if d.Name == "" {
		return fmt.Errorf("device validation: name does not exist")
	} else if utf8.RuneCountInString(d.Name) > MaxDeviceNameLen {
		return fmt.Errorf("device validation: name length is longer than %d", MaxDeviceNameLen)
	}

	if d.ManagementIPv4 == "" {
		return fmt.Errorf("device validation: management ipv4 does not exist")
	} else if utf8.RuneCountInString(d.ManagementIPv4) > MaxDeviceNameLen {
		return fmt.Errorf("device validation: management ipv4 length is longer than %d", MaxDeviceManagementIPv4Len)
	}

	if d.Vendor == "" {
		return fmt.Errorf("device validation: vendor does not exist")
	}

	switch strings.ToLower(d.Vendor) {
	case Cisco:
	case Juniper:

	default:
		{
			return fmt.Errorf("device validation: unrecognised vendor %s", d.Vendor)
		}
	}

	if d.Version == "" {
		return fmt.Errorf("device validation: version does not exist")
	}

	return nil
}

type DeviceStore interface {
	DeviceExists(id int) bool
	FindDeviceByID(id int) (*Device, error)
	FindDevices() ([]*Device, int, error)
	CreateDevice(device *Device) error
	UpdateDevice(id int, device *Device) error
	DeleteDevice(id int) error
}
