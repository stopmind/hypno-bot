package judgment

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils"
)

func (c *content) checkJudge(send *discordgo.MessageCreate) bool {
	if send.Author.ID != c.Config.Judge {
		err := c.Reply(send, "Вы не судья!")

		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}

		return true
	}

	return false
}

type side struct {
	Name   string
	Points int
}

type process struct {
	Side1 side
	Side2 side
}

func (c *content) getProcess() *process {
	session, _ := c.Sessions.GetSession(c.Config.Judge)
	if session == nil {
		return nil
	}

	return session.Data.(*process)
}

func (c *content) begin(send *discordgo.MessageCreate) {
	if c.checkJudge(send) {
		return
	}

	args, err := utils.NewArgsParser().
		AddString().
		AddString().
		AddString().
		Parse(2, send.Content)

	if err != nil {
		utils.ReplyIncorrectArgsError(send, err)
		return
	}

	proc := c.getProcess()

	if proc != nil {
		err = c.Reply(send, "Процесс уже идет.")

		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}

		return
	}

	proc = &process{
		Side1: side{
			Name:   args.Get(1).(string),
			Points: 0,
		},
		Side2: side{
			Name:   args.Get(2).(string),
			Points: 0,
		},
	}

	_, err = c.Sessions.NewSession(c.Config.Judge, proc)
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
		return
	}

	err = c.Reply(send, "Встаньте! Суд идет!")

	if err != nil {
		utils.ReplyUnexpectedError(send, err)
	}
}

func (c *content) stop(send *discordgo.MessageCreate) {
	if c.checkJudge(send) {
		return
	}

	var err error
	proc := c.getProcess()
	if proc != nil {
		var session *core.Session
		session, err = c.Sessions.GetSession(c.Config.Judge)

		if session != nil {
			session.Close()
			err = c.Reply(send, "Судебный процесс завершен досрочно.")
		}
	} else {
		err = c.Reply(send, "Никакого судебного процесса не идет.")
	}

	if err != nil {
		utils.ReplyUnexpectedError(send, err)
	}
}

func (c *content) point(send *discordgo.MessageCreate) {
	if c.checkJudge(send) {
		return
	}

	args, err := utils.NewArgsParser().
		AddString().
		AddString().
		Parse(2, send.Content)

	if err != nil {
		utils.ReplyIncorrectArgsError(send, err)
		return
	}

	proc := c.getProcess()

	if proc == nil {
		err = c.Reply(send, "Никакого судебного процесса не идет.")
		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}
		return
	}

	name := args.Get(1).(string)

	if proc.Side1.Name == name {
		proc.Side1.Points += 1
	} else if proc.Side2.Name == name {
		proc.Side2.Points += 1
	} else {
		err = c.Reply(send, "Не удалось найти данную сторону.")
		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}

		return
	}

	end := false

	if proc.Side1.Points == 3 {
		err = c.Reply(send, fmt.Sprintf("**%v** проиграл суд.", proc.Side1.Name))
		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}
		end = true
	}

	if proc.Side2.Points == 3 {
		err = c.Reply(send, fmt.Sprintf("**%v** проиграл суд.", proc.Side2.Name))
		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}
		end = true
	}

	if end {
		session, _ := c.Sessions.GetSession(c.Config.Judge)
		session.Close()
		return
	}

	err = c.Reply(send, "Очко было учтено.")
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
	}
}
