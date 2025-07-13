package config

import (
	"fmt"
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

type SwaggerConfig struct {
	Key string
}

type HashConfig struct {
	Salt string
}

type JWTConfig struct {
	SecretKey string
}

type Config struct {
	Mongo      MongoConfig
	Encryption EncryptionConfig
	Swagger    SwaggerConfig
	Hash       HashConfig
	JWT        JWTConfig
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.Mongo.URI = os.Getenv("MONGO_CONNECTION")
	cfg.Mongo.Host = os.Getenv("MONGO_HOST")
	cfg.Mongo.Port = os.Getenv("MONGO_PORT")
	cfg.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_PASSWORD")
	cfg.Mongo.Database = os.Getenv("MONGO_DATABASE")
	cfg.Encryption.Key = os.Getenv("ENCRYPTION_KEY")
	cfg.Swagger.Key = os.Getenv("SWAGGER_KEY")
	cfg.Hash.Salt = os.Getenv("HASH_SALT")
	cfg.JWT.SecretKey = os.Getenv("JWT_SECRET")

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

		cfg.Mongo.URI = viper.GetString("mongo.connection")
		cfg.Mongo.Host = viper.GetString("mongo.host")
		cfg.Mongo.Port = viper.GetString("mongo.port")
		cfg.Mongo.User = viper.GetString("mongo.user")
		cfg.Mongo.Password = viper.GetString("mongo.password")
		cfg.Mongo.Database = viper.GetString("mongo.database")
		cfg.Encryption.Key = viper.GetString("encryption.key")
		cfg.Swagger.Key = viper.GetString("swagger.key")
		cfg.Hash.Salt = viper.GetString("hash.salt")
		cfg.JWT.SecretKey = viper.GetString("jwt.secret")

	}

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

	return cfg, nil
}
