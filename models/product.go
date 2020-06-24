package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

const (
	PRODUCT_STATUS_SELECTED    = iota + 1 // Отобран
	PRODUCT_STATUS_NOT_FOUND              // Не найден
	PRODUCT_STATUS_REPLACEMENT            // Замена
)

type Product struct {
	gorm.Model
	UID                 uuid.UUID `gorm:"type:uuid;unique:true"`
	Department          int32
	Code                string
	Barcode             string `gorm:"size:50"`
	Count               int32
	Ratio               float64
	Measure             string `gorm:"size:255"`
	Quantity            float64
	OrderID             uint
	IsWeight            bool
	Order               Order
	Vat                 float64 // НДС
	Name                string
	Price               float64
	FullPrice           float64
	DiscountPrice       float64
	SelectedStatus      int32      `gorm:"default:1"`                      // статус отборки
	Replacements        []*Product `gorm:"foreignkey:ParentReplacementID"` // Товары на замену
	ParentReplacementID uint
	ReplacementStatus   bool `gorm:"default:false"` //Статус замены (подтверждени или нет)
}
