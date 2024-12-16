package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config config

type config struct {
	Postgres      postgresqlInfo
	Server        serverInfo
	ElasticSearch elasticSearchInfo
}

type serverInfo struct {
	Host    string
	Port    int
	TimeOut int
}

type postgresqlInfo struct {
	Host     string
	Port     int
	UserName string
	Password string
	Database string
}

type elasticSearchInfo struct {
	Host     string
	Port     int
	Username string
	Password string
}

func init() {
	load()
}

func load() {
	viper.SetConfigName("config")
	viper.AddConfigPath("app/config")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Config can't be read")
		return
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Panic("Config can't be load")
		return
	}
	log.Println(Config)
}
