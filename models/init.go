package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func init() {
	conn, err := gorm.Open("sqlite3", "webmine.db")
	if err != nil {
		log.Println("It's open db error:", err)
	}

	db = conn
	db.Debug().AutoMigrate(&User{}, &Mod{}) //Миграция базы данных
	initModes()
	initUsers()
}

// возвращает дескриптор объекта DB
func GetDB() *gorm.DB {
	return db
}
