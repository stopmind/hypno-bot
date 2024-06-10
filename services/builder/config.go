package builder

import "hypno-bot/core"

type WithConfig[TConfig any] struct {
	Config TConfig
}

func (c *WithConfig[TConfig]) initConfig(storage *core.Storage) error {
	return storage.ReadTOML("config.toml", &c.Config)

}

type configInit interface {
	initConfig(storage *core.Storage) error
}
