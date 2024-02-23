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
	defaultConnString = "postgres://postgres:password@localhost:5432/rinha?pool_max_conns=6&pool_min_conns=6&pool_max_conn_lifetime=330s"
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

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(fmt.Sprintf("failed to parse connection string: %s", err.Error()))
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Sprintf("unable to create connection pool: %s", err.Error()))
	}
	defer dbpool.Close()

	r := postgres.NewPostgreSQL(dbpool)
	service := customer.NewService(r)

	app := fiber.Handlers(service)
	app.Listen(serverPort)
}
