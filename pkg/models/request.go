package models

// CompletionRequest represents the request body for the completions endpoint
type CompletionRequest struct {
	Prompt string `json:"prompt"`
}

// AnalysisRequest represents the request body for the analysis endpoint
type AnalysisRequest struct {
	ID string `json:"id"`
}

// FeedbackRequest represents the request body for the feedback endpoint
type FeedbackRequest struct {
	ID        string `json:"id"`
	IsCorrect bool   `json:"is_correct"`
	Feedback  string `json:"feedback"`
}

// TestLlamaRequest represents the request body for the test-llama endpoint
type TestLlamaRequest struct {
	Message string `json:"message"`
}

// CompletionResponse represents the response for the completions endpoint
type CompletionResponse struct {
	Text     string `json:"text"`
	Analysis struct {
		IsCorrect bool   `json:"is_correct"`
		Feedback  string `json:"feedback"`
	} `json:"analysis"`
}

// AnalysisResponse represents the response for the analysis endpoint
type AnalysisResponse struct {
	Analysis string `json:"analysis"`
}

// TestLlamaResponse represents the response for the test-llama endpoint
type TestLlamaResponse struct {
	Response string `json:"response"`
}
