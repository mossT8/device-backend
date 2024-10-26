package entity

import (
	"database/sql"
	"errors"

	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

func (u *Units) GetUnitByID(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT u.name, u.symbol
        FROM units u
        WHERE u.ID = ?;
    `, u.ID).Scan(
		&u.Name,
		&u.Symbol,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundUnitByID
		}
		return qErr
	}

	return nil
}

func (u *Units) GetUnitByName(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT u.ID, u.symbol
        FROM units u
        WHERE u.name = ?;
    `, u.Name).Scan(
		&u.ID,
		&u.Symbol,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundUnitByName
		}
		return qErr
	}

	return nil
}

func (u *Units) CountUnits(conn datastore.MySqlDataStore) (*int64, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	var count int64
	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT COUNT(u.ID)
        FROM units u;
    `).Scan(
		&count,
	); qErr != nil {
		return nil, qErr
	}

	return &count, nil
}

func (u *Units) ListUnits(conn datastore.MySqlDataStore, page, pageSize int64) ([]Units, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	rows, qErr := conn.ReaderDB.QueryContext(ctx, `
        SELECT u.ID, u.name, u.symbol
        FROM units u
        ORDER BY u.ID
        LIMIT ? OFFSET ?;
    `, pageSize, page*pageSize)
	if qErr != nil {
		return nil, qErr
	}

	defer func() {
		conn.CloseRows(rows)
	}()

	units := make([]Units, 0)
	for rows.Next() {
		unit := Units{}
		if sErr := rows.Scan(
			&unit.ID,
			&unit.Name,
			&unit.Symbol,
		); sErr != nil {
			return nil, sErr
		}
		units = append(units, unit)
	}

	return units, nil
}
