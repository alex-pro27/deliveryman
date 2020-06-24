package models

import (
	"git.samberi.com/dois/delivery_api/utils"
	"github.com/jinzhu/gorm"
)

type TokenManager struct {
	*gorm.DB
	self *Token
}

func (manager *TokenManager) NewToken(user *User) {
	if user.TokenID != 0 {
		manager.Delete(&Token{}, user.TokenID)
	}
	key := utils.GenerateHash()
	t := Token{}
	manager.First(&t, "key = ?", key)
	if t.ID != 0 {
		manager.NewToken(user)
	} else {
		manager.self.Key = key
		manager.Create(manager.self)
		manager.NewRecord(*manager.self)
		user.TokenID = manager.self.ID
		user.Token = *manager.self
	}
}
