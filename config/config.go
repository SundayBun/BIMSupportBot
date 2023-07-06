package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Postgres          PostgresConfig
	Server            ServerConfig
	Metrics           Metrics
	TelegramApiConfig TelegramBotApiConfig
	Mongo             MongoDbCongig
}

type MongoDbCongig struct {
	Url        string
	DataBase   string
	Collection string
}

type TelegramBotApiConfig struct {
	Token  string
	Domain string
	Secret string
}

// Server config struct
type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         bool
	Driver          string
	MaxOpenConns    int
	ConnMaxLifetime int
	MaxIdleConns    int
	ConnMaxIdleTime int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
