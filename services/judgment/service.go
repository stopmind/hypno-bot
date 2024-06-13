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
	case "справка":
		c.help(send)
		break
	default:
		utils.ReplyError(send, "Я чет не понял", fmt.Sprintf("Не получилось найти колманду `%s`.", args.Get(0)))
		break
	}

}

const helpMessage = "# Справка\n\n" +
	"## Суд\n" +
	"`?суд начать <сторона1> <сторона2>` - начать суд\n" +
	"`?суд стоп` - закончить суд досрочно\n" +
	"`?суд поинт <сторона>` - добавить очко на одну из сторон, если сторона получит 3 очка, она выиграет\n" +
	"## Отзывы\n" +
	"`?суд отзыв <пользователь> <очки> <комментарий>` - оставить отзыв на человека\n" +
	"`?суд карма <пользователь>` - посмотреть карму пользователя\n"

func (c *content) help(send *discordgo.MessageCreate) {
	err := c.Reply(send, helpMessage)
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
	}
}

func BuildService() core.Service {
	c := new(content)
	return builder.BuildService(c).
		AddCommand("?суд", c.command).
		Finish()
}
