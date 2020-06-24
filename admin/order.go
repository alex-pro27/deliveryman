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

func (s *Server) GetAllOrders(ctx context.Context, in *gen.PaginateRequest) (*gen.GetAllOrdersResponse, error) {
	db := ctx.Value("DB").(*gorm.DB)
	var orders []models.Order
	queryset := db.Preload(
		"DeliveryMan",
	).Preload(
		"Customer",
	).Find(&orders).Order("created_at desc")
	data := common.Paginate(
		&orders,
		queryset,
		int(in.Page),
		int(in.PageSize),
		[]string{"Customer", "DeliveryMan"},
		func(obj interface{}) interface{} {
			order := obj.(models.Order)
			return serializers.ToOrderMessage(&order)
		},
	)
	var orderMessage []*gen.OrderMessage
	for _, it := range data.Data.([]interface{}) {
		orderMessage = append(orderMessage, it.(*gen.OrderMessage))
	}
	resp := &gen.GetAllOrdersResponse{
		Paginate: data.Paginate,
		Data:     orderMessage,
	}
	return resp, nil
}

func (s *Server) SetOrderUser(ctx context.Context, in *gen.SetOrderUserRequest) (*gen.OrderMessage, error) {
	order := new(models.Order)
	user := new(models.User)
	db := ctx.Value("DB").(*gorm.DB)
	if res := db.Preload(
		"DeliveryMan",
	).Preload(
		"Customer",
	).Preload(
		"Products",
	).First(
		order, "uid = ?", in.OrderUID,
	); res.Error != nil {
		return nil, status.Error(codes.Aborted, "Order not found")
	}
	if res := db.First(user, "uid = ?", in.UserUID); res.Error != nil {
		return nil, status.Error(codes.Aborted, "User not found")
	}
	resp := serializers.ToOrderMessage(order)
	return resp, nil
}
