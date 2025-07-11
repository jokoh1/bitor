#!/bin/bash

# Fix Cloudflare Pages Configuration Script
# This script helps resolve the build output directory mismatch

set -e

echo "🔧 Fixing Cloudflare Pages Configuration"
echo "========================================"

echo ""
echo "The issue is that Cloudflare Pages is looking for 'frontend/dist' but the build outputs to 'build'"
echo ""
echo "To fix this, you need to update the Cloudflare Pages dashboard settings:"
echo ""
echo "1. Go to https://pages.cloudflare.com/"
echo "2. Select your Bitor project"
echo "3. Go to Settings → Builds & deployments"
echo "4. Update these settings:"
echo "   - Root directory: frontend"
echo "   - Build command: pnpm install && pnpm build"
echo "   - Build output directory: build"
echo "   - Node.js version: 18"
echo ""
echo "Alternatively, you can use the wrangler.toml file that was just created:"
echo "The wrangler.toml file in frontend/ now has the correct configuration."
echo ""
echo "After updating the settings, trigger a new deployment by:"
echo "1. Making a small change to frontend/"
echo "2. Committing and pushing to main"
echo "3. Or manually triggering a deployment from the dashboard"
echo ""

# Check if we're in the right directory
if [ ! -f "frontend/package.json" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Verify the wrangler.toml file exists
if [ -f "frontend/wrangler.toml" ]; then
    echo "✅ wrangler.toml file is present with correct configuration"
else
    echo "❌ Error: wrangler.toml file is missing"
    exit 1
fi

# Show current build output
echo ""
echo "Current build configuration:"
echo "==========================="
echo "SvelteKit outputs to: build/"
echo "Cloudflare expects: build/ (after fix)"
echo ""

echo "🎯 Next steps:"
echo "1. Update Cloudflare Pages dashboard settings"
echo "2. Test with a small frontend change"
echo "3. Verify deployment works correctly"
echo "" 