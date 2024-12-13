package main

import (
	"fmt"
	version "github.com/littlebluewhite/Account"
	"github.com/littlebluewhite/Account/api/group"
	"github.com/littlebluewhite/Account/app/dbs"
	"github.com/littlebluewhite/Account/app/servers"
	"github.com/littlebluewhite/Account/app/servers/user_server"
	"github.com/littlebluewhite/Account/app/servers/workspace_server"
	"github.com/littlebluewhite/Account/app/websocket_hub"
	"github.com/littlebluewhite/Account/docs"
	"github.com/littlebluewhite/Account/util/my_log"
	"path/filepath"
	"runtime"

	"context"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	_ "github.com/littlebluewhite/Account/docs"
	"github.com/littlebluewhite/Account/util/config"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	rootPath string
)

// 初始化配置
func init() {
	// log配置
	_, b, _, _ := runtime.Caller(0)
	rootPath = filepath.Dir(filepath.Dir(filepath.Dir(b)))
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

// @host      127.0.0.1:9600
// @BasePath  /api

func main() {
	mainLog := my_log.NewLog("main")
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mainLog.Infoln("command module start")

	// read config
	Config := config.NewConfig[config.Config](rootPath, "config", "config", config.Yaml)

	ServerConfig := Config.Server

	// swagger docs host
	docsHost := fmt.Sprintf("%s:%s", ServerConfig.SwaggerHost, ServerConfig.Port)
	docs.SwaggerInfo.Host = docsHost
	docs.SwaggerInfo.Version = version.Version

	// DBs start includes SQL Cache
	DBS := dbs.NewDbs(mainLog, false, Config)
	defer func() {
		DBS.GetIdb().Close()
		mainLog.Infoln("influxDB Disconnect")
	}()

	// create websocket manager
	hm := websocket_hub.NewHubManager()

	// create servers
	userServer := user_server.NewUserServer(DBS)
	workspaceServer := workspace_server.NewWorkspaceServer(DBS)
	s := servers.NewServers(userServer, workspaceServer)

	// start server
	go func() {
		s.Start(ctx)
	}()

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

	group.Inject(apiServer, DBS, hm, s)

	// for api server shout down gracefully
	serverShutdown := make(chan struct{})
	go func() {
		_ = <-ctx.Done()
		mainLog.Infoln("Gracefully shutting down api server")
		_ = apiServer.Shutdown()
		serverShutdown <- struct{}{}
	}()

	if err := apiServer.Listen(":9600"); err != nil {
		mainLog.Errorf("listen: %s\n", err)
	}

	// Listen for the interrupt signal.
	<-serverShutdown

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	DBS.Close()
	s.Close()
	mainLog.Infoln("Server exiting")

}
