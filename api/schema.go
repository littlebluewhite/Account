package api

import (
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/littlebluewhite/Account/app/dbs/influxdb"
	"github.com/littlebluewhite/Account/entry/e_module"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dbs interface {
	GetSql() *gorm.DB
	GetCache() *cache.Cache
	GetRdb() redis.UniversalClient
	GetIdb() *influxdb.Influx
	Close()
}

type HubManager interface {
	RegisterHub(module e_module.Module)
	Broadcast(module e_module.Module, message []byte)
	WsConnect(module e_module.Module, conn *websocket.Conn) error
}

type Servers interface {
	Start(ctx context.Context)
	Login(username string, password string) error
}

type Logger interface {
	Infoln(args ...interface{})
	Infof(s string, args ...interface{})
	Errorln(args ...interface{})
	Errorf(s string, args ...interface{})
	Warnln(args ...interface{})
	Warnf(s string, args ...interface{})
}
