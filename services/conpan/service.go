package conpan

import (
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils"
	"hypno-bot/utils/builder"
	"slices"
	"strings"
)

type content struct {
	builder.WithContainer
	builder.WithState[struct {
		Codes []string `json:"codes"`
	}]
	builder.WithConfig[struct {
		Admins []string `toml:"admins"`
	}]
}

func (c *content) command(send *discordgo.MessageCreate) {
	args, err := utils.NewArgsParser().
		AddString().
		Parse(1, send.Content)

	if err != nil {
		utils.ReplyIncorrectArgsError(send, err)
	}

	session, _ := c.Sessions.GetSession(send.Author.ID)
	if session == nil {
		session, _ = c.Sessions.NewSession(send.Author.ID, &shellSession{
			isAdmin: c.Config.Admins == nil || len(c.Config.Admins) == 0 || slices.Contains(c.Config.Admins, send.Author.ID),
		})
	}

	session.Extend()

	shell := session.Data.(*shellSession)

	c.execute(shell, send, strings.Split(args.Get(0).(string), " "))
}

func BuildService() core.Service {
	c := new(content)
	return builder.BuildService(c).
		AddCommand("!conpan", c.command).
		Finish()
}
