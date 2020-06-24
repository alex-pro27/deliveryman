package serializers

import (
	"git.samberi.com/dois/delivery_api/models"
	"git.samberi.com/dois/delivery_api/proto/gen"
)

func ToCustomerMessage(customer *models.User) *gen.CustomerMessage {
	return &gen.CustomerMessage{
		UID:      customer.UID.String(),
		Email:    customer.Email,
		Phone:    customer.Phone,
		Address:  customer.Address,
		FullName: customer.GetFullName(),
	}
}

func ToProductMessage(product *models.Product) *gen.ProductMessage {
	replacementsUID := make([]string, 0)
	for _, r := range product.Replacements {
		replacementsUID = append(replacementsUID, r.UID.String())
	}
	return &gen.ProductMessage{
		UID:               product.UID.String(),
		Name:              product.Name,
		Count:             product.Count,
		Ratio:             product.Ratio,
		Measure:           product.Measure,
		Quantity:          product.Quantity,
		IsWeight:          product.IsWeight,
		Price:             product.Price,
		FullPrice:         product.FullPrice,
		DiscountPrice:     product.DiscountPrice,
		SelectedStatus:    product.SelectedStatus,
		ReplacementsUID:   replacementsUID,
		ReplacementStatus: product.ReplacementStatus,
	}
}

func ToRepeatProductMessage(products []models.Product) []*gen.ProductMessage {
	resp := make([]*gen.ProductMessage, 0)
	for _, product := range products {
		resp = append(resp, ToProductMessage(&product))
	}
	return resp
}

func ToOrderMessage(order *models.Order) *gen.OrderMessage {
	return &gen.OrderMessage{
		UID:            order.UID.String(),
		DeliverymanUID: order.DeliveryMan.UID.String(),
		Customer:       ToCustomerMessage(&order.Customer),
		Products:       ToRepeatProductMessage(order.Products),
		Address:        order.Address,
		Latitude:       order.Latitude,
		Longitude:      order.Longitude,
		Status:         order.Status,
		OriginalSum:    order.OriginalSum,
		ConfirmSum:     order.ConfirmSum,
	}
}

func ToUserMessage(user *models.User) *gen.UserMessage {
	return &gen.UserMessage{
		UID:         user.UID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Phone:       user.Phone,
		Active:      user.Active,
		IsSuperUser: user.IsSuperUser,
		UserTypes:   user.UserTypes,
		Token:       user.Token.Key,
	}
}
