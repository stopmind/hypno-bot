package core

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
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

	err = toml.Unmarshal(source, &config)
	if err != nil {
		return Config{}, fmt.Errorf("config: %s", err.Error())
	}

	return config, nil
}
