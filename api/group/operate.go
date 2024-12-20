package group

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/littlebluewhite/Account/app/dbs/influxdb"
	"github.com/littlebluewhite/Account/entry/e_log"
	"time"
)

type Operate struct {
	idb *influxdb.Influx
}

type dbs interface {
	GetIdb() *influxdb.Influx
}

func NewOperate(dbs dbs) *Operate {
	o := &Operate{
		idb: dbs.GetIdb(),
	}
	return o
}

func (o *Operate) WriteLog(c *fiber.Ctx) (err error) {
	now := time.Now()
	response := c.Response()
	l := e_log.Log{
		Timestamp:     float64(now.UnixMilli()) / 1000000,
		Account:       c.Get("Account"),
		ContentLength: len(response.Body()),
		Datetime:      now,
		IP:            c.IP(),
		Referer:       c.Get("Referer"),
		StatusCode:    response.StatusCode(),
		Token:         c.Get("Authorization"),
		UserAgent:     c.Get("User-Agent"),
		WebPath:       c.Get("Web-Path"),
	}
	jL, err := json.Marshal(l)
	if err != nil {
		return
	}
	p := influxdb2.NewPoint("log",
		map[string]string{},
		map[string]interface{}{"data": jL},
		now,
	)
	o.idb.Writer().WritePoint(p)
	return
}

func (o *Operate) ReadLog(start, stop string) (logs []e_log.Log, err error) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	stmt := fmt.Sprintf(`from(bucket:"account")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "log")
|> filter(fn: (r) => r._field == "data")
`, start, stopValue)
	result, err := o.idb.Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var l e_log.Log
			v := result.Record().Value()
			vString, ok := v.(string)
			if !ok {
				fmt.Printf("value: %v is not string", v)
				continue
			}
			if e := json.Unmarshal([]byte(vString), &l); e != nil {
				fmt.Printf(e.Error())
				continue
			}
			logs = append(logs, l)
		}
	} else {
		return
	}
	// send empty []
	if logs == nil {
		logs = make([]e_log.Log, 0)
	}
	return
}
