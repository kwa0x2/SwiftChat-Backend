package config

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection(){
	db, err := gorm.Open(postgres.Open("postgres://nettasec:nettaseclocal@localhost:5437/nettasec_global_db?sslmode=disable"), &gorm.Config{})


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