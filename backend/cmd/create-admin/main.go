package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	host := os.Getenv("USER_DB_HOST")
	port := os.Getenv("USER_DB_PORT")
	user := os.Getenv("USER_DB_USER")
	password := os.Getenv("USER_DB_PASSWORD")
	dbname := os.Getenv("USER_DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5434" // Default port mapped in dev
	}
	// If run from host against docker mapped port, user needs to check docker-compose.dev.yml
	// user_db maps 5434:5432
	port = "5434"

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	email := "sanin.tony.dev@gmail.com"
	pass := "iamtheroot"
	name := "Admin"

	// Check if exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking user: %v", err)
	}
	if exists {
		log.Println("User already exists")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	id := uuid.New()
	role := "admin"

	query := `INSERT INTO users (id, name, email, password, role_id) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.ExecContext(context.Background(), query, id, name, email, string(hashed), role)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}

	fmt.Printf("Admin user created: %s / %s\n", email, pass)
}
