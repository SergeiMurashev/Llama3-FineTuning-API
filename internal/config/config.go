package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GigaChat struct {
		APIURL       string `json:"api_url"`
		APIScope     string `json:"api_scope"`
		Model        string `json:"model"`
		ClientID     string `json:"client_id"`
		RQUID        string `json:"rq_uid"`
		AuthKey      string `json:"-"`
		ClientSecret string `json:"-"`
		SessionID    string `json:"-"`
	} `json:"gigachat"`
	Llama struct {
		APIURL string `json:"api_url"`
		APIKey string `json:"-"`
	} `json:"llama"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		SSLMode  string `json:"sslmode"`
	} `json:"database"`
}

func LoadConfig(configPath string) (*Config, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Read config file
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Load sensitive data from environment variables
	cfg.GigaChat.AuthKey = os.Getenv("GIGACHAT_AUTH_KEY")
	cfg.GigaChat.ClientSecret = os.Getenv("GIGACHAT_CLIENT_SECRET")
	cfg.GigaChat.SessionID = os.Getenv("GIGACHAT_SESSION_ID")
	cfg.Llama.APIKey = os.Getenv("LLAMA_API_KEY")

	return &cfg, nil
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
