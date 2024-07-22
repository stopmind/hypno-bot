package builder

import (
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"strings"
)

type slashCommandInfo struct {
	handler core.InteractionCreateHandler
	command *discordgo.ApplicationCommand
}

type builtService struct {
	container          *core.ServiceContainer
	content            any
	handlers           []any
	slashCommands      []slashCommandInfo
	componentsHandlers []core.InteractionCreateHandler
}

func (b *builtService) Init(container *core.ServiceContainer) error {
	b.container = container

	if contain, ok := b.content.(containerInit); ok {
		contain.initContainer(container)
	}

	if config, ok := b.content.(configInit); ok {
		if err := config.initConfig(container.Storage); err != nil {
			return err
		}
	}
	if state, ok := b.content.(stateInit); ok {
		if err := state.initState(container.Storage); err != nil {
			container.Logger.Print(err)
		}
	}

	for _, handler := range b.handlers {
		container.Handlers.AddHandler(handler)
	}

	for _, slashCommand := range b.slashCommands {
		err := b.container.Interactions.AddCommand(slashCommand.command, slashCommand.handler)
		if err != nil {
			return err
		}
	}

	for _, handler := range b.componentsHandlers {
		b.container.Interactions.AddComponentHandler(handler)
	}

	return nil
}

func (b *builtService) Stop() {
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
		content:       content,
		handlers:      make([]any, 0),
		slashCommands: make([]slashCommandInfo, 0),
	}}
}

func (b *ServiceBuilder) AddHandler(handler any) *ServiceBuilder {
	b.service.handlers = append(b.service.handlers, handler)
	return b
}

type CommandAction func(send *discordgo.MessageCreate)

func (b *ServiceBuilder) AddCommand(name string, action CommandAction) *ServiceBuilder {
	return b.AddHandler(func(send *discordgo.MessageCreate) {
		parts := strings.SplitN(send.Content, " ", 2)
		if !send.Author.Bot && parts[0] == name {
			action(send)
		}
	})
}

func (b *ServiceBuilder) AddSlashCommand(name string, description string, handler core.InteractionCreateHandler) *SlashCommandBuilder {
	builder, command := buildSlashCommand(b, name, description)

	b.service.slashCommands = append(b.service.slashCommands, slashCommandInfo{
		handler: handler,
		command: command,
	})

	return builder
}

func (b *ServiceBuilder) AddComponentHandler(handler core.InteractionCreateHandler) *ServiceBuilder {
	b.service.componentsHandlers = append(b.service.componentsHandlers, handler)
	return b
}

func (b *ServiceBuilder) Finish() core.Service {
	return b.service
}
