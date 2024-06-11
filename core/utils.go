package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
}

func NewLogger(name string) (*log.Logger, error) {
	path := fmt.Sprintf("log/%s.log", name)
	file, err := openFile(path)

	if err != nil {
		return nil, err
	}

	return log.New(file, "", 0), nil
}

func (c *ServiceContainer) Reply(send *discordgo.MessageCreate, content string) error {
	_, err := c.Bot.ChannelMessageSendReply(send.ChannelID, content, send.Reference())

	return err
}

func MemberName(member *discordgo.Member) string {
	name := member.Nick

	if name == "" {
		name = member.User.GlobalName
	}
	if name == "" {
		name = member.User.Username
	}

	return name
}
