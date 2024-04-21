package services

import (
	"hypno-bot/core"
	"log"
)

func addService(bot *core.Bot, logger *log.Logger, name string, service core.Service) {
	err := bot.AddService(name, service)
	if err != nil {
		logger.Println(err)
	}
}

func Init(bot *core.Bot, logger *log.Logger) {
}
