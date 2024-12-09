package mocks

import (
	"github.com/littlebluewhite/Account/api" // Adjust the import path as necessary
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infoln(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Infof(s string, args ...interface{}) {
	m.Called(append([]interface{}{s}, args...)...)
}

func (m *MockLogger) Errorln(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(s string, args ...interface{}) {
	m.Called(append([]interface{}{s}, args...)...)
}

func (m *MockLogger) Warnln(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(s string, args ...interface{}) {
	m.Called(append([]interface{}{s}, args...)...)
}

// Compile-time assertion to ensure MockLogger implements api.Logger
var _ api.Logger = (*MockLogger)(nil)
