# Peeingdog Server

A simple Go server with PostgreSQL database integration and REST APIs for managing users and location-based messages.

## Features

- **User Management** - Create and manage users
- **Location-Based Messages** - Send messages tied to geographic locations
- **Message Expiration** - Messages expire after 24 hours and are automatically archived
- **Nearby Query** - Query messages by location and radius
- **Graceful Shutdown** - Handles SIGINT/SIGTERM signals properly
- **Error Group** - Concurrent goroutine management with error handling

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker & Docker Compose (for local database)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Start Database

```bash
docker-compose up -d postgres
```

Or using task:
```bash
task db:up
```

### 3. Configure Environment

Create a `.env` file (or set environment variables):

```
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/peeingdog?sslmode=disable
```

### 4. Run Server

```bash
go run main.go
```

Or using task:
```bash
task dev
```

The server will start on `http://localhost:8080`

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Messages Table (Active)
```sql
CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  text VARCHAR(180) NOT NULL,
  latitude DECIMAL(10, 8) NOT NULL,
  longitude DECIMAL(11, 8) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours')
);
```

### Archived Messages Table
```sql
CREATE TABLE archived_messages (
  id SERIAL PRIMARY KEY,
  original_message_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  text VARCHAR(180) NOT NULL,
  latitude DECIMAL(10, 8) NOT NULL,
  longitude DECIMAL(11, 8) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  expired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API Endpoints

### Health Check
```bash
GET /health
```

### User Management

#### List All Users
```bash
GET /api/users
```

#### Get Single User
```bash
GET /api/users/{id}
```

#### Create User
```bash
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

### Messages

#### Create Message
```bash
POST /api/users/{id}/messages
Content-Type: application/json

{
  "text": "Hello from NYC!",
  "latitude": 40.7128,
  "longitude": -74.0060
}
```

- Text: Max 180 characters
- Latitude: -90 to 90
- Longitude: -180 to 180

#### Get User Messages
```bash
GET /api/users/{id}/messages
```

Returns all active (non-expired) messages from a user.

#### Get Nearby Messages
```bash
GET /api/messages/nearby?lat=40.7128&lon=-74.0060&radius=1
```

Query parameters:
- `lat` - Latitude (required)
- `lon` - Longitude (required)
- `radius` - Search radius in kilometers (required)

Returns all active messages within the specified radius of the location.

## Message Lifecycle

1. **Created** - Message is created with 24-hour expiration time
2. **Active** - Message is visible in API responses for 24 hours
3. **Expired** - After 24 hours, message is no longer visible
4. **Archived** - Expired messages are moved to `archived_messages` table for audit purposes

Messages can be manually archived by calling:
```go
messageService.ArchiveExpiredMessages()
```

Or by triggering the database function directly (scheduled via application or pg_cron if enabled).

## Project Structure

```
peeingdog-server/
├── main.go           # Entry point with graceful shutdown
├── go.mod            # Go module file
├── config/           # Configuration
│   └── config.go
├── db/               # Database connection
│   └── db.go
├── service/          # Business logic
│   ├── user.go
│   └── message.go
├── handlers/         # HTTP handlers
│   └── handlers.go
├── Dockerfile.postgres   # PostgreSQL Docker image
├── docker-compose.yml    # Container orchestration
├── Taskfile.yml          # Task automation
├── schema.sql            # Database schema
├── .gitignore
└── README.md
```

## Useful Tasks

```bash
task dev              # Start database and run server
task dev:fresh        # Fresh start with clean database
task db:up            # Start PostgreSQL
task db:down          # Stop PostgreSQL
task db:clear         # Clear and reinitialize database
task db:reset         # Full reset
task deps             # Update Go dependencies
task build            # Build binary
task test             # Run tests
task clean            # Clean build artifacts
task help             # Show all available tasks
```

## Graceful Shutdown

The server listens for SIGINT (Ctrl+C) and SIGTERM signals:

```bash
# Start server
go run main.go

# Press Ctrl+C to gracefully shutdown
# The server will:
# 1. Stop accepting new connections
# 2. Wait up to 10 seconds for existing requests to finish
# 3. Close database connections
# 4. Exit cleanly
```

