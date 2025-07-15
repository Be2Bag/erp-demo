package config

import (
	"fmt"
	"os"
	"strconv"
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

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type WasabiConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	Endpoint  string
}

type Config struct {
	Mongo      MongoConfig
	Encryption EncryptionConfig
	Swagger    SwaggerConfig
	Hash       HashConfig
	JWT        JWTConfig
	Email      EmailConfig
	Wasabi     WasabiConfig
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
	cfg.Email.Host = os.Getenv("EMAIL_HOST")
	if portStr := os.Getenv("EMAIL_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid EMAIL_PORT: %v", err)
		}
		cfg.Email.Port = port
	}
	cfg.Email.Username = os.Getenv("EMAIL_USERNAME")
	cfg.Email.Password = os.Getenv("EMAIL_PASSWORD")
	cfg.Email.From = os.Getenv("EMAIL_FROM")
	cfg.Wasabi.AccessKey = os.Getenv("WASABI_ACCESS_KEY")
	cfg.Wasabi.SecretKey = os.Getenv("WASABI_SECRET_KEY")
	cfg.Wasabi.Bucket = os.Getenv("WASABI_BUCKET")
	cfg.Wasabi.Region = os.Getenv("WASABI_REGION")
	cfg.Wasabi.Endpoint = os.Getenv("WASABI_ENDPOINT")

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
		cfg.Email.Host = viper.GetString("email.host")
		cfg.Email.Port = viper.GetInt("email.port")
		cfg.Email.Username = viper.GetString("email.username")
		cfg.Email.Password = viper.GetString("email.password")
		cfg.Email.From = viper.GetString("email.from")
		cfg.Wasabi.AccessKey = viper.GetString("wasabi.access_key")
		cfg.Wasabi.SecretKey = viper.GetString("wasabi.secret_key")
		cfg.Wasabi.Bucket = viper.GetString("wasabi.bucket")
		cfg.Wasabi.Region = viper.GetString("wasabi.region")
		cfg.Wasabi.Endpoint = viper.GetString("wasabi.endpoint")

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
