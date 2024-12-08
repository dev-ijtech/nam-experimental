package namsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (s *DeviceStore) FindDeviceByID(ctx context.Context, id int) (*nam.Device, error) {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("find device by id: %v", err)
	}

	defer tx.Rollback()

	device, err := findDialByID(ctx, tx, id)

	if err != nil {
		return &nam.Device{}, fmt.Errorf("find device by id: %v", err)
	}

	return device, nil
}

func (s *DeviceStore) FindDevices(ctx context.Context, filter nam.DeviceFilter) ([]*nam.Device, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, 0, fmt.Errorf("find devices: %v", err)
	}

	defer tx.Rollback()

	devices, n, err := findDevices(ctx, tx, filter)

	if err != nil {
		return devices, 0, fmt.Errorf("find devices: %v", err)
	}

	return devices, n, nil
}

func (s *DeviceStore) CreateDevice(ctx context.Context, d *nam.Device) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("create device: %v", err)
	}

	defer tx.Rollback()

	err = createDevice(ctx, tx, d)

	if err != nil {
		return fmt.Errorf("create device: %v", err)
	}

	return tx.Commit()
}

func (s *DeviceStore) UpdateDevice(ctx context.Context, id int, update nam.DeviceUpdate) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("update device: %v", err)
	}

	defer tx.Rollback()

	err = updateDevice(ctx, tx, id, update)

	if err != nil {
		return fmt.Errorf("update device: %v", err)
	}

	return tx.Commit()
}

func (s *DeviceStore) DeleteDevice(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("delete device: %v", err)
	}

	defer tx.Rollback()

	err = deleteDevice(ctx, tx, id)

	if err != nil {
		return fmt.Errorf("delete device: %v", err)
	}

	return tx.Commit()
}

func checkDeviceExists(ctx context.Context, tx *sql.Tx, id int) (bool, error) {
	var n int

	row := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM Device WHERE id = ?", id)

	err := row.Scan(&n)

	if err != nil {
		return false, err
	}

	return n != 0, nil
}

func findDialByID(ctx context.Context, tx *sql.Tx, id int) (*nam.Device, error) {
	dials, _, err := findDevices(ctx, tx, nam.DeviceFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(dials) == 0 {
		return nil, errors.New("device not found")
	}
	return dials[0], nil
}

func findDevices(ctx context.Context, tx *sql.Tx, filter nam.DeviceFilter) (devices []*nam.Device, n int, err error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	limit_offset := ""

	if v := filter.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}

	if filter.Limit > 0 && filter.Offset > 0 {
		limit_offset = fmt.Sprintf(`LIMIT %d OFFSET %d`, filter.Limit, filter.Offset)
	} else if filter.Limit > 0 {
		limit_offset = fmt.Sprintf(`LIMIT %d`, filter.Limit)
	} else if filter.Offset > 0 {
		limit_offset = fmt.Sprintf(`OFFSET %d`, filter.Offset)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT
		    id,
		    name,
		    managementIPv4,
		    vendor,
		    version,
		    createdAt,
		    updatedAt
		FROM Device
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+limit_offset,
		args...,
	)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var device nam.Device
		var createdAt, updatedAt string

		if err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.ManagementIPv4,
			&device.Vendor,
			&device.Version,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, 0, err
		}

		if device.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			return nil, 0, err
		}

		if device.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, 0, err
		}

		devices = append(devices, &device)
		n++
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return
}

func createDevice(ctx context.Context, tx *sql.Tx, device *nam.Device) error {
	if device.CreatedAt.IsZero() {
		device.CreatedAt = time.Now()
		device.UpdatedAt = time.Now()
	}

	if device.UpdatedAt.IsZero() {
		device.UpdatedAt = time.Now()
	}

	_, err := tx.ExecContext(ctx, `
		INSERT INTO Device (
			name,
			managementIPv4,
			vendor,
			version,
			createdAt,
			updatedAt
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		device.Name,
		device.ManagementIPv4,
		device.Vendor,
		device.Version,
		device.CreatedAt.Format(time.RFC3339),
		device.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return err
	}

	return nil
}

func updateDevice(ctx context.Context, tx *sql.Tx, id int, update nam.DeviceUpdate) error {
	exists, err := checkDeviceExists(ctx, tx, id)

	if err != nil {
		return fmt.Errorf("check device exists: %v", err)
	}

	if !exists {
		return errors.New("device does not exist")
	}

	fields, args := []string{"updatedAt = ?"}, []interface{}{time.Now().Format(time.RFC3339)}

	if v := update.Name; v != nil {
		fields, args = append(fields, "name = ?"), append(args, *v)
	}

	if v := update.ManagementIPv4; v != nil {
		fields, args = append(fields, "managementIPv4 = ?"), append(args, *v)
	}

	if v := update.Vendor; v != nil {
		fields, args = append(fields, "vendor = ?"), append(args, *v)
	}

	if v := update.Version; v != nil {
		fields, args = append(fields, "version = ?"), append(args, *v)
	}

	args = append(args, id)

	_, err = tx.ExecContext(ctx, `
		UPDATE Device
		SET

		`+strings.Join(fields, ",")+`
		WHERE id = ?`,
		args...)

	if err != nil {
		return err
	}

	return nil
}

func deleteDevice(ctx context.Context, tx *sql.Tx, id int) error {
	exists, err := checkDeviceExists(ctx, tx, id)

	if err != nil {
		return fmt.Errorf("check device exists: %v", err)
	}

	if !exists {
		return errors.New("device does not exist")
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM Device WHERE id=?", id)

	if err != nil {
		return err
	}

	return nil
}
