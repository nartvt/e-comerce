package config

import (
	"component-master/util"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig        `mapstructure:"server" json:"server,omitempty"`
	Database       DatabaseConfig      `mapstructure:"database" json:"database,omitempty"`
	Redis          RedisConfig         `mapstructure:"redis" json:"redis,omitempty"`
	Elastic        ElasticSearchConfig `mapstructure:"elastic" json:"elastic,omitempty"`
	Log            LogConfig           `mapstructure:"log" json:"log,omitempty"`
	Middleware     MiddlewareConfig    `mapstructure:"middleware" json:"middleware,omitempty"`
	GrpcPromotiion GrpcConfigClient    `mapstructure:"grpcPromotion" json:"grpc_promotiion,omitempty"`
}

type ServerConfig struct {
	Http ServerInfo `mapstructure:"http" json:"http,omitempty"`
	Grpc ServerInfo `mapstructure:"grpc" json:"grpc,omitempty"`
}

type ServerInfo struct {
	Host           string `mapstructure:"host" json:"host,omitempty"`
	Port           int    `mapstructure:"port" json:"port,omitempty"`
	EnableTLS      bool   `mapstructure:"enableTLS" json:"enable_tls,omitempty"`
	ConnectTimeOut int    `mapstructure:"connectTimeOut" json:"connect_time_out,omitempty"`
}

type DatabaseConfig struct {
	DriverName         string        `mapstructure:"driverName" json:"driver_name,omitempty"`
	Host               string        `mapstructure:"host" json:"host,omitempty"`
	Port               int           `mapstructure:"port" json:"port,omitempty"`
	UserName           string        `mapstructure:"userName" json:"user_name,omitempty"`
	Password           string        `mapstructure:"password" json:"password,omitempty"`
	DBName             string        `mapstructure:"dbName" json:"db_name,omitempty"`
	SSLMode            string        `mapstructure:"sslMode" json:"ssl_mode,omitempty"`
	MaxOpenConnections int           `mapstructure:"maxOpenConnections" json:"max_open_connections,omitempty"`
	MaxIdleConnections int           `mapstructure:"maxIdleConnections" json:"max_idle_connections,omitempty"`
	MaxConnLifetime    time.Duration `mapstructure:"maxConnLifetime" json:"max_conn_lifetime,omitempty"`
	MaxConnIdleTime    time.Duration `mapstructure:"maxConnIdleTime" json:"max_conn_idle_time,omitempty"`
}

type MiddlewareConfig struct {
	BasicAuth BasicAuthConfig `mapstructure:"basicAuth" json:"basic_auth,omitempty"`
	Casbin    CasbinConfig    `mapstructure:"casbin" json:"casbin,omitempty"`
	Token     TokenConfig     `mapstructure:"token" json:"token,omitempty"`
	Static    StaticConfig    `mapstructure:"static" json:"static,omitempty"`
}

type TokenConfig struct {
	AccessTokenSecret  string        `mapstructure:"accessTokenSecretKey" json:"access_token_secret,omitempty"`
	AccessTokenExp     time.Duration `mapstructure:"accessTokenExp" json:"access_token_exp,omitempty"`
	RefreshTokenSecret string        `mapstructure:"refreshTokenSecretKey" json:"refresh_token_secret,omitempty"`
	RefreshTokenExp    time.Duration `mapstructure:"refreshTokenExp" json:"refresh_token_exp,omitempty"`
}

type CasbinConfig struct {
	ModelPath  string `mapstructure:"modelPath" json:"model_path,omitempty"`
	PolicyPath string `mapstructure:"policyPath" json:"policy_path,omitempty"`
}

type BasicAuthConfig struct {
	Username string `mapstructure:"username" json:"username,omitempty"`
	Password string `mapstructure:"password" json:"password,omitempty"`
}

type StaticConfig struct {
	Unauthorized string `mapstructure:"unauthorized" json:"unauthorized,omitempty"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host" json:"host,omitempty"`
	Port         int           `mapstructure:"port" json:"port,omitempty"`
	Password     string        `mapstructure:"password" json:"password,omitempty"`
	DB           int           `mapstructure:"db" json:"db,omitempty"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout" json:"dial_timeout,omitempty"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout" json:"read_timeout,omitempty"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"write_timeout,omitempty"`
	MaxIdle      int           `mapstructure:"maxIdle" json:"max_idle,omitempty"`
}

type ElasticSearchConfig struct {
	Host               string `mapstructure:"host" json:"host,omitempty"`
	Port               int    `mapstructure:"port" json:"port,omitempty"`
	User               string `mapstructure:"user" json:"user,omitempty"`
	Password           string `mapstructure:"password" json:"password,omitempty"`
	Index              string `mapstructure:"index" json:"index,omitempty"`
	MaxRetry           int    `mapstructure:"maxRetry" json:"max_retry,omitempty"`
	RetryOnStatusCodes []int  `mapstructure:"retryOnStatusCodes" json:"retry_on_status_codes,omitempty"`
}

type GrpcConfigClient struct {
	Host         string
	Port         int
	ReadTimeOut  int
	WriteTimeOut int
}

type LogConfig struct {
	Environment string     `mapstructure:"env" json:"environment,omitempty"`
	LogLevel    slog.Level `mapstructure:"level" json:"log_level,omitempty"`
	JSONOutput  bool       `mapstructure:"jsonOutput" json:"json_output,omitempty"`
	AddSource   bool       `mapstructure:"addSource" json:"add_source,omitempty"`
}

func loadConfig(configfile string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(configfile)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	if err := v.ReadInConfig(); err != nil {
		log.Println("config file not found. Using exists environment variables")
		return nil, err
	}
	overrideconfig(v)
	return v, nil
}

func overrideconfig(v *viper.Viper) {
	for _, key := range v.AllKeys() {
		envKey := "APP_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		envValue := os.Getenv(envKey)
		if envValue != "" {
			v.Set(key, envValue)
		}

	}
}

func LoadConfig(pathToFile string, env string, config any) error {
	configPath := pathToFile + "/" + "application"
	if len(env) > 0 {
		configPath = configPath + "-" + env
	}

	pwd, err := os.Getwd()
	if err == nil && len(pwd) > 0 {
		configPath = pwd + "/" + configPath
	}

	confFile := fmt.Sprintf("%s.yaml", configPath)
	slog.Info(fmt.Sprintf("Config file path: %s", confFile))
	v, err := loadConfig(confFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	slog.Info(util.StructToJson(config))
	return nil
}

func (r *DatabaseConfig) BuildDatabaseConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		r.Host, r.Port, r.UserName, r.Password, r.DBName, func() string {
			if len(r.SSLMode) == 0 {
				return "disable"
			}
			return r.SSLMode
		}(),
	)
}

func (r *RedisConfig) BuildRedisConnectionString() string {
	return fmt.Sprintf("redis://%s:%s@%s:%d", "", "", r.Host, r.Port)
}

func (r *ElasticSearchConfig) BuildElasticSearchConnectionString() []string {
	return []string{
		fmt.Sprintf("http://%s:%d", r.Host, r.Port),
	}
}
