package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml/v2"
	"hypno-bot/core"
	"strings"
	"time"
)

type SpectacleService struct {
	config struct {
		Default   string `toml:"default"`
		End       string `toml:"end"`
		Scenarios map[string][]struct {
			Content string  `toml:"c"`
			Time    float32 `toml:"t"`
		} `toml:"scenarios"`
	}
}

func (s *SpectacleService) Stop() {

}

func (y *SpectacleService) Init(container *core.ServiceContainer) error {

	data, err := container.Storage.ReadFile("config.toml")
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, &y.config)
	if err != nil {
		return err
	}

	container.Bot.AddHandler(func(session *discordgo.Session, send *discordgo.MessageCreate) {
		if !strings.HasPrefix(send.Content, "!spectacle") || send.Author.Bot {
			return
		}
		scenario := y.config.Scenarios[y.config.Default]

		go func() {
			first := scenario[0]

			var message *discordgo.Message
			message, err = session.ChannelMessageSendReply(send.ChannelID, first.Content, send.Reference())

			if err != nil {
				container.Logger.Print(err)
				return
			}

			time.Sleep(time.Duration(first.Time * float32(time.Second)))

			for i := 1; i < len(scenario); i++ {
				block := scenario[i]
				_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, block.Content)
				if err != nil {
					container.Logger.Print(err)
					return
				}

				time.Sleep(time.Duration(block.Time * float32(time.Second)))
			}

			_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, y.config.End)
			if err != nil {
				container.Logger.Print(err)
				return
			}
		}()
	})

	return nil
}
