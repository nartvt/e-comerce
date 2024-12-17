package client

import (
	"component-master/config"
	"component-master/infra/logging"
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	doOne sync.Once
	conn  *grpc.ClientConn

	grpcReadTimeout  time.Duration
	grpcWriteTimeout time.Duration
)

const (
	// second unit
	WriteTimeOutDefault = 120
	ReadTimeOutDefault  = 120
	ContextTimeout      = 120 * time.Second
)

type ClientConfig struct {
	Host            string
	Port            int
	EnableTLS       bool
	Timeout         time.Duration
	KeepAliveParams keepalive.ClientParameters
}

func ContextwithTimeout() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	go cancelContext(ContextTimeout, cancel)
	return ctx
}

func cancelContext(timeout time.Duration, cancel context.CancelFunc) {
	time.Sleep(timeout)
	cancel()
	slog.Info("context canceled")
}

func InitConnection(conf config.GrpcConfigClient, logConf config.LogConfig) (*grpc.ClientConn, error) {
	cfg := mapClientConfig(conf)
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(cfg.KeepAliveParams),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(logging.UnaryClientInterceptor(logging.InitLogger(logConf))),
	}

	if cfg.EnableTLS {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}

	grpcAddress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(grpcAddress, opts...)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func mapClientConfig(conf config.GrpcConfigClient) ClientConfig {
	if conf.ReadTimeOut < 1000 || conf.WriteTimeOut < 1000 {
		conf.ReadTimeOut = ReadTimeOutDefault
		conf.WriteTimeOut = WriteTimeOutDefault
	}
	return ClientConfig{
		Host:      conf.Host,
		Port:      conf.Port,
		EnableTLS: false,
		Timeout:   time.Duration(conf.ReadTimeOut) * time.Second,
		KeepAliveParams: keepalive.ClientParameters{
			Time:                time.Duration(conf.WriteTimeOut) * time.Second,
			Timeout:             time.Duration(conf.WriteTimeOut) * time.Second,
			PermitWithoutStream: true,
		}}
}
