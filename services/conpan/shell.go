package conpan

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand/v2"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type shellSession struct {
	isAdmin bool
}

func (c *content) execute(session *shellSession, send *discordgo.MessageCreate, args []string) {
	checkAdmin := func() bool {
		if !session.isAdmin {
			err := c.Reply(send, "You are not a admin")
			if err != nil {
				c.Logger.Print(err)
			}
		}

		return session.isAdmin
	}

	checkUpdatesAccepted := func() bool {
		if c.State.UpdatesBlocked {
			err := c.Reply(send, "Updates blocked")
			if err != nil {
				c.Logger.Print(err)
			}
		}

		return !c.State.UpdatesBlocked
	}

	switch args[0] {
	case "code":
		c.useCode(session, send, args[1])
		break
	case "request-code":
		if checkAdmin() {
			c.requestCode(send)
		}
		break
	case "restart":
		if checkAdmin() {
			c.restart()
		}
		break
	case "update":
		if checkAdmin() && checkUpdatesAccepted() {
			c.update()
		}
		break
	case "read":
		if checkAdmin() {
			c.read(send, args[1])
		}
		break
	case "log":
		if checkAdmin() {
			c.log(send, args[1])
		}
		break
	case "write":
		if checkAdmin() && checkUpdatesAccepted() {
			c.write(send, args)
		}
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

func (c *content) requestCode(send *discordgo.MessageCreate) {
	code := strconv.Itoa(rand.Int())
	c.State.Codes = append(c.State.Codes, code)

	err := c.Reply(send, code)
	if err != nil {
		c.Logger.Print(err)
	}
}

func (c *content) update() {
	err := exec.Command("sh", "-c", "cd .. && ./bot update && ./bot restart").Run()
	if err != nil {
		c.Logger.Print(err)
	}
}

func (c *content) restart() {
	err := exec.Command("sh", "-c", "cd .. && ./bot restart").Run()
	if err != nil {
		c.Logger.Print(err)
	}
}

func (c *content) read(send *discordgo.MessageCreate, path string) {
	data, err := os.ReadFile(fmt.Sprintf("storage/%s", path))
	message := fmt.Sprintf("```%s\n```", string(data))
	if err != nil {
		message = "Failed to read file"
	}

	err = c.Reply(send, message)
	if err != nil {
		c.Logger.Print(err)
	}
}

func (c *content) log(send *discordgo.MessageCreate, logName string) {
	data, err := os.ReadFile(fmt.Sprintf("log/%s.log", logName))
	message := fmt.Sprintf("```%s\n```", string(data))
	if err != nil {
		message = "Failed to read log file"
	}

	err = c.Reply(send, message)
	if err != nil {
		c.Logger.Print(err)
	}
}

func (c *content) write(send *discordgo.MessageCreate, args []string) {
	path := args[1]
	data := strings.Join(args[2:], " ")

	err := os.WriteFile(fmt.Sprintf("storage/%s", path), []byte(data), 0644)
	if err != nil {
		err = c.Reply(send, err.Error())
		if err != nil {
			c.Logger.Print(err)
		}
	}
}
