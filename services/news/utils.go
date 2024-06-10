package news

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/utils"
	"slices"
)

func (s *Service) checkRole(send *discordgo.MessageCreate) bool {
	if !slices.Contains(send.Member.Roles, s.config.EditorsRole) {
		utils.ReplyError(send, ":pinching_hand: Недостаточно прав", fmt.Sprintf("У тебя нет роли, челедь."))
		return false
	}

	return true
}
