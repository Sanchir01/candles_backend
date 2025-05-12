package config

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env"`
	Domain     string `yaml:"domain"`
	HttpServer `yaml:"http_server"`
	Errors     `yaml:"errors"`
	Grpc       GrpcOrder
	RedisDB    Redis      `yaml:"redis"`
	DB         DataBase   `yaml:"database"`
	S3Store    S3Store    `yaml:"s3store"`
	Prometheus Prometheus `yaml:"prometheus"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	DBNumber int    `yaml:"dbnumber"`
}
type Prometheus struct {
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Port        string        `yaml:"port"  env-default:"8081"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}
type DataBase struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Database    string `yaml:"dbname"`
	SSL         string `yaml:"ssl"`
	MaxAttempts int    `yaml:"max_attempts"`
}
type GrpcOrder struct {
	Timeout int `yaml:"timeout"`
	Port    int `yaml:"port"`
}
type S3Store struct {
	Key        string `yaml:"key"`
	BucketName string `yaml:"bucketname"`
	Region     string `yaml:"region"`
	URL        string `yaml:"url"`
}
type HttpServer struct {
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Port        string        `yaml:"port"  env-default:"5000"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}

type Errors struct {
	Unauthorized ErrorsBody `yaml:"unauthorized"  env-default:"Unauthorized"`
	NotFound     ErrorsBody `yaml:"not_found"  env-default:"Not found"`
}
type ErrorsBody struct {
	message string `yaml:"message"`
	code    int    `yaml:"code"`
}

func InitConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("ошибка при инициализации переменных окружения", err.Error())
	}
	configPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist:%s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
