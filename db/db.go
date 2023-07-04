package db

import (
	"fmt"

	"github.com/arvnd79/ie-proj/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Setup(databaseName string) *gorm.DB {
	db := newDB(databaseName)
	migrate(db)
	db.LogMode(true)
	return db
}

func newDB(name string) *gorm.DB {
	db, err := gorm.Open("sqlite3", "./"+name) //dbms is SQLite
	if err != nil {
		fmt.Println("Error in creating database : ", err)
		return nil
	}
	return db
}
func migrate(db *gorm.DB) {
	db.AutoMigrate(&model.User{}, &model.Request{}, &model.URL{})
}
