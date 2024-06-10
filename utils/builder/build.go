package builder

import (
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"strings"
)

type builtService struct {
	container *core.ServiceContainer
	content   any
	handlers  []any
}

func (b builtService) Init(container *core.ServiceContainer) error {
	b.container = container

	if contain, ok := b.content.(containerInit); ok {
		contain.initContainer(b.container)
	}

	if config, ok := b.content.(configInit); ok {
		if err := config.initConfig(container.Storage); err != nil {
			return err
		}
	}
	if state, ok := b.content.(stateInit); ok {
		if err := state.initState(container.Storage); err != nil {
			return err
		}
	}

	for _, handler := range b.handlers {
		container.Bot.AddHandler(handler)
	}

	return nil
}

func (b builtService) Stop() {
	if state, ok := b.content.(stateInit); ok {
		if err := state.saveState(b.container.Storage); err != nil {
			b.container.Logger.Print(err)
		}
	}
}

type ServiceBuilder struct {
	service *builtService
}

func BuildService(content any) *ServiceBuilder {
	return &ServiceBuilder{&builtService{
		content:  content,
		handlers: make([]any, 0),
	}}
}

func (b *ServiceBuilder) AddHandler(handler any) *ServiceBuilder {
	b.service.handlers = append(b.service.handlers, handler)

	return b
}

type CommandAction func(send *discordgo.MessageCreate)

func (b *ServiceBuilder) AddCommand(name string, action CommandAction) *ServiceBuilder {
	return b.AddHandler(func(session *discordgo.Session, send *discordgo.MessageCreate) {
		if !send.Author.Bot && strings.HasPrefix(send.Content, name) {
			action(send)
		}
	})
}

func (b *ServiceBuilder) Finish() core.Service {
	return b.service
}
