package database

import (
	"common-component/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/zerolog/log"
)

func InitDatabaseConnect(conf *config.DatabaseConfig) (*sqlx.DB, error) {
	sqlConnectUrl := conf.BuildDatabaseConnectionString()
	log.Info().Msgf("connecting to database %s, %s", "url", sqlConnectUrl)
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
	return db, nil
}
