package namsql

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/dev-ijtech/nam-experimental"
)

type DeviceStore struct {
	store []*nam.Device

	db *sql.DB
}

func NewDeviceService(db *sql.DB) *DeviceStore {
	return &DeviceStore{
		store: make([]*nam.Device, 0),
		db:    db,
	}
}

func (s *DeviceStore) DeviceExists(id int) bool {
	index := slices.IndexFunc(s.store, func(d *nam.Device) bool { return (d.ID == id) })

	return index == -1
}

func (s *DeviceStore) FindDeviceByID(ctx context.Context, id int) (*nam.Device, error) {
	var device nam.Device

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("find device by id: %v", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT * WHERE ID=? FROM Device;", id)
	if err != nil {
		return nil, fmt.Errorf("find device by id: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&device.ID, &device.Name, &device.ManagementIPv4, &device.Vendor, &device.Version, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("find device by id: %v", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("find device by id: %v", err)
	}

	return &device, nil
}

func (s *DeviceStore) FindDevices() ([]*nam.Device, error) {
	return s.store, nil
}

// TODO change to accommodate the Validator interface
func (s *DeviceStore) MakeDevice(d *nam.Device) error {
	if err := d.Valid(); err != nil {
		return fmt.Errorf("create device: %v", err)
	}

	s.store = append(s.store, d)

	return nil
}

// TODO change to accommodate the Validator interface
func (s *DeviceStore) UpdateDevice(id int, d *nam.Device) error {
	index := slices.IndexFunc(s.store, func(d *nam.Device) bool { return (d.ID == id) })

	if index == -1 {
		return fmt.Errorf("update device: device with id %d does not exist", id)
	}

	if err := d.Valid(); err != nil {
		return fmt.Errorf("update device: %v", err)
	}

	s.store[index] = d

	return nil
}

func (s *DeviceStore) DeleteDevice(id int) error {
	index := slices.IndexFunc(s.store, func(d *nam.Device) bool { return (d.ID == id) })

	if index == -1 {
		return fmt.Errorf("delete device: device with id %d does not exist", id)
	}

	s.store = slices.Delete(s.store, index, index+1)

	return nil
}
