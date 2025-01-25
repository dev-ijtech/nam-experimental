package southbound

import (
	"errors"
	"log"

	"github.com/dev-ijtech/nam-experimental"
)

type SouthboundImpl struct {
	LoginUser     string
	LoginPassword string
	Logger        *log.Logger
}

func NewSouthboundService(loginUser string, loginPassword string, logger *log.Logger) *SouthboundImpl {
	return &SouthboundImpl{LoginUser: loginUser, LoginPassword: loginPassword, Logger: logger}
}

func (s SouthboundImpl) DeviceFactory(device *nam.Device) (nam.SouthboundOps, error) {
	switch device.Vendor {
	case nam.Juniper:
		return &JuniperDevice{Device: device, SouthboundImpl: &s}, nil
	default:
		return nil, errors.New("southbound service not supported for vendor " + device.Vendor)
	}
}
