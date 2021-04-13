package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var db *gorm.DB //база данных

func init() {

	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}

	conn, err := gorm.Open("sqlite3", "webmine.db")
	if err != nil {
		log.Println("It's error", err)
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
