package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/Account/api"
)

type Group struct {
	app fiber.Router
	dbs api.Dbs
	hm  api.HubManager
	s   api.Servers
}

func NewAPIGroup(app fiber.Router, dbs api.Dbs, hm api.HubManager, s api.Servers) *Group {
	return &Group{
		app: app,
		dbs: dbs,
		hm:  hm,
		s:   s,
	}
}

func (g *Group) GetApp() fiber.Router {
	return g.app
}

func (g *Group) GetDbs() api.Dbs {
	return g.dbs
}

func (g *Group) GetWebsocketManager() api.HubManager {
	return g.hm
}

func (g *Group) GetServers() api.Servers {
	return g.s
}
