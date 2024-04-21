package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"path/filepath"
)

type Bot struct {
	*discordgo.Session

	services map[string]*ServiceContainer
}

func (b *Bot) Init() error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	b.Session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	return nil
}

func (b *Bot) Run() error {
	defer b.Close()

	err := b.Open()
	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	for {

	}

	return nil
}

func (b *Bot) AddService(name string, service Service) error {
	_, ok := b.services[name]

	if ok {
		return fmt.Errorf("bot: service with name \"%s\" already exist", name)
	}

	container := new(ServiceContainer)

	container.Name = name
	container.Service = service
	container.Bot = b
	container.Sessions = newSessionsManager(container)
	container.Storage = newStorage(filepath.Join("storage", name))
	var err error
	container.Logger, err = NewLogger(name)

	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	service.Init(container)

	return nil
}
