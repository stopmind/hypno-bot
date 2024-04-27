package main

import (
	"hypno-bot/core"
	"hypno-bot/services"
	"os"
	"os/signal"
)

func main() {
	logger, err := core.NewLogger("core")
	if err != nil {
		panic(err)
	}

	bot := core.NewBot()
	err = bot.Start()
	if err != nil {
		logger.Panic(err)
	}

	services.Init(bot, logger)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	bot.Stop()
}
