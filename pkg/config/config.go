package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"gorm.io/gorm/logger"
)

type ServerConfig struct {
	Host         string
	Port         string
	IdleTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	LogLevel        logger.LogLevel
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

func Load() (Config, error) {
	_ = gotenv.Load()

	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("SERVER_PORT", "3000")
	viper.SetDefault("IDLE_TIMEOUT_SEC", 5)
	viper.SetDefault("WRITE_TIMEOUT_SEC", 5)
	viper.SetDefault("READ_TIMEOUT_SEC", 5)

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "ishari")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_CONN_MAX_LIFETIME_MIN", 60)
	viper.SetDefault("DB_CONN_MAX_IDLE_TIME_MIN", 10)
	viper.SetDefault("DB_LOG_LEVEL", "warn")

	// JWT defaults
	viper.SetDefault("JWT_ACCESS_TOKEN_TTL_MIN", 15)
	viper.SetDefault("JWT_REFRESH_TOKEN_TTL_DAYS", 7)

	viper.AutomaticEnv()

	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, errors.New("JWT_SECRET environment variable is not set")
	}
	if jwtSecret == "your-super-secret-key-change-in-production" {
		return Config{}, errors.New("JWT_SECRET is set to the default insecure value. Please change it in production")
	}

	logLevel := logger.Warn
	switch viper.GetString("DB_LOG_LEVEL") {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	return Config{
		Server: ServerConfig{
			Host:         viper.GetString("SERVER_HOST"),
			Port:         viper.GetString("SERVER_PORT"),
			IdleTimeout:  time.Duration(viper.GetInt("IDLE_TIMEOUT_SEC")) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt("WRITE_TIMEOUT_SEC")) * time.Second,
			ReadTimeout:  time.Duration(viper.GetInt("READ_TIMEOUT_SEC")) * time.Second,
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetString("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			DBName:          viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSLMODE"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime: time.Duration(viper.GetInt("DB_CONN_MAX_LIFETIME_MIN")) * time.Minute,
			ConnMaxIdleTime: time.Duration(viper.GetInt("DB_CONN_MAX_IDLE_TIME_MIN")) * time.Minute,
			LogLevel:        logLevel,
		},
		JWT: JWTConfig{
			Secret:          jwtSecret,
			AccessTokenTTL:  time.Duration(viper.GetInt("JWT_ACCESS_TOKEN_TTL_MIN")) * time.Minute,
			RefreshTokenTTL: time.Duration(viper.GetInt("JWT_REFRESH_TOKEN_TTL_DAYS")) * 24 * time.Hour,
		},
	}, nil
}
