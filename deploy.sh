#!/bin/bash

# Bitor Deployment Script
# This script helps you deploy the Bitor application

set -e

echo "🚀 Bitor Deployment Script"
echo "=========================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "frontend/package.json" ] || [ ! -f "backend/go.mod" ]; then
    print_error "Please run this script from the root directory of the Bitor project"
    exit 1
fi

# Function to test frontend build
test_frontend_build() {
    print_status "Testing frontend build..."
    cd frontend
    
    if ! command -v pnpm &> /dev/null; then
        print_error "pnpm is not installed. Please install it first: npm install -g pnpm"
        exit 1
    fi
    
    pnpm install
    pnpm build
    
    if [ $? -eq 0 ]; then
        print_success "Frontend build successful!"
    else
        print_error "Frontend build failed!"
        exit 1
    fi
    
    cd ..
}

# Function to test backend build
test_backend_build() {
    print_status "Testing backend build..."
    cd backend
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install it first."
        exit 1
    fi
    
    go mod download
    go build -o bitor main.go
    
    if [ $? -eq 0 ]; then
        print_success "Backend build successful!"
    else
        print_error "Backend build failed!"
        exit 1
    fi
    
    cd ..
}

# Function to show deployment options
show_deployment_options() {
    echo ""
    echo "📋 Deployment Options:"
    echo "======================"
    echo "1. Frontend (Cloudflare Pages):"
    echo "   - Manual: Go to https://pages.cloudflare.com/"
    echo "   - Automated: Set up GitHub Actions secrets"
    echo ""
    echo "2. Backend (Railway):"
    echo "   - Manual: Go to https://railway.app/"
    echo "   - Automated: Set up Railway CLI and GitHub Actions"
    echo ""
    echo "3. Both (Recommended):"
    echo "   - Use the automated GitHub Actions workflows"
    echo ""
}

# Function to show environment variables needed
show_env_vars() {
    echo ""
    echo "🔧 Required Environment Variables:"
    echo "=================================="
    echo ""
    echo "Frontend (Cloudflare Pages):"
    echo "- VITE_BACKEND_URL: Your backend URL (e.g., https://your-backend.railway.app)"
    echo "  → Set this in Cloudflare Pages project settings"
    echo ""
    echo "Backend (Railway):"
    echo "- PORT: 8090"
    echo "- POCKETBASE_ADMIN_EMAIL: Admin email for PocketBase"
    echo "- POCKETBASE_ADMIN_PASSWORD: Admin password for PocketBase"
    echo "  → Set these in Railway project variables"
    echo ""
    echo "GitHub Actions Secrets (only if using automated deployment):"
    echo "- CLOUDFLARE_API_TOKEN: Your Cloudflare API token"
    echo "- CLOUDFLARE_ACCOUNT_ID: Your Cloudflare account ID"
    echo "- RAILWAY_TOKEN: Your Railway token"
    echo ""
}

# Function to show next steps
show_next_steps() {
    echo ""
    echo "📝 Next Steps:"
    echo "=============="
    echo "1. Push your code to GitHub:"
    echo "   git add ."
    echo "   git commit -m 'Add deployment configuration'"
    echo "   git push origin main"
    echo ""
    echo "2. Set up your hosting accounts:"
    echo "   - Cloudflare Pages: https://pages.cloudflare.com/"
    echo "   - Railway: https://railway.app/"
    echo ""
    echo "3. Configure environment variables in your hosting platforms:"
    echo "   - Cloudflare Pages: Project Settings → Environment variables"
    echo "   - Railway: Project → Variables"
    echo ""
    echo "4. Deploy manually through the hosting platform dashboards"
    echo ""
    echo "5. Monitor your deployments and check the logs"
    echo ""
}

# Main script logic
case "${1:-}" in
    "test")
        print_status "Running build tests..."
        test_frontend_build
        test_backend_build
        print_success "All builds successful! Ready for deployment."
        ;;
    "frontend")
        test_frontend_build
        show_deployment_options
        ;;
    "backend")
        test_backend_build
        show_deployment_options
        ;;
    "env")
        show_env_vars
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  test     - Test both frontend and backend builds"
        echo "  frontend - Test frontend build and show deployment options"
        echo "  backend  - Test backend build and show deployment options"
        echo "  env      - Show required environment variables"
        echo "  help     - Show this help message"
        echo ""
        echo "Examples:"
        echo "  $0 test     # Test both builds"
        echo "  $0 frontend # Test frontend only"
        echo "  $0 env      # Show environment variables"
        ;;
    *)
        print_status "Welcome to Bitor deployment!"
        echo ""
        print_status "Testing builds..."
        test_frontend_build
        test_backend_build
        print_success "All builds successful!"
        
        show_deployment_options
        show_env_vars
        show_next_steps
        
        print_success "Ready for deployment! Check DEPLOYMENT.md for detailed instructions."
        ;;
esac 