package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigReader interface {
	ReadConfig(configFilePath string, config *Config) error
}

func ReadConfig(configFilePath string, config *Config) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("ユーザー情報を取得できませんでした: %w", err)
	}
	path := strings.Replace(configFilePath, "~", currentUser.HomeDir, 1)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return fmt.Errorf("failed to decode config file as YAML %s: %w", configFilePath, err)
	}
	return nil
}
