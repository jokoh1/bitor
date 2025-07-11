#!/bin/bash
set -e

echo "🚀 Starting build process..."

# Check if we're in the backend directory
if [ ! -f "main.go" ]; then
    echo "❌ Error: This script must be run from the backend directory"
    exit 1
fi

# Check if frontend directory exists
if [ ! -d "../frontend" ]; then
    echo "❌ Error: Frontend directory not found at ../frontend"
    exit 1
fi

echo "📦 Building frontend..."
cd ../frontend

# Install dependencies
if [ -f "pnpm-lock.yaml" ]; then
    echo "Installing dependencies with pnpm..."
    pnpm install --frozen-lockfile
else
    echo "Installing dependencies with npm..."
    npm ci
fi

# Build frontend
echo "Building frontend..."
npm run build

# Create pb_public directory in backend
echo "📁 Creating pb_public directory..."
mkdir -p ../backend/pb_public

# Copy build output to pb_public
echo "📋 Copying build files to pb_public..."
cp -r build/* ../backend/pb_public/

# Go back to backend directory
cd ../backend

echo "🔧 Building backend..."
# Build the Go application
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bitor main.go

echo "✅ Build completed successfully!" 