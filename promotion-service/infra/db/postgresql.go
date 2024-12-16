package db

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"

	"promotion-service/config"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func InitDB() {
	conf := config.Config
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		conf.Postgres.UserName,
		conf.Postgres.Password,
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.Database,
	)

	var err error
	dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	once.Do(func() {
		if dbInstance == nil {
			InitDB()
		}
	})
	return dbInstance
}

func CloseDB() {
	if db, _ := dbInstance.DB(); db != nil {
		if err := db.Close(); err != nil {
			fmt.Println("[ERROR] Cannot close mysql connection, err:", err)
		}
	}
}
