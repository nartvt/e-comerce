package http

import (
	"component-master/config"
	"component-master/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HttpServer struct {
	appName string
	conf    *config.ServerInfo
}

func (r *HttpServer) Start() {

}

func (r *HttpServer) InitHttpServer() {
	app := fiber.New(r.ConfigFiber(r.conf))
	app.Use(middleware.CorsFilter())
}

func (r *HttpServer) ConfigFiber(conf *config.ServerInfo) fiber.Config {
	return fiber.Config{
		AppName:           "Fiber App",
		EnablePrintRoutes: true,
		ReadTimeout:       time.Duration(conf.ConnectTimeOut) * time.Millisecond,
		WriteTimeout:      time.Duration(conf.ConnectTimeOut) * time.Millisecond,
	}
}
