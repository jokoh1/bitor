# Platform Configuration for Path-Based Deployment

This guide shows how to configure Cloudflare Pages and Railway to only deploy when their relevant files change.

## 🎯 Goal

- **Frontend changes** → Only Cloudflare Pages deploys
- **Backend changes** → Only Railway deploys
- **No cross-deployment confusion**

## 📋 Cloudflare Pages Configuration

### Step 1: Access Cloudflare Pages Dashboard
1. Go to [Cloudflare Pages](https://pages.cloudflare.com/)
2. Select your project

### Step 2: Configure Build Settings
1. Go to **Settings** → **Builds & deployments**
2. Set **Root directory** to: `frontend`
3. Set **Build command** to: `pnpm install && pnpm build`
4. Set **Build output directory** to: `build`
5. Set **Node.js version** to: `18`

### Step 3: Verify Configuration
- ✅ **Triggers**: Only when `frontend/` files change
- ✅ **Builds**: Only the SvelteKit application
- ✅ **Deploys**: Only frontend assets

## 🚂 Railway Configuration

### Step 1: Access Railway Dashboard
1. Go to [Railway](https://railway.app/)
2. Select your project

### Step 2: Configure Service Settings
1. Go to **Settings** → **General**
2. Set **Root Directory** to: `backend`
3. Set **Build Command** to: `go build -o bitor main.go`
4. Set **Start Command** to: `./bitor serve --http 0.0.0.0:$PORT`

### Step 3: Verify Configuration
- ✅ **Triggers**: Only when `backend/` files change
- ✅ **Builds**: Only the Go application
- ✅ **Deploys**: Only backend service

## 🔧 Current Repository Configuration

### Railway Configuration File
```json
// backend/railway.json
{
  "$schema": "https://railway.app/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "backend/Dockerfile.railway",
    "context": ".."
  },
  "deploy": {
    "numReplicas": 1,
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 10
  }
}
```

### What This Does
- Uses the root directory as build context
- Allows Dockerfile to access both `frontend/` and `backend/`
- Builds frontend and embeds it in backend

## 🧪 Testing the Configuration

### Test 1: Frontend-Only Change
```bash
# Make a change to frontend only
echo "// test" >> frontend/src/app.html
git add frontend/src/app.html
git commit -m "test: frontend only change"
git push origin main
```

**Expected Result:**
- ✅ Cloudflare Pages deploys
- ❌ Railway does NOT deploy

### Test 2: Backend-Only Change
```bash
# Make a change to backend only
echo "// test" >> backend/main.go
git add backend/main.go
git commit -m "test: backend only change"
git push origin main
```

**Expected Result:**
- ❌ Cloudflare Pages does NOT deploy
- ✅ Railway deploys

### Test 3: Both Changes
```bash
# Make changes to both
echo "// test" >> frontend/src/app.html
echo "// test" >> backend/main.go
git add .
git commit -m "test: both changes"
git push origin main
```

**Expected Result:**
- ✅ Cloudflare Pages deploys
- ✅ Railway deploys

## 🚨 Troubleshooting

### If Both Platforms Still Deploy on Every Push

1. **Check Dashboard Settings**
   - Verify root directory is set correctly in both platforms
   - Ensure settings are saved

2. **Check Platform Logs**
   - Look at deployment logs to see what triggered the build
   - Verify which files were detected as changed

3. **Platform Limitations**
   - Some platforms may not support path-based triggers
   - In that case, use the branch-based workflow we created

### Alternative: Branch-Based Workflow

If path-based triggers don't work, use our scripts:

```bash
# Frontend changes only
./scripts/deploy-frontend.sh my-feature

# Backend changes only
./scripts/deploy-backend.sh my-fix
```

## 📊 Benefits of Path-Based Triggers

✅ **Automatic**: No manual intervention needed  
✅ **Fast**: Only relevant platform builds  
✅ **Efficient**: No unnecessary deployments  
✅ **Simple**: Push to main, platform decides what to deploy  

## 🔄 Migration Steps

1. **Configure Cloudflare Pages** with root directory `frontend/`
2. **Configure Railway** with root directory `backend/`
3. **Test with small changes** to verify behavior
4. **Monitor deployments** to ensure correct triggers

## 📞 Support

If you encounter issues:

1. Check platform documentation for path-based triggers
2. Verify dashboard settings are correct
3. Test with isolated changes
4. Use branch-based workflow as fallback 