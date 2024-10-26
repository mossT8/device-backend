package entity

import (
	"database/sql"
	"errors"

	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

func (d *Device) AddDevice(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO devices (account_id, device_name, serial_number, model_id, model_config, created_at, modified_at)
        VALUES (?, ?, ?, ?, ?, ?, ?);
    `)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	result, err := stmt.ExecContext(ctx,
		d.AccountId,
		d.Name,
		d.SerialNumber,
		d.ModelId,
		d.ModelConfig,
		d.CreatedAt,
		d.ModifiedAt,
	)
	if err != nil {
		return err
	}

	cErr = tx.Commit()
	if cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	lastId, cErr := result.LastInsertId()
	if cErr != nil {
		return cErr
	}

	d.SetID(lastId)

	return nil
}

func (d *Device) GetDeviceByID(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT d.account_id, d.device_name, d.serial_number, d.model_id, d.model_config, d.created_at, d.modified_at
        FROM devices d
        WHERE d.ID = ?;
    `, d.ID).Scan(
		&d.AccountId,
		&d.Name,
		&d.SerialNumber,
		&d.ModelId,
		&d.ModelConfig,
		&d.CreatedAt,
		&d.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundDeviceByID
		}
		return qErr
	}

	return nil
}

func (d *Device) GetDeviceBySerialNumber(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT d.ID, d.account_id, d.device_name, d.model_id, d.model_config, d.created_at, d.modified_at
        FROM devices d
        WHERE d.serial_number = ?;
    `, d.SerialNumber).Scan(
		&d.ID,
		&d.AccountId,
		&d.Name,
		&d.ModelId,
		&d.ModelConfig,
		&d.CreatedAt,
		&d.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundDeviceBySerialNumber
		}
		return qErr
	}

	return nil
}

func (d *Device) CountDevices(conn datastore.MySqlDataStore) (*int64, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	var count int64
	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT COUNT(d.ID)
        FROM devices d
		WHERE d.account_id = ?;
    `, d.AccountId).Scan(
		&count,
	); qErr != nil {
		return nil, qErr
	}

	return &count, nil
}

func (d *Device) ListDevices(conn datastore.MySqlDataStore, page, pageSize int64) ([]Device, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	rows, qErr := conn.ReaderDB.QueryContext(ctx, `
        SELECT d.ID, d.device_name, d.serial_number, d.model_id, d.model_config, d.created_at, d.modified_at
        FROM devices d
		WHERE d.account_id = ?
        ORDER BY d.ID
        LIMIT ? OFFSET ?;
    `, d.AccountId, pageSize, page*pageSize)
	if qErr != nil {
		return nil, qErr
	}

	defer func() {
		conn.CloseRows(rows)
	}()

	devices := make([]Device, 0)
	for rows.Next() {
		device := Device{AccountId: d.AccountId}
		if sErr := rows.Scan(
			&device.ID,
			&device.Name,
			&device.SerialNumber,
			&device.ModelId,
			&device.ModelConfig,
			&device.CreatedAt,
			&device.ModifiedAt,
		); sErr != nil {
			return nil, sErr
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (d *Device) UpdateDevice(conn datastore.MySqlDataStore, tx *sql.Tx) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if tx == nil {
		var cErr error
		tx, cErr = conn.WriterDB.BeginTx(ctx, nil)
		if cErr != nil {
			return cErr
		}
	}

	stmt, tErr := tx.PrepareContext(ctx, `
        UPDATE devices d
        SET d.device_name = ?, d.serial_number = ?, d.model_id = ?, d.model_config = ?, d.modified_at = ?
        WHERE d.ID = ? AND d.account_id = ?;
    `)
	if tErr != nil {
		return tErr
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	if _, sErr := stmt.ExecContext(ctx,
		d.Name,
		d.SerialNumber,
		d.ModelId,
		d.ModelConfig,
		d.ModifiedAt,
		d.ID,
		d.AccountId,
	); sErr != nil {
		return sErr
	}

	if cErr := tx.Commit(); cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	return nil
}

func (d *Device) DeleteDevice(conn datastore.MySqlDataStore, tx *sql.Tx) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if tx == nil {
		var cErr error
		tx, cErr = conn.WriterDB.BeginTx(ctx, nil)
		if cErr != nil {
			return cErr
		}
	}

	stmt, tErr := tx.PrepareContext(ctx, `
        DELETE FROM devices
        WHERE ID = ? AND account_id = ?;
    `)
	if tErr != nil {
		return tErr
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	if _, sErr := stmt.ExecContext(ctx, d.ID, d.AccountId); sErr != nil {
		return sErr
	}

	if cErr := tx.Commit(); cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	return nil
}
