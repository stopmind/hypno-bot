package conpan

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand/v2"
	"os"
	"os/exec"
	"slices"
	"strconv"
)

type shellSession struct {
	isAdmin bool
}

func (c *content) execute(session *shellSession, send *discordgo.MessageCreate, args []string) {
	switch args[0] {
	case "code":
		c.useCode(session, send, args[1])
		break
	case "request-code":
		c.requestCode(session, send)
		break
	}
}

func (c *content) useCode(session *shellSession, send *discordgo.MessageCreate, code string) {
	if !slices.Contains(c.State.Codes, code) {
		err := c.Reply(send, "Code is incorrect")
		if err != nil {
			c.Logger.Print(err)
		}
		return
	}

	i := slices.Index(c.State.Codes, code)
	c.State.Codes = slices.Delete(c.State.Codes, i, i+1)

	session.isAdmin = true
}

func (c *content) requestCode(session *shellSession, send *discordgo.MessageCreate) {
	if !session.isAdmin {
		err := c.Reply(send, "You are not a admin")
		if err != nil {
			c.Logger.Print(err)
		}
		return
	}

	code := strconv.Itoa(rand.Int())
	c.State.Codes = append(c.State.Codes, code)

	err := c.Reply(send, code)
	if err != nil {
		c.Logger.Print(err)
	}
}

func (c *content) update(session *shellSession, send *discordgo.MessageCreate) {
	if !session.isAdmin {
		err := c.Reply(send, "You are not a admin")
		if err != nil {
			c.Logger.Print(err)
		}
		return
	}
}

func (c *content) restart(session *shellSession, send *discordgo.MessageCreate) {
	if !session.isAdmin {
		err := c.Reply(send, "You are not a admin")
		if err != nil {
			c.Logger.Print(err)
		}
		return
	}

	_ = exec.Command("sh", "-c", "../bot restart").Run()
}

func (c *content) read(session *shellSession, send *discordgo.MessageCreate, path string) {
	if !session.isAdmin {
		err := c.Reply(send, "You are not a admin")
		if err != nil {
			c.Logger.Print(err)
		}
		return
	}

	data, err := os.ReadFile(fmt.Sprintf("storage/%s", path))
	message := string(data)
	if err != nil {
		message = "Failed to read file"
	}

	err = c.Reply(send, message)
	if err != nil {
		c.Logger.Print(err)
	}
}

func (c *content) log(session *shellSession, send *discordgo.MessageCreate, logName string) {
	if !session.isAdmin {
		err := c.Reply(send, "You are not a admin")
		if err != nil {
			c.Logger.Print(err)
		}
		return
	}

	data, err := os.ReadFile(fmt.Sprintf("log/%s", logName))
	message := string(data)
	if err != nil {
		message = "Failed to read log file"
	}

	err = c.Reply(send, message)
	if err != nil {
		c.Logger.Print(err)
	}
}
