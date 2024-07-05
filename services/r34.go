package services

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/services/achievements"
	"hypno-bot/utils/builder"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type R34Service struct {
	builder.WithContainer
	builder.WithConfig[struct {
		Channel string
	}]
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

	if len(posts) < count {
		images := make([]string, len(posts))
		for i, post := range posts {
			images[i] = post.FileUrl
		}

		return images, nil
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

func (r *R34Service) checkTimer(send *discordgo.MessageCreate) bool {
	session, _ := r.Sessions.GetSession(send.Author.ID)
	if session != nil && time.Now().Sub(session.Data.(time.Time)) < time.Second*1 {
		r.replyFile(send, "assets/toofast.txt")
		return true
	}

	if session == nil {
		session, _ = r.Sessions.NewSession(send.Author.ID, time.Now())
	} else {
		session.Data = time.Now()
		session.Extend()
	}

	return false
}

func (r *R34Service) command(send *discordgo.MessageCreate) {
	if send.ChannelID != r.Config.Channel {
		r.replyFile(send, "assets/channelmessage.txt")
		return
	}

	if r.checkTimer(send) {
		return
	}

	var tags []string
	var err error
	var count int

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

	if count > 100 {
		r.replyFile(send, "assets/cowboy.txt")
		return
	}

	images, err := getImages(count, tags)

	if len(images) == 0 || err != nil {
		r.replyFile(send, "assets/cantfind.txt")
		return
	}

	if len(images) < count {
		r.replyFile(send, "assets/toofew.txt")
	}

	const ipm = 5
	messageCount := int(math.Ceil(float64(len(images)) / ipm))

	for i := 0; i < messageCount; i++ {
		if _, err = r.Bot.ChannelMessageSendReply(send.ChannelID, strings.Join(images[i*ipm:int(math.Min(float64(len(images)), float64(i+1)*ipm))], "\n"), send.Reference()); err != nil {
			r.Logger.Print(err)
		}
	}

	achievements.OnR34(send.Author.ID)
}

func (r *R34Service) commandTo(send *discordgo.MessageCreate) {
	if send.ChannelID != r.Config.Channel {
		r.replyFile(send, "assets/channelmessage.txt")
		return
	}

	var tags []string
	var err error
	var count int
	target := ""

	parts := strings.Split(send.Content, " ")
	count = 5

	if len(parts) >= 2 {
		targetRaw := parts[1]
		if strings.HasPrefix(targetRaw, "<@") && strings.HasSuffix(targetRaw, ">") {
			target = strings.TrimSuffix(strings.TrimPrefix(targetRaw, "<@"), ">")
		}
	}

	if target == "" {
		r.replyFile(send, "assets/comu.txt")
		return
	}

	if len(parts) >= 3 {
		count, err = strconv.Atoi(parts[2])
		if err != nil {
			count = 5
		}
	}

	if r.checkTimer(send) {
		return
	}

	tags = make([]string, 0)
	if len(parts) >= 4 {
		tags = parts[3:]
	}

	if count > 100 {
		r.replyFile(send, "assets/cowboy.txt")
		return
	}

	images, err := getImages(count, tags)

	if len(images) == 0 || err != nil {
		r.replyFile(send, "assets/cantfind.txt")
		return
	}

	if len(images) < count {
		r.replyFile(send, "assets/toofew.txt")
	}

	for _, image := range images {
		_, err = r.Bot.ChannelMessageSend(send.ChannelID, fmt.Sprintf("<@%v> %v", target, image))
		if err != nil {
			r.Logger.Print(err)
		}
	}

	achievements.OnR34To(send.Author.ID)
}

func (r *R34Service) replyFile(send *discordgo.MessageCreate, path string) {
	content, err := r.Storage.ReadFile(path)

	if err != nil {
		r.Logger.Print(err)
		return
	}

	err = r.Reply(send, string(content))
	if err != nil {
		r.Logger.Print(err)
	}
}

func BuildR34Service() core.Service {
	c := new(R34Service)
	return builder.BuildService(c).
		AddCommand("!r34", c.command).
		AddCommand("!r34to", c.commandTo).
		Finish()
}
