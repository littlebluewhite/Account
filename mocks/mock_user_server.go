package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockUserServer struct {
	mock.Mock
}

func (m *MockUserServer) ReloadCacheByIDs(ids []int32) error {
	args := m.Called(ids)
	return args.Error(0)
}

func (m *MockUserServer) Close() {
	m.Called()
}
