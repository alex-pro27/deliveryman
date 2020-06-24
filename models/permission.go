package models

import "github.com/jinzhu/gorm"

type PermissionAccess uint

const (
	FORBIDDEN PermissionAccess = 0
	READ      PermissionAccess = 3
	WRITE     PermissionAccess = 5
	DELETE    PermissionAccess = 7
)

var PermissionAccessChoices = map[PermissionAccess]string{
	FORBIDDEN: "Не разрешено",
	READ:      "Только для чтения",
	WRITE:     "Доступ на запись",
	DELETE:    "Доступ на удаление",
}

type Permission struct {
	gorm.Model
	ContentTypeID uint
	ContentType   ContentType      `gorm:"foreignkey: ContentTypeID"`
	Access        PermissionAccess `gorm:"default:3"`
}

func (Permission) GetChoiceAccess() map[PermissionAccess]string {
	return PermissionAccessChoices
}

func (permission Permission) GetPermissionName() string {
	return PermissionAccessChoices[permission.Access]
}
