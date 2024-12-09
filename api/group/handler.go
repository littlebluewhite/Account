package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/littlebluewhite/Account/api"
	"github.com/littlebluewhite/Account/entry/e_log"
	"github.com/littlebluewhite/Account/util"
)

type hOperate interface {
	ReadLog(start, stop string) ([]e_log.Log, error)
}

type Handler struct {
	o hOperate
	l api.Logger
}

func NewHandler(o hOperate, l api.Logger) *Handler {
	return &Handler{
		o: o,
		l: l,
	}
}

// GetHistory swagger
// @Summary get logs history
// @Tags    Log
// @Accept  json
// @Produce json
// @Param       start  query     string true "start time"
// @Param       stop  query     string false "stop time"
// @Success 200 {array} e_log.Log
// @Router  /api/logs [get]
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	start := c.Query("start")
	if start == "" {
		return util.Err(c, util.MyErr("No start time input"), 0)
	}
	stop := c.Query("stop")
	data, err := h.o.ReadLog(start, stop)
	if err != nil {
		return util.Err(c, err, 0)
	}
	return c.Status(200).JSON(data)
}
