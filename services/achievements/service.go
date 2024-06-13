package achievements

import (
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils"
	"hypno-bot/utils/builder"
	"math/rand/v2"
	"slices"
	"time"
)

type userInfo struct {
	Achievements []string `json:"achievements"`
	Progress     progress `json:"progress"`
}

type content struct {
	builder.WithContainer
	builder.WithConfig[struct {
		InfoChannel  string            `toml:"info_channel"`
		Achievements map[string]string `toml:"achievements"`
	}]
	builder.WithState[struct {
		Users map[string]*userInfo `json:"users"`
	}]
}

var c *content

func (c *content) getUserInfo(id string) *userInfo {
	user, ok := c.State.Users[id]
	if ok {
		return user
	}

	if c.State.Users == nil {
		c.State.Users = make(map[string]*userInfo)
	}

	user = &userInfo{
		Achievements: make([]string, 0),
		Progress: progress{
			AristocratCount: 0,
			Counters:        make(map[string]int),
		},
	}
	c.State.Users[id] = user
	return user
}

func (c *content) giveAchievement(userID string, achievement string) {
	user := c.getUserInfo(userID)
	if slices.Contains(user.Achievements, achievement) {
		return
	}

	user.Achievements = append(user.Achievements, achievement)

	go func() {
		time.Sleep(time.Duration(6+rand.IntN(5)) * time.Second)

		notify, err := c.Storage.GetTemplate("assets/notify.tmp")
		if err != nil {
			c.Logger.Print(err)
			return
		}

		text, err := utils.ExecuteTemplate(notify, map[string]any{
			"UserID": userID,
			"Name":   c.Config.Achievements[achievement],
		})
		if err != nil {
			c.Logger.Print(err)
			return
		}

		_, err = c.Bot.ChannelMessageSend(c.Config.InfoChannel, text)
		if err != nil {
			c.Logger.Print(err)
		}
	}()
}

func (c *content) profile(send *discordgo.MessageCreate) {
	template, err := c.Storage.GetTemplate("assets/profile.tmp")
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
		return
	}

	result, err := utils.ExecuteTemplate(template, map[string]any{
		"User":         c.getUserInfo(send.Author.ID),
		"Achievements": c.Config.Achievements,
	})
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
		return
	}

	err = c.Reply(send, result)
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
	}
}

func BuildService() core.Service {
	c = new(content)
	return builder.BuildService(c).
		AddCommand("?ачивки", c.profile).
		AddCommand("test", func(send *discordgo.MessageCreate) {
			c.giveAchievement(send.Author.ID, "test")
		}).
		AddHandler(c.aristocratCheck).
		Finish()
}
