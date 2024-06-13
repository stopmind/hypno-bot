package services

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml/v2"
	"hypno-bot/core"
	"hypno-bot/services/achievements"
	"hypno-bot/utils"
	"slices"
	"strings"
	"time"
)

type RankService struct {
	*core.ServiceContainer

	config struct {
		Blacklist []string `toml:"blacklist"`
	}
	data struct {
		UsersStat    map[string]int
		UsersCookies map[string]int
		ChannelsStat map[string]int
	}
}

func (r *RankService) saveData() error {
	file, err := r.Storage.OpenFile("data.json")
	if err != nil {
		return err
	}

	data, err := json.Marshal(r.data)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	_ = file.Close()

	return err
}

func (r *RankService) loadData() error {
	data, err := r.Storage.ReadFile("data.json")
	if err == nil {
		if err = json.Unmarshal(data, &r.data); err != nil {
			r.Logger.Print(err)
		}
	} else {
		r.Logger.Print(err)
	}

	if r.data.UsersStat == nil {
		r.data.UsersStat = make(map[string]int)
	}
	if r.data.ChannelsStat == nil {
		r.data.ChannelsStat = make(map[string]int)
	}
	if r.data.UsersCookies == nil {
		r.data.UsersCookies = make(map[string]int)
	}

	return nil
}

func (r *RankService) topSend(topType string, send *discordgo.MessageCreate) error {
	message := "Че ты от меня хочешь?"
	switch topType {
	case "users":
		members := make([]*discordgo.Member, 0, len(r.data.UsersStat))

		for k := range r.data.UsersStat {
			member, err := r.Bot.GuildMember(send.GuildID, k)
			if err != nil {
				continue
			}
			members = append(members, member)
		}

		slices.SortFunc(members, func(a, b *discordgo.Member) int {
			return r.data.UsersStat[b.User.ID] - r.data.UsersStat[a.User.ID]
		})

		template, err := r.Storage.GetTemplate("assets/users.tmp")
		if err == nil {
			message, err = utils.ExecuteTemplate(template, map[string]any{
				"users": members,
				"stat":  r.data.UsersStat,
			})
		}

		if err != nil {
			return err
		}
	case "channels":
		channels := make([]*discordgo.Channel, 0, len(r.data.ChannelsStat))

		for k := range r.data.ChannelsStat {
			channel, err := r.Bot.Channel(k)
			if err != nil {
				continue
			}
			channels = append(channels, channel)
		}
		slices.SortFunc(channels, func(a, b *discordgo.Channel) int {
			return r.data.ChannelsStat[b.ID] - r.data.ChannelsStat[a.ID]
		})

		template, err := r.Storage.GetTemplate("assets/channels.tmp")
		if err == nil {
			message, err = utils.ExecuteTemplate(template, map[string]any{
				"channels": channels,
				"stat":     r.data.ChannelsStat,
			})
		}

		if err != nil {
			return err
		}

		break
	case "cookies":
		users := make([]*discordgo.Member, 0, len(r.data.UsersStat))

		for k := range r.data.UsersCookies {
			member, err := r.Bot.GuildMember(send.GuildID, k)
			if err != nil {
				continue
			}
			users = append(users, member)
		}

		slices.SortFunc(users, func(a, b *discordgo.Member) int {
			return r.data.UsersCookies[b.User.ID] - r.data.UsersCookies[a.User.ID]
		})

		template, err := r.Storage.GetTemplate("assets/cookies.tmp")
		if err == nil {
			message, err = utils.ExecuteTemplate(template, map[string]any{
				"users": users,
				"stat":  r.data.UsersCookies,
			})
		}

		if err != nil {
			return err
		}

		break
	}

	_, err := r.Bot.ChannelMessageSendReply(send.ChannelID, message, send.Reference())

	return err
}

func (r *RankService) Init(container *core.ServiceContainer) error {
	r.ServiceContainer = container

	data, err := r.Storage.ReadFile("config.toml")
	if err != nil {
		return err
	}

	if err = toml.Unmarshal(data, &r.config); err != nil {
		return err
	}

	if err = r.loadData(); err != nil {
		return err
	}

	go func() {
		for {
			year, month, day := time.Now().Date()
			nextTime := time.Date(year, month, day+1, 0, 0, 0, 0, time.Local)
			time.Sleep(time.Until(nextTime))

			if err := r.saveData(); err != nil {
				r.Logger.Print(err)
			}
		}
	}()

	const cookieEmoji = "🍪"

	r.Bot.AddHandler(func(session *discordgo.Session, send *discordgo.MessageCreate) {
		if send.Author.Bot {
			return
		}

		if !slices.Contains(r.config.Blacklist, send.ChannelID) {
			r.data.UsersStat[send.Author.ID] += 1
			r.data.ChannelsStat[send.ChannelID] += 1
		}

		if !strings.HasPrefix(send.Content, "!top") {
			return
		}

		topType := "users"

		parts := strings.SplitN(send.Content, " ", 2)
		if len(parts) == 2 {
			topType = parts[1]
		}

		if err = r.topSend(topType, send); err != nil {
			r.Logger.Print(err)
		}
	})

	r.Bot.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
		if event.Member.User.Bot || event.Emoji.Name != cookieEmoji {
			return
		}

		message, err := r.Bot.ChannelMessage(event.ChannelID, event.MessageID)
		if err != nil {
			return
		}

		if message.Author.ID == event.UserID {
			return
		}

		if _, ok := r.data.UsersCookies[event.UserID]; !ok {
			r.data.UsersCookies[message.Author.ID] = 0
		}
		r.data.UsersCookies[message.Author.ID] += 1
		if r.data.UsersCookies[message.Author.ID] > 100 {
			achievements.CheckCookies(message.Author.ID)
		}
	})

	r.Bot.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionRemove) {
		if event.Emoji.Name != cookieEmoji {
			return
		}

		user, err := r.Bot.User(event.UserID)
		if user.Bot || err != nil {
			return
		}

		message, err := r.Bot.ChannelMessage(event.ChannelID, event.MessageID)
		if err != nil {
			return
		}

		if message.Author.ID == event.UserID {
			return
		}

		if _, ok := r.data.UsersCookies[event.UserID]; !ok {
			r.data.UsersCookies[message.Author.ID] = 0
		}
		r.data.UsersCookies[message.Author.ID] -= 1
	})

	return nil
}

func (r *RankService) Stop() {
	if err := r.saveData(); err != nil {
		r.Logger.Print(err)
	}
}
