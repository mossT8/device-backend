package entity

import (
	"database/sql"
	"errors"

	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

func (s *Sensor) AddSensor(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO sensors (unit_id, code, name, config_required, config_default)
        VALUES (?, ?, ?, ?, ?);
    `)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	result, err := stmt.ExecContext(ctx,
		s.UnitId,
		s.Code,
		s.Name,
		s.ConfigRequried,
		s.DefaultConfig,
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

	s.SetID(lastId)

	return nil
}

func (s *Sensor) GetSensorByID(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT s.unit_id, s.code, s.name, s.config_required, s.config_default
        FROM sensors s
        WHERE s.ID = ?;
    `, s.ID).Scan(
		&s.UnitId,
		&s.Code,
		&s.Name,
		&s.ConfigRequried,
		&s.DefaultConfig,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundSensorByID
		}
		return qErr
	}

	return nil
}

func (s *Sensor) GetSensorByCode(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT s.ID, s.unit_id, s.name, s.config_required, s.config_default
        FROM sensors s
        WHERE s.code = ?;
    `, s.Code).Scan(
		&s.ID,
		&s.UnitId,
		&s.Name,
		&s.ConfigRequried,
		&s.DefaultConfig,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundSensorByCode
		}
		return qErr
	}

	return nil
}

func (s *Sensor) CountSensors(conn datastore.MySqlDataStore) (*int64, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	var count int64
	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT COUNT(s.ID)
        FROM sensors s;
    `).Scan(
		&count,
	); qErr != nil {
		return nil, qErr
	}

	return &count, nil
}

func (s *Sensor) ListSensors(conn datastore.MySqlDataStore, page, pageSize int64) ([]Sensor, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	rows, qErr := conn.ReaderDB.QueryContext(ctx, `
        SELECT s.ID, s.unit_id, s.code, s.name, s.config_required, s.config_default
        FROM sensors s
        ORDER BY s.ID
        LIMIT ? OFFSET ?;
    `, pageSize, page*pageSize)
	if qErr != nil {
		return nil, qErr
	}

	defer func() {
		conn.CloseRows(rows)
	}()

	sensors := make([]Sensor, 0)
	for rows.Next() {
		sensor := Sensor{}
		if sErr := rows.Scan(
			&sensor.ID,
			&sensor.UnitId,
			&sensor.Code,
			&sensor.Name,
			&sensor.ConfigRequried,
			&sensor.DefaultConfig,
		); sErr != nil {
			return nil, sErr
		}
		sensors = append(sensors, sensor)
	}

	return sensors, nil
}
