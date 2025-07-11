# Bitor Deployment Workflow

This document describes the professional branch-based deployment workflow for Bitor.

## 🎯 Overview

Instead of deploying on every push to main, we use feature branches to:
- **Isolate changes** by component (frontend/backend)
- **Target deployments** to the correct platform
- **Enable code review** before deployment
- **Prevent cross-deployment confusion**

## 🚀 Quick Start

### Frontend Changes Only
```bash
./scripts/deploy-frontend.sh [optional-branch-name]
# Make changes to frontend/
git add .
git commit -m "your message"
git push origin feature/frontend-[name]
# Create PR → Merge → Cloudflare Pages deploys
```

### Backend Changes Only
```bash
./scripts/deploy-backend.sh [optional-branch-name]
# Make changes to backend/
git add .
git commit -m "your message"
git push origin feature/backend-[name]
# Create PR → Merge → Railway deploys
```

### Full Stack Changes
```bash
./scripts/deploy-full.sh [optional-branch-name]
# Make changes to both frontend/ and backend/
git add .
git commit -m "your message"
git push origin feature/full-[name]
# Create PR → Merge → Both platforms deploy
```

## 📋 Detailed Workflow

### 1. Frontend Deployment

**When to use**: Changes to `frontend/` directory only

```bash
# Start frontend workflow
./scripts/deploy-frontend.sh dashboard-update

# Make your changes
# Edit files in frontend/

# Commit and push
git add .
git commit -m "feat: add new dashboard component"
git push origin feature/frontend-dashboard-update

# Create PR on GitHub
# Review and merge → Cloudflare Pages auto-deploys
```

**Result**: Only Cloudflare Pages deploys, Railway doesn't trigger

### 2. Backend Deployment

**When to use**: Changes to `backend/` directory only

```bash
# Start backend workflow
./scripts/deploy-backend.sh api-fix

# Make your changes
# Edit files in backend/

# Commit and push
git add .
git commit -m "fix: resolve API authentication issue"
git push origin feature/backend-api-fix

# Create PR on GitHub
# Review and merge → Railway auto-deploys
```

**Result**: Only Railway deploys, Cloudflare Pages doesn't trigger

### 3. Full Stack Deployment

**When to use**: Changes to both `frontend/` and `backend/` directories

```bash
# Start full stack workflow
./scripts/deploy-full.sh major-update

# Make your changes
# Edit files in both frontend/ and backend/

# Commit and push
git add .
git commit -m "feat: implement new user authentication system"
git push origin feature/full-major-update

# Create PR on GitHub
# Review and merge → Both platforms deploy
```

**Result**: Both Cloudflare Pages and Railway deploy

## 🔧 Script Details

### Available Scripts

| Script | Purpose | Branch Prefix | Deploys To |
|--------|---------|---------------|------------|
| `deploy-frontend.sh` | Frontend changes only | `feature/frontend-` | Cloudflare Pages |
| `deploy-backend.sh` | Backend changes only | `feature/backend-` | Railway |
| `deploy-full.sh` | Full stack changes | `feature/full-` | Both platforms |

### Script Features

- ✅ **Automatic branch creation** with timestamp or custom name
- ✅ **Pre-flight checks** (uncommitted changes, directory validation)
- ✅ **Colored output** for better UX
- ✅ **Clear next steps** guidance
- ✅ **Error handling** and validation

### Usage Examples

```bash
# Auto-generated branch name
./scripts/deploy-frontend.sh
# Creates: feature/frontend-20241211-143022

# Custom branch name
./scripts/deploy-frontend.sh new-feature
# Creates: feature/frontend-new-feature

# Backend with custom name
./scripts/deploy-backend.sh bugfix-123
# Creates: feature/backend-bugfix-123
```

## 📊 Benefits

### Before (Current Issues)
- ❌ Both platforms deploy on every push
- ❌ Cross-deployment confusion
- ❌ No code review process
- ❌ Hard to isolate issues

### After (New Workflow)
- ✅ **Targeted deployments** - only relevant platform deploys
- ✅ **Clean separation** - frontend/backend changes isolated
- ✅ **Code review** - PR process before deployment
- ✅ **Better debugging** - know exactly what caused issues
- ✅ **Professional workflow** - industry standard approach

## 🛠️ Setup Requirements

### 1. Make Scripts Executable
```bash
chmod +x scripts/deploy-*.sh
```

### 2. GitHub Branch Protection (Recommended)
1. Go to GitHub repository → Settings → Branches
2. Add rule for `main` branch
3. Enable:
   - ✅ Require pull request reviews
   - ✅ Require status checks to pass
   - ✅ Include administrators

### 3. Conventional Commits (Recommended)
Use conventional commit messages:
```bash
git commit -m "feat: add new dashboard"
git commit -m "fix: resolve API authentication"
git commit -m "docs: update deployment guide"
```

## 🔍 Monitoring Deployments

### Cloudflare Pages
- **URL**: https://pages.cloudflare.com/
- **Triggers**: Merges to main with frontend changes
- **Build Time**: ~2-3 minutes

### Railway
- **URL**: https://railway.app/
- **Triggers**: Merges to main with backend changes
- **Build Time**: ~3-5 minutes

### GitHub Actions
- **URL**: https://github.com/jokoh1/bitor/actions
- **Purpose**: Build testing and validation
- **Triggers**: All pushes and PRs

## 🚨 Troubleshooting

### Common Issues

1. **Script not executable**
   ```bash
   chmod +x scripts/deploy-*.sh
   ```

2. **Uncommitted changes**
   ```bash
   git add . && git commit -m "WIP: save current work"
   # or
   git stash
   ```

3. **Wrong directory**
   ```bash
   # Make sure you're in the project root
   ls frontend/package.json backend/main.go
   ```

4. **Branch already exists**
   ```bash
   # Delete local branch
   git branch -D feature/frontend-name
   # or use different name
   ./scripts/deploy-frontend.sh new-name
   ```

### Getting Help

- Check script output for error messages
- Verify you're in the project root directory
- Ensure all changes are committed before running scripts
- Check GitHub Actions for build status

## 📈 Best Practices

1. **Use descriptive branch names**
   ```bash
   ./scripts/deploy-frontend.sh dashboard-redesign
   ./scripts/deploy-backend.sh user-auth-fix
   ```

2. **Write clear commit messages**
   ```bash
   git commit -m "feat: implement new user dashboard"
   git commit -m "fix: resolve API rate limiting issue"
   ```

3. **Test locally before pushing**
   ```bash
   # Frontend
   cd frontend && pnpm build
   
   # Backend
   cd backend && go build
   ```

4. **Review PRs thoroughly**
   - Check for unintended changes
   - Verify deployment target is correct
   - Test functionality if possible

5. **Monitor deployments**
   - Watch build logs for errors
   - Verify deployment success
   - Test functionality after deployment

## 🔄 Migration from Current Workflow

### For Existing Development

1. **Complete current work** on main branch
2. **Switch to new workflow** for future changes
3. **Use scripts** for all new features/fixes
4. **Gradually adopt** branch-based workflow

### Team Adoption

1. **Share this documentation** with team members
2. **Set up branch protection** rules
3. **Use scripts consistently** across team
4. **Monitor and improve** workflow over time

---

## 📞 Support

If you encounter issues with this workflow:

1. Check this documentation first
2. Review script error messages
3. Check GitHub Actions for build status
4. Verify deployment platform logs
5. Create an issue on GitHub if needed 