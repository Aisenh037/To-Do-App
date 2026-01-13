# Netlify Deployment Guide for Angular Frontend

## 🚀 Quick Deploy to Netlify

### **Option 1: Netlify UI (Recommended for First Time)**

#### Step 1: Sign Up / Login
1. Go to **https://netlify.com**
2. Click **"Sign up"** or **"Log in"**
3. Use **GitHub** for easy integration

#### Step 2: New Site from Git
1. Click **"Add new site"** → **"Import an existing project"**
2. Choose **"Deploy with GitHub"**
3. Authorize Netlify to access your GitHub
4. Select repository: **`Aisenh037/To-Do-App`**

#### Step 3: Configure Build Settings

Netlify should auto-detect these from `netlify.toml`, but verify:

| Setting | Value |
|---------|-------|
| **Base directory** | `web` |
| **Build command** | `npm install && npm run build` |
| **Publish directory** | `dist/web/browser` |
| **Production branch** | `main` |

#### Step 4: Set Environment Variables

Before deploying, add environment variable:

1. Click **"Site configuration"** → **"Environment variables"**
2. Add variable:
   ```
   Key: API_URL
   Value: https://go-todo-api.onrender.com
   ```
   (Change this to your actual Render backend URL after you deploy it)

#### Step 5: Deploy!

1. Click **"Deploy site"**
2. Wait 2-3 minutes for build
3. Your site will be live at: `https://random-name-123456.netlify.app`

#### Step 6: Custom Domain (Optional)

1. Go to **"Domain management"**
2. Click **"Add custom domain"** or **"Change site name"**
3. Rename to something like: `todo-app-yourname.netlify.app`

---

### **Option 2: Netlify CLI (For Advanced Users)**

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Login to Netlify
netlify login

# Navigate to your project
cd c:\Users\ASUS\Desktop\Golang\go-todo-api

# Deploy (from project root)
netlify deploy --prod

# Follow prompts:
# - Create & configure new site: Yes
# - Publish directory: web/dist/web/browser
```

---

## 🔧 Update Angular App to Use Environment Variables

You'll need to update your Angular services to use the Render backend URL.

### Find Your API Service Files

Look for files that make HTTP requests, typically:
- `web/src/app/services/auth.service.ts`
- `web/src/app/services/todo.service.ts`

### Update API Base URL

Change from:
```typescript
private apiUrl = 'http://localhost:8080/api';
```

To:
```typescript
private apiUrl = 'https://go-todo-api.onrender.com/api';
```

Or better yet, create an environment file for production.

---

## ✅ After Deployment Checklist

Once both backend (Render) and frontend (Netlify) are deployed:

1. **Update CORS in Backend**:
   Add your Netlify URL to allowed origins:
   ```
   ALLOWED_ORIGINS=https://your-app.netlify.app
   ```

2. **Test Your Live App**:
   - Open your Netlify URL
   - Try registering a new user
   - Try logging in
   - Create some todos

3. **Update GitHub README**:
   Add live demo links:
   ```markdown
   ## 🌐 Live Demo
   
   - **Frontend**: https://your-app.netlify.app
   - **API Documentation**: https://go-todo-api.onrender.com/swagger/index.html
   ```

---

## 🐛 Troubleshooting

### Build Fails
- Check build logs in Netlify
- Ensure `package.json` has all dependencies
- Verify Node version (v18+)

### API Requests Fail (CORS Error)
- Add Netlify URL to `ALLOWED_ORIGINS` in Render
- Redeploy backend on Render

### Routes Don't Work (404 on refresh)
- The `netlify.toml` should handle this
- Verify the redirects rule is in place

### Environment Variable Not Working
- Set `API_URL` in Netlify dashboard
- Redeploy after adding env vars

---

## 📊 Deployment Status

**Current Setup:**
- ✅ Dockerfile ready for Render
- ✅ render.yaml for backend
- ✅ netlify.toml for frontend
- ⏳ Backend deployment (Render) - Pending
- ⏳ Frontend deployment (Netlify) - Pending

**Next Steps:**
1. Deploy backend to Render first
2. Get the Render URL (e.g., `https://go-todo-api.onrender.com`)
3. Update Angular services with Render URL
4. Deploy frontend to Netlify
5. Update README with live demo links

---

## 🎯 Pro Tips

- **Free Tier Limits**: Render free tier sleeps after 15 min of inactivity (first request takes ~30s to wake up)
- **Build Time**: Netlify builds are usually faster than Render (2-3 min vs 8-10 min)
- **Custom Domains**: Both platforms support custom domains for free
- **Auto Deploy**: Both auto-deploy when you push to GitHub main branch 🚀
