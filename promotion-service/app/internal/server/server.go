package main

import (
	"promotion-service/app/infra/db"
	"promotion-service/app/infra/grpc"
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
