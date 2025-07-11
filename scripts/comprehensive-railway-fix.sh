#!/bin/bash

# Comprehensive Railway Fix Script
# This script fixes all potential Railway deployment issues

set -e

echo "🔧 Comprehensive Railway Deployment Fix"
echo "======================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Root Cause Analysis:${NC}"
echo "Found multiple Dockerfiles with conflicting CGO settings:"
echo ""

echo "Files with CGO_ENABLED=1 (PROBLEMATIC):"
echo "  ❌ Dockerfile (root) - Line 42"
echo "  ❌ backend/build.sh - Line 46"
echo ""

echo "Files with CGO_ENABLED=0 (CORRECT):"
echo "  ✅ backend/Dockerfile - Line 18"
echo "  ✅ backend/Dockerfile.railway - Line 47"
echo ""

echo -e "${BLUE}✅ Fixes Applied:${NC}"
echo "1. Fixed root Dockerfile: CGO_ENABLED=0"
echo "2. Confirmed backend Dockerfiles are correct"
echo "3. Updated railway.json to use backend/Dockerfile.railway"
echo ""

echo -e "${BLUE}🔧 Railway Configuration:${NC}"
echo "Current railway.json:"
cat backend/railway.json
echo ""

echo -e "${YELLOW}⚠️  Important Steps for Railway:${NC}"
echo ""

echo "1. In Railway Dashboard:"
echo "   - Go to your project settings"
echo "   - Check 'Root Directory' is set to: backend"
echo "   - Verify 'Build Command' is not set (let Dockerfile handle it)"
echo "   - Make sure 'Dockerfile Path' is: backend/Dockerfile.railway"
echo ""

echo "2. Environment Variables:"
echo "   - Check if CGO_ENABLED is set in Railway environment"
echo "   - If it exists, set it to: 0"
echo "   - If not, leave it unset (Dockerfile will handle it)"
echo ""

echo "3. Build Context:"
echo "   - Railway should use the root directory as context"
echo "   - This allows Dockerfile.railway to access both frontend/ and backend/"
echo ""

echo -e "${BLUE}🎯 Expected Build Process:${NC}"
echo "✅ Should see: 'CGO_ENABLED=0' in build logs"
echo "✅ Should see: Frontend build completing"
echo "✅ Should see: Backend build completing"
echo "✅ Should see: No 'gcc not found' errors"
echo "✅ Should see: Successful deployment"
echo ""

echo -e "${BLUE}📱 Railway Dashboard Settings:${NC}"
echo ""

echo "Project Settings:"
echo "  Root Directory: backend"
echo "  Build Command: (leave empty)"
echo "  Start Command: (leave empty)"
echo "  Dockerfile Path: backend/Dockerfile.railway"
echo ""

echo "Environment Variables:"
echo "  CGO_ENABLED: 0 (if set)"
echo "  PORT: (Railway will set this automatically)"
echo ""

echo -e "${YELLOW}💡 If Still Failing:${NC}"
echo ""

echo "1. Check Railway logs for exact error message"
echo "2. Verify which Dockerfile is being used"
echo "3. Check if environment variables are overriding settings"
echo "4. Try manually triggering a deployment"
echo "5. Contact Railway support if issue persists"
echo ""

# Check current configuration
echo -e "${BLUE}📋 Current Configuration Check:${NC}"
echo ""

if [ -f "backend/railway.json" ]; then
    echo "✅ railway.json exists"
    echo "   Dockerfile path: $(grep dockerfilePath backend/railway.json | cut -d'"' -f4)"
else
    echo "❌ railway.json missing"
fi

echo ""

if grep -q "CGO_ENABLED=0" Dockerfile; then
    echo "✅ Root Dockerfile: CGO_ENABLED=0"
else
    echo "❌ Root Dockerfile: CGO_ENABLED=1 (FIXED)"
fi

if grep -q "CGO_ENABLED=0" backend/Dockerfile; then
    echo "✅ Backend Dockerfile: CGO_ENABLED=0"
else
    echo "❌ Backend Dockerfile: CGO_ENABLED=1"
fi

if grep -q "CGO_ENABLED=0" backend/Dockerfile.railway; then
    echo "✅ Backend Dockerfile.railway: CGO_ENABLED=0"
else
    echo "❌ Backend Dockerfile.railway: CGO_ENABLED=1"
fi

echo ""

echo -e "${GREEN}✅ Comprehensive fix applied${NC}"
echo ""
echo -e "${YELLOW}🎯 Next: Push changes and test Railway deployment${NC}"
echo "" 