package admin

import (
	"context"
	"git.samberi.com/dois/delivery_api/common"
	"git.samberi.com/dois/delivery_api/models"
	"git.samberi.com/dois/delivery_api/proto/gen"
	"git.samberi.com/dois/delivery_api/serializers"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAllUsers(ctx context.Context, in *gen.PaginateRequest) (*gen.GetAllUsersResponse, error) {
	db := ctx.Value("DB").(*gorm.DB)
	var users []models.User
	queryset := db.Find(&users).Order("created_at desc")
	data := common.Paginate(
		&users,
		queryset,
		int(in.Page),
		int(in.PageSize),
		[]string{"Token"},
		func(i interface{}) interface{} {
			user := i.(models.User)
			return serializers.ToUserMessage(&user)
		},
	)
	userMessage := make([]*gen.UserMessage, 0)
	for _, it := range data.Data.([]interface{}) {
		userMessage = append(userMessage, it.(*gen.UserMessage))
	}
	resp := &gen.GetAllUsersResponse{
		Paginate: data.Paginate,
		Data:     userMessage,
	}
	return resp, nil
}

func (s *Server) CreateUser(ctx context.Context, in *gen.CreateUserRequest) (*gen.UserMessage, error) {
	user := new(models.User)
	db := ctx.Value("DB").(*gorm.DB)
	err := user.Manager(db).Create(in)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	resp := serializers.ToUserMessage(user)
	return resp, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *gen.UpdateUserRequest) (*gen.UserMessage, error) {
	user := new(models.User)
	db := ctx.Value("DB").(*gorm.DB)
	err := user.Manager(db).Update(in)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	resp := serializers.ToUserMessage(user)
	return resp, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *gen.DeleteUserRequest) (*gen.SimpleMessage, error) {
	user := new(models.User)
	db := ctx.Value("DB").(*gorm.DB)
	if res := db.First(user, "uid = ?", in.UserUID); res.Error != nil {
		return nil, status.Error(codes.Aborted, "user not found")
	}
	err := user.Manager(db).Delete()
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	resp := &gen.SimpleMessage{
		Success: true,
	}
	return resp, err
}

func (s *Server) UserTypes(ctx context.Context, request *gen.EmptyRequest) (*gen.UserTypesResponse, error) {
	panic("implement me")
}
