package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"path/filepath"
)

type Bot struct {
	*discordgo.Session

	Services map[string]*ServiceContainer
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
	_, ok := b.Services[name]

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
	container.Slash = newSlashCommandsManager(b, container.Handlers, container.Logger)

	if err != nil {
		return fmt.Errorf("bot: %s", err.Error())
	}

	if err = service.Init(container); err != nil {
		container.Logger.Print(err)
		return fmt.Errorf("bot: service init: %s", err.Error())
	}

	b.Services[name] = container

	return nil
}

func (b *Bot) Update() {
	for _, service := range b.Services {
		service.Sessions.Update()
	}
}

func (b *Bot) Stop() {
	for _, c := range b.Services {
		c.Stop()
	}

	commands, err := b.ApplicationCommands(b.State.User.ID, "")

	if err != nil {
		print(err.Error())
	} else {
		for _, c := range commands {
			if err = b.ApplicationCommandDelete(c.ApplicationID, c.GuildID, c.ID); err != nil {
				print(err.Error())
			}
		}
	}

	if err != nil {
		print(err.Error())
	}

	_ = b.Close()
}

func NewBot() *Bot {
	result := new(Bot)

	result.Services = make(map[string]*ServiceContainer)

	return result
}
