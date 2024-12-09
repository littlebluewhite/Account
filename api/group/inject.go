package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/littlebluewhite/Account/api"
	"github.com/littlebluewhite/Account/api/group/user"
	"github.com/littlebluewhite/Account/util/my_log"
	"io"
	"os"
)

func Inject(app *fiber.App, dbs api.Dbs, hm api.HubManager, s api.Servers) {
	// Middleware
	log := my_log.NewLog("api/inject.log")
	fiberLog := getFiberLogFile(log)
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Output: fiberLog,
	}))

	//swagger routes
	app.Get("/swagger/*", swagger.HandlerDefault)

	// api group add cors middleware
	Api := app.Group("/api/account", cors.New())

	// use middleware to write log
	o := NewOperate(dbs)
	h := NewHandler(o, log)
	Api.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		err = o.WriteLog(c)
		if err != nil {
			log.Errorln(err)
		}
		return err
	})
	Api.Get("/logs", h.GetHistory)

	// create new group
	g := NewAPIGroup(Api, dbs, hm, s)

	// model registration
	user.RegisterRouter(g)
}

func getFiberLogFile(log api.Logger) io.Writer {
	fiberFile, err := os.OpenFile("./log/fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("can not open log file: " + err.Error())
	}
	return io.MultiWriter(fiberFile, os.Stdout)
}
