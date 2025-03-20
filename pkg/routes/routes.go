package routes

import (
	"net/http"

	"llama3-server/pkg/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(handler *handlers.CompletionHandler) {
	http.HandleFunc("/v1/completions", handler.HandleCompletions)
	http.HandleFunc("/v1/analyze", handler.HandleAnalysis)
	http.HandleFunc("/v1/feedback", handler.HandleFeedback)
	http.HandleFunc("/v1/test-llama", handler.HandleTestLlama)
}
