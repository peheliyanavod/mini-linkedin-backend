package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // The MySQL Driver
)

func main() {
	// Format: "username:password@tcp(127.0.0.1:3306)/database_name"
	dsn := "root:@tcp(127.0.0.1:3306)/skillhub_db"

	// Open a connection pool
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	fmt.Println("Successfully connected to the MySQL Database!")
}