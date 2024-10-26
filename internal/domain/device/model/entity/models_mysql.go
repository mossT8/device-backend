package entity

import (
	"database/sql"
	"errors"

	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

func (u *Models) GetModelByID(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT m.name, m.code, m.created_at, m.modified_at
        FROM models m
        WHERE m.ID = ?;
    `, u.ID).Scan(
		&u.Name,
		&u.Code,
		&u.CreatedAt,
		&u.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundModelByID
		}
		return qErr
	}

	return nil
}

func (u *Models) CountModels(conn datastore.MySqlDataStore) (*int64, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	var count int64
	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT COUNT(m.ID)
        FROM models m;
    `).Scan(
		&count,
	); qErr != nil {
		return nil, qErr
	}

	return &count, nil
}

func (u *Models) ListModels(conn datastore.MySqlDataStore, page, pageSize int64) ([]Models, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	rows, qErr := conn.ReaderDB.QueryContext(ctx, `
        SELECT m.ID, m.name, m.code, m.created_at, m.modified_at
        FROM models m
        ORDER BY m.ID
        LIMIT ? OFFSET ?;
    `, pageSize, page*pageSize)
	if qErr != nil {
		return nil, qErr
	}

	defer func() {
		conn.CloseRows(rows)
	}()

	models := make([]Models, 0)
	for rows.Next() {
		tempModel := Models{}
		if sErr := rows.Scan(
			&tempModel.ID,
			&tempModel.Name,
			&tempModel.Code,
			&tempModel.CreatedAt,
			&tempModel.ModifiedAt,
		); sErr != nil {
			return nil, sErr
		}
		models = append(models, tempModel)
	}

	return models, nil
}
