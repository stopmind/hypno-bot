package core

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type InteractionCreateHandler func(*discordgo.InteractionCreate)

type InteractionsManager struct {
	bot                *Bot
	handlers           *HandlersManager
	logger             *log.Logger
	commandsHandlers   map[string]InteractionCreateHandler
	componentsHandlers []InteractionCreateHandler
	initialized        bool
}

func (s *InteractionsManager) init() {
	s.initialized = true
	s.handlers.AddHandler(func(i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			handler, ok := s.commandsHandlers[i.Interaction.ApplicationCommandData().Name]

			if !ok {
				s.logger.Printf("unknown slash command %s", i.Interaction.ApplicationCommandData().Name)
				return
			}

			handler(i)
			break
		case discordgo.InteractionMessageComponent:
			for _, handler := range s.componentsHandlers {
				handler(i)
			}
			break
		}
	})
}

func (s *InteractionsManager) AddCommand(command *discordgo.ApplicationCommand, handler InteractionCreateHandler) error {
	if !s.initialized {
		s.init()
	}

	cmd, err := s.bot.ApplicationCommandCreate(s.bot.State.User.ID, "", command)

	if err != nil {
		return err
	}
	s.commandsHandlers[cmd.Name] = handler

	return nil
}

func (s *InteractionsManager) AddComponentHandler(handler InteractionCreateHandler) {
	if !s.initialized {
		s.init()
	}

	s.componentsHandlers = append(s.componentsHandlers, handler)
}

func newSlashCommandsManager(bot *Bot, handlers *HandlersManager, logger *log.Logger) *InteractionsManager {
	return &InteractionsManager{
		bot:              bot,
		handlers:         handlers,
		logger:           logger,
		commandsHandlers: make(map[string]InteractionCreateHandler),
		initialized:      false,
	}
}
