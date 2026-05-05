package main

import (
	"database/sql"
	"encoding/json" // New: To translate our data into JSON
	"fmt"
	"log"
	"net/http"      // New: To build the web server

	_ "github.com/go-sql-driver/mysql"
)

// skill data structure
type Skill struct {
	ID   int    `json:"id"`
	Name string `json:"skill_name"`
}

// We make db a global variable here so our HTTP handlers can easily access it
var db *sql.DB

func main() {

	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/skillhub_db"

	// Open a connection pool
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	fmt.Println("Successfully connected to the MySQL Database!")

	// Our updated ordering window
	http.HandleFunc("/skills", func(w http.ResponseWriter, r *http.Request) {
		// --- THE NEW CORS SECURITY CLEARANCE ---
		// Tell the browser: "I allow requests from Angular (port 4200)"
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		// Tell the browser: "Angular is allowed to use these methods"
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
		// Tell the browser: "Angular is allowed to send JSON headers"
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// The browser sends a quick "OPTIONS" request first to check the rules. 
		// If it's an OPTIONS request, we just say "OK" and stop here.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		
		// get method for getting all skills	
		case http.MethodGet:
			// Query the database
			rows, err := db.Query("SELECT id, skill_name FROM skills")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var skills []Skill
			for rows.Next() {
				var s Skill
				// Scan reads the columns from the database row into our Go struct
				if err := rows.Scan(&s.ID, &s.Name); err != nil {
					log.Println(err)
					continue
				}
				skills = append(skills, s)
			}
			json.NewEncoder(w).Encode(skills)

		// post method for adding a new skill	
		case http.MethodPost:
			var newSkill Skill
			// Decode the JSON order ticket coming from Angular into our Go struct
			if err := json.NewDecoder(r.Body).Decode(&newSkill); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			// Insert into the database (Using '?' prevents SQL injection attacks!)
			result, err := db.Exec("INSERT INTO skills (skill_name) VALUES (?)", newSkill.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			// Grab the new ID and send it back to confirm success
			newID, _ := result.LastInsertId()
			newSkill.ID = int(newID)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newSkill)


		// --- 3. THE DELETE REQUEST (Delete) ---
		case http.MethodDelete:
			// Read the "?id=" part from the URL
			idStr := r.URL.Query().Get("id")
			
			// If they forgot to send an ID, tell them it's a bad request
			if idStr == "" {
				http.Error(w, "Missing skill ID", http.StatusBadRequest)
				return
			}

			// Delete the skill from the database where the ID matches
			_, err := db.Exec("DELETE FROM skills WHERE id = ?", idStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Send back a success message
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Skill deleted successfully"})
			
			
		// --- 4. THE PUT REQUEST (Update) ---
		case http.MethodPut:
			// 1. Get the ID from the URL (e.g., ?id=5)
			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				http.Error(w, "Missing skill ID", http.StatusBadRequest)
				return
			}

			// 2. Read the new JSON data (the updated skill name)
			var updatedSkill Skill
			if err := json.NewDecoder(r.Body).Decode(&updatedSkill); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			// 3. Run the SQL UPDATE command
			_, err := db.Exec("UPDATE skills SET skill_name = ? WHERE id = ?", updatedSkill.Name, idStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// 4. Send a success response
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Skill updated successfully"})
			

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}