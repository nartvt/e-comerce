package config_test

import (
	"common-component/config"
	"common-component/infra/database"
	"fmt"
	"log"
	"testing"
)

func TestInitconfig(t *testing.T) {
	conf := &config.Config{}
	err := config.LoadConfig("config", "dev", conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf.Server.Port)

	dbConf := conf.Database
	fmt.Println(dbConf.DriverName)

	db, err := database.InitDatabaseConnect(&dbConf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(db != nil)

	defer db.Close()
}
