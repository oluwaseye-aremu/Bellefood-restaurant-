# 🚀 Production Deployment Guide

Deploy your Lumière Restaurant website to production.

## 📋 Pre-Deployment Checklist

- [ ] Test all features locally
- [ ] Have Neon PostgreSQL database ready
- [ ] Have Paystack live API keys (not test keys)
- [ ] Domain name (optional but recommended)
- [ ] SSL certificate (handled automatically by most platforms)

## 🎯 Recommended Stack

**Backend:** Railway, Render, or Heroku  
**Frontend:** Netlify or Vercel  
**Database:** Neon PostgreSQL (already using)  

## Part 1: Deploy Backend

### Option A: Deploy to Railway (Recommended - Easiest)

1. **Sign up at [Railway.app](https://railway.app/)**

2. **Create New Project**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Connect your GitHub account
   - Select your repository

3. **Configure Environment Variables**
   
   In Railway dashboard, add these variables:
   ```
   DATABASE_URL=your_neon_postgresql_connection_string
   PAYSTACK_SECRET_KEY=sk_live_your_live_secret_key
   PORT=8080
   ```

4. **Create railway.toml** in backend folder:
   ```toml
   [build]
   builder = "NIXPACKS"
   
   [deploy]
   startCommand = "go run main.go"
   ```

5. **Deploy**
   - Railway will automatically deploy
   - Note your backend URL (e.g., `https://your-app.railway.app`)

### Option B: Deploy to Render

1. **Sign up at [Render.com](https://render.com/)**

2. **Create New Web Service**
   - Connect GitHub
   - Select repository
   - Choose "Go" as environment

3. **Configure**
   - **Build Command:** `go build -o main main.go`
   - **Start Command:** `./main`
   - **Add Environment Variables:**
     ```
     DATABASE_URL=your_neon_connection
     PAYSTACK_SECRET_KEY=sk_live_xxx
     PORT=8080
     ```

4. **Deploy** - Render will build and deploy automatically

### Option C: Deploy to Heroku

1. **Install Heroku CLI**
   ```bash
   # Mac
   brew tap heroku/brew && brew install heroku
   
   # Windows
   # Download from: https://devcenter.heroku.com/articles/heroku-cli
   ```

2. **Login and Create App**
   ```bash
   heroku login
   cd backend
   heroku create your-restaurant-api
   ```

3. **Set Environment Variables**
   ```bash
   heroku config:set DATABASE_URL="your_neon_url"
   heroku config:set PAYSTACK_SECRET_KEY="sk_live_xxx"
   ```

4. **Create Procfile** in backend folder:
   ```
   web: go run main.go
   ```

5. **Deploy**
   ```bash
   git add .
   git commit -m "Deploy backend"
   git push heroku main
   ```

## Part 2: Deploy Frontend

### Option A: Deploy to Netlify (Recommended)

1. **Sign up at [Netlify.com](https://www.netlify.com/)**

2. **Prepare Frontend**
   
   Update all API_URL references to your production backend:
   
   In these files, change:
   ```javascript
   const API_URL = 'http://localhost:8080/api';
   ```
   
   To:
   ```javascript
   const API_URL = 'https://your-backend.railway.app/api';
   ```
   
   Files to update:
   - `frontend/js/menu-dynamic.js`
   - `frontend/cart.html`
   - `frontend/checkout.html`
   - `frontend/tracking.html`
   - `admin/admin.html`

3. **Deploy via Netlify UI**
   - Drag and drop your `frontend` folder
   - Or connect GitHub repo
   - Netlify will provide a URL like: `https://your-site.netlify.app`

4. **Custom Domain (Optional)**
   - In Netlify: Settings → Domain Management
   - Add your custom domain
   - Follow DNS instructions

### Option B: Deploy to Vercel

1. **Sign up at [Vercel.com](https://vercel.com/)**

2. **Update API URLs** (same as Netlify step 2)

3. **Deploy**
   ```bash
   npm i -g vercel
   cd frontend
   vercel
   ```

4. **Configure**
   - Follow prompts
   - Vercel will provide URL

## Part 3: Update CORS Settings

After deploying frontend, update backend CORS:

In `main.go`, find:
```go
AllowedOrigins: []string{"*"},
```

Change to:
```go
AllowedOrigins: []string{
    "https://your-site.netlify.app",
    "https://www.yourdomain.com",
},
```

Redeploy backend after this change.

## Part 4: Configure Paystack for Production

1. **Get Live API Key**
   - Go to Paystack Dashboard
   - Switch from Test to Live mode
   - Get your live secret key (starts with `sk_live_`)

2. **Update Backend**
   - Set `PAYSTACK_SECRET_KEY` to live key in Railway/Render/Heroku

3. **Frontend**
   - No Paystack public key is required in the current redirect-based checkout flow.
   - The frontend calls `/api/payment/initialize` and redirects to the Paystack `authorization_url`.

## Part 5: Setup Custom Domain (Optional)

### For Frontend (Netlify)

1. Buy domain from Namecheap, GoDaddy, etc.
2. In Netlify: Settings → Domain Management → Add Custom Domain
3. Update DNS records at your domain provider:
   ```
   Type: A
   Name: @
   Value: [Netlify IP from dashboard]
   
   Type: CNAME
   Name: www
   Value: [your-site].netlify.app
   ```

### For Backend (Railway)

1. In Railway: Settings → Domains → Add Custom Domain
2. Update DNS:
   ```
   Type: CNAME
   Name: api
   Value: [your-app].railway.app
   ```

3. Update frontend API_URL to: `https://api.yourdomain.com/api`

## Part 6: Admin Dashboard Security

### Important: Add Real Authentication

The current admin uses simple client-side auth. For production, implement:

1. **Backend Authentication**

Add to `main.go`:

```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Hash password check (use bcrypt in production)
    if req.Username == "admin" && checkPassword(req.Password) {
        token := generateJWT() // Implement JWT
        json.NewEncoder(w).Encode(map[string]string{
            "token": token,
        })
        return
    }
    
    http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}
```

2. **Protect Admin Routes**

Add middleware to verify JWT tokens on admin endpoints.

3. **Update Admin Login**

Make it POST to backend instead of client-side check.

## 🔒 Security Checklist

- [ ] Use HTTPS (handled by platforms)
- [ ] Use live Paystack keys in production
- [ ] Implement proper admin authentication (JWT)
- [ ] Restrict CORS to your domains only
- [ ] Never expose .env file
- [ ] Use strong admin passwords
- [ ] Enable Paystack webhook signature verification
- [ ] Regular database backups (Neon handles this)
- [ ] Monitor error logs
- [ ] Rate limiting on API endpoints

## 📊 Monitoring & Maintenance

### Setup Monitoring

**Backend (Railway/Render):**
- Check deployment logs
- Monitor CPU/Memory usage
- Setup alerts for downtime

**Frontend (Netlify/Vercel):**
- Analytics built-in
- Monitor build/deploy success

**Database (Neon):**
- Monitor connection pool
- Check query performance
- Regular backups

### Paystack Webhooks (Advanced)

To receive real-time payment updates:

1. In Paystack Dashboard: Developers → Webhooks
2. Add endpoint: `https://your-backend.com/api/webhooks/paystack`
3. Select Paystack events such as `charge.success`
4. Implement webhook handler in backend

## 🐛 Troubleshooting

**CORS Errors:**
```
Update AllowedOrigins in main.go to include your frontend URL
```

**Images not loading:**
```
Check that uploads/ folder has write permissions
Verify image URLs are correct in database
```

**Payment failing:**
```
Verify live Paystack keys are set
Check webhook configurations
Test with Paystack test payments first
```

**Database connection issues:**
```
Verify DATABASE_URL includes ?sslmode=require
Check Neon dashboard for connection limits
Ensure IP allowlist is configured (if enabled)
```

## 💰 Cost Estimate

**Neon PostgreSQL:** Free tier (up to 512MB)  
**Railway:** Free $5 credit/month, then ~$5-20/month  
**Netlify:** Free tier (100GB bandwidth)  
**Paystack:** fees vary by country, card type, and settlement rules  

**Total:** ~$5-20/month for small-medium traffic

## 📈 Scaling Tips

1. **Database:** Upgrade Neon plan as needed
2. **Backend:** Railway auto-scales, or add more instances
3. **CDN:** Use Cloudflare for static assets
4. **Caching:** Implement Redis for frequently accessed data
5. **Load Balancing:** Add nginx if traffic is very high

## 🎉 Post-Deployment

After successful deployment:

1. **Test Everything:**
   - [ ] Browse menu
   - [ ] Add to cart
   - [ ] Complete checkout with Paystack test payment
   - [ ] Track order
   - [ ] Login to admin
   - [ ] Add/delete menu items
   - [ ] Update order tracking

2. **Share with Users:**
   - Announce on social media
   - Add Google Analytics
   - Setup SEO metadata
   - Submit to Google Search Console

3. **Monitor First Orders:**
   - Watch Paystack dashboard
   - Check order logs
   - Verify email notifications work
   - Test customer tracking links

## 🔗 Useful Resources

- [Railway Docs](https://docs.railway.app/)
- [Netlify Docs](https://docs.netlify.com/)
- [Neon Docs](https://neon.tech/docs)
- [Paystack Production Checklist](https://paystack.com/docs/development/checklist)
- [Go Deployment Best Practices](https://golang.org/doc/articles/wiki/)

---

**Congratulations! Your restaurant is now live! 🎊**

Need help? Check the troubleshooting section or reach out to the platform's support.
