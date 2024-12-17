package database

import (
	"component-master/config"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DatabaseClient struct {
	db *sqlx.DB
}

func InitDatabaseConnect(conf *config.DatabaseConfig) (*DatabaseClient, error) {
	sqlConnectUrl := conf.BuildDatabaseConnectionString()
	slog.Info(fmt.Sprintf("connecting to database %s, %s", "url", sqlConnectUrl))
	db, err := sqlx.Connect(conf.DriverName, sqlConnectUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database, %v", err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConnections)
	db.SetMaxIdleConns(conf.MaxIdleConnections)

	if conf.MaxConnLifetime > 0 {
		db.SetConnMaxLifetime(conf.MaxConnLifetime)
	}

	if conf.MaxConnIdleTime > 0 {
		db.SetConnMaxIdleTime(conf.MaxConnIdleTime)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database, %v", err)
	}
	return &DatabaseClient{
		db: db,
	}, nil
}

func (db *DatabaseClient) Close() {
	slog.Info("closing database connection")
	err := db.db.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to close database connection, %v", err))
	}
}

func (db *DatabaseClient) BeginTx() (*sqlx.Tx, error) {
	tx, err := db.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction, %v", err)
	}
	return tx, nil
}

func (db *DatabaseClient) RollbackTx(tx *sqlx.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return fmt.Errorf("failed to rollback transaction, %v", err)
	}
	return nil
}

func (db *DatabaseClient) CommitTx(tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction, %v", err)
	}
	return nil
}

func (db *DatabaseClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query, %v", err)
	}
	return result, nil
}

func (db *DatabaseClient) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := db.db.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query database, %v", err)
	}
	return rows, nil
}

func (db *DatabaseClient) QueryRow(query string, args ...interface{}) *sqlx.Row {
	row := db.db.QueryRowx(query, args...)
	return row
}

func (db *DatabaseClient) Prepare(query string) (*sqlx.Stmt, error) {
	stmt, err := db.db.Preparex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement, %v", err)
	}
	return stmt, nil
}

func (db *DatabaseClient) GetDB() *sqlx.DB {
	return db.db
}
