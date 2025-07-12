package config

import (
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

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
