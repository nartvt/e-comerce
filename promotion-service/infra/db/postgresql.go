package db

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"

	"promotion-service/config"
)

var instance *gorm.DB

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
	instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	if instance == nil {
		InitDB()
	}
	return instance
}

func CloseDB() {
	if db, _ := instance.DB(); db != nil {
		if err := db.Close(); err != nil {
			fmt.Println("[ERROR] Cannot close mysql connection, err:", err)
		}
	}
}
