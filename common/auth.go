package common

import (
	"context"
	"git.samberi.com/dois/delivery_api/models"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckAuthByToken(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Token")
	if err != nil {
		return nil, err
	}
	db := ctx.Value("DB").(*gorm.DB)
	user, err := new(models.User).Manager(db).GetByToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}
	newCtx := context.WithValue(ctx, "user", user)
	return newCtx, nil
}
