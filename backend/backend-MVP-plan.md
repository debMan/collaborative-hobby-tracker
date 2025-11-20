# Backend MVP Implementation Plan

## Tech Stack Confirmed

- **Database:** MongoDB
- **Framework:** Gin
- **Config:** koanf (YAML + env vars)
- **Logging:** zap
- **AI:** OpenAI-compatible (Ollama)
- **OAuth:** Google + GitHub
- **Testing:** TDD with unit tests (mocks) + integration tests (real DB)

## Implementation Phases

### Phase 1: Project Setup & Foundation (TDD Setup)

1. Initialize Go module and project structure
2. Set up configuration management (koanf)
3. Set up logging (zap)
4. Set up MongoDB connection and health check
5. Create base test utilities (testcontainers for integration tests)
6. Set up Makefile for common commands

### Phase 2: Domain Models & Repository Layer

**TDD approach: Write repository interface tests first**

1. Define domain models (User, Item, Category, Circle, Tag)
2. Create repository interfaces
3. Write repository tests (unit + integration)
4. Implement MongoDB repositories
5. Create indexes for performance

### Phase 3: Authentication System

**TDD: Write auth tests first**

1. JWT token generation/validation
2. Password hashing (bcrypt)
3. Auth service (register, login)
4. Auth middleware for protected routes
5. OAuth integration (Google, then GitHub)

### Phase 4: Core API - Items CRUD

**TDD: Write handler tests first**

1. Item service (business logic)
2. Item handlers (HTTP layer)
3. Routes setup with Gin
4. Authorization checks (user owns item)
5. Integration tests (end-to-end)

### Phase 5: Categories & Circles

**TDD: Tests first**

1. Category service + handlers
2. Circle service + handlers
3. Circle membership & access control
4. Sharing logic (invite links, email invites)

### Phase 6: AI Integration

**TDD: Mock AI responses in tests**

1. AI client interface (OpenAI-compatible)
2. Categorization service
3. Tag suggestion service
4. Configuration for different providers (Ollama)

### Phase 7: Content Fetchers

**TDD: Mock external HTTP calls**

1. Fetcher interface (strategy pattern)
2. YouTube scraper
3. Instagram scraper (basic)
4. Generic web scraper (Open Graph)
5. Fetcher registry

### Phase 8: Import Endpoint

**TDD: Integration tests**

1. Import service (orchestrates AI + fetchers)
2. Import handler
3. End-to-end import flow test

### Phase 9: Tags & Final Endpoints

1. Tag service + handlers
2. Tag learning logic
3. Final integration tests

### Phase 10: Polish & Documentation

1. Error handling improvements
2. API documentation
3. Docker setup
4. README with setup instructions

## Deliverables

- Complete backend API with all MVP endpoints
- MongoDB schema with indexes
- Comprehensive test suite (unit + integration)
- Docker Compose for local development
- Configuration examples
- API documentation
