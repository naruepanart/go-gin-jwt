package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() (*gorm.DB) {
	var err error
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database : error=%v", err)
		return nil
	}

	return db
}