package mine

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type userData struct {
	Balance int `json:"balance"`
}

func (c *content) getUserData(id string) *userData {
	if c.State.Users == nil {
		c.State.Users = make(map[string]*userData)
	}

	user, ok := c.State.Users[id]
	if !ok {
		user = &userData{}
		c.State.Users[id] = user
	}

	return user
}

func (c *content) balance(interaction *discordgo.Interaction) error {
	return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Баланс: %d", c.getUserData(interaction.Member.User.ID).Balance),
		},
	})
}
