package core

import "log"

type Service interface {
	Init(container *ServiceContainer) error
	Stop()
}

type ServiceContainer struct {
	Service  Service
	Storage  *Storage
	Logger   *log.Logger
	Sessions *SessionsManager
	Bot      *Bot
	Name     string
}

func (c *ServiceContainer) Stop() {
	c.Service.Stop()
}