package servers

import (
	"context"
	"github.com/littlebluewhite/Account/dal/model"
	"github.com/littlebluewhite/Account/entry/domain"
)

type userServer interface {
	Start(ctx context.Context)
	Login(username string, password string) (model.User, error)
	Register(register domain.Register) error
	LoginWithToken(token string) (model.User, error)
	Close()
}

type workspaceServer interface {
	Start(ctx context.Context)
	Close()
}

type Servers struct {
	UserServer      userServer
	WorkspaceServer workspaceServer
}

func NewServers(us userServer, wss workspaceServer) *Servers {
	return &Servers{
		UserServer:      us,
		WorkspaceServer: wss,
	}
}

func (s *Servers) Start(ctx context.Context) {
	go func() { s.UserServer.Start(ctx) }()
	go func() { s.WorkspaceServer.Start(ctx) }()
}

func (s *Servers) Close() {
	s.UserServer.Close()
	s.WorkspaceServer.Close()
}

func (s *Servers) Login(username string, password string) (model.User, error) {
	user, err := s.UserServer.Login(username, password)
	return user, err
}

func (s *Servers) Register(register domain.Register) error {
	err := s.UserServer.Register(register)
	return err
}

func (s *Servers) LoginWithToken(token string) (model.User, error) {
	user, err := s.UserServer.LoginWithToken(token)
	return user, err
}
