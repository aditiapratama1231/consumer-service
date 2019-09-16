package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //for import mysql
	"github.com/joho/godotenv"
)

// DBInit create connection to database
func DBInit() *gorm.DB {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", username, password, host, port, dbName)
	db, err := gorm.Open("mysql", dbURI)
	if err != nil {
		log.Panic("failed to connect to database")
	}
	db.DB().SetConnMaxLifetime(time.Minute * 5)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(5)

	debug, _ := strconv.ParseBool(os.Getenv("GORM_DEBUG"))
	db.LogMode(debug)
	return db
}
