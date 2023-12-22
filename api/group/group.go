package group

import (
	"account/api"
	"account/app/dbs"
	"github.com/gofiber/fiber/v2"
)

type Group struct {
	app fiber.Router
	dbs dbs.Dbs
	wm  api.WebsocketManager
}

func NewAPIGroup(app fiber.Router, dbs dbs.Dbs, wm api.WebsocketManager) *Group {
	return &Group{
		app: app,
		dbs: dbs,
		wm:  wm,
	}
}

func (g *Group) GetApp() fiber.Router {
	return g.app
}

func (g *Group) GetDbs() dbs.Dbs {
	return g.dbs
}

func (g *Group) GetWebsocketManager() api.WebsocketManager {
	return g.wm
}
