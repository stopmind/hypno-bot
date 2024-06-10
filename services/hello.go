package services

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml/v2"
	"hypno-bot/core"
)

type HelloService struct {
	*core.ServiceContainer

	config struct {
		HiMsg      string `toml:"hi_message"`
		ByeMsg     string `toml:"bye_message"`
		ChannelMsg string `toml:"channel"`
	}
}

func (h *HelloService) Init(container *core.ServiceContainer) error {
	h.ServiceContainer = container

	data, err := h.Storage.ReadFile("WithConfig.toml")
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, &h.config)
	if err != nil {
		return err
	}

	h.Bot.AddHandler(func(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
		_, err = session.ChannelMessageSend(h.config.ChannelMsg, fmt.Sprintf(h.config.HiMsg, event.User.Username))

		if err != nil {
			h.Logger.Print(err)
		}
	})

	h.Bot.AddHandler(func(session *discordgo.Session, event *discordgo.GuildMemberRemove) {
		_, err = session.ChannelMessageSend(h.config.ChannelMsg, fmt.Sprintf(h.config.ByeMsg, event.User.Username))

		if err != nil {
			h.Logger.Print(err)
		}
	})

	return nil
}

func (h *HelloService) Stop() {

}
