package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Email         string `json:"email"`
	EmailPassword string `json:"email_password"`
}

func NewConfig() (*Config, error) {
	buff, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.ErrNotExist == err {
			return &Config{
				Host: "127.0.0.1",
				Port: 8080,
			}, nil
		}
		return nil, err
	}

	var cfg = &Config{}
	err = json.Unmarshal(buff, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
