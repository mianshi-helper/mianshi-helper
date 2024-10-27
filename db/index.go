package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getDSN() string {
	return "root:123456@tcp(127.0.0.1:3306)/mianshi_helper?charset=utf8mb4&parseTime=True&loc=Local"
}

func ConnectDB() *gorm.DB {
	dsn := getDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db failed: %v\n", err)
	}
	return db
}
