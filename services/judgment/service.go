package judgment

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils"
	"hypno-bot/utils/builder"
)

type content struct {
	builder.WithContainer
	builder.WithConfig[struct {
		Judge string `toml:"judge"`
	}]
	builder.WithState[struct {
		Users map[string]*userInfo `json:"users"`
	}]
}

func (c *content) command(send *discordgo.MessageCreate) {
	args, err := utils.NewArgsParser().
		AddString().
		AddString().
		Parse(1, send.Content)

	if err != nil {
		utils.ReplyIncorrectArgsError(send, err)
		return
	}

	switch args.Get(0) {
	case "начать":
		c.begin(send)
		break
	case "стоп":
		c.stop(send)
		break
	case "поинт":
		c.point(send)
		break
	case "карма":
		c.karma(send)
		break
	case "отзыв":
		c.review(send)
		break
	default:
		utils.ReplyError(send, "Я чет не понял", fmt.Sprintf("Не получилось найти колманду `%s`.", args.Get(0)))
		break
	}

}

func BuildService() core.Service {
	c := new(content)
	return builder.BuildService(c).
		AddCommand("?суд", c.command).
		Finish()
}
