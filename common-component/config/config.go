package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig        `mapstructure:"server"`
	Database DatabaseConfig      `mapstructure:"database"`
	Redis    RedisConfig         `mapstructure:"redis"`
	Elastic  ElasticSearchConfig `mapstructure:"elastic"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	DriverName         string        `mapstructure:"driverName"`
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	UserName           string        `mapstructure:"userName"`
	Password           string        `mapstructure:"password"`
	DBName             string        `mapstructure:"dbName"`
	SSLMode            string        `mapstructure:"sslMode"`
	MaxOpenConnections int           `mapstructure:"maxOpenConnections"`
	MaxIdleConnections int           `mapstructure:"maxIdleConnections"`
	MaxConnLifetime    time.Duration `mapstructure:"maxConnLifetime"`
	MaxConnIdleTime    time.Duration `mapstructure:"maxConnIdleTime"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	MaxIdle      int           `mapstructure:"maxIdle"`
}

type ElasticSearchConfig struct {
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	Index              string `mapstructure:"index"`
	MaxRetry           int    `mapstructure:"maxRetry"`
	RetryOnStatusCodes []int  `mapstructure:"retryOnStatusCodes"`
}

func loadConfig(configfile string) (*viper.Viper, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found. Using exists environment variables")
		return nil, err
	}
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
	v, err := loadConfig(configPath + ".yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
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
