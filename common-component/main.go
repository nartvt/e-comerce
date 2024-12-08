package main

import (
	"common-component/config"
	"fmt"
	"log"
)

func main() {
	conf, err := config.LoadConfig("config", "dev")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf.Server.Port)
}
