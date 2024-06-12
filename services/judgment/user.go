package judgment

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils"
	"math"
	"strconv"
)

type review struct {
	Comment string `json:"comment"`
	Points  int    `json:"points"`
	Author  string `json:"author"`
}

type userInfo struct {
	Karma       int      `json:"karma"`
	LastReviews []review `json:"last_reviews"`
}

func (c *content) getUser(userId string) *userInfo {
	user, ok := c.State.Users[userId]

	if !ok {
		if c.State.Users == nil {
			c.State.Users = make(map[string]*userInfo)
		}

		user = &userInfo{
			Karma:       0,
			LastReviews: make([]review, 0),
		}
		c.State.Users[userId] = user
	}

	return user
}

func (u *userInfo) addReview(newReview review) {
	u.LastReviews = append([]review{newReview}, u.LastReviews...)

	if len(u.LastReviews) > 10 {
		u.LastReviews = u.LastReviews[:5]
	}

	u.Karma = int(math.Max(math.Min(float64(u.Karma+newReview.Points), 10), -10))
}

func (c *content) karma(send *discordgo.MessageCreate) {
	args, err := utils.NewArgsParser().
		AddString().
		AddUser().
		Parse(2, send.Content)

	if err != nil {
		utils.ReplyIncorrectArgsError(send, err)
		return
	}

	var intToStr = func(a int) string {
		if a > 0 {
			return fmt.Sprintf("+%d", a)
		}

		return strconv.Itoa(a)
	}

	userId := args.Get(1).(string)
	info := c.getUser(userId)

	result := fmt.Sprintf("Карма <@%v>: %v\n\n", userId, intToStr(info.Karma))

	for _, r := range info.LastReviews {
		member, err := c.Bot.GuildMember(send.GuildID, r.Author)
		name := ""

		if err == nil {
			name = core.MemberName(member)
		} else {
			user, err := c.Bot.User(r.Author)
			if err == nil {
				name = user.Username
			} else {
				name = "неизвестного"
			}
		}

		result = fmt.Sprintf("%v%v (*%v*) от %v\n", result, r.Comment, intToStr(r.Points), name)

	}

	if len(info.LastReviews) == 0 {
		result += "*Отзывов нет*"
	}

	err = c.Reply(send, result)
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
	}
}

func (c *content) review(send *discordgo.MessageCreate) {
	args, err := utils.NewArgsParser().
		AddString().
		AddUser().
		AddInt().
		AddString().
		Parse(4, send.Content)

	if err != nil {
		utils.ReplyIncorrectArgsError(send, err)
		return
	}

	userId := args.Get(1).(string)

	if userId == send.Author.ID {
		err = c.Reply(send, "Оставлять отзыв на самого себя нельзя!")
		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}
		return
	}

	info := c.getUser(userId)
	r := review{
		Comment: args.Get(3).(string),
		Points:  args.Get(2).(int),
		Author:  send.Author.ID,
	}

	if math.Abs(float64(r.Points)) > 5 {
		err = c.Reply(send, "Ваш отзыв слишком мощный. Ставьте отметку в диапозоне от -5 до 5.")
		if err != nil {
			utils.ReplyUnexpectedError(send, err)
		}
		return
	}

	info.addReview(r)
	err = c.Reply(send, "Ваш отзыв был принят.")
	if err != nil {
		utils.ReplyUnexpectedError(send, err)
	}
}
