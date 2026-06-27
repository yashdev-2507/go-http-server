# Go HTTP Server

A lightweight HTTP server built from scratch using only Go's standard library — 
no frameworks. Built to understand request handling, middleware patterns, and 
concurrency primitives at a low level.

## Endpoints
| Method | Path    | Description                        |
|--------|---------|------------------------------------|
| GET    | /hello  | Returns a plain text greeting      |
| GET    | /health | Returns server status as JSON      |
| POST   | /echo   | Echoes back the JSON request body  |
| POST   | /set    | Stores a key-value pair in memory  |
| GET    | /get    | Retrieves a value by key           |

## Implementation notes
- Middleware via closure/wrapper pattern — handlers are wrapped, not modified
- In-memory KV store protected with sync.RWMutex for concurrent read safety
- Structured JSON responses with correct HTTP status semantics (4xx for client errors)

## Run
go run main.go
