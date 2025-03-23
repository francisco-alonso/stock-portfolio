package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/caarlos0/env/v11"
	"github.com/francisco-alonso/stock-portfolio/user-service/internal/adapters/api"
	postgres "github.com/francisco-alonso/stock-portfolio/user-service/internal/adapters/pg"
	"github.com/francisco-alonso/stock-portfolio/user-service/internal/services"
	_ "github.com/lib/pq"
)

type Config struct {
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
    
    encodedPassword := url.QueryEscape(cfg.DBPassword)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, encodedPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := postgres.NewPostgresUserRepository(db)
	service := services.NewUserService(repo)
	handler := api.NewUserHandler(service)

	http.HandleFunc("/create-user", handler.CreateUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
