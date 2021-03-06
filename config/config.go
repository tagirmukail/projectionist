package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Host            string         `json:"host"`
	Port            int            `json:"port"`
	GrpcPort        int            `json:"grpc_port"`
	GrpcApiPort     int            `json:"grpc_api_port"`
	TokenSecretKey  string         `json:"token_secret_key"`
	AccessAddresses []string       `json:"access_addresses"`
	Email           string         `json:"email"`
	EmailPassword   string         `json:"email_password"`
	HealthCheck     HealthCheckCfg `json:"health_check"`
	NotifierConfig  NotifierConfig `json:"notifier"`
}

type HealthCheckCfg struct {
	ConnCount   int `json:"conn_count"`
	ConnTimeout int `json:"conn_timeout"`
}

type NotifierConfig struct {
	Enable   bool   `json:"enable"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	From     string `json:"from"`
	Password string `json:"password"`
}

func NewConfig() (*Config, error) {
	buff, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.ErrNotExist == err {
			return &Config{
				Host:            "127.0.0.1",
				Port:            8080,
				GrpcPort:        8081,
				GrpcApiPort:     8082,
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

	if cfg.GrpcPort == 0 {
		cfg.GrpcPort = 8081
	}

	if cfg.GrpcApiPort == 0 {
		cfg.GrpcApiPort = 8082
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
