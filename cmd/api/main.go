package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/customer"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/customer/postgres"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/internal/httpserver/fiber"
)

const (
	defaultConnString = "postgres://postgres:password@localhost:5432/rinha?sslmode=disable&pool_min_conns=0&pool_max_conns=15"
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

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	r := postgres.NewPostgreSQL(dbpool)
	service := customer.NewService(r)

	app := fiber.Handlers(service)
	app.Listen(serverPort)
}
