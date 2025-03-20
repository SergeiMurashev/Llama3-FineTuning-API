package main

import (
	"fmt"
	"log"
	"net/http"

	"llama3-server/db"
	"llama3-server/gigachat"
	"llama3-server/internal/config"
	"llama3-server/llama"
	"llama3-server/pkg/handlers"
	"llama3-server/pkg/routes"
	"llama3-server/pkg/services"
)

func main() {
	// Load configuration
	configPath := "configs.json"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize clients
	gigaChatClient := gigachat.NewClient(
		cfg.GigaChat.APIURL,
		cfg.GigaChat.AuthKey,
		cfg.GigaChat.ClientID,
		cfg.GigaChat.RQUID,
		cfg.GigaChat.APIScope,
		cfg.GigaChat.Model,
	)
	llamaClient := llama.NewClient(cfg.Llama.APIURL, cfg.Llama.APIKey)

	// Initialize database
	if err := db.InitDB(cfg.DatabaseDSN()); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize service and handler
	completionService := services.NewCompletionService(gigaChatClient, llamaClient)
	completionHandler := handlers.NewCompletionHandler(completionService)

	// Setup routes
	routes.SetupRoutes(completionHandler)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	log.Printf("GigaChat API URL: %s", cfg.GigaChat.APIURL)
	log.Printf("Llama API URL: %s", cfg.Llama.APIURL)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
