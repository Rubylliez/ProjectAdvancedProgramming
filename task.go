package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID          uint   `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Gender      string `json:"gender"`
}

var (
	db *gorm.DB
)

func main() {
	var err error
	dsn := "host=localhost user=postgres password=1959 dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Could not migrate table:", err)
	}

	fmt.Println("Connected to the database")

	http.HandleFunc("/uploadjson", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var newUser User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			http.Error(w, "Failed to insert user into database", http.StatusInternalServerError)
			return
		}

		successMsg := map[string]string{"message": "User successfully inserted into the database"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successMsg)
	})

	http.HandleFunc("/getusers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var users []User
		db.Find(&users)

		csvData := "ID,Full Name,Username,Email,Phone Number,Password,Gender\n"
		for _, user := range users {
			csvData += fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s\n", user.ID, user.FullName, user.Username, user.Email, user.PhoneNumber, user.Password, user.Gender)
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=users.csv")
		_, err := w.Write([]byte(csvData))
		if err != nil {
			http.Error(w, "Failed to write CSV response", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/createuser", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var newUser User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			http.Error(w, "Failed to insert user into database", http.StatusInternalServerError)
			return
		}

		successMsg := map[string]string{"message": "User successfully inserted into the database"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successMsg)
	})

	http.HandleFunc("/getuser/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		params := mux.Vars(r)
		id := params["id"]

		var user User
		result := db.First(&user, id)
		if result.Error != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("/updateuser/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPut {
			http.Error(w, "Only PUT requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		params := mux.Vars(r)
		id := params["id"]

		var user User
		result := db.First(&user, id)
		if result.Error != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		var updatedName struct {
			FullName string `json:"full_name"`
		}
		err := json.NewDecoder(r.Body).Decode(&updatedName)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		db.Model(&user).Update("full_name", updatedName.FullName)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User name updated successfully"})
	})

	http.HandleFunc("/deleteuser/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodDelete {
			http.Error(w, "Only DELETE requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		params := mux.Vars(r)
		id := params["id"]

		var user User
		result := db.First(&user, id)
		if result.Error != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		db.Delete(&user)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	})

	fmt.Println("Server running on http://localhost:5050")
	log.Fatal(http.ListenAndServe(":5050", nil))
}
