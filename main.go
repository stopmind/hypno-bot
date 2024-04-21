package main

import (
	"hypno-bot/core"
	"hypno-bot/services"
)

func main() {
	logger, err := core.NewLogger("core")
	if err != nil {
		panic(err)
	}

	bot := core.Bot{}
	err = bot.Init()
	if err != nil {
		logger.Panic(err)
	}

	services.Init(&bot, logger)

	err = bot.Run()
	if err != nil {
		logger.Panic(err)
	}
}
