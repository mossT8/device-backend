package entity

import (
	"database/sql"
	"errors"

	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

func (a *Account) AddAccount(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO accounts (email, password_hash, salt, name, receive_updates, verified, created_at, modified_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	result, err := stmt.ExecContext(ctx,
		a.Email,
		a.PasswordHash,
		a.Salt,
		a.Name,
		a.ReceivesUpdates,
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

func (a *Account) GetAccountByEmail(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT a.ID, a.password_hash, a.salt, a.name, a.receive_updates, a.verified, a.created_at, a.modified_at
        FROM accounts a
        WHERE a.email = ? AND a.active = 1;
    `, a.Email).Scan(
		&a.ID,
		&a.PasswordHash,
		&a.Salt,
		&a.Name,
		&a.ReceivesUpdates,
		&a.Verified,
		&a.CreatedAt,
		&a.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundAccountByEmail
		}
		return qErr
	}

	return nil
}

func (a *Account) GetAccountByID(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT a.email, a.password_hash, a.salt, a.name, a.receive_updates, a.verified, a.created_at, a.modified_at
        FROM accounts a
        WHERE a.ID = ? AND a.active = 1;
    `, a.ID).Scan(
		&a.Email,
		&a.PasswordHash,
		&a.Salt,
		&a.Name,
		&a.ReceivesUpdates,
		&a.Verified,
		&a.CreatedAt,
		&a.ModifiedAt,
	); qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return domain.ErrNotFoundAccountByID
		}
		return qErr
	}

	return nil
}

func (a *Account) CountAccounts(conn datastore.MySqlDataStore) (*int64, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	var count int64
	if qErr := conn.ReaderDB.QueryRowContext(ctx, `
        SELECT COUNT(a.ID)
        FROM accounts a
		WHERE a.active = 1;
    `).Scan(
		&count,
	); qErr != nil {
		return nil, qErr
	}

	return &count, nil
}

func (a *Account) ListAccounts(conn datastore.MySqlDataStore, page, pageSize int64) ([]Account, error) {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	rows, qErr := conn.ReaderDB.QueryContext(ctx, `
        SELECT a.ID, a.email, a.password_hash, a.salt, a.name, a.receive_updates, a.verified, a.created_at, a.modified_at
        FROM accounts a
		WHERE active = 1
        LIMIT ?
        OFFSET ?;
    `, pageSize, page*pageSize)
	if qErr != nil {
		return nil, qErr
	}

	defer func() {
		conn.CloseRows(rows)
	}()

	accounts := make([]Account, 0)
	for rows.Next() {
		account := Account{}
		if sErr := rows.Scan(
			&account.ID,
			&account.Email,
			&account.PasswordHash,
			&account.Salt,
			&account.Name,
			&account.ReceivesUpdates,
			&account.Verified,
			&account.CreatedAt,
			&account.ModifiedAt,
		); sErr != nil {
			return nil, sErr
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
func (a *Account) UpdateAccount(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE accounts a
		SET a.email = ?, a.name = ?, a.receive_updates = ?, a.verified = ?, a.modified_at = ?
		WHERE a.ID = ? AND a.active = 1;
	`)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	_, err = stmt.ExecContext(ctx,
		a.Email,
		a.Name,
		a.ReceivesUpdates,
		a.Verified,
		a.ModifiedAt,
		a.ID,
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

func (a *Account) DeleteAccount(conn datastore.MySqlDataStore) error {
	ctx, cancel := conn.NewSqlContext()
	defer cancel()

	tx, cErr := conn.WriterDB.BeginTx(ctx, nil)
	if cErr != nil {
		return cErr
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE accounts a
		SET a.active = 0
		WHERE a.ID = ?;
	`)
	if err != nil {
		return err
	}

	defer func() {
		conn.CloseStatement(stmt)
	}()

	_, err = stmt.ExecContext(ctx, a.ID)
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
