package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:255"`
	Permissions []Permission `gorm:"many2many:roles_permissions;" form:"label:Разрешение"`
	Users       []User       `gorm:"many2many:users_roles;" form:"label:Пользователи"`
}
