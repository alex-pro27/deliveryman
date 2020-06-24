package deliveryman

import (
	"context"
	"git.samberi.com/dois/delivery_api/common"
)

type Server struct {
}

func (server *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	switch fullMethodName {
	case "/deliveryman.Api/RequestTwoFactAuth", "/deliveryman.Api/ConfirmTwoFactAuth":
		return ctx, nil
	default:
		return common.CheckAuthByToken(ctx)
	}
}

func NewServer() *Server {
	return new(Server)
}
