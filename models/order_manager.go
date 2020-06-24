package models

import (
	"errors"
	"git.samberi.com/dois/delivery_api/logger"
	"git.samberi.com/dois/delivery_api/proto/gen"
	"git.samberi.com/dois/delivery_api/services/payment_service"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"math"
)

type OrderManager struct {
	*gorm.DB
	self *Order
}

func (manager *OrderManager) ConfirmOrder(data *gen.ConfirmOrderRequest) error {
	orderProducts := koazee.StreamOf(manager.self.Products)
	tx := manager.Begin()
	for _, item := range data.Body.Products {
		find := orderProducts.Filter(func(p Product) bool {
			return p.UID.String() == item.UID
		}).Out().Val().([]Product)
		if len(find) == 0 {
			continue
		}
		product := find[0]
		product.Count = item.Count
		tx.Save(&product)
		for _, _item := range item.Replacements {
			_find := orderProducts.Filter(func(p Product) bool {
				return p.UID.String() == item.UID
			}).Out().Val().([]Product)
			if len(_find) == 0 {
				continue
			}
			rp := _find[0]
			rp.Count = _item.Count
			rp.ReplacementStatus = true
			tx.Save(&rp)
		}
	}
	products := make([]payment_service.Product, 0)
	for _, item := range manager.self.Products {
		quantity := float64(item.Count) * item.Quantity * item.Ratio
		// Цена за одну позицию
		priceDiscount := math.Round((item.DiscountPrice/quantity)*100) / 100
		// Сумма за всю позицию
		sumPos := math.Round(priceDiscount*quantity*100) / 100
		// Скидка за одну позицию
		discount := math.Round((item.FullPrice-item.DiscountPrice)*100) / 100
		// ндс за одну позицю
		taxSum := math.Round((priceDiscount/(1+item.Vat/100))*(item.Vat/100)*100) / 100
		products = append(products, payment_service.Product{
			Name:          item.Name,
			Code:          item.Code,
			Quantity:      item.Quantity,
			Measure:       item.Measure,
			Barcode:       item.Barcode,
			Department:    item.Department,
			Price:         item.Price,
			Discount:      discount,
			TaxSum:        taxSum,
			TaxRate:       item.Vat,
			Amount:        sumPos,
			PriceDiscount: priceDiscount,
		})
	}
	if err := payment_service.ConfirmOrder(payment_service.Order{
		UID:      manager.self.UID.String(),
		Products: products,
		Customer: payment_service.Customer{
			UID:   manager.self.Customer.UID.String(),
			Phone: manager.self.Customer.Phone,
		},
	}); err == nil {
		if err := tx.Commit().Error; err != nil {
			logger.HandleError(err)
			return errors.New("ошибка подтверждения заказа")
		}
		manager.self.Status = ORDER_COMPLETE
		manager.DB.Save(manager.self)
		return nil
	} else {
		logger.HandleError(err)
		return err
	}
}
