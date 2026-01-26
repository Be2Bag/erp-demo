package config

import (
	"errors"
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
	Username string
	Password string
	From     string
	Port     int
}

type SupabaseConfig struct {
	AccessKey  string
	SecretKey  string
	Bucket     string
	Region     string
	Endpoint   string
	PublicBase string
}

type CloudflareConfig struct {
	AccessKey     string
	SecretKey     string
	Bucket        string
	Region        string
	Endpoint      string
	PublicBaseURL string
}

type Config struct {
	Mongo      MongoConfig
	Encryption EncryptionConfig
	Swagger    SwaggerConfig
	Hash       HashConfig
	JWT        JWTConfig
	Email      EmailConfig
	Supabase   SupabaseConfig
	Cloudflare CloudflareConfig
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
			return nil, fmt.Errorf("invalid EMAIL_PORT: %w", err)
		}
		cfg.Email.Port = port
	}
	cfg.Email.Username = os.Getenv("EMAIL_USERNAME")
	cfg.Email.Password = os.Getenv("EMAIL_PASSWORD")
	cfg.Email.From = os.Getenv("EMAIL_FROM")
	cfg.Supabase.AccessKey = os.Getenv("SUPABASE_ACCESS_KEY")
	cfg.Supabase.SecretKey = os.Getenv("SUPABASE_SECRET_KEY")
	cfg.Supabase.Bucket = os.Getenv("SUPABASE_BUCKET")
	cfg.Supabase.Region = os.Getenv("SUPABASE_REGION")
	cfg.Supabase.Endpoint = os.Getenv("SUPABASE_ENDPOINT")
	cfg.Supabase.PublicBase = os.Getenv("SUPABASE_PUBLIC_BASE")
	cfg.Cloudflare.AccessKey = os.Getenv("CLOUDFLARE_ACCESS_KEY")
	cfg.Cloudflare.SecretKey = os.Getenv("CLOUDFLARE_SECRET_KEY")
	cfg.Cloudflare.Bucket = os.Getenv("CLOUDFLARE_BUCKET")
	cfg.Cloudflare.Region = os.Getenv("CLOUDFLARE_REGION")
	cfg.Cloudflare.Endpoint = os.Getenv("CLOUDFLARE_ENDPOINT")
	cfg.Cloudflare.PublicBaseURL = os.Getenv("CLOUDFLARE_PUBLIC_BASE_URL")

	if cfg.Mongo.URI == "" && cfg.Mongo.Host == "" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if !errors.As(err, &configFileNotFoundError) {
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
		cfg.Supabase.AccessKey = viper.GetString("supabase.access_key")
		cfg.Supabase.SecretKey = viper.GetString("supabase.secret_key")
		cfg.Supabase.Bucket = viper.GetString("supabase.bucket")
		cfg.Supabase.Region = viper.GetString("supabase.region")
		cfg.Supabase.Endpoint = viper.GetString("supabase.endpoint")
		cfg.Supabase.PublicBase = viper.GetString("supabase.public_base")
		cfg.Cloudflare.AccessKey = viper.GetString("cloudflare.access_key")
		cfg.Cloudflare.SecretKey = viper.GetString("cloudflare.secret_key")
		cfg.Cloudflare.Bucket = viper.GetString("cloudflare.bucket")
		cfg.Cloudflare.Region = viper.GetString("cloudflare.region")
		cfg.Cloudflare.Endpoint = viper.GetString("cloudflare.endpoint")
		cfg.Cloudflare.PublicBaseURL = viper.GetString("cloudflare.public_base_url")
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
