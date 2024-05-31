package news

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/utils"
	"slices"
)

func (s *Service) replyError(send *discordgo.MessageCreate, title string, content string) {
	temp, err := s.Storage.GetTemplate("assets/error.jet")

	var text string
	if err == nil {
		text, err = utils.ExecuteTemplate(temp, &struct {
			Title, Content string
		}{title, content})
	}

	if err != nil {
		s.Logger.Print(err)
		if err := s.Reply(send, "Absolutely failed..."); err != nil {
			s.Logger.Print(err)
		}

		return
	}

	err = s.Reply(send, text)
	if err != nil {
		s.Logger.Print(err)
	}
}

func (s *Service) replyFile(send *discordgo.MessageCreate, path string) {
	content, err := s.Storage.ReadFile(path)

	if err != nil {
		s.replyUnexpectedError(send, err)
		return
	}

	err = s.Reply(send, string(content))
	if err != nil {
		s.Logger.Print(err)
	}
}

func (s *Service) replyIncorrectArgsError(send *discordgo.MessageCreate, err error) {
	s.replyError(send, ":stop_sign: Некоректные аргументы", fmt.Sprintf("Подробнее:\n`%v`", err.Error()))
}

func (s *Service) replyUnexpectedError(send *discordgo.MessageCreate, err error) {
	s.replyError(send, ":stop_sign: Неожиданная ошибка", err.Error())
}

func (s *Service) checkRole(send *discordgo.MessageCreate) bool {
	if !slices.Contains(send.Member.Roles, s.config.EditorsRole) {
		s.replyError(send, ":pinching_hand: Недостаточно прав", fmt.Sprintf("У тебя нет роли, челедь."))
		return false
	}

	return true
}
