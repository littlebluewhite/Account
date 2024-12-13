package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/Account/entry/domain"
)

type Group struct {
	app fiber.Router
	dbs domain.Dbs
	hm  domain.HubManager
	s   domain.Servers
}

func NewAPIGroup(app fiber.Router, dbs domain.Dbs, hm domain.HubManager, s domain.Servers) *Group {
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

func (g *Group) GetDbs() domain.Dbs {
	return g.dbs
}

func (g *Group) GetWebsocketManager() domain.HubManager {
	return g.hm
}

func (g *Group) GetServers() domain.Servers {
	return g.s
}
