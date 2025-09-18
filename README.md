# RESTful API Project with Go Fiber

A simple RESTful API built with Go Fiber that returns a "hello world" JSON response.

## Features

- Built with [Go Fiber](https://gofiber.io/) framework
- Returns JSON response: `{"message": "hello world"}` on GET /
- Runs on port 3000

## Prerequisites

- Go 1.17 or higher
- Go Fiber v2

## Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd /Users/ekaphon.m/Documents/Workspace/Temp_AI_RestAPI
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Run the server

```bash
go run main.go
```

The server will start on port 3000.

### Test the API

Open your browser or use curl to test the endpoint:

```bash
curl http://localhost:3000/
```

Response:
```json
{"message": "hello world"}
```

## Build

To build the executable:

```bash
go build -o server main.go
```

Then run:
```bash
./server
```

## Project Structure

```
.
├── main.go       # Main application file
├── go.mod        # Go module file
├── go.sum        # Go dependencies checksums
└── README.md     # This file
```

## API Endpoints

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| GET    | /        | Get hello world message | `{"message": "hello world"}` |
