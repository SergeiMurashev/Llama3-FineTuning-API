package services

import (
	"fmt"
	"time"

	"llama3-server/db"
	"llama3-server/gigachat"
	"llama3-server/llama"
	"llama3-server/pkg/models"
)

type CompletionService struct {
	gigaChatClient *gigachat.Client
	llamaClient    *llama.Client
}

func NewCompletionService(gigaChatClient *gigachat.Client, llamaClient *llama.Client) *CompletionService {
	return &CompletionService{
		gigaChatClient: gigaChatClient,
		llamaClient:    llamaClient,
	}
}

func (s *CompletionService) GetCompletion(prompt string) (*models.CompletionResponse, error) {
	// Get response from GigaChat
	gigaChatResp, err := s.gigaChatClient.GetCompletion(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get GigaChat response: %v", err)
	}

	// Get Llama's analysis
	analysis, err := s.llamaClient.AnalyzeResponse(prompt, gigaChatResp)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze with Llama: %v", err)
	}

	response := &models.CompletionResponse{
		Text: gigaChatResp,
	}
	response.Analysis.IsCorrect = analysis.IsCorrect
	response.Analysis.Feedback = analysis.Feedback

	// Store in database
	entry := db.DatabaseEntry{
		ID:               generateID(),
		Prompt:           prompt,
		GigaChatResponse: gigaChatResp,
		LlamaResponse:    analysis.Analysis,
		IsCorrect:        analysis.IsCorrect,
		Feedback:         analysis.Feedback,
	}
	if err := db.StoreEntry(entry); err != nil {
		return nil, fmt.Errorf("failed to store in database: %v", err)
	}

	return response, nil
}

func (s *CompletionService) GetAnalysis(id string) (*models.AnalysisResponse, error) {
	// Get entry from database
	entry, err := db.GetEntry(id)
	if err != nil {
		return nil, fmt.Errorf("entry not found: %v", err)
	}

	// Perform Llama analysis
	analysis, err := s.llamaClient.AnalyzeResponse(entry.Prompt, entry.GigaChatResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze: %v", err)
	}

	return &models.AnalysisResponse{
		Analysis: analysis.Analysis,
	}, nil
}

func (s *CompletionService) UpdateFeedback(id string, isCorrect bool, feedback string) error {
	// Update database entry with feedback
	if err := db.UpdateFeedback(id, isCorrect, feedback); err != nil {
		return fmt.Errorf("failed to update feedback: %v", err)
	}

	// Get entry for training
	entry, err := db.GetEntry(id)
	if err != nil {
		return fmt.Errorf("entry not found: %v", err)
	}

	// Train Llama model with feedback
	if err := s.llamaClient.TrainModel(entry.Prompt, entry.GigaChatResponse, isCorrect, feedback); err != nil {
		return fmt.Errorf("failed to train model: %v", err)
	}

	return nil
}

func (s *CompletionService) TestLlama(message string) (*models.TestLlamaResponse, error) {
	// Get response from Llama
	response, err := s.llamaClient.SimpleChat(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get Llama response: %v", err)
	}

	return &models.TestLlamaResponse{
		Response: response,
	}, nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
