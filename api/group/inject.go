package group

import (
	"account/api"
	"account/app/dbs"
	"account/util/logFile"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"io"
	"os"
)

func Inject(app *fiber.App, dbs dbs.Dbs, wm api.WebsocketManager) {
	// Middleware
	log := logFile.NewLogFile("api", "inject.log")
	fiberLog := getFiberLogFile(log)
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Output: fiberLog,
	}))

	//swagger routes
	app.Get("/swagger/*", swagger.HandlerDefault)

	// api group add cors middleware
	Api := app.Group("/api", cors.New())

	// use middleware to write log
	o := NewOperate(dbs)
	h := NewHandler(o, log)
	Api.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		err = o.WriteLog(c)
		if err != nil {
			log.Error().Println(err)
		}
		return err
	})
	Api.Get("/logs", h.GetHistory)

	// create new group
	_ = NewAPIGroup(Api, dbs, wm)

	// model registration
}

func getFiberLogFile(log logFile.LogFile) io.Writer {
	fiberFile, err := os.OpenFile("./log/fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Fatal("can not open log file: " + err.Error())
	}
	return io.MultiWriter(fiberFile, os.Stdout)
}
