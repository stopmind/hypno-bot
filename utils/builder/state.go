package builder

import (
	"hypno-bot/core"
)

type WithState[TState any] struct {
	State TState
}

func (c *WithState[TState]) initState(storage *core.Storage) error {
	err := storage.ReadJson("state.json", &c.State)

	if err != nil {
		c.State = *new(TState)
	}

	return err
}

func (c *WithState[TState]) saveState(storage *core.Storage) error {
	return storage.WriteJson("state.json", &c.State)
}

type stateInit interface {
	initState(storage *core.Storage) error
	saveState(storage *core.Storage) error
}
