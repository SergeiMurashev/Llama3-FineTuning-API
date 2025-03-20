package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseEntry struct {
	ID               string    `json:"id"`
	Prompt           string    `json:"prompt"`
	GigaChatResponse string    `json:"gigachat_response"`
	LlamaResponse    string    `json:"llama_response"`
	IsCorrect        bool      `json:"is_correct"`
	Feedback         string    `json:"feedback"`
	Timestamp        time.Time `json:"timestamp"`
}

var db *sql.DB

func InitDB(dbURL string) error {
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS interactions (
			id TEXT PRIMARY KEY,
			prompt TEXT NOT NULL,
			gigachat_response TEXT NOT NULL,
			llama_response TEXT NOT NULL,
			is_correct BOOLEAN DEFAULT false,
			feedback TEXT,
			timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func StoreEntry(entry DatabaseEntry) error {
	_, err := db.Exec(`
		INSERT INTO interactions (id, prompt, gigachat_response, llama_response, is_correct, feedback, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		entry.ID,
		entry.Prompt,
		entry.GigaChatResponse,
		entry.LlamaResponse,
		entry.IsCorrect,
		entry.Feedback,
		entry.Timestamp,
	)
	return err
}

func GetEntry(id string) (DatabaseEntry, error) {
	var entry DatabaseEntry
	err := db.QueryRow(`
		SELECT id, prompt, gigachat_response, llama_response, is_correct, feedback, timestamp
		FROM interactions
		WHERE id = $1
	`,
		id,
	).Scan(
		&entry.ID,
		&entry.Prompt,
		&entry.GigaChatResponse,
		&entry.LlamaResponse,
		&entry.IsCorrect,
		&entry.Feedback,
		&entry.Timestamp,
	)
	if err != nil {
		return DatabaseEntry{}, err
	}
	return entry, nil
}

func UpdateFeedback(id string, isCorrect bool, feedback string) error {
	_, err := db.Exec(`
		UPDATE interactions
		SET is_correct = $1, feedback = $2
		WHERE id = $3
	`,
		isCorrect,
		feedback,
		id,
	)
	return err
}

func GetIncorrectResponses() ([]DatabaseEntry, error) {
	rows, err := db.Query(`
		SELECT id, prompt, gigachat_response, llama_response, is_correct, feedback, timestamp
		FROM interactions
		WHERE is_correct = false
		ORDER BY timestamp DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []DatabaseEntry
	for rows.Next() {
		var entry DatabaseEntry
		err := rows.Scan(
			&entry.ID,
			&entry.Prompt,
			&entry.GigaChatResponse,
			&entry.LlamaResponse,
			&entry.IsCorrect,
			&entry.Feedback,
			&entry.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
