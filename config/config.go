package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	e := godotenv.Load()

	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	host := os.Getenv("db_host")
	dbName := os.Getenv("db_name")

	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, username, password, dbName)
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)

	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&todoModel{})
}

func GetDB() *gorm.DB {
	return db
}
