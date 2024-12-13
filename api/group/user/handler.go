package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/entry/domain"
	"github.com/littlebluewhite/Account/util"
)

type server interface {
	Login(username, password string) (model.User, error)
	Register(register domain.Register) error
	LoginWithToken(token string) (model.User, error)
}

type Handler struct {
	s server
	l domain.Logger
}

func NewHandler(s server, l domain.Logger) *Handler {
	return &Handler{
		s: s,
		l: l,
	}
}

// Login swagger
// @Summary login with username and password
// @Tags    User
// @Accept  json
// @Produce json
// @Param   login  body  domain.Login  true  "username and password"
// @Success 200 {string} string
// @Router  /api/account/user/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	entry := domain.Login{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("Login: ", err)
		return util.Err(c, err, 0)
	}
	loginUser, err := h.s.Login(entry.Username, entry.Password)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(loginUser)
}

// Register swagger
// @Summary register user
// @Tags    User
// @Accept  json
// @Produce json
// @Param   register  body  domain.Login  true  "username and password"
// @Success 200 {string} string
// @Router  /api/account/user/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	entry := domain.Register{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("Register: ", err)
		return util.Err(c, err, 0)
	}
	if err := h.s.Register(entry); err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON("register success")
}

// LoginWithToken swagger
// @Summary login with token
// @Tags    User
// @Accept  json
// @Produce json
// @Param   loginWithToken  body  domain.LoginWithToken true  "token"
// @Success 200 {string} string
// @Router  /api/account/user/loginWithToken [post]
func (h *Handler) LoginWithToken(c *fiber.Ctx) error {
	entry := domain.LoginWithToken{}
	if err := c.BodyParser(&entry); err != nil {
		h.l.Errorln("Register: ", err)
		return util.Err(c, err, 0)
	}
	loginUser, err := h.s.LoginWithToken(entry.Token)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(loginUser)
}
