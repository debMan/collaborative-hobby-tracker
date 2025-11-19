# Hobby Tracker Backend

Go backend API for the Collaborative Hobby Tracker application.

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: MongoDB
- **Configuration**: koanf (YAML + Environment Variables)
- **Logging**: zap (Uber's structured logger)
- **Authentication**: JWT + OAuth 2.0 (Google, GitHub)
- **AI Integration**: OpenAI-compatible API (Ollama, LocalAI, OpenAI)

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── api/              # HTTP handlers and routes
│   ├── domain/           # Business logic and models
│   ├── repository/       # Database access layer
│   ├── ai/               # AI service integration
│   ├── fetchers/         # Content fetchers (YouTube, Instagram, etc.)
│   └── auth/             # Authentication utilities
├── pkg/                  # Reusable packages
│   ├── logger/           # Logging utilities
│   ├── validator/        # Request validation
│   └── errors/           # Custom error types
├── config/               # Configuration files
└── migrations/           # Database migrations (indexes)
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB 7.0 or higher
- Docker & Docker Compose (for local development)
- Ollama (for AI features) or access to OpenAI API

### Installation

1. **Clone the repository** (if not already cloned)

2. **Install dependencies**:

   ```bash
   make install
   ```

3. **Set up configuration**:

   ```bash
   cp config.example.yaml config.yaml
   ```

   Edit `config.yaml` with your settings.

4. **Set environment variables** (optional, overrides config.yaml):

   ```bash
   export HT_AUTH_JWT_SECRET="your-secret-key"
   export HT_DATABASE_URI="mongodb://localhost:27017"
   export HT_AI_BASE_URL="http://localhost:11434"
   # ... more variables as needed
   ```

5. **Start MongoDB**:

   ```bash
   make docker-up
   # or just MongoDB:
   make mongo-up
   ```

6. **Run the application**:

   ```bash
   make run
   ```

   The API will be available at `http://localhost:8080`

### Development

For hot-reload during development:

```bash
make dev
```

This uses [Air](https://github.com/cosmtrek/air) for automatic reloading.

## Configuration

Configuration is loaded from `config.yaml` and can be overridden with environment variables.

Environment variables use the `HT_` prefix and follow the structure:

```
HT_<SECTION>_<KEY>

Example:
HT_SERVER_PORT=8080
HT_DATABASE_URI=mongodb://localhost:27017
HT_AUTH_JWT_SECRET=my-secret
```

See `config.example.yaml` for all available options.

## API Endpoints

### Health Check

- `GET /health` - Server health check

### API v1

- Base URL: `/api/v1`

#### Authentication (Coming Soon)

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user
- `GET /api/v1/auth/oauth/{provider}` - OAuth login (Google, GitHub)

#### Items (Coming Soon)

- `GET /api/v1/items` - List items
- `POST /api/v1/items` - Create item
- `GET /api/v1/items/:id` - Get item
- `PATCH /api/v1/items/:id` - Update item
- `DELETE /api/v1/items/:id` - Delete item

#### Categories, Circles, Tags, Import (Coming Soon)

## Testing

Run all tests:

```bash
make test
```

Run unit tests only:

```bash
make test-unit
```

Run integration tests only:

```bash
make test-integration
```

## Development Commands

```bash
make help              # Show available commands
make install           # Install dependencies
make build             # Build application
make run               # Run application
make dev               # Run with hot reload
make test              # Run all tests
make test-unit         # Run unit tests
make test-integration  # Run integration tests
make clean             # Clean build artifacts
make docker-up         # Start Docker services
make docker-down       # Stop Docker services
make lint              # Run linter
make fmt               # Format code
```

## Docker

Start all services (MongoDB):

```bash
docker-compose up -d
```

Stop services:

```bash
docker-compose down
```

## OAuth Setup

### Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/api/v1/auth/oauth/google/callback`
6. Add client ID and secret to config

### GitHub OAuth

1. Go to GitHub Settings > Developer settings > OAuth Apps
2. Create a new OAuth App
3. Set authorization callback URL: `http://localhost:8080/api/v1/auth/oauth/github/callback`
4. Add client ID and secret to config

## AI Setup (Ollama)

1. Install Ollama from [ollama.ai](https://ollama.ai)
2. Pull a model:

   ```bash
   ollama pull llama3:8b
   ```

3. Start Ollama (usually runs automatically)
4. Configure in `config.yaml`:

   ```yaml
   ai:
     provider: "ollama"
     base_url: "http://localhost:11434"
     model: "llama3:8b"
   ```

## License

See parent project LICENSE file.
