package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/Account/entry/domain"
	"github.com/littlebluewhite/Account/util/my_log"
)

func RegisterRouter(g group) {
	log := my_log.NewLog("router/user")
	app := g.GetApp()

	c := app.Group("/user")

	c.Use(func(c *fiber.Ctx) error {
		c.Locals("Module", "account")
		return c.Next()
	})

	h := NewHandler(g.GetServers(), log)

	c.Post("/login", h.Login)
	c.Post("/register", h.Register)
	c.Post("/loginWithToken", h.LoginWithToken)
}

type group interface {
	GetApp() fiber.Router
	GetServers() domain.Servers
}
