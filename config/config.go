package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Host            string         `json:"host"`
	Port            int            `json:"port"`
	TokenSecretKey  string         `json:"token_secret_key"`
	AccessAddresses []string       `json:"access_addresses"`
	Email           string         `json:"email"`
	EmailPassword   string         `json:"email_password"`
	HealthCheck     HealthCheckCfg `json:"health_check"`
}

type HealthCheckCfg struct {
	ConnCount   int `json:"conn_count"`
	ConnTimeout int `json:"conn_timeout"`
}

func NewConfig() (*Config, error) {
	buff, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.ErrNotExist == err {
			return &Config{
				Host:            "127.0.0.1",
				Port:            8080,
				TokenSecretKey:  "Secret",
				AccessAddresses: []string{"*"},
			}, nil
		}
		return nil, err
	}

	var cfg = &Config{}
	err = json.Unmarshal(buff, cfg)
	if err != nil {
		return nil, err
	}

	cfg.checkDefault()

	return cfg, nil
}

func (cfg *Config) checkDefault() {
	if cfg.TokenSecretKey == "" {
		cfg.TokenSecretKey = "Secret"
	}

	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	if cfg.AccessAddresses == nil || len(cfg.AccessAddresses) == 0 {
		cfg.AccessAddresses = []string{"*"}
	}

	if cfg.HealthCheck.ConnCount == 0 {
		cfg.HealthCheck.ConnCount = 10
	}

	if cfg.HealthCheck.ConnTimeout == 0 {
		cfg.HealthCheck.ConnTimeout = 30
	}
}
