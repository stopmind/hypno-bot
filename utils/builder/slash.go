package builder

import "github.com/bwmarrin/discordgo"

type SlashCommandBuilder struct {
	command        *discordgo.ApplicationCommand
	serviceBuilder *ServiceBuilder
}

func (s *SlashCommandBuilder) EndCommand() *ServiceBuilder {
	return s.serviceBuilder
}

func buildSlashCommand(serviceBuilder *ServiceBuilder, name string, description string) (*SlashCommandBuilder, *discordgo.ApplicationCommand) {
	command := &discordgo.ApplicationCommand{
		Name:        name,
		Description: description,
		Options:     []*discordgo.ApplicationCommandOption{},
	}
	return &SlashCommandBuilder{
		command:        command,
		serviceBuilder: serviceBuilder,
	}, command
}

func (s *SlashCommandBuilder) addOption(option *discordgo.ApplicationCommandOption) *SlashCommandBuilder {
	s.command.Options = append(s.command.Options, option)
	return s
}

func (s *SlashCommandBuilder) AddStringOption(name, description string, required bool) *SlashCommandBuilder {
	return s.addOption(&discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        name,
		Description: description,
		Required:    required,
	})
}

func (s *SlashCommandBuilder) AddSubCommand(name, description string) *SlashCommandBuilder {
	return s.addOption(&discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        name,
		Description: description,
		Options:     []*discordgo.ApplicationCommandOption{},
	})
}
