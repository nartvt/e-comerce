package main

import (
	"component-master/config"
	"component-master/infra/database"
	"component-master/infra/logging"
	"component-master/infra/redis"
	"component-master/infra/search"
	"component-master/util"
	"fmt"
	"log/slog"
)

func init() {
	util.LoadEnv()
}

// func TestInitconfig(t *testing.T) {
func main() {
	conf := &config.Config{}
	err := config.LoadConfig("config", "dev", conf)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config, %v", err))
		return
	}
	logging.InitLogger(conf.Log)

	pringLogconfig(&conf.Log)
	printMiddlewareConfig(&conf.Middleware)

	fmt.Println(conf.Server.Port)

	dbConf := conf.Database
	fmt.Println(dbConf.DriverName)

	db, err := database.InitDatabaseConnect(&dbConf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	slog.Info(fmt.Sprintf("DATABASE CONNECCT SUCCESS: %v", db != nil && db.GetDB() != nil))

	rd, err := redis.InitRedis(&conf.Redis)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init redis, %v", err))
		return
	}
	slog.Info(fmt.Sprintf("REDIS CONNECT: %v", rd != nil))

	defer db.Close()

	elasticClient, err := search.InitElasticSearch(&conf.Elastic)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init elastic search, %v", err))
		return
	}
	elasticClient.CreateIndex()
	fmt.Println("ELASTIC SEARCH: ", elasticClient != nil)
}

func pringLogconfig(conf *config.LogConfig) {
	slog.Info(fmt.Sprintf("ENV: %s", conf.Environment))
	slog.Info(fmt.Sprintf("LEVEL: %s", conf.LogLevel))
	slog.Info(fmt.Sprintf("JSON OUTPUT: %t", conf.JSONOutput))
	slog.Info(fmt.Sprintf("ADD SOURCE: %t", conf.AddSource))
}

func printMiddlewareConfig(conf *config.MiddlewareConfig) {
	slog.Info(fmt.Sprintf("TOKEN SECRET: %s", conf.Token.AccessTokenSecret))
	slog.Info(fmt.Sprintf("TOKEN EXP: %s", conf.Token.AccessTokenExp))
	slog.Info(fmt.Sprintf("REFRESH TOKEN SECRET: %s", conf.Token.RefreshTokenSecret))
	slog.Info(fmt.Sprintf("REFRESH TOKEN EXP: %s", conf.Token.RefreshTokenExp))
}
