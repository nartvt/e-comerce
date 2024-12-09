package config_test

import (
	"common-component/config"
	"common-component/infra/database"
	"common-component/infra/logging"
	"common-component/infra/redis"
	"common-component/infra/search"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestInitconfig(t *testing.T) {
	logging.InitLogger()
	conf := &config.Config{}
	err := config.LoadConfig("config", "dev", conf)
	if err != nil {
		log.Error().Msgf("failed to load config, %v", err)
		return
	}

	fmt.Println(conf.Server.Port)

	dbConf := conf.Database
	fmt.Println(dbConf.DriverName)

	db, err := database.InitDatabaseConnect(&dbConf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("DATABASE CONNECCT SUCCESS: ", db != nil)

	rd, err := redis.InitRedis(&conf.Redis)
	if err != nil {
		log.Error().Msgf("failed to init redis, %v", err)
		return
	}
	fmt.Println("REDIS CONNECT: ", rd != nil)

	defer db.Close()

	elasticClient, err := search.InitElasticSearch(&conf.Elastic)
	if err != nil {
		log.Error().Msgf("failed to init elastic search, %v", err)
		return
	}
	elasticClient.CreateIndex()
	fmt.Println("ELASTIC SEARCH: ", elasticClient != nil)
}
