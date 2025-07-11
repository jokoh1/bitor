# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy frontend files
COPY frontend/package*.json ./
COPY frontend/pnpm-lock.yaml ./

# Install pnpm
RUN npm install -g pnpm

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN pnpm build

# Build stage for backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Copy frontend build output
COPY --from=frontend-builder /app/build ./pb_public

# Build the application (CGO disabled for Railway compatibility)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bitor main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=backend-builder /app/bitor .

# Create necessary directories
RUN mkdir -p pb_data ansible nuclei-templates

# Expose port
EXPOSE 8090

# Run the application
CMD ["./bitor", "serve", "--http", "0.0.0.0:8090"]
