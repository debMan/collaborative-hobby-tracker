#!/bin/bash

# Setup script for Hobby Tracker Backend

set -e

echo "🚀 Setting up Hobby Tracker Backend..."
echo ""

# Create config directory if it doesn't exist
mkdir -p config

# Copy example config if config.yaml doesn't exist
if [ ! -f "config/config.yaml" ]; then
    echo "📝 Creating config.yaml from example..."
    cp config/config.example.yaml config/config.yaml
    echo "✅ config.yaml created"
else
    echo "⚠️  config.yaml already exists, skipping..."
fi

# Check if MongoDB is running
echo ""
echo "🔍 Checking MongoDB connection..."
if ! command -v mongosh &> /dev/null && ! command -v mongo &> /dev/null; then
    echo "⚠️  MongoDB client not found. Trying to connect anyway..."
fi

# Try to connect to MongoDB
if nc -z localhost 27017 2>/dev/null; then
    echo "✅ MongoDB is running on localhost:27017"
else
    echo ""
    echo "⚠️  MongoDB is not running on localhost:27017"
    echo ""
    echo "To start MongoDB with Docker:"
    echo "  make docker-up"
    echo ""
    echo "Or install MongoDB locally:"
    echo "  https://www.mongodb.com/docs/manual/installation/"
    echo ""
fi

# Install Go dependencies
echo ""
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy
echo "✅ Dependencies installed"

# Build the application
echo ""
echo "🔨 Building application..."
make build
echo "✅ Application built successfully"

echo ""
echo "✨ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Start MongoDB: make docker-up"
echo "2. Edit config/config.yaml with your settings"
echo "3. Run the application: make run"
echo "4. Visit http://localhost:8080/health"
echo ""
echo "For development with hot reload: make dev"
echo ""
