package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func PostgreConnection(){
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", os.Getenv("POSTGRE_USER"), os.Getenv("POSTGRE_PASSWORD"), os.Getenv("POSTGRE_HOST"), os.Getenv("POSTGRE_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})


	if err != nil{
		panic(err)
	}
	sqlDb, err := db.DB()

	if err != nil{
		panic(err)
	}

	start := time.Now()

	for sqlDb.Ping() != nil{
		if start.After(start.Add(10 * time.Second)) {
			fmt.Println("Failed to connection database after 10 seconds")
			break
		}
	}

	fmt.Println("Connected to database: ",sqlDb.Ping() == nil)
	DB = db
}