package builder

import "hypno-bot/core"

type WithContainer struct {
	*core.ServiceContainer
}

func (w *WithContainer) initContainer(container *core.ServiceContainer) {
	w.ServiceContainer = container
}

type containerInit interface {
	initContainer(container *core.ServiceContainer)
}
