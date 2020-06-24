package database

import (
	"fmt"
	"git.samberi.com/dois/delivery_api/config"
	"git.samberi.com/dois/delivery_api/logger"
	"git.samberi.com/dois/delivery_api/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	if DB == nil {
		dbConf := config.Config.Database
		params := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			dbConf.Host,
			dbConf.Port,
			dbConf.User,
			dbConf.Database,
			dbConf.Password,
		)
		var err error
		DB, err = gorm.Open("postgres", params)
		logger.HandleError(err)
		if config.Config.System.Debug {
			DB.LogMode(true)
		}
	}
	return DB
}

var Models = []interface{}{
	models.User{},
	models.Token{},
	models.Order{},
	models.Product{},
	models.UserDeviceInfo{},
	models.Role{},
	models.Permission{},
	models.ContentType{},
}

func MigrateDB() {
	db := ConnectDB()
	db.AutoMigrate(Models...)
	db.Model(models.User{}).AddForeignKey("token_id", "tokens(id)", "CASCADE", "CASCADE")
	db.Model(models.Order{}).AddForeignKey("delivery_man_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(models.Order{}).AddForeignKey("customer_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(models.Product{}).AddForeignKey("order_id", "orders(id)", "CASCADE", "CASCADE")
	db.Model(models.Product{}).AddForeignKey("parent_replacement_id", "products(id)", "CASCADE", "CASCADE")
	db.Model(models.UserDeviceInfo{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(models.Permission{}).AddForeignKey("content_type_id", "content_types(id)", "CASCADE", "CASCADE")
	db.Table("users_roles").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Table("users_roles").AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")
	db.Table("roles_permissions").AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")
	db.Table("roles_permissions").AddForeignKey("permission_id", "permissions(id)", "CASCADE", "CASCADE")
	access := []models.PermissionAccess{
		models.READ,
		models.WRITE,
		models.DELETE,
	}
	tx := db.Begin()
	for _, model := range Models {
		tableName := tx.NewScope(model).GetModelStruct().TableName(db)
		contentType := new(models.ContentType)
		tx.FirstOrCreate(contentType, models.ContentType{
			Model: tableName,
		})
		for _, a := range access {
			tx.FirstOrCreate(new(models.Permission), models.Permission{
				Access:        a,
				ContentTypeID: contentType.ID,
			})
		}
	}
	tx.Commit()
}
