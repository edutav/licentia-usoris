package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SMTP     SMTPConfig
	JWT      JWTConfig
	Env      Environment
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type JWTConfig struct {
	SecretKey string
}

type Environment struct {
	Env string
}

func Load() (*Config, error) {
	dir, _ := os.Getwd()
	log.Println("Current directory: ", dir)

	env := os.Getenv("APP_ENV")
	if env == "" {
		return nil, fmt.Errorf("APP_ENV not set")
	}

	file := fmt.Sprintf("config.%s", env)
	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
		return nil, err
	}

	log.Printf("Using config: %s", viper.ConfigFileUsed())

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		return nil, err
	}

	config.Env.Env = env

	return &config, nil
}
