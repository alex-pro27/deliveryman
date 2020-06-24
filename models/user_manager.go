package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.samberi.com/dois/delivery_api/helpers"
	"git.samberi.com/dois/delivery_api/proto/gen"
	"git.samberi.com/dois/delivery_api/utils"
	"github.com/jinzhu/gorm"
	"time"
)

type UserManager struct {
	*gorm.DB
	self *User
}

func (manager *UserManager) SetNotificationToken(deviceUID, notificationToken string) error {
	deviceInfo := new(UserDeviceInfo)
	if err := deviceInfo.SetUID(deviceUID); err != nil {
		return err
	}
	manager.DB.First(deviceInfo, "device_uid = ?", deviceInfo.DeviceUID)
	deviceInfo.NotificationToken = notificationToken
	deviceInfo.UserID = manager.self.ID
	manager.DB.Save(deviceInfo)
	return nil
}

func (manager *UserManager) Create(in *gen.CreateUserRequest) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, in.Body); err != nil {
		return err
	}
	errs := make(map[string]string)
	if res := manager.First(&User{}, "user_name = ?", manager.self.UserName); !res.RecordNotFound() {
		errs["UserName"] = fmt.Sprintf("Имя пользователя %s уже занято", manager.self.UserName)
	}
	if manager.self.Email != "" {
		if res := manager.First(&User{}, "email = ?", manager.self.Email); !res.RecordNotFound() {
			errs["Email"] = fmt.Sprintf("Email %s уже занят", manager.self.Email)
		}
	}
	if len(errs) > 0 {
		message, _ := json.Marshal(errs)
		return errors.New(string(message))
	}
	(&Token{}).Manager(manager.DB).NewToken(manager.self)
	manager.DB.Create(manager.self)
	if manager.self.UID, err = utils.GenerateUUID(manager.self.ID); err == nil {
		if res := manager.DB.Save(manager.self); res.Error != nil {
			err = res.Error
		}
	}
	return err
}

func (manager *UserManager) Update(in *gen.UpdateUserRequest) (err error) {
	if res := manager.First(manager.self, "uid = ?", in.UserUID); res.Error != nil {
		return errors.New("Пользователь не найден")
	}
	if err = helpers.SetFieldsForModel(manager.self, in.Body); err != nil {
		return err
	}
	errs := make(map[string]string)
	if manager.self.Email != "" {
		res := manager.First(&User{}, "email = ? and not id = ?", manager.self.Email, manager.self.ID)
		if !res.RecordNotFound() {
			errs["Email"] = fmt.Sprintf("Email %s занят", manager.self.Email)
		}
	}
	if len(errs) > 0 {
		message, _ := json.Marshal(errs)
		return errors.New(string(message))
	}
	(&Token{}).Manager(manager.DB).NewToken(manager.self)
	manager.Save(manager.self)
	return nil
}

func (manager *UserManager) Delete() (err error) {
	now := time.Now()
	manager.self.DeletedAt = &now
	manager.self.Active = false
	manager.Save(manager.self)
	return nil
}

func (manager *UserManager) GetByUID(uid string) (*User, error) {
	if res := manager.Preload(
		"Token",
	).First(
		manager.self, "uid = ?", uid,
	); res.Error != nil {
		return nil, res.Error
	}
	return manager.self, nil
}

func (manager *UserManager) GetByUserName(username string) (*User, error) {
	if res := manager.Preload(
		"Token",
	).First(
		manager.self, "active = true AND user_name = ? OR email = ? OR phone = ?",
		username, username, username,
	); res.Error != nil {
		return nil, res.Error
	}
	return manager.self, nil
}

func (manager *UserManager) GetByToken(token string) (*User, error) {
	manager.First(&manager.self.Token, "key = ?", token)
	if res := manager.Find(
		manager.self, "token_id = ?", manager.self.Token.ID,
	); res.Error != nil {
		return nil, res.Error
	}
	return manager.self, nil
}
