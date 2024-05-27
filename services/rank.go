package services

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml/v2"
	"hypno-bot/core"
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

	if err == nil {
		return nil
	}

	r.data.UsersStat = make(map[string]int)
	r.data.ChannelsStat = make(map[string]int)

	return nil
}

func (r *RankService) topSend(topType string, send *discordgo.MessageCreate) error {
	message := ""
	switch topType {
	case "users":
		message += "Ах, как я люблю задротов этого сервера, вот они, слево направо:\n\n"

		users := make([]string, 0, len(r.data.UsersStat))

		for k := range r.data.UsersStat {
			users = append(users, k)
		}

		slices.SortFunc(users, func(a, b string) int {
			return r.data.UsersStat[b] - r.data.UsersStat[a]
		})

		for i, user := range users {
			member, err := r.Bot.GuildMember(send.GuildID, user)
			if err != nil {
				return err
			}

			name := member.Nick

			if name == "" {
				name = member.User.GlobalName
			}
			if name == "" {
				name = member.User.Username
			}

			message += fmt.Sprintf("**%v. %s** *(%v сообщений)*\n", i+1, name, r.data.UsersStat[user])
		}
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

		if err = r.topSend(topType, send); err != nil {
			r.Logger.Print(err)
		}
	})

	return nil
}

func (r *RankService) Stop() {
	if err := r.saveData(); err != nil {
		r.Logger.Print(err)
	}
}
