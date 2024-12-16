package main

import (
	"promotion-service/infra/db"
	elasticsearch "promotion-service/infra/elastic-search"
	"promotion-service/infra/grpc"
)

func main() {
	setupInfra()
	defer closeInfra()
}

func setupInfra() {
	db.InitDB()
	elasticsearch.InitES()
	grpc.InitGrpcServer()
}

func closeInfra() {
	db.CloseDB()
	elasticsearch.CloseES()
}
