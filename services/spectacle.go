package services

import (
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/services/builder"
	"time"
)

type Spectacle struct {
	builder.WithContainer
	builder.WithConfig[struct {
		Default   string `toml:"default"`
		End       string `toml:"end"`
		Scenarios map[string][]struct {
			Content string  `toml:"c"`
			Time    float32 `toml:"t"`
		} `toml:"scenarios"`
	}]
}

func (s *Spectacle) command(send *discordgo.MessageCreate) {
	scenario := s.Config.Scenarios[s.Config.Default]

	go func() {
		first := scenario[0]

		message, err := s.Bot.ChannelMessageSendReply(send.ChannelID, first.Content, send.Reference())

		if err != nil {
			s.Logger.Print(err)
			return
		}

		time.Sleep(time.Duration(first.Time * float32(time.Second)))

		for i := 1; i < len(scenario); i++ {
			block := scenario[i]
			_, err = s.Bot.ChannelMessageEdit(message.ChannelID, message.ID, block.Content)
			if err != nil {
				s.Logger.Print(err)
				return
			}

			time.Sleep(time.Duration(block.Time * float32(time.Second)))
		}

		_, err = s.Bot.ChannelMessageEdit(message.ChannelID, message.ID, s.Config.End)
		if err != nil {
			s.Logger.Print(err)
			return
		}
	}()
}

func BuildSpectacleService() core.Service {
	content := new(Spectacle)
	return builder.BuildService(content).
		AddCommand("!spectacle", content.command).
		Finish()
}
