package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"mossT8.github.com/device-backend/internal/application/config/types"
)

type Store struct {
	WriterDB *sql.DB
	ReaderDB *sql.DB
	config   types.EngineDB
}

// NewStore will create a variable that represent the Repository struct
func NewStore(config types.EngineDB) *Store {
	return &Store{WriterDB: nil, ReaderDB: nil, config: config}
}

func (r *Store) Start() error {
	writerDB, err := createConnection(r.config.Writer)
	if err != nil {
		return err
	}
	readerDB, err := createConnection(r.config.Reader)
	if err != nil {
		writerDB.Close()
		return err
	}

	r.WriterDB = writerDB
	r.ReaderDB = readerDB

	return nil
}

func createConnection(dbConfig *types.DBConfig) (*sql.DB, error) {
	if dbConfig.Port == 0 {
		dbConfig.Port = 3306
	}
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	db, err := sql.Open(dbConfig.Dialect, connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	if dbConfig.MaxConnections > 0 {
		db.SetMaxOpenConns(dbConfig.MaxConnections)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Close attaches the provider and close the connection
func (r *Store) Close() (error, error) {
	wErr := r.WriterDB.Close()
	rErr := r.ReaderDB.Close()
	return wErr, rErr
}

// Ping both writer and reader
func (r *Store) Ping() error {
	err := r.WriterDB.Ping()
	if err != nil {
		return err
	}
	err = r.ReaderDB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (r *Store) CreateWriterTransaction() (*sql.Tx, error) {
	return r.WriterDB.Begin()
}
