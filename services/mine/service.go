package mine

import (
	"github.com/bwmarrin/discordgo"
	"hypno-bot/core"
	"hypno-bot/utils/builder"
	"strconv"
)

type content struct {
	builder.WithContainer
	builder.WithConfig[struct {
		Ores   map[string]*ore  `toml:"ores"`
		Items  map[string]*item `toml:"items"`
		MaxOre int              `toml:"maxOre"`
		MinOre int              `toml:"minOre"`
	}]
	builder.WithState[struct {
		Users map[string]*userData `json:"users"`
	}]
}

func (c *content) command(create *discordgo.InteractionCreate) {
	var err error

	switch create.ApplicationCommandData().Options[0].Name {
	case "копать":
		err = c.start(create.Interaction)
	case "баланс":
		err = c.balance(create.Interaction)
		break
	}

	if err != nil {
		c.Logger.Print(err.Error())
	}
}

func (c *content) button(create *discordgo.InteractionCreate) {
	var err error
	id := create.MessageComponentData().CustomID
	if id[0] == 'p' {
		path, _ := strconv.Atoi(id[1:])
		err = c.next(create.Interaction, path)
	} else if id == "exit" {
		err = c.exit(create.Interaction)
	} else if id == "support" {
		err = c.support(create.Interaction)
	}

	if err != nil {
		c.Logger.Print(err.Error())
	}
}

func Build() core.Service {
	c := &content{}
	return builder.BuildService(c).
		AddSlashCommand("шахта", "шахта", c.command).
		AddSubCommand("копать", "копать").
		AddSubCommand("баланс", "баланс").
		EndCommand().
		AddComponentHandler(c.button).
		Finish()
}
