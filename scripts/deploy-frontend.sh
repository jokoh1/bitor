#!/bin/bash

# Frontend Deployment Workflow Script
# Usage: ./scripts/deploy-frontend.sh [branch-name]

set -e

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
if [ ! -f "frontend/package.json" ]; then
    print_error "This script must be run from the project root directory"
    exit 1
fi

# Get branch name from argument or generate one
if [ -n "$1" ]; then
    BRANCH_NAME="feature/frontend-$1"
else
    BRANCH_NAME="feature/frontend-$(date +%Y%m%d-%H%M%S)"
fi

print_status "Starting frontend deployment workflow..."
print_status "Branch name: $BRANCH_NAME"

# Check if we have uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    print_warning "You have uncommitted changes. Please commit or stash them first."
    git status --short
    exit 1
fi

# Create and switch to new branch
print_status "Creating feature branch..."
git checkout -b "$BRANCH_NAME"

print_success "Created branch: $BRANCH_NAME"
print_status "Make your frontend changes in the frontend/ directory"
print_status "Then run: git add . && git commit -m 'your message' && git push origin $BRANCH_NAME"
print_status "Finally, create a PR on GitHub to merge into main"

print_status ""
print_status "Next steps:"
echo "1. Make your frontend changes"
echo "2. git add ."
echo "3. git commit -m 'your commit message'"
echo "4. git push origin $BRANCH_NAME"
echo "5. Create PR on GitHub: https://github.com/jokoh1/bitor/pull/new/$BRANCH_NAME"
echo "6. Merge PR → Cloudflare Pages will auto-deploy" 