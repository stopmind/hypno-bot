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
	addService(bot, logger, "stopmind.fun.spectacle", new(SpectacleService))
	addService(bot, logger, "stopmind.fun.games", new(GamesService))
	//addService(bot, logger, "stopmind.fun.mine", new(MineService))
	addService(bot, logger, "stopmind.serv.rank", new(RankService))
	addService(bot, logger, "stopmind.serv.hello", new(HelloService))
	//addService(bot, logger, "stopmind.serv.conpan", new(ConpanService))
}
