package main

import (
	"log"
	"net/http"
	"os"

	"github.com/francisco-alonso/trade-engine-service/internal/adapter/api"
	"github.com/francisco-alonso/trade-engine-service/internal/application"
	"github.com/francisco-alonso/trade-engine-service/internal/infra/pubsub"
)

func main() {
    projectId := os.Getenv("PROJECT_ID")
    topicId := os.Getenv("PUB_SUB_TOPIC")
    
    if projectId == "" || topicId == "" {
        log.Fatalf("Env variables projectId aor topicId are missing.")
    }
    
    publisher, err := pubsub.NewPublisher(projectId, topicId)
    if err != nil {
        log.Fatalf("Failed to create publisher: %v", err)
    }
    
    tradeService := application.NewTradeService(publisher)
    handler := api.NewHandler(tradeService)

    log.Println("Trade Engine Service is running on port 8081...")
    log.Fatal(http.ListenAndServe(":8081", handler.Router()))
}
