package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type MongoConfig struct {
	URI      string
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type EncryptionConfig struct {
	Key string
}

type Config struct {
	Mongo      MongoConfig
	Encryption EncryptionConfig
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Load from ENV first
	cfg.Mongo.URI = os.Getenv("MONGO_CONNECTION")
	cfg.Mongo.Host = os.Getenv("MONGO_HOST")
	cfg.Mongo.Port = os.Getenv("MONGO_PORT")
	cfg.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_PASSWORD")
	cfg.Mongo.Database = os.Getenv("MONGO_DATABASE")
	cfg.Encryption.Key = os.Getenv("ENCRYPTION_KEY")

	// If ENV not found, fallback to config.yaml
	if cfg.Mongo.URI == "" && cfg.Mongo.Host == "" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}

		// Load from file
		cfg.Mongo.URI = viper.GetString("mongo.connection")
		cfg.Mongo.Host = viper.GetString("mongo.host")
		cfg.Mongo.Port = viper.GetString("mongo.port")
		cfg.Mongo.User = viper.GetString("mongo.user")
		cfg.Mongo.Password = viper.GetString("mongo.password")
		cfg.Mongo.Database = viper.GetString("mongo.database")
		cfg.Encryption.Key = viper.GetString("encryption.key")
	}

	// Fallback: Build URI if not set
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

	// Debug log
	log.Println("[Mongo URI]", cfg.Mongo.URI)
	log.Println("[Mongo Host]", cfg.Mongo.Host)
	log.Println("[Mongo User]", cfg.Mongo.User)
	log.Println("[Mongo Password]", cfg.Mongo.Password)
	log.Println("[Mongo Database]", cfg.Mongo.Database)
	log.Println("[Encryption Key]", cfg.Encryption.Key)

	return cfg, nil
}
