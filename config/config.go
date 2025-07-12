package config

import (
	"strings"

	"github.com/spf13/viper"
)

type MongoConfig struct {
	URI      string `mapstructure:"connection"`
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
	return &cfg, nil
}
