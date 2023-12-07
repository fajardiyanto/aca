package database

import (
	"fmt"
	"log"

	"github.com/aca/permit/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	host := "localhost"
	port := 3306
	user := "root"
	password := "root"
	dbname := "permit"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err.Error())
	}
	db.AutoMigrate(&models.Permit{}, &models.Auth{}, &models.Position{}, &models.Simper{}, &models.Department{})

	return db.Debug()
}
