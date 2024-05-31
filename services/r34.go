package services

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type R34Service struct {
	*core.ServiceContainer

	config struct {
		TrollMod    bool
		TrollTarget string
		TrollTags   []string
		TrollCount  int
	}
}

func getImages(count int, tags []string) ([]string, error) {
	var posts []struct {
		FileUrl string `json:"file_url"`
	}

	response, err := http.Get(fmt.Sprintf("https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&json=1&limit=300&tags=%v", strings.Join(tags, "+")))
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err = response.Body.Close(); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(raw, &posts); err != nil {
		return nil, err
	}

	images := make([]string, count)

	usedPosts := make([]int, 0)
	for i := 0; i < count; i++ {
		var id int
		for {
			id = rand.IntN(len(posts))
			if !slices.Contains(usedPosts, id) {
				break
			}
		}

		usedPosts = append(usedPosts, id)
		images[i] = posts[id].FileUrl
	}

	return images, nil
}

func (r *R34Service) Init(container *core.ServiceContainer) error {
	r.ServiceContainer = container

	err := r.Storage.ReadTOML("config.toml", &r.config)
	if err != nil {
		return err
	}

	r.Bot.AddHandler(func(session *discordgo.Session, send *discordgo.MessageCreate) {
		r.command(send)
	})

	return nil
}

func (r *R34Service) command(send *discordgo.MessageCreate) {
	if !strings.HasPrefix(send.Content, "!r34") || send.Author.Bot {
		return
	}

	var tags []string
	var err error
	var count int

	if r.config.TrollMod && send.Author.ID == r.config.TrollTarget {
		count = r.config.TrollCount
		tags = r.config.TrollTags
	} else {
		parts := strings.Split(send.Content, " ")

		count = 5

		if len(parts) >= 2 {
			count, err = strconv.Atoi(parts[1])
			if err != nil {
				count = 5
			}
		}

		tags = make([]string, 0)
		if len(parts) >= 3 {
			tags = parts[2:]
		}
	}

	images, err := getImages(count, tags)

	if err != nil {
		if _, err = r.Bot.ChannelMessageSendReply(send.ChannelID, err.Error(), send.Reference()); err != nil {
			r.Logger.Print(err)
		}

		return
	}

	const ipm = 5
	messageCount := int(math.Ceil(float64(len(images)) / ipm))

	for i := 0; i < messageCount; i++ {
		if _, err = r.Bot.ChannelMessageSendReply(send.ChannelID, strings.Join(images[i*ipm:int(math.Min(float64(len(images)), float64(i+1)*ipm))], "\n"), send.Reference()); err != nil {
			r.Logger.Print(err)
		}
	}
}

func (r *R34Service) Stop() {

}
