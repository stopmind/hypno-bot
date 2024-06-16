package core

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"reflect"
)

type HandlersManager struct {
	logger      *log.Logger
	bot         *Bot
	handlers    map[reflect.Type][]func(any)
	initialized bool
}

func (h *HandlersManager) AddHandler(handler any) {
	if !h.initialized {
		h.init()
	}

	t := reflect.TypeOf(handler).In(0)

	slice, ok := h.handlers[t]
	if !ok {
		slice = make([]func(any), 0, 1)
	}

	h.handlers[t] = append(slice, handler.(func(any)))
}

func (h *HandlersManager) init() {
	h.bot.AddHandler(h.handle)
	h.initialized = true
}

func (h *HandlersManager) handle(_ *discordgo.Session, event any) {
	handlers, ok := h.handlers[reflect.TypeOf(event)]

	if !ok {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			h.logger.Println(err)
		}
	}()

	for _, handler := range handlers {
		handler(event)
	}
}

func newHandlersManager(logger *log.Logger, bot *Bot) *HandlersManager {
	return &HandlersManager{
		logger:      logger,
		bot:         bot,
		initialized: false,
		handlers:    make(map[reflect.Type][]func(any)),
	}
}
