package llama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiURL string
	apiKey string
}

type AnalysisRequest struct {
	Prompt           string `json:"prompt"`
	GigaChatResponse string `json:"gigachat_response"`
}

type AnalysisResponse struct {
	IsCorrect bool   `json:"is_correct"`
	Feedback  string `json:"feedback"`
	Analysis  string `json:"analysis"`
}

type TrainingRequest struct {
	Prompt           string `json:"prompt"`
	GigaChatResponse string `json:"gigachat_response"`
	IsCorrect        bool   `json:"is_correct"`
	Feedback         string `json:"feedback"`
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func NewClient(apiURL, apiKey string) *Client {
	return &Client{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

func (c *Client) AnalyzeResponse(prompt, response string) (*AnalysisResult, error) {
	analysisPrompt := fmt.Sprintf("Analyze the following response to the prompt '%s':\n\nResponse: %s\n\nIs this response correct and helpful? Provide feedback.", prompt, response)

	result, err := c.SimpleChat(analysisPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze response: %v", err)
	}

	// Simple analysis logic - can be improved
	isCorrect := true // This should be determined based on the analysis
	return &AnalysisResult{
		Analysis:  result,
		IsCorrect: isCorrect,
		Feedback:  "Response analyzed successfully",
	}, nil
}

func (c *Client) TrainModel(prompt, response string, isCorrect bool, feedback string) error {
	trainingPrompt := fmt.Sprintf("Learn from this interaction:\nPrompt: %s\nResponse: %s\nCorrect: %v\nFeedback: %s",
		prompt, response, isCorrect, feedback)

	_, err := c.SimpleChat(trainingPrompt)
	if err != nil {
		return fmt.Errorf("failed to train model: %v", err)
	}

	return nil
}

func (c *Client) SimpleChat(message string) (string, error) {
	reqBody := GenerateRequest{
		Model:  "llama2",
		Prompt: message,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Response, nil
}

type AnalysisResult struct {
	Analysis  string
	IsCorrect bool
	Feedback  string
}
