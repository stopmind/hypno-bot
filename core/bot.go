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

func (b *Bot) Start() error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	b.Session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	b.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildPresences

	err = b.Open()
	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
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

	container.Handlers = newHandlersManager(container.Logger, b)

	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	if err = service.Init(container); err != nil {
		container.Logger.Print(err)
		return fmt.Errorf("bot: service init: %s", err.Error())
	}

	b.services[name] = container

	return nil
}

func (b *Bot) Update() {
	for _, service := range b.services {
		service.Sessions.Update()
	}
}

func (b *Bot) Stop() {
	for _, c := range b.services {
		c.Stop()
	}
	_ = b.Close()
}

func NewBot() *Bot {
	result := new(Bot)

	result.services = make(map[string]*ServiceContainer)

	return result
}
