package middleware

import (
	"context"
	"git.samberi.com/dois/delivery_api/database"
	"google.golang.org/grpc"
)

func dbSetContext(ctx context.Context) context.Context {
	db := database.ConnectDB()
	return context.WithValue(ctx, "DB", db)
}

func DbSetContextUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		nextCtx := dbSetContext(ctx)
		return handler(nextCtx, req)
	}
}
