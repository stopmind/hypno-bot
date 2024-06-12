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

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

L:
	for {
		select {
		case <-c:
			break L
		default:
			time.Sleep(1 * time.Second)
			bot.Update()
			break
		}
	}

	bot.Stop()
}
