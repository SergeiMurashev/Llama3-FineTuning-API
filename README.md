# Llama 3 Server

A Go-based HTTP server that provides an API interface for Llama 3 model interactions.

## Prerequisites

- Go 1.16 or later
- PostgreSQL 14 or later
- Llama 3 model files (to be configured)

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod tidy
```

3. Initialize the database:
```bash
./scripts/init_db.sh
```

## Configuration

The server configuration is stored in `configs.json`. Make sure to update the following settings according to your environment:

- GigaChat API URL and key
- Llama API URL
- Database connection details

## Running the Server

You can run the server in two ways:

1. Using make (recommended):
```bash
make run
```

2. Directly with go:
```bash
go run cmd/server/main.go
```

The server will start on the configured host and port (default: localhost:8080).

## API Usage

### Completions Endpoint

Send a POST request to `/v1/completions` with a JSON body:

```json
{
    "prompt": "Your prompt here"
}
```

Example using curl:
```bash
curl -X POST http://localhost:8080/v1/completions \
     -H "Content-Type: application/json" \
     -d '{"prompt": "What is the capital of France?"}'
```

### Test Llama Endpoint

Test direct communication with Llama:

```bash
curl -X POST http://localhost:8080/v1/test-llama \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello, how are you?"}'
```

### Analysis Endpoint

Get analysis of a specific response:

```bash
curl -X POST http://localhost:8080/v1/analyze \
     -H "Content-Type: application/json" \
     -d '{"id": "ID_FROM_PREVIOUS_RESPONSE"}'
```

### Feedback Endpoint

Provide feedback for model training:

```bash
curl -X POST http://localhost:8080/v1/feedback \
     -H "Content-Type: application/json" \
     -d '{
       "id": "ID_FROM_PREVIOUS_RESPONSE",
       "is_correct": true,
       "feedback": "The response was accurate"
     }'
```

## Project Structure

```
Llama3-FineTuning-API/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   └── config/
│       └── config.go
├── db/
│   └── db.go
├── gigachat/
│   └── client.go
├── llama/
│   └── client.go
├── configs.json
├── go.mod
└── README.md
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
