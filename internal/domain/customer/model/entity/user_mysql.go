package entity

import (
	"database/sql"
	"errors"

	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

func (u *User) AddUser(conn datastore.MySqlDataStore, tx *sql.Tx) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if tx == nil {
		var cErr error
		if tx, cErr = conn.WriterDB.BeginTx(ctx, nil); cErr != nil {
			return cErr
		}
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO users (
			account_id,
			email,
			cell,
			first_name,
			last_name,
			verified,
			receive_updates,
			created_at,
			modified_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	result, err := stmt.ExecContext(ctx,
		u.AccountId,
		u.Email,
		u.Cell,
		u.FirstName,
		u.LastName,
		u.Verified,
		u.ReceivesUpdates,
		u.CreatedAt,
		u.ModifiedAt,
	)
	if err != nil {
		return err
	}

	if cErr := tx.Commit(); cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	lastId, cErr := result.LastInsertId()
	if cErr != nil {
		return cErr
	}
	u.SetID(lastId)

	return nil
}

func (u *User) GetUserByEmail(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
		SELECT u.ID, u.account_ID, u.cell, u.first_name, u.last_name, u.receive_updates, u.verified, u.created_at, u.modified_at
		FROM users u
		WHERE u.email = ? AND u.active = 1;
	`, u.Email).Scan(
		&u.ID,
		&u.AccountId,
		&u.Cell,
		&u.FirstName,
		&u.LastName,
		&u.ReceivesUpdates,
		&u.Verified,
		&u.CreatedAt,
		&u.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundUserByEmail
		}
		return qErr
	}

	return nil
}

func (u *User) GetUserByID(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
		SELECT u.email, u.cell, u.first_name, u.last_name, u.receive_updates, u.verified, u.created_at, u.modified_at
		FROM users u
		WHERE u.account_ID = ? AND u.ID = ? AND u.active = 1;
	`, u.AccountId, u.ID).Scan(
		&u.Email,
		&u.Cell,
		&u.FirstName,
		&u.LastName,
		&u.ReceivesUpdates,
		&u.Verified,
		&u.CreatedAt,
		&u.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundUserByID
		}
		return qErr
	}

	return nil
}

func (u *User) CountUsers(conn datastore.MySqlDataStore) (*int64, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	var count int64
	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM users u
		WHERE u.account_ID = ? AND u.active = 1;
	`, u.AccountId).Scan(
		&count,
	); qErr != nil {
		return nil, qErr
	}

	return &count, nil
}

func (u *User) ListUsers(conn datastore.MySqlDataStore, page, pageSize int64) ([]User, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	rows, qErr := conn.ReaderDB.QueryContext(ctx, `
		SELECT u.ID, u.email, u.cell, u.first_name, u.last_name, u.receive_updates, u.verified, u.created_at, u.modified_at
		FROM users u
		WHERE u.account_ID = ? AND u.active = 1
		LIMIT ?
		OFFSET ?;
	`, u.AccountId, pageSize, page*pageSize)
	if qErr != nil {
		return nil, qErr
	}

	defer func() {
		conn.CloseRows(rows)
	}()

	users := make([]User, 0)
	for rows.Next() {
		user := User{AccountId: u.AccountId}
		if sErr := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Cell,
			&user.FirstName,
			&user.LastName,
			&user.ReceivesUpdates,
			&user.Verified,
			&user.CreatedAt,
			&user.ModifiedAt,
		); sErr != nil {
			return nil, sErr
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *User) UpdateUser(conn datastore.MySqlDataStore, tx *sql.Tx) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if tx == nil {
		var cErr error
		tx, cErr = conn.WriterDB.BeginTx(ctx, nil)
		if cErr != nil {
			return cErr
		}
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE users
		SET email = ?, cell = ?, first_name = ?, last_name = ?, receive_updates = ?, verified = ?, modified_at = ?
		WHERE ID = ? AND account_ID = ?;
	`)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	if _, err = stmt.ExecContext(ctx,
		u.Email,
		u.Cell,
		u.FirstName,
		u.LastName,
		u.ReceivesUpdates,
		u.Verified,
		u.ModifiedAt,
		u.ID,
		u.AccountId,
	); err != nil {
		return err
	}

	if cErr := tx.Commit(); cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	return nil
}

func (u *User) DeleteUser(conn datastore.MySqlDataStore, tx *sql.Tx) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if tx == nil {
		var cErr error
		if tx, cErr = conn.WriterDB.BeginTx(ctx, nil); cErr != nil {
			return cErr
		}
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE users
		SET active = 0
		WHERE ID = ? AND account_ID = ?;
	`)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	if _, err = stmt.ExecContext(ctx, u.ID, u.AccountId); err != nil {
		return err
	}

	if cErr := tx.Commit(); cErr != nil {
		conn.RollbackAndJoinErrorIfAny(tx)
		return cErr
	}

	return nil
}
