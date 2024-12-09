package servers

import (
	"context"
)

type userServer interface {
	Start(ctx context.Context)
	Login(username string, password string) error
}

type workspaceServer interface {
	Start(ctx context.Context)
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

func (s *Servers) Login(username string, password string) error {
	err := s.UserServer.Login(username, password)
	return err
}
