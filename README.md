# Learning API

This project is a Go REST API using Gin (web framework), GORM (ORM), and MySQL. It supports development and production profiles, and provides CRUD APIs for topics, questions, and answers, as well as an endpoint to fetch all questions with answers for a specific topic.

## Features
- Gin web framework
- GORM ORM with MySQL
- Dev/Production profiles
- Topics, Questions, Answers tables (with relationships)
- CRUD APIs for all tables
- Endpoint: Get all questions with answers for a topic
- Unit tests for APIs

## Getting Started
1. Install Go and MySQL
2. Clone this repo
3. Configure your database connection in `config.yaml` or environment variables
4. Run: `go run main.go`

## Testing
Run unit tests with:
```
go test ./...
```
