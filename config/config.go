package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type MongoConfig struct {
	URI      string `mapstructure:"connection"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

type Config struct {
	Mongo      MongoConfig      `mapstructure:"mongo"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// map nested keys (dot) to underscores for ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// ignore missing config file, but error on other issues
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	// build URI from individual parts if not provided
	if cfg.Mongo.URI == "" {
		auth := ""
		if cfg.Mongo.User != "" && cfg.Mongo.Password != "" {
			auth = fmt.Sprintf("%s:%s@", cfg.Mongo.User, cfg.Mongo.Password)
		}
		addr := cfg.Mongo.Host
		if cfg.Mongo.Port != "" {
			addr = fmt.Sprintf("%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)
		}
		cfg.Mongo.URI = fmt.Sprintf("mongodb://%s%s/%s?retryWrites=true&w=majority",
			auth, addr, cfg.Mongo.Database)
	}
	return &cfg, nil
}
