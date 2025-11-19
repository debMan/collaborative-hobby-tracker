# Collaborative Hobby Tracker

A full-stack web application for tracking and organizing hobbies with AI-powered categorization and collaborative features.

## 🎯 Overview

This application helps users track and organize their hobbies with intelligent features:

- **AI-Powered Categorization**: Automatically categorize items (movies, restaurants, travel, activities, etc.) using AI
- **Multi-Source Import**: Import from Instagram, YouTube, Twitter/X, TikTok, Telegram, Wikipedia, or plain text
- **Smart Tagging**: AI learns your tagging preferences and suggests tags for new items
- **Collaborative Circles**: Share categories with different groups (Partner, Friends, Family, etc.)
- **Access Control**: Fine-grained permissions (private, view, edit, admin)
- **Calendar Integration**: Sync with Google Calendar and Apple Calendar for planning
- **Metadata Enrichment**: Fetch details from IMDB, Google Maps, and other sources

## 🏗️ Architecture

This is a monorepo containing both frontend and backend:

```
collaborative-hobby-tracker/
├── frontend/          # React + TypeScript web application
│   ├── src/
│   ├── public/
│   └── package.json
├── backend/           # Go API server
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── go.mod
├── .gitignore
├── LICENSE
└── README.md          # This file
```

## 🚀 Quick Start

### Prerequisites

- **Frontend**: Node.js 18+ and npm
- **Backend**: Go 1.21+, MongoDB 7.0+
- **AI Features**: Ollama (or access to OpenAI API)

### 1. Start MongoDB

```bash
# Using Docker (easiest)
cd backend
make docker-up

# Or install MongoDB locally
# https://www.mongodb.com/docs/manual/installation/
```

### 2. Start Backend API

```bash
cd backend

# First time setup
cp config.example.yaml config.yaml
# Edit config.yaml with your settings (MongoDB URI, JWT secret, etc.)

# Install dependencies and run
make install
make run

# For development with hot reload
make dev
```

Backend will be available at **<http://localhost:8080>**

Test the API:

```bash
curl http://localhost:8080/health
```

### 3. Start Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will be available at **<http://localhost:5173>**

## 📚 Documentation

- **[Frontend README](frontend/README.md)** - React application details, components, state management
- **[Backend README](backend/README.md)** - API documentation, deployment, testing

## 🛠️ Tech Stack

### Frontend

- **Framework**: React 18
- **Language**: TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Routing**: React Router v6
- **Icons**: Lucide React

### Backend

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: MongoDB
- **Configuration**: koanf (YAML + environment variables)
- **Logging**: zap (structured logging)
- **Authentication**: JWT + OAuth 2.0 (Google, GitHub)
- **AI**: OpenAI-compatible API (Ollama, LocalAI, OpenAI)

## 🧪 Development

### Frontend Commands

```bash
cd frontend

npm run dev          # Start dev server (http://localhost:5173)
npm run build        # Build for production
npm run preview      # Preview production build
npm run lint         # Run ESLint
```

### Backend Commands

```bash
cd backend

make dev             # Run with hot reload
make build           # Build binary
make run             # Run built binary
make test            # Run all tests
make test-unit       # Run unit tests only
make test-integration # Run integration tests
make lint            # Run linter
make fmt             # Format code
make docker-up       # Start MongoDB
make docker-down     # Stop MongoDB
```

## 🌟 Features

### Implemented (MVP)

- ✅ Clean, distraction-free UI
- ✅ User authentication (email + password)
- ✅ OAuth login (Google, GitHub)
- ✅ Item management (create, read, update, delete)
- ✅ Categories and organization
- ✅ Circles for sharing
- ✅ Basic access control
- ✅ MongoDB persistence
- ✅ Health monitoring

### In Progress

- 🔄 AI-powered categorization
- 🔄 Content fetching from external sources
- 🔄 Tag learning and suggestions
- 🔄 Metadata enrichment (IMDB, Google Maps)

### Planned

- 📅 Calendar integration (Google, Apple)
- 🔔 Notifications and reminders
- 🔄 Real-time collaboration
- 📊 Analytics and insights
- 📱 Mobile apps (iOS, Android)
- 🌐 Internationalization (i18n)

## 🎨 Design Philosophy

- **Simplicity First**: Clean interface without visual clutter
- **Content-First**: Focus on your hobbies, not decorative UI
- **List-Based**: Traditional list view instead of cards
- **Progressive Disclosure**: Advanced features available but not overwhelming
- **Material Design**: Familiar and consistent UI patterns

## 📄 License

See [LICENSE](LICENSE) file for details.

## 🤝 Contributing

This is a personal learning project. Suggestions and feedback are welcome via GitHub issues!

## 💬 Support

For questions or issues:

- Open an issue on GitHub
- Check the documentation in `frontend/` and `backend/` directories

---

Built with ❤️ as a learning project to practice full-stack development, Go, React, and cloud deployment.
