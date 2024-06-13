package main

import (
	"hypno-bot/core"
	i "hypno-bot/init"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	go func() {
		for {
			time.Sleep(time.Millisecond * 200)
			bot.Update()
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	bot.Stop()
	os.Exit(0)
}
