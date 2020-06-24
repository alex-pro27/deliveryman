package admin

import (
	"context"
)

type Server struct {
}

func NewSever() *Server {
	return new(Server)
}

func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	switch fullMethodName {
	case "/admin.Admin/Login":
		return ctx, nil
	default:
		return CheckAuth(ctx)
	}
}
