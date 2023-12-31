package main

import (
	"account/api/group"
	"account/app/dbs"
	"account/app/websocket_manager"
	_ "account/docs"
	"account/util/config"
	"account/util/logFile"
	"context"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	mainLog logFile.LogFile
)

// 初始化配置
func init() {
	// log配置
	mainLog = logFile.NewLogFile("", "main.log")
}

// @title           Schedule-Task-Command swagger API
// @version         2.0.0
// @description     This is a account server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Wilson
// @contact.url    https://github.com/littlebluewhite
// @contact.email  wwilson008@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:5487

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mainLog.Info().Println("command module start")

	// DBs start includes SQL Cache
	DBS := dbs.NewDbs(mainLog, false)
	defer func() {
		DBS.GetIdb().Close()
		mainLog.Info().Println("influxDB Disconnect")
	}()

	// create websocket manager
	wm := websocket_manager.NewWebsocketManager()
	// run websocket manager
	go wm.Run()

	ServerConfig := config.NewConfig[config.ServerConfig](".", "env", "server")

	var sb strings.Builder
	sb.WriteString(":")
	sb.WriteString(ServerConfig.Port)
	//addr := sb.String()
	apiServer := fiber.New(
		fiber.Config{
			ReadTimeout:  ServerConfig.ReadTimeout * time.Minute,
			WriteTimeout: ServerConfig.WriteTimeout * time.Minute,
			AppName:      "account",
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		},
	)

	group.Inject(apiServer, DBS, wm)

	// for api server shout down gracefully
	serverShutdown := make(chan struct{})
	go func() {
		_ = <-ctx.Done()
		mainLog.Info().Println("Gracefully shutting down api server")
		_ = apiServer.Shutdown()
		serverShutdown <- struct{}{}
	}()

	if err := apiServer.Listen(":5487"); err != nil {
		mainLog.Error().Fatalf("listen: %s\n", err)
	}

	// Listen for the interrupt signal.
	<-serverShutdown

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	time.Sleep(1 * time.Second)
	mainLog.Info().Println("Server exiting")

}
