package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/albar2305/payment-app/utils/common"
)

type ApiConfig struct {
	ApiPort string
}

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   string
}

type FileConfig struct {
	FilePath string
}

type TokenConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	TokenSymetricKey     string
}

type Config struct {
	ApiConfig
	DbConfig
	FileConfig
	TokenConfig
}

// Method
func (c *Config) ReadConfig() error {
	err := common.LoadEnv()
	if err != nil {
		return err
	}

	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{
		ApiPort: os.Getenv("API_PORT"),
	}

	c.FileConfig = FileConfig{
		FilePath: os.Getenv("FILE_PATH"),
	}

	appAccessTokenDuration, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return err
	}
	accessTokenDuration := time.Duration(appAccessTokenDuration) * time.Minute

	appRefreshTokenDuration, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		return err
	}
	refreshTokenDuration := time.Duration(appRefreshTokenDuration) * time.Hour

	c.TokenConfig = TokenConfig{
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
		TokenSymetricKey:     os.Getenv("TOKEN_SYMMETRIC_KEY"),
	}

	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Name == "" ||
		c.DbConfig.User == "" || c.DbConfig.Password == "" || c.DbConfig.Driver == "" ||
		c.ApiConfig.ApiPort == "" || c.FileConfig.FilePath == "" {
		return fmt.Errorf("missing required environment variables")
	}
	return nil
}

// constructor
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.ReadConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
