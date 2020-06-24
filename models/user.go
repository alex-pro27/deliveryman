package models

import (
	"errors"
	"fmt"
	"git.samberi.com/dois/delivery_api/helpers"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/wesovilabs/koazee"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

const (
	USERTYPE_CUSTOMER int64 = iota + 1
	USERTYPE_DELIVERYMAN
	USERTYPE_ADMIN
)

var UserTypeChoices = map[int64]string{
	USERTYPE_CUSTOMER:    "Покупатель",
	USERTYPE_DELIVERYMAN: "Доставщик",
	USERTYPE_ADMIN:       "Администратор",
}

type User struct {
	gorm.Model
	UID         uuid.UUID        `gorm:"type:uuid;"`
	FirstName   string           `gorm:"size:255;not null"`
	LastName    string           `gorm:"size:255;"`
	UserName    string           `gorm:"size:255;unique_index;not null"`
	Password    string           `gorm:"size:60"`
	TmpPassword string           `gorm:"size:60"`
	Email       string           `gorm:"type:varchar(100)"`
	Phone       string           `gorm:"type:varchar(17)"`
	Active      bool             `gorm:"default:true"`
	UserTypes   pq.Int64Array    `gorm:"type:integer[]"`
	Roles       []Role           `gorm:"many2many:users_roles;"`
	IsSuperUser bool             `gorm:"default:false"`
	Orders      []Order          `gorm:"foreignkey:UserID"`
	DeviceInfo  []UserDeviceInfo `gorm:"foreignkey:UserID"`
	Address     string
	TokenID     uint
	Token       Token
}

func (user User) GetChoicesUserType() map[int64]string {
	return UserTypeChoices
}

func (user User) GetUserTypesNames() string {
	userTypesNames := make([]string, 0)
	for _, userType := range user.UserTypes {
		userTypesNames = append(userTypesNames, UserTypeChoices[userType])
	}
	return strings.Join(userTypesNames, ", ")
}

func (user *User) SetUserTypes(userTypes []int64) error {
	allUserTypesStream := koazee.StreamOf([]int64{
		USERTYPE_CUSTOMER,
		USERTYPE_DELIVERYMAN,
		USERTYPE_ADMIN,
	})
	_userTypes := koazee.StreamOf(userTypes).RemoveDuplicates().Filter(
		func(i int64) bool {
			index, _ := allUserTypesStream.IndexOf(i)
			return index > -1
		},
	)
	user.UserTypes = _userTypes.Out().Val().([]int64)
	return nil
}

func (user User) GetFullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user User) String() string {
	return user.GetFullName()
}

func (user User) IsAdmin() bool {
	if !user.IsSuperUser {
		userRoleTypes := koazee.StreamOf(user.UserTypes).Filter(func(x int64) bool {
			return x == USERTYPE_ADMIN
		}).Out().Val().([]int64)
		return len(userRoleTypes) > 0
	}
	return true
}

func (user *User) SetPhone(phone string) error {
	phone = helpers.PreparePhone(phone)
	if phone != "" {
		if matched, _ := regexp.MatchString("^7\\d{10}", phone); !matched {
			return errors.New("Неверно заполнен номер телефона")
		}
	}
	user.Phone = phone
	return nil
}

func (user *User) SetEmail(email string) error {
	email = strings.ToLower(strings.Trim(email, ""))
	if email == "" {
		return nil
	}
	matched, _ := regexp.MatchString(
		"^[a-z0-9_][a-z0-9.\\-_]{1,100}@[a-z0-9\\-_]{1,100}\\.[a-z0-9\\-_]{1,50}[a-z0-9_]$",
		email,
	)
	if !matched {
		return errors.New("Неверно заполнен email")
	}
	user.Email = email
	return nil
}

func (user *User) SetPassword(password string) error {
	if password = strings.Trim(password, ""); len(password) < 3 {
		return errors.New("пароль должен состоять минимум из 3х символов")
	}
	user.Password = password
	user.HashPassword()
	return nil
}

func (user *User) HashPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	user.Password = string(hash)
}

func (user *User) HashTmpPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.TmpPassword), bcrypt.MinCost)
	user.TmpPassword = string(hash)
}

func (user User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user User) CheckTmpPassword(tmpPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.TmpPassword), []byte(tmpPassword))
	return err == nil
}

func (user *User) Manager(db *gorm.DB) *UserManager {
	return &UserManager{DB: db, self: user}
}
