package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils/builder"
)

type service struct {
	builder.WithContainer
}

var instance *service

func InitUtilsService() core.Service {
	instance = new(service)
	return builder.BuildService(instance).Finish()
}

func ReplyError(send *discordgo.MessageCreate, title string, content string) {
	temp, err := instance.Storage.GetTemplate("assets/error.jet")

	var text string
	if err == nil {
		text, err = ExecuteTemplate(temp, &struct {
			Title, Content string
		}{title, content})
	}

	if err != nil {
		instance.Logger.Print(err)
		if err := instance.Reply(send, "Absolutely failed..."); err != nil {
			instance.Logger.Print(err)
		}

		return
	}

	err = instance.Reply(send, text)
	if err != nil {
		instance.Logger.Print(err)
	}
}

func ReplyFile(send *discordgo.MessageCreate, path string) {
	content, err := instance.Storage.ReadFile(path)

	if err != nil {
		ReplyUnexpectedError(send, err)
		return
	}

	err = instance.Reply(send, string(content))
	if err != nil {
		instance.Logger.Print(err)
	}
}

func ReplyIncorrectArgsError(send *discordgo.MessageCreate, err error) {
	ReplyError(send, ":stop_sign: Некоректные аргументы", fmt.Sprintf("Подробнее:\n`%v`", err.Error()))
}

func ReplyUnexpectedError(send *discordgo.MessageCreate, err error) {
	ReplyError(send, ":stop_sign: Неожиданная ошибка", err.Error())
}
