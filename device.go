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
	ManagementIPv4 string    `json:"managementIPv4"`
	Vendor         string    `json:"vendor"`
	Version        string    `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (d Device) Valid() ProblemSet {
	problems := ProblemSet{}

	if d.Name == "" {
		problems.Add("Device name", "is empty")
	} else if utf8.RuneCountInString(d.Name) > MaxDeviceNameLen {
		problems.Add("Device name", fmt.Sprintf("length of device name is longer than %d characters.", MaxDeviceNameLen))
	}

	if d.ManagementIPv4 == "" {
		problems.Add("Device management IPv4 address", "does not exist")
	} else if utf8.RuneCountInString(d.ManagementIPv4) > MaxDeviceNameLen {
		problems.Add("Device management IPv4 address", "is invalid")
	}

	if d.Vendor == "" {
		problems.Add("Device vendor", "does not exist")
	}

	switch strings.ToLower(d.Vendor) {
	case Cisco:
	case Juniper:

	default:
		{
			problems.Add("Device vendor", fmt.Sprintf("unrecognised device vendor %s.", d.Vendor))
		}
	}

	if d.Version == "" {
		problems.Add("Device version", "is empty")
	}

	return problems
}

type DeviceStore interface {
	DeviceExists(id int) bool
	FindDeviceByID(id int) (*Device, error)
	FindDevices() ([]*Device, int, error)
	CreateDevice(device *Device) error
	UpdateDevice(id int, device *Device) error
	DeleteDevice(id int) error
}

type DeviceFilter struct {
	ID *int `json:"id"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type DeviceUpdate struct {
	Name           *string `json:"name"`
	ManagementIPv4 *string `json:"managementIPv4"`
	Vendor         *string `json:"vendor"`
	Version        *string `json:"version"`
}
