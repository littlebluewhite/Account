package mocks

import (
	"github.com/littlebluewhite/Account/app/dbs/influxdb"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDbs struct {
	mock.Mock
}

func (m *MockDbs) GetCache() *cache.Cache {
	args := m.Called()
	return args.Get(0).(*cache.Cache)
}

func (m *MockDbs) GetSql() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDbs) GetRdb() redis.UniversalClient {
	args := m.Called()
	return args.Get(0).(redis.UniversalClient)
}

func (m *MockDbs) GetIdb() *influxdb.Influx {
	args := m.Called()
	return args.Get(0).(*influxdb.Influx)
}

func (m *MockDbs) Close() {
	m.Called()
}
