package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// User структура для хранения данных о пользователе
type User struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Gender      string `json:"gender"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1959"
	dbname   = "postgres"
)

func main() {
	// Строка подключения к базе данных PostgreSQL
	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Подключение к базе данных
	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	fmt.Println("Connected to the database")

	// Обработчик POST запросов
	http.HandleFunc("/uploadjson", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Если это предзапрос OPTIONS, просто отправляем пустой ответ с корректными заголовками
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Проверка метода запроса - должен быть POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// Декодирование JSON из тела запроса
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		// Преобразование JSON данных в нужный формат (User в данном случае)
		// Здесь нужно использовать вашу структуру данных, соответствующую вашим данным
		var newUser User
		// Пример преобразования JSON в структуру User
		newUser.FullName = data["full_name"].(string)
		newUser.Username = data["username"].(string)
		newUser.Email = data["email"].(string)
		newUser.PhoneNumber = data["phone_number"].(string)
		newUser.Password = data["password"].(string)
		newUser.Gender = data["gender"].(string)

		// SQL запрос для вставки данных нового пользователя в базу данных
		sqlStatement := `
		INSERT INTO users (full_name, username, email, phone_number, password, gender)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

		// Выполнение SQL запроса
		var userID int
		err = db.QueryRow(sqlStatement, newUser.FullName, newUser.Username, newUser.Email, newUser.PhoneNumber, newUser.Password, newUser.Gender).Scan(&userID)
		if err != nil {
			http.Error(w, "Failed to insert user into database", http.StatusInternalServerError)
			return
		}

		// Ответ об успешной вставке пользователя в базу данных
		successMsg := map[string]string{"message": "User successfully inserted into the database"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successMsg)
	})

	http.HandleFunc("/getusers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Проверка метода запроса - должен быть GET
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// SQL запрос для выборки всех пользователей из базы данных
		sqlStatement := `SELECT * FROM users`

		// Выполнение SQL запроса
		rows, err := db.Query(sqlStatement)
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Создаем CSV-заголовок
		csvData := "ID,Full Name,Username,Email,Phone Number,Password,Gender\n"

		// Формирование CSV данных
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.PhoneNumber, &user.Password, &user.Gender)
			if err != nil {
				http.Error(w, "Failed to scan user row", http.StatusInternalServerError)
				return
			}
			csvData += fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s\n", user.ID, user.FullName, user.Username, user.Email, user.PhoneNumber, user.Password, user.Gender)
		}

		// Отправляем CSV в ответе
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=users.csv")
		_, err = w.Write([]byte(csvData))
		if err != nil {
			http.Error(w, "Failed to write CSV response", http.StatusInternalServerError)
			return
		}
	})

	// Запуск сервера на порту 8080
	fmt.Println("Server running on http://localhost:5050")
	log.Fatal(http.ListenAndServe(":5050", nil))
}
