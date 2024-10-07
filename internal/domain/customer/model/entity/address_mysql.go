package entity

import (
	"database/sql"
	"errors"

	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

func (a *Address) AddAddress(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO addresses (account_ID, name, address_line1, address_line2, city, state, postal_code, country, verified, created_at, modified_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	result, err := stmt.ExecContext(ctx,
		a.AccountId,
		a.Name,
		a.AddressLine1,
		a.AddressLine2,
		a.City,
		a.State,
		a.PostalCode,
		a.Country,
		a.Verified,
		a.CreatedAt,
		a.ModifiedAt,
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

	a.SetID(lastId)

	return nil
}

func (a *Address) GetAddressByID(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT a.account_ID, a.name, a.address_line1, a.address_line2, a.city, a.state, a.postal_code, a.country, a.verified, a.created_at, a.modified_at
        FROM addresses a
        WHERE a.ID = ? AND a.active = 1;
    `, a.ID).Scan(
		&a.AccountId,
		&a.Name,
		&a.AddressLine1,
		&a.AddressLine2,
		&a.City,
		&a.State,
		&a.PostalCode,
		&a.Country,
		&a.Verified,
		&a.CreatedAt,
		&a.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundAddressByID
		}
		return qErr
	}

	return nil
}

func (a *Address) CountAddresses(conn datastore.MySqlDataStore) (*int64, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	var count int64
	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM addresses a
        WHERE a.account_ID = ? AND a.active = 1;
    `, a.AccountId).Scan(
		&count,
	); qErr != nil {
		return nil, qErr
	}

	return &count, nil
}

func (a *Address) ListAddresses(conn datastore.MySqlDataStore, page, pageSize int64) ([]Address, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	rows, qErr := conn.ReaderDB.QueryContext(ctx, `
        SELECT a.ID, a.account_ID, a.name, a.address_line1, a.address_line2, a.city, a.state, a.postal_code, a.country, a.verified, a.created_at, a.modified_at
        FROM addresses a
        WHERE a.account_ID = ? AND a.active = 1
        LIMIT ?
        OFFSET ?;
    `, a.AccountId, pageSize, page*pageSize)
	if qErr != nil {
		return nil, qErr
	}

	defer func() {
		conn.CloseRows(rows)
	}()

	addresses := make([]Address, 0)
	for rows.Next() {
		address := Address{AccountId: a.AccountId}
		if sErr := rows.Scan(
			&address.ID,
			&address.AccountId,
			&address.Name,
			&address.AddressLine1,
			&address.AddressLine2,
			&address.City,
			&address.State,
			&address.PostalCode,
			&address.Country,
			&address.Verified,
			&address.CreatedAt,
			&address.ModifiedAt,
		); sErr != nil {
			return nil, sErr
		}

		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (a *Address) UpdateAddress(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE addresses
		SET name = ?, address_line1 = ?, address_line2 = ?, city = ?, state = ?, postal_code = ?, country = ?, verified = ?, modified_at = ?
		WHERE ID = ? AND account_ID = ? AND active = 1;
	`)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	_, err = stmt.ExecContext(ctx,
		a.Name,
		a.AddressLine1,
		a.AddressLine2,
		a.City,
		a.State,
		a.PostalCode,
		a.Country,
		a.Verified,
		a.ModifiedAt,
		a.ID,
		a.AccountId,
	)
	if err != nil {
		return err
	}

	cErr = tx.Commit()
	if cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	return nil
}

func (a *Address) DeleteAddress(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE addresses
		SET active = 0
		WHERE ID = ? AND account_ID = ? AND active = 1;
	`)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	_, err = stmt.ExecContext(ctx, a.ID, a.AccountId)
	if err != nil {
		return err
	}

	cErr = tx.Commit()
	if cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	return nil
}
