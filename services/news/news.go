package news

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils"
	"strings"
	"time"
)

type Service struct {
	*core.ServiceContainer

	config struct {
		EditorsChannel string
		PublishChannel string
		EditorsRole    string
	}

	state struct {
		NextID int
	}

	currentRelease release
}

type block struct {
	Title   string
	Content string
	Checker string
	Author  string
}

type release struct {
	Tagline string
	Summary string
	Index   int
	Date    string

	Blocks []*block
}

func (s *Service) Init(container *core.ServiceContainer) error {
	s.ServiceContainer = container

	s.currentRelease = release{
		Blocks: make([]*block, 0),
	}

	err := s.Storage.ReadTOML("config.toml", &s.config)

	if err != nil {
		return err
	}

	_ = s.Storage.ReadJson("state.json", &s.state)

	s.Bot.AddHandler(func(session *discordgo.Session, send *discordgo.MessageCreate) {
		if !strings.HasPrefix(send.Content, "?газета") || send.Author.Bot {
			return
		}

		args, err := utils.NewArgsParser().
			AddString().
			AddString().
			Parse(0, send.Content)

		if err != nil {
			s.replyIncorrectArgsError(send, err)
			return
		}

		switch args.Get(0) {
		case nil:
			s.replyFile(send, "assets/index.txt")
		case "справка":
			s.replyFile(send, "assets/help.txt")
		case "предложить":
			s.propose(send)
		case "опубликовать":
			s.publish(send)
		case "одобрить":
			s.accept(send)
		default:
			s.replyError(send, ":stop_sign: Некоректная команда", fmt.Sprintf("Неизвестная команда: `%v`\nПопробуйте `?газета справка`", args.Get(0)))
		}
	})

	return nil
}

func (s *Service) Stop() {
	if err := s.Storage.WriteJson("state.json", &s.state); err != nil {
		s.Logger.Print(err)
	}
}

func (s *Service) propose(send *discordgo.MessageCreate) {
	args, err := utils.NewArgsParser().
		AddString().
		AddString().
		AddString().
		Parse(3, send.Content)

	if err != nil {
		s.replyIncorrectArgsError(send, err)
		return
	}

	session, _ := s.Sessions.GetSession(send.Author.ID)

	if session == nil {
		session, err = s.Sessions.NewSession(send.Author.ID, make([]*block, 0))
		if err != nil {
			s.replyUnexpectedError(send, err)
			return
		}
	}

	newBlock := &block{
		Title:   strings.ReplaceAll(args.Get(1).(string), "_", " "),
		Content: args.Get(2).(string),
		Author:  send.Author.ID,
	}

	id := len(session.Data.([]*block))
	session.Data = append(session.Data.([]*block), newBlock)

	templ, err := s.Storage.GetTemplate("assets/notify.jet")
	if err != nil {
		s.replyUnexpectedError(send, err)
		return
	}

	text, err := utils.ExecuteTemplate(templ, &struct {
		NewBlock *block
		Id       int
		Service  string
	}{newBlock, id, s.config.EditorsRole})
	if err != nil {
		s.replyUnexpectedError(send, err)
		return
	}

	_, err = s.Bot.ChannelMessageSend(s.config.EditorsChannel, text)
	if err != nil {
		s.Logger.Print(err)
		return
	}
}

func (s *Service) publish(send *discordgo.MessageCreate) {
	if !s.checkRole(send) {
		return
	}

	args, err := utils.NewArgsParser().
		AddString().
		AddString().
		AddString().
		Parse(3, send.Content)

	if err != nil {
		s.replyIncorrectArgsError(send, err)
		return
	}

	s.currentRelease.Tagline = strings.ReplaceAll(args.Get(1).(string), "_", " ")
	s.currentRelease.Summary = args.Get(2).(string)
	s.currentRelease.Date = time.Now().Format("02.01.2006")
	s.currentRelease.Index = s.state.NextID

	tmpl, err := s.Storage.GetTemplate("assets/begin.tmp")
	if err == nil {
		text, err := utils.ExecuteTemplate(tmpl, s.currentRelease)
		if err == nil {
			_, err = s.Bot.ChannelMessageSend(s.config.PublishChannel, text)
		}
	}
	if err != nil {
		s.replyUnexpectedError(send, err)
		return
	}

	tmpl, err = s.Storage.GetTemplate("assets/block.tmp")
	if err != nil {
		s.replyUnexpectedError(send, err)
		return
	}
	for _, block := range s.currentRelease.Blocks {
		text, err := utils.ExecuteTemplate(tmpl, block)
		if err == nil {
			_, err = s.Bot.ChannelMessageSend(s.config.PublishChannel, text)
		}
		if err != nil {
			s.replyUnexpectedError(send, err)
			return
		}
	}

	tmpl, err = s.Storage.GetTemplate("assets/end.tmp")
	if err == nil {
		text, err := utils.ExecuteTemplate(tmpl, s.currentRelease)
		if err == nil {
			_, err = s.Bot.ChannelMessageSend(s.config.PublishChannel, text)
		}
	}
	if err != nil {
		s.replyUnexpectedError(send, err)
		return
	}

	err = s.Storage.WriteJson(fmt.Sprintf("archive/%v.json", s.state.NextID), &s.currentRelease)
	if err != nil {
		s.replyUnexpectedError(send, err)
	}

	s.state.NextID += 1
	err = s.Storage.WriteJson("state.json", &s.state)
	if err != nil {
		s.replyUnexpectedError(send, err)
	}

	s.currentRelease = release{
		Blocks: make([]*block, 0),
	}
}

func (s *Service) accept(send *discordgo.MessageCreate) {
	if !s.checkRole(send) {
		return
	}

	agrs, err := utils.NewArgsParser().
		AddString().
		AddString().
		AddInt().
		Parse(3, send.Content)

	if err != nil {
		s.replyIncorrectArgsError(send, err)
		return
	}

	author := agrs.Get(1).(string)
	id := agrs.Get(2).(int)

	session, err := s.Sessions.GetSession(author)
	var aBlock *block
	if err == nil && len(session.Data.([]*block)) > id {
		aBlock = session.Data.([]*block)[id]
	}
	if aBlock == nil {
		s.replyError(send, ":bangbang: Не найдено", "Не удалось найти данное предложение")
		return
	}

	session.Data.([]*block)[id] = nil

	aBlock.Checker = send.Author.ID
	s.currentRelease.Blocks = append(s.currentRelease.Blocks, aBlock)
}
