package core

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type SlashCommandHandler func(*discordgo.InteractionCreate)

type SlashCommandsManager struct {
	bot              *Bot
	handlers         *HandlersManager
	logger           *log.Logger
	commandsHandlers map[string]SlashCommandHandler
	initialized      bool
}

func (s *SlashCommandsManager) init() {
	s.initialized = true
	s.handlers.AddHandler(func(i *discordgo.InteractionCreate) {
		handler, ok := s.commandsHandlers[i.Interaction.ApplicationCommandData().Name]

		if !ok {
			s.logger.Printf("unknown slash command %s", i.Interaction.ApplicationCommandData().Name)
			return
		}

		handler(i)
	})
}

func (s *SlashCommandsManager) AddCommand(command *discordgo.ApplicationCommand, handler SlashCommandHandler) error {
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

func newSlashCommandsManager(bot *Bot, handlers *HandlersManager, logger *log.Logger) *SlashCommandsManager {
	return &SlashCommandsManager{
		bot:              bot,
		handlers:         handlers,
		logger:           logger,
		commandsHandlers: make(map[string]SlashCommandHandler),
		initialized:      false,
	}
}
