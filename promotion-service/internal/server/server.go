package main

import (
	"promotion-service/infra/db"
	"promotion-service/infra/grpc"
)

func main() {
	setupInfra()
	defer closeInfra()
}

func setupInfra() {
	db.InitDB()
	grpc.InitGrpcServer()
}

func closeInfra() {
	db.CloseDB()
}
