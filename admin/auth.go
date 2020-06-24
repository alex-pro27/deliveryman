package admin

import (
	"context"
	"git.samberi.com/dois/delivery_api/models"
	"git.samberi.com/dois/delivery_api/proto/gen"
	"git.samberi.com/dois/delivery_api/serializers"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckAuth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Token")
	if err != nil {
		return nil, err
	}
	db := ctx.Value("DB").(*gorm.DB)
	user, err := new(models.User).Manager(db).GetByToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}
	if user.IsAdmin() {
		newCtx := context.WithValue(ctx, "user", user)
		return newCtx, nil
	}
	return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
}

func (s *Server) Login(ctx context.Context, req *gen.LoginRequest) (*gen.UserMessage, error) {
	user := new(models.User)
	db := ctx.Value("DB").(*gorm.DB)
	if res := db.Preload(
		"Token",
	).First(
		user,
		"user_name = ? and (? = ANY (user_types) or is_super_user = true) and active = true",
		req.Body.Username,
		models.USERTYPE_ADMIN,
	); res.Error != nil || !user.CheckPassword(req.Body.Password) {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	resp := serializers.ToUserMessage(user)
	return resp, nil
}
