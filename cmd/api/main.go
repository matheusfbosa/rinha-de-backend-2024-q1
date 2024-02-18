package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/customer"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/customer/postgres"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/internal/httpserver/fiber"
)

const (
	defaultConnString = "postgres://postgres:password@db:5432/rinha?sslmode=disable"
	defaultServerPort = ":3000"
)

func main() {
	connString := os.Getenv("DB_URL")
	if connString == "" {
		connString = defaultConnString
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	r := postgres.NewPostgreSQL(db)
	service := customer.NewService(r)

	app := fiber.Handlers(service)
	app.Listen(serverPort)
}
