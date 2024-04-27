package services

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
)

type GamesSession struct {
}

type GamesService struct {
	*core.ServiceContainer
}

func genMap(width int, height int, bombsCount int) []int {
	result := make([]int, width*height)
	bombs := make([]int, bombsCount)
	for i := 0; i < bombsCount; i++ {
		bombs[i] = -1
	}

	for i := 0; i < bombsCount; i++ {
		var pos int

		for {
			pos = rand.IntN(len(result))
			if !slices.Contains[[]int](bombs, pos) {
				break
			}
		}

		bombs[i] = pos

		x, y := pos%width, pos/width

		if x -= 1; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}
		if x += 2; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}

		y -= 1
		if x -= 2; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}
		if x += 1; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}
		if x += 1; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}

		y += 2
		if x -= 2; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}
		if x += 1; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}
		if x += 1; x >= 0 && x < width && y >= 0 && y < height {
			result[x+height*y] += 1
		}
	}

	for _, pos := range bombs {
		result[pos] = -1
	}

	return result
}

var toEmoji = map[int]string{
	-1: ":red_square:",
	0:  ":blue_square:",
	1:  ":one:",
	2:  ":two:",
	3:  ":three:",
	4:  ":four:",
	5:  ":five:",
	6:  ":six:",
	7:  ":seven:",
	8:  ":eight:",
	9:  ":nine:",
}

func (g *GamesService) Sapper(send *discordgo.MessageCreate) {
	reply := "Неправильные аргументы"

	args := strings.Split(send.Content, " ")

	var width, height, bombsCount int
	var m []int

	switch len(args) {
	case 1:
		width, height, bombsCount = 5, 5, 5
	case 4:
		var err error
		if width, err = strconv.Atoi(args[1]); err != nil {
			goto end
		}
		if height, err = strconv.Atoi(args[2]); err != nil {
			goto end
		}
		if bombsCount, err = strconv.Atoi(args[3]); err != nil {
			goto end
		}

		if bombsCount > width*height {
			goto end
		}

	default:
		goto end
	}

	m = genMap(width, height, bombsCount)

	reply = ""
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			reply += fmt.Sprintf("||%s||", toEmoji[m[x+y*width]])
		}

		reply += "\n"
	}

end:
	_, err := g.Bot.ChannelMessageSendReply(send.ChannelID, reply, send.Reference())

	if err != nil {
		g.Logger.Print(err)
	}
}

func (g *GamesService) Init(container *core.ServiceContainer) error {
	g.ServiceContainer = container

	container.Bot.AddHandler(func(session *discordgo.Session, send *discordgo.MessageCreate) {
		if !strings.HasPrefix(send.Content, "!sapper") || send.Author.Bot {
			return
		}
		g.Sapper(send)
	})

	return nil
}

func (g *GamesService) Stop() {

}
