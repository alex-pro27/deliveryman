package deliveryman

import (
	"context"
	"fmt"
	"git.samberi.com/dois/delivery_api/helpers"
	"git.samberi.com/dois/delivery_api/models"
	"git.samberi.com/dois/delivery_api/proto/gen"
	"git.samberi.com/dois/delivery_api/serializers"
	"git.samberi.com/dois/delivery_api/utils"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

func (s *Server) Login(ctx context.Context, in *gen.LoginRequest) (*gen.UserMessage, error) {
	db := ctx.Value("DB").(*gorm.DB)
	user, err := new(models.User).Manager(db).GetByUserName(in.Body.Username)
	if err == nil && user.CheckPassword(in.Body.Password) {
		resp := serializers.ToUserMessage(user)
		return resp, nil
	}
	return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
}

func (s *Server) RequestTwoFactAuth(ctx context.Context, in *gen.RequestTwoFactAuthMessage) (*gen.SimpleMessage, error) {
	db := ctx.Value("DB").(*gorm.DB)
	user, err := new(models.User).Manager(db).GetByUserName(helpers.PreparePhone(in.Body.Phone))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	confirmCode := fmt.Sprintf("%v", rand.Intn(9999))
	if err := utils.SendSms(user.Phone, fmt.Sprintf("Код подтверждения %s", confirmCode)); err != nil {
		return nil, status.Error(codes.Aborted, "Ошибка отправки смс")
	}
	println("CONFIRM CODE:", confirmCode)
	user.TmpPassword = confirmCode
	user.HashTmpPassword()
	db.Save(user)
	return &gen.SimpleMessage{Success: true}, nil
}

func (s *Server) ConfirmTwoFactAuth(ctx context.Context, in *gen.ConfirmTwoFactAuthMessage) (*gen.UserMessage, error) {
	db := ctx.Value("DB").(*gorm.DB)
	user, err := new(models.User).Manager(db).GetByUserName(helpers.PreparePhone(in.Body.Phone))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	if !user.CheckTmpPassword(in.Body.Password) {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	user.TmpPassword = ""
	db.Save(user)
	return serializers.ToUserMessage(user), nil
}
