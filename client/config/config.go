package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Hosts []struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"hosts"`
	AppToken string `yaml:"app_token"`
	BotToken string `yaml:"bot_token"`
}

var Cfg *Config

func LoadConfig(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &Cfg); err != nil {
		return err
	}
	return nil
}
