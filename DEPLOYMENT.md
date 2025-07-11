# Bitor Deployment Guide

This guide will help you deploy the Bitor application to production. The application consists of:
- **Frontend**: SvelteKit (deployed to Cloudflare Pages)
- **Backend**: Go + PocketBase (deployed to Railway)

## How Deployment Works

### When You Push to GitHub

1. **GitHub Actions runs build tests** (no automatic deployment)
2. **Frontend build test**: Ensures your SvelteKit app builds successfully
3. **Backend build test**: Ensures your Go app builds successfully
4. **Manual deployment required**: You deploy through hosting platform dashboards

### What Gets Deployed Where

| Component | Location | Deployment Method | What Happens |
|-----------|----------|-------------------|--------------|
| **Frontend** | Cloudflare Pages | Manual via dashboard | Only `frontend/` directory is deployed |
| **Backend** | Railway | Manual via dashboard | Only `backend/` directory is deployed |
| **Full Repository** | GitHub | Automatic on push | All code is stored, but not deployed |

### Why This Approach?

- ✅ **No secrets in GitHub** - Environment variables set in hosting platforms
- ✅ **Better security** - Sensitive data stays in hosting platforms
- ✅ **Flexible deployment** - Deploy frontend and backend independently
- ✅ **Build validation** - GitHub Actions ensures code builds before deployment

## Prerequisites

- GitHub account
- Cloudflare account
- Railway account (or alternative backend hosting)
- Git installed locally

## Step 1: Frontend Deployment (Cloudflare Pages)

### Manual Deployment via Cloudflare Dashboard (Recommended)

1. **Go to [Cloudflare Pages](https://pages.cloudflare.com/)**
2. **Click "Create a project"**
3. **Connect your GitHub repository**
4. **Configure build settings:**
   - **Project name**: `bitor-frontend`
   - **Production branch**: `main`
   - **Root directory**: `frontend` ← **This is key!**
   - **Build command**: `pnpm install && pnpm build`
   - **Build output directory**: `build`
   - **Node.js version**: `18`

5. **Add Environment Variables in Cloudflare Pages:**
   - Go to your project → Settings → Environment variables
   - Add the following variables:
     - `VITE_BACKEND_URL`: Your backend URL (e.g., `https://your-backend.railway.app`)
   - **Note**: These are set directly in Cloudflare, not in GitHub secrets

6. **Click "Save and Deploy"**

### What Cloudflare Pages Sees

When you connect your repository to Cloudflare Pages:
- ✅ **Sees**: `frontend/` directory only
- ❌ **Doesn't see**: `backend/`, `wrangler.toml`, `.github/`, etc.
- ✅ **Builds**: Only the SvelteKit application
- ✅ **Deploys**: Only the built frontend files

## Step 2: Backend Deployment (Railway)

### Manual Deployment via Railway Dashboard (Recommended)

1. **Go to [Railway](https://railway.app/)**
2. **Click "New Project"**
3. **Select "Deploy from GitHub repo"**
4. **Choose your repository**
5. **Configure the service:**
   - **Root Directory**: `backend` ← **This is key!**
   - **Build Command**: `go build -o bitor main.go`
   - **Start Command**: `./bitor serve --http 0.0.0.0:$PORT`

6. **Add Environment Variables in Railway:**
   - Go to your Railway project → Variables
   - Add the following variables:
     - `PORT`: `8090`
     - `POCKETBASE_ADMIN_EMAIL`: Admin email for PocketBase
     - `POCKETBASE_ADMIN_PASSWORD`: Admin password for PocketBase
     - Any other environment variables your app needs

### What Railway Sees

When you connect your repository to Railway:
- ✅ **Sees**: `backend/` directory only
- ❌ **Doesn't see**: `frontend/`, `wrangler.toml`, `.github/`, etc.
- ✅ **Builds**: Only the Go application
- ✅ **Deploys**: Only the backend service

## Step 3: Connect Frontend to Backend

1. **Get your backend URL from Railway**
2. **Update your Cloudflare Pages environment variables:**
   - Go to your Cloudflare Pages project → Settings → Environment variables
   - Set `VITE_BACKEND_URL` to your Railway backend URL

3. **Redeploy your frontend** (this will happen automatically if you're using GitHub Actions)

## Step 4: Custom Domain (Optional)

### Frontend (Cloudflare Pages)
1. Go to your Cloudflare Pages project → Custom domains
2. Add your domain
3. Update DNS records as instructed

### Backend (Railway)
1. Go to your Railway project → Settings → Domains
2. Add your custom domain
3. Update DNS records as instructed

## Step 5: Environment Variables

### Frontend Environment Variables (Cloudflare Pages)
Set these directly in your Cloudflare Pages project settings:
- `VITE_BACKEND_URL`: Your backend API URL

### Backend Environment Variables (Railway)
Set these directly in your Railway project variables:
- `PORT`: Port number (usually 8090)
- `POCKETBASE_ADMIN_EMAIL`: Admin email for PocketBase
- `POCKETBASE_ADMIN_PASSWORD`: Admin password for PocketBase
- Any other environment variables your Go app needs

## Step 6: Testing Your Deployment

1. **Test the frontend**: Visit your Cloudflare Pages URL
2. **Test the backend**: Visit your Railway backend URL
3. **Test the connection**: Try logging in or using features that require backend communication

## Troubleshooting

### Common Issues

1. **Build fails on Cloudflare Pages:**
   - Check that the root directory is set to `frontend`
   - Verify that `pnpm` is available (it should be automatically detected)
   - Check the build logs for specific errors

2. **Backend deployment fails on Railway:**
   - Ensure the root directory is set to `backend`
   - Check that all Go dependencies are properly specified in `go.mod`
   - Verify the build and start commands

3. **Frontend can't connect to backend:**
   - Check that `VITE_BACKEND_URL` is set correctly in Cloudflare Pages
   - Verify CORS settings on your backend
   - Check that the backend is running and accessible

4. **Database issues:**
   - PocketBase will create its database automatically
   - Ensure the `pb_data` directory is writable
   - Check Railway logs for database-related errors

### Getting Help

- Check the build logs in Cloudflare Pages and Railway
- Review the application logs in Railway
- Check the browser console for frontend errors
- Verify all environment variables are set correctly

## Alternative Backend Hosting Options

If Railway doesn't work for you, consider these alternatives:

1. **Fly.io**: Great for Go applications
2. **Render**: Easy deployment with good free tier
3. **DigitalOcean App Platform**: Reliable and scalable
4. **Heroku**: Classic choice (requires credit card)
5. **AWS/GCP/Azure**: For more advanced users

## Security Considerations

1. **Environment Variables**: Never commit sensitive data to your repository
2. **HTTPS**: Both Cloudflare Pages and Railway provide HTTPS by default
3. **CORS**: Configure CORS properly on your backend
4. **Authentication**: Ensure your authentication system is properly configured

## Monitoring and Maintenance

1. **Set up monitoring**: Use Railway's built-in monitoring or add external services
2. **Regular updates**: Keep your dependencies updated
3. **Backup strategy**: Consider backing up your PocketBase database
4. **Logs**: Monitor application logs for errors and performance issues

## Cost Optimization

1. **Cloudflare Pages**: Free tier includes 500 builds/month
2. **Railway**: Free tier includes $5 credit/month
3. **Monitor usage**: Keep track of your resource usage
4. **Optimize builds**: Minimize build times and resource usage

---

## Quick Start Commands

```bash
# Clone and setup
git clone <your-repo>
cd bitor

# Test frontend locally
cd frontend
pnpm install
pnpm dev

# Test backend locally
cd ../backend
go mod download
go run main.go serve

# Deploy (if using GitHub Actions)
git add .
git commit -m "Ready for deployment"
git push origin main
```

Your application should now be deployed and accessible via the provided URLs! 