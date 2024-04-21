package core

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Token string
}

func LoadConfig() (Config, error) {
	config := Config{}

	source, err := os.ReadFile("./config.toml")
	if err != nil {
		return Config{}, fmt.Errorf("config: %s", err.Error())
	}

	_, err = toml.Decode(string(source), &config)
	if err != nil {
		return Config{}, fmt.Errorf("config: %s", err.Error())
	}

	return config, nil
}
