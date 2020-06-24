package deliveryman

import (
	"context"
	"git.samberi.com/dois/delivery_api/models"
	"git.samberi.com/dois/delivery_api/proto/gen"
	"git.samberi.com/dois/delivery_api/serializers"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOrder(ctx context.Context, in *gen.GetOrderRequest) (*gen.OrderMessage, error) {
	order := new(models.Order)
	db := ctx.Value("DB").(*gorm.DB)
	user := ctx.Value("user").(*models.User)
	if res := db.Preload(
		"Products.Replacements",
	).Preload(
		"DeliveryMan",
	).Preload(
		"Customer",
	).Find(order, "uid = ? and delivery_man_id = ?", in.OrderUID, user.ID); res.Error != nil {
		return nil, status.Error(codes.NotFound, "Заказ не найден")
	}
	return serializers.ToOrderMessage(order), nil
}

func (s *Server) GetOrders(ctx context.Context, request *gen.EmptyRequest) (resp *gen.GetOrdersResponse, err error) {
	db := ctx.Value("DB").(*gorm.DB)
	user := ctx.Value("user").(*models.User)
	var orders []*models.Order
	db.Preload(
		"DeliveryMan",
	).Preload(
		"Customer",
	).Find(
		&orders,
		"delivery_man_id = ? "+
			"and (status IN (?) "+
			"or (status IN (?) "+
			"and created_at BETWEEN current_date and (current_date + '1 day'::interval)))",
		user.ID,
		[]models.OrderStatus{
			models.ORDER_ACTIVE,
			models.ORDER_SHIPPING,
		},
		[]models.OrderStatus{
			models.ORDER_COMPLETE,
			models.ORDER_CANCELED,
		},
	)
	ordersMess := make([]*gen.OrderMessage, 0)
	for _, order := range orders {
		ordersMess = append(ordersMess, serializers.ToOrderMessage(order))
	}
	resp = &gen.GetOrdersResponse{Orders: ordersMess}
	return resp, nil
}

func (s *Server) ChangeOrderStatus(ctx context.Context, in *gen.ChangeOrderStatusRequest) (*gen.OrderMessage, error) {
	db := ctx.Value("DB").(*gorm.DB)
	user := ctx.Value("user").(*models.User)
	order := new(models.Order)
	if res := db.Find(order, "uid = ? and delivery_man_id = ?", in.OrderUID, user.ID); res.Error != nil {
		return nil, status.Error(codes.NotFound, "Заказ не найден")
	}
	if in.Body.Status == models.ORDER_COMPLETE {
		return nil, status.Error(codes.Aborted, "Ошибка изменния статуса заказа")
	}
	_ = order.SetStatus(in.Body.Status)
	db.Save(order)
	return serializers.ToOrderMessage(order), nil
}

func (s *Server) ConfirmOrder(ctx context.Context, in *gen.ConfirmOrderRequest) (*gen.OrderMessage, error) {
	db := ctx.Value("DB").(*gorm.DB)
	user := ctx.Value("user").(*models.User)
	order := new(models.Order)
	if res := db.Preload(
		"Products.Replacements",
	).Preload(
		"DeliveryMan",
	).Preload(
		"Customer",
	).Find(order, "uid = ? and delivery_man_id =?", in.OrderUID, user.ID); res.Error != nil {
		return nil, status.Error(codes.NotFound, "Заказ не найден")
	}
	if err := order.Manager(db).ConfirmOrder(in); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return serializers.ToOrderMessage(order), nil
}

func (s *Server) SetNotificationToken(ctx context.Context, in *gen.NotificationTokenRequest) (resp *gen.SimpleMessage, err error) {
	db := ctx.Value("DB").(*gorm.DB)
	user := ctx.Value("user").(*models.User)
	err = user.Manager(db).SetNotificationToken(in.Body.DeviceUID, in.Body.NotificationToken)
	if err != nil {
		return nil, status.Error(codes.DataLoss, err.Error())
	}
	resp = &gen.SimpleMessage{
		Success: true,
	}
	return resp, nil
}
