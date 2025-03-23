package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/francisco-alonso/stock-portfolio/portfolio-service/api"
	"github.com/francisco-alonso/stock-portfolio/portfolio-service/infra"
	"github.com/francisco-alonso/stock-portfolio/portfolio-service/services"
)

func main() {
    ctx := context.Background()

    // Obtener datos de entorno
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbName := os.Getenv("DB_NAME")
    secretName := os.Getenv("DB_SECRET_NAME")

    if dbHost == "" || dbUser == "" || dbName == "" || secretName == "" {
        log.Fatal("Environment variables DB_HOST, DB_USER, DB_NAME, and DB_SECRET_NAME must be set")
    }

    password, err := infra.GetSecret(ctx, secretName)
    if err != nil {
        log.Fatalf("Failed to access secret: %v", err)
    }

    database, err := infra.NewPostgresDB(dbHost, dbUser, password, dbName)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    portfolioService := services.NewPortfolioService(database)
    handler := api.NewHandler(portfolioService)
    
    http.HandleFunc("/portfolio", handler.GetPortfolio)
    http.HandleFunc("/position", handler.AddPosition)

    log.Println("Portfolio Service is running on port 8082...")
    log.Fatal(http.ListenAndServe(":8082", nil))
}
