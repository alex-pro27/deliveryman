package models

import (
	"crypto/md5"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"io"
)

type UserDeviceInfo struct {
	gorm.Model
	DeviceUID         uuid.UUID `gorm:"type:uuid;"`
	NotificationToken string    `gorm:"size:255"`
	UserID            uint
	User              User
}

func (deviceInfo *UserDeviceInfo) PrepareUID(uid string) (uuid.UUID, error) {
	h := md5.New()
	io.WriteString(h, uid)
	return uuid.FromBytes(h.Sum(nil))
}

func (deviceInfo *UserDeviceInfo) SetUID(uid string) (err error) {
	deviceInfo.DeviceUID, err = deviceInfo.PrepareUID(uid)
	if err != nil {
		return err
	}
	return nil
}
