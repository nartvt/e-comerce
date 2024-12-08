package main

import (
	"common-component/config"
	"fmt"
	"log"
)

func main() {
	conf := &config.Config{}
	err := config.LoadConfig("config", "dev", conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf.Server.Port)
}
