package main

import (
	"database/sql"
	"encoding/json" // New: To translate our data into JSON
	"fmt"
	"log"
	"net/http"      // New: To build the web server

	_ "github.com/go-sql-driver/mysql"
)

// dummy data structure
type Skill struct {
	ID   int    `json:"id"`
	Name string `json:"skill_name"`
}

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

	http.HandleFunc("/skills", func(w http.ResponseWriter, r *http.Request) {
		
		dummySkills := []Skill{
			{ID: 10, Name: "Go"},
			{ID: 20, Name: "Angular"},
		}

		// Set the content type to JSON so Angular knows how to read it
		w.Header().Set("Content-Type", "application/json")
		
		// Translate our Go data into JSON and send it back
		json.NewEncoder(w).Encode(dummySkills)
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}