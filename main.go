package main

import (
	"hypno-bot/core"
	i "hypno-bot/init"
	"os"
	"os/signal"
	"syscall"
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

	i.Init(bot, logger)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	bot.Stop()
}
