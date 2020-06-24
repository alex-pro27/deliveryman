package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/wesovilabs/koazee"
)

type OrderStatus int32

const (
	ORDER_ACTIVE = iota + 1
	ORDER_SHIPPING
	ORDER_COMPLETE
	ORDER_CANCELED
)

type Order struct {
	gorm.Model
	UID           uuid.UUID `gorm:"type:uuid;"`
	DeliveryManID uint      `gorm:"default:null"`
	DeliveryMan   User
	CustomerID    uint
	Customer      User
	Products      []Product `gorm:"foreignkey:OrderID"`
	OriginalSum   float64   // Изначальная сумма заказа
	ConfirmSum    float64   // Новая подтвержденна сумма заказа
	Address       string
	Latitude      float64
	Longitude     float64
	Status        int32
}

func (order *Order) SetStatus(status int32) error {
	statuses := koazee.StreamOf([]int32{
		ORDER_ACTIVE,
		ORDER_SHIPPING,
		ORDER_COMPLETE,
		ORDER_CANCELED,
	})
	index, _ := statuses.IndexOf(status)
	if index > -1 {
		order.Status = status
		return nil
	}
	return errors.New("status not supported")
}

func (order *Order) Manager(db *gorm.DB) *OrderManager {
	return &OrderManager{DB: db, self: order}
}
