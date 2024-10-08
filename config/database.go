package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database instance.
var DB *gorm.DB

// region "PostgreConnection" initializes the PostgreSQL database connection.
func PostgreConnection() {
	// Build the Data Source Name (DSN) for connecting to PostgreSQL.
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", os.Getenv("POSTGRE_USER"), os.Getenv("POSTGRE_PASSWORD"), os.Getenv("POSTGRE_HOST"), os.Getenv("POSTGRE_DB"))

	// Open a connection to the database using GORM.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NowFunc: func() time.Time {
		return time.Now().UTC()
	}})

	if err != nil {
		panic(err) // Panic if there is an error while connecting.
	}

	// Get the generic database object from GORM.
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	start := time.Now()

	// Attempt to ping the database until a successful connection is established or 10 seconds have passed.
	for sqlDb.Ping() != nil {
		if start.After(start.Add(10 * time.Second)) {
			fmt.Println("Failed to connection database after 10 seconds")
			break
		}
	}

	fmt.Println("Connected to database: ", sqlDb.Ping() == nil)
	DB = db // Assign the DB instance to the global variable.
}

// endregion
