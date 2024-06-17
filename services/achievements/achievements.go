package achievements

import (
	"github.com/bwmarrin/discordgo"
	"unicode"
)

type progress struct {
	AristocratCount int            `json:"aristocrat_count"`
	Counters        map[string]int `json:"counters"`
}

func (c *content) aristocratCheck(_ *discordgo.Session, send *discordgo.MessageCreate) {
	info := c.getUserInfo(send.Author.ID)

	if info.Progress.AristocratCount >= 10 {
		c.giveAchievement(send.Author.ID, "aristocrat")
		return
	}

	if !unicode.IsUpper(rune(send.Content[0])) || send.Content[len(send.Content)-1] != '.' {
		return
	}

	info.Progress.AristocratCount += 1

	if info.Progress.AristocratCount == 10 {
		c.giveAchievement(send.Author.ID, "aristocrat")
	}
}

func CheckCookies(userId string) {
	c.giveAchievement(userId, "cookies")
}

func counterCheck(achievement string, requiredCount int) func(userId string) {
	return func(userId string) {
		info := c.getUserInfo(userId)

		if info.Progress.Counters[achievement] >= requiredCount {
			c.giveAchievement(userId, achievement)
			return
		}

		info.Progress.Counters[achievement] += 1

		if info.Progress.Counters[achievement] == requiredCount {
			c.giveAchievement(userId, achievement)
		}
	}
}

var (
	OnR34    = counterCheck("r34fan", 40)
	OnR34To  = counterCheck("r34friend", 10)
	OnBlock  = counterCheck("writer", 10)
	OnReview = counterCheck("snitch", 10)
)
