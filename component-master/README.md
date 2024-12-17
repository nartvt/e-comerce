## COMPONENT MASTER

Keep it simple

#### Tech stack
- Golang(1.23.0)
- Docker
- Elastic search(7.17.0)
- Viper
- Grpc
- Redis
- Casbin
- Fiber V2
- Jwt Token
- PostgresSQL
- Slog (Build-in)

### Run test

### Grpc generate

```
make grpc

```

```bash
go test -run config/config.go config_test.go
```

### Setup elastic search
1. https://www.elastic.co/blog/getting-started-with-the-elastic-stack-and-docker-compose


### Docker commands
1. Docker stop all containers: `docker stop $(docker ps -a -q)`
2. Docker remove all containers: `docker rm $(docker ps -a -q)`
