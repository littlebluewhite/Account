package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/Account/api"
	"github.com/littlebluewhite/Account/util"
)

type server interface {
	Login(username, password string) error
}

type Handler struct {
	s server
	l api.Logger
}

func NewHandler(s server, l api.Logger) *Handler {
	return &Handler{
		s: s,
		l: l,
	}
}

// Login swagger
// @Summary login with username and password
// @Tags    Login
// @Accept  json
// @Produce json
// @Param   login  body  user.Login  true  "username and password"
// @Success 200 {string} string
// @Router  /api/account/user/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	entry := Login{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("Login: ", err)
		return util.Err(c, err, 0)
	}
	if err := h.s.Login(entry.Username, entry.Password); err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("login success")
}
