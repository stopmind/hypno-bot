package init

import (
	"hypno-bot/core"
	"hypno-bot/services"
	"hypno-bot/services/achievements"
	"hypno-bot/services/conpan"
	"hypno-bot/services/judgment"
	"hypno-bot/services/news"
	"hypno-bot/utils"
	"log"
)

func addService(bot *core.Bot, logger *log.Logger, name string, service core.Service) {
	err := bot.AddService(name, service)
	if err != nil {
		logger.Println(err)
	}
}

func addCriticalService(bot *core.Bot, logger *log.Logger, name string, service core.Service) {
	err := bot.AddService(name, service)
	if err != nil {
		_ = bot.Close()
		logger.Panic(err)
	}
}

func Init(bot *core.Bot, logger *log.Logger) {
	addCriticalService(bot, logger, "serv.utils", utils.InitUtilsService())
	addCriticalService(bot, logger, "serv.conpan", conpan.BuildService())

	addService(bot, logger, "fun.spectacle", services.BuildSpectacleService())
	addService(bot, logger, "fun.games", new(services.GamesService))
	addService(bot, logger, "serv.rank", new(services.RankService))
	//addService(bot, logger, "serv.hello", new(HelloService))
	addService(bot, logger, "fun.r34", services.BuildR34Service())
	addService(bot, logger, "fun.news", new(news.Service))
	addService(bot, logger, "fun.achievements", achievements.BuildService())
	addService(bot, logger, "fun.judgment", judgment.BuildService())
}
