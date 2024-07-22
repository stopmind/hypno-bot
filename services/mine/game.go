package mine

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand/v2"
)

type ore struct {
	Emoji       string `toml:"emoji"`
	Label       string `toml:"label"`
	Cost        int    `toml:"cost"`
	MinChance   int    `toml:"min_chance"`
	MaxChance   int    `toml:"max_chance"`
	Coefficient int    `toml:"coefficient"`
	Start       int    `toml:"start"`
}

type item struct {
	Emoji    string `toml:"emoji"`
	Label    string `toml:"label"`
	Cost     int    `toml:"cost"`
	Duration int    `toml:"duration"`
}

type userGame struct {
	started         bool
	ores            map[string]int
	next            []string
	collapseCounter int
	step            int
	items           map[string]int
}

func (u *userGame) addOre(ore string, count int) {
	_, ok := u.ores[ore]
	if !ok {
		u.ores[ore] = count
		return
	}
	u.ores[ore] += count
}

func (u *userGame) finish() {
	u.started = false
	u.collapseCounter = 0
	u.ores = make(map[string]int)
	u.step = 0
	u.items = make(map[string]int)
}

func (c *content) getUserGame(id string) *userGame {
	session, _ := c.Sessions.GetSession(id)

	if session != nil {
		session.Extend()
		return session.Data.(*userGame)
	}

	game := &userGame{
		ores:  make(map[string]int),
		items: make(map[string]int),
	}

	_, _ = c.Sessions.NewSession(id, game)

	return game

}

func (c *content) genPaths(count int, step int) []string {
	result := make([]string, count)

	ores := make(map[string]int, len(c.Config.Ores))
	randRange := 0
	for oreID, data := range c.Config.Ores {
		chance := max(data.MinChance, min(data.MaxChance, data.Start+data.Coefficient*step))
		ores[oreID] = chance
		randRange += chance
	}

	for i := 0; i < count; i++ {
		index := 1 + rand.IntN(randRange-1)
		for oreID, chance := range ores {
			index -= chance
			if index <= 0 {
				result[i] = oreID
				break
			}
		}
	}

	return result
}

func (c *content) start(interaction *discordgo.Interaction) error {
	uGame := c.getUserGame(interaction.Member.User.ID)
	uGame.next = c.genPaths(3, 0)
	uGame.started = true

	return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c.locationTemplate("Добро пожаловать в шахту", "Добро пожаловать в шахту №1 в Чертевозеве.\nКопмания \"StopBusiness Mines Ltd.\" предоставляет вам базовые инструменты для выполнения своей работы.\nИх вы должны вернуть после окончания вашей работы. (в ином случае предусмотрен штраф в размере 50 (c))\nТак же мы пригнудительно купим все добытое вами в шахтах по самой демократичной цене на рынке.", map[string]int{}),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: c.Config.Ores[uGame.next[0]].Label,
							Emoji: &discordgo.ComponentEmoji{
								Name: c.Config.Ores[uGame.next[0]].Emoji,
							},
							Style:    discordgo.SecondaryButton,
							CustomID: "p0",
						},
						discordgo.Button{
							Label: c.Config.Ores[uGame.next[1]].Label,
							Emoji: &discordgo.ComponentEmoji{
								Name: c.Config.Ores[uGame.next[1]].Emoji,
							},
							Style:    discordgo.SecondaryButton,
							CustomID: "p1",
						},
						discordgo.Button{
							Label: c.Config.Ores[uGame.next[2]].Label,
							Emoji: &discordgo.ComponentEmoji{
								Name: c.Config.Ores[uGame.next[2]].Emoji,
							},
							Style:    discordgo.SecondaryButton,
							CustomID: "p2",
						},
					},
				},
			},
		},
	})
}

func (c *content) noStarted(interaction *discordgo.Interaction) error {
	return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Это не твоя игра или она закончилась",
		},
	})
}

const (
	cave int = iota
	oldCollapse
	oldCollapseWithCorpse
	oldMine
	none
)

func (c *content) next(interaction *discordgo.Interaction, selectedPath int) error {
	game := c.getUserGame(interaction.Member.User.ID)

	if !game.started {
		return c.noStarted(interaction)
	}

	game.collapseCounter += 1 + rand.IntN(2)
	if game.collapseCounter >= 24 {
		return c.collapse(interaction)
	}

	game.step++

	oreCount := c.Config.MinOre + rand.IntN(c.Config.MaxOre-c.Config.MinOre)

	if d, ok := game.items["pickaxe"]; ok && d > 0 {
		game.items["pickaxe"]--
		oreCount += oreCount / 2
	}

	oreData := c.Config.Ores[game.next[selectedPath]]

	items := map[string]int{
		fmt.Sprintf("%s %s", oreData.Emoji, oreData.Label): oreCount,
	}

	game.addOre(game.next[selectedPath], oreCount)

	location := none

	if game.collapseCounter >= 18 && rand.IntN(10) < 4 {
		location = oldCollapse
		if rand.IntN(10) < 7 {
			location = oldCollapseWithCorpse
		}
	} else if rand.IntN(10) < 2 {
		location = oldMine
	} else if rand.IntN(10) < 2 {
		location = cave
	}

	pathsCount := 3
	message := ""

	addItem := func(id string) {
		itemData := c.Config.Items[id]
		items[fmt.Sprintf("%s %s", itemData.Emoji, itemData.Label)] = 1
		_, ok := game.items[id]
		if !ok {
			game.items[id] = itemData.Duration
			return
		}
		game.items[id] += itemData.Duration
	}

	switch location {
	case none:
		message = c.locationTemplate("Продолжение пути", "Ты просто продолжил свой путь и добыл руды.", items)
		break
	case cave:
		pathsCount = 11
		message = c.locationTemplate("Пещера", "Ты очутился в какой-то большой пещере.", items)
		break
	case oldCollapse:
		message = c.locationTemplate("Старый обвал", "Тут когда-то произошел обвал.\nСтоит быть осторожным.", items)
		break
	case oldCollapseWithCorpse:
		addItem("pickaxe")
		if rand.IntN(3) == 0 {
			addItem("lamp")
		}
		message = c.locationTemplate("Старый обвал", "Тут когда-то произошел обвал.\nНо тут осталось тело.\nТы решил забрать у него часть вещей.", items)
		break
	case oldMine:
		if rand.IntN(2) == 0 {
			addItem("support")
		} else {
			addItem("lamp")
		}
		message = c.locationTemplate("Старая шахта", "Старая шахта, тут точно можно что-то найти.", items)
		break
	}

	if d, ok := game.items["lamp"]; ok && d > 0 {
		game.items["lamp"]--
		pathsCount++
	}

	game.next = c.genPaths(pathsCount, game.step)

	row := discordgo.ActionsRow{Components: []discordgo.MessageComponent{}}
	pathsButtonsRows := []discordgo.MessageComponent{}
	for i := 0; i < pathsCount; i++ {
		if len(row.Components) == 5 {
			pathsButtonsRows = append(pathsButtonsRows, row)
			row = discordgo.ActionsRow{Components: []discordgo.MessageComponent{}}
		}
		row.Components = append(row.Components, discordgo.Button{
			Label: c.Config.Ores[game.next[i]].Label,
			Emoji: &discordgo.ComponentEmoji{
				Name: c.Config.Ores[game.next[i]].Emoji,
			},
			Style:    discordgo.SecondaryButton,
			CustomID: fmt.Sprintf("p%d", i),
		})
	}
	pathsButtonsRows = append(pathsButtonsRows, row)

	mainRow := discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.Button{
			Label: "Покинуть",
			Emoji: &discordgo.ComponentEmoji{
				Name: "⬅️",
			},
			Style:    discordgo.DangerButton,
			CustomID: "exit",
		},
	}}

	if d, ok := game.items["support"]; ok && d > 0 {
		mainRow.Components = append(mainRow.Components, discordgo.Button{
			Label: c.Config.Items["support"].Label,
			Emoji: &discordgo.ComponentEmoji{
				Name: c.Config.Items["support"].Emoji,
			},
			Style:    discordgo.PrimaryButton,
			CustomID: "support",
		})
	}

	return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    message,
			Components: append(pathsButtonsRows, mainRow),
		},
	})
}

func (c *content) collapse(interaction *discordgo.Interaction) error {
	c.getUserData(interaction.Member.User.ID).Balance -= 90
	c.getUserGame(interaction.Member.User.ID).finish()

	return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c.endTemplate(map[string]int{
				"Штраф за потерю инструментов": -50,
				"Спасательные работы":          -40,
			}, "Произошел обвал", "Тебя завалило, но были проведены спасательные работы.\nВсе твои вещи были утеряны."),
		},
	})
}

func (c *content) exit(interaction *discordgo.Interaction) error {
	game := c.getUserGame(interaction.Member.User.ID)

	if !game.started {
		return c.noStarted(interaction)
	}

	data := c.getUserData(interaction.Member.User.ID)
	bills := map[string]int{}

	for oreID, count := range game.ores {
		oreData := c.Config.Ores[oreID]
		data.Balance += oreData.Cost * count
		bills[fmt.Sprintf("%s %s x%d", oreData.Emoji, oreData.Label, count)] = oreData.Cost * count
	}

	game.finish()

	return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c.endTemplate(bills, "Выход", "Ты вышел из шахты и вернул свои инструменты.\nСдав добытую руду, ты получил деньги за нее."),
		},
	})
}

func (c *content) support(interaction *discordgo.Interaction) error {
	game := c.getUserGame(interaction.Member.User.ID)

	if !game.started {
		return c.noStarted(interaction)
	}

	if game.items["support"] == 0 {
		return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "У тебя нет опор.",
			},
		})
	}

	game.items["support"]--
	game.collapseCounter -= 7

	return c.Bot.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ты поставил опору.",
		},
	})
}
