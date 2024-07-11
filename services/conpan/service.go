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
		Codes          []string `json:"codes"`
		UpdatesBlocked bool     `json:"updates_blocked"`
	}]
	builder.WithConfig[struct {
		Admins             []string `toml:"admins"`
		UpdatesControllers []string `toml:"updates_controllers"`
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

	checkIsController := func(send *discordgo.MessageCreate) bool {
		if !slices.Contains(c.Config.UpdatesControllers, send.Author.ID) {
			err := c.Reply(send, "You are not a updates controller")
			if err != nil {
				c.Logger.Print(err)
			}
		}

		return slices.Contains(c.Config.UpdatesControllers, send.Author.ID)
	}

	return builder.BuildService(c).
		AddCommand("!conpan", c.command).
		AddCommand("!block-updates", func(send *discordgo.MessageCreate) {
			if checkIsController(send) {
				c.State.UpdatesBlocked = true
				err := c.Reply(send, "Updates blocked")
				if err != nil {
					c.Logger.Print(err)
				}
			}
		}).
		AddCommand("!accept-updates", func(send *discordgo.MessageCreate) {
			if checkIsController(send) {
				c.State.UpdatesBlocked = false
				err := c.Reply(send, "Updates accepted")
				if err != nil {
					c.Logger.Print(err)
				}
			}
		}).
		Finish()
}
