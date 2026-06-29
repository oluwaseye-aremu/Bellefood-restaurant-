# 🔧 Troubleshooting Guide

Common issues and their solutions.

## 🚨 Backend Issues

### "Database connection failed"

**Symptoms:**
```
Error: pq: SSL is not enabled on the server
```

**Solutions:**
1. Add `?sslmode=require` to your DATABASE_URL:
   ```
   postgresql://user:pass@host/db?sslmode=require
   ```

2. Verify Neon connection string is correct:
   - Login to neon.tech
   - Copy connection string from dashboard
   - Make sure it's the full URL, not just hostname

3. Check if Neon database is active:
   - Free tier databases may pause after inactivity
   - Visit Neon dashboard to wake it up

### "Port already in use"

**Symptoms:**
```
Error: listen tcp :8080: bind: address already in use
```

**Solutions:**
```bash
# Find what's using port 8080
lsof -i :8080
# Or on Windows
netstat -ano | findstr :8080

# Kill the process
kill -9 <PID>
# Or on Windows
taskkill /PID <PID> /F

# Or use a different port
export PORT=8081
go run main.go
```

### "Module not found"

**Symptoms:**
```
Error: cannot find module github.com/gorilla/mux
```

**Solutions:**
```bash
# Install dependencies
go mod tidy

# If still failing, manually install
go get github.com/gorilla/mux
go get github.com/lib/pq
go get github.com/rs/cors
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv

# Clear module cache if needed
go clean -modcache
go mod tidy
```

### "Table does not exist"

**Symptoms:**
```
Error: pq: relation "menu_items" does not exist
```

**Solutions:**
1. The tables should be created automatically. Restart the backend:
   ```bash
   # Stop server (Ctrl+C)
   # Start again
   go run main.go
   ```

2. Check database logs for table creation:
   ```
   Tables created successfully!
   ```

3. If tables still not created, manually run SQL:
   ```sql
   -- Connect to Neon PostgreSQL and run these commands
   CREATE TABLE IF NOT EXISTS menu_items (
       id SERIAL PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       description TEXT,
       price DECIMAL(10, 2) NOT NULL,
       category VARCHAR(100),
       image_url TEXT,
       ingredients TEXT,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   -- Repeat for orders and order_tracking tables (see main.go)
   ```

## 🌐 Frontend Issues

### Menu Not Loading

**Symptoms:**
- Empty menu page
- Console error: "Failed to fetch"

**Solutions:**
1. Check backend is running:
   ```bash
   # Should see "Server starting on port 8080"
   curl http://localhost:8080/api/menu
   ```

2. Verify API_URL in `menu-dynamic.js`:
   ```javascript
   const API_URL = 'http://localhost:8080/api';
   // Make sure this matches your backend port
   ```

3. Check browser console (F12) for CORS errors

4. Try accessing API directly in browser:
   ```
   http://localhost:8080/api/menu
   ```

### CORS Error

**Symptoms:**
```
Access to fetch at 'http://localhost:8080/api/menu' from origin 
'http://localhost:3000' has been blocked by CORS policy
```

**Solutions:**
1. Update CORS in `main.go`:
   ```go
   AllowedOrigins: []string{
       "http://localhost:3000",
       "http://127.0.0.1:3000",
       "*", // Allow all (for development only)
   },
   ```

2. Restart backend after changes

3. Clear browser cache (Ctrl+Shift+Delete)

### Images Not Displaying

**Symptoms:**
- Menu items show but no images
- 404 errors in console

**Solutions:**
1. Check uploads folder exists:
   ```bash
   cd backend
   ls -la uploads/
   ```

2. Verify image URLs in database:
   ```sql
   SELECT id, title, image_url FROM menu_items;
   ```

3. Check image paths in frontend:
   ```javascript
   // Should be:
   `${API_URL.replace('/api', '')}${item.image_url}`
   // Results in: http://localhost:8080/uploads/image.jpg
   ```

4. Ensure backend serves uploads:
   ```go
   router.PathPrefix("/uploads/").Handler(
       http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads")))
   )
   ```

### Cart Not Persisting

**Symptoms:**
- Cart empties after page reload
- Items don't stay in cart

**Solutions:**
1. Check browser's LocalStorage:
   ```javascript
   // Open console (F12) and run:
   localStorage.getItem('cart')
   // Should show JSON array
   ```

2. Clear LocalStorage if corrupted:
   ```javascript
   localStorage.clear()
   ```

3. Check for JavaScript errors in console

4. Verify cart is being saved:
   ```javascript
   // In menu-dynamic.js
   console.log('Cart saved:', cart);
   localStorage.setItem('cart', JSON.stringify(cart));
   ```

## 💳 Payment Issues

### Paystack Payment Failing

**Symptoms:**
- "Payment failed" message
- Checkout button doesn't work

**Solutions:**
1. Verify Paystack keys are correct:
   ```bash
   # In backend .env
   PAYSTACK_SECRET_KEY=sk_test_...
   # Should start with sk_test_ for testing
   ```

2. Use Paystack test payment details from your Paystack dashboard or Paystack documentation.

3. Check Paystack dashboard for errors:
   - Go to https://dashboard.paystack.com/test/logs
   - Look for failed payment attempts

4. Confirm checkout receives an `authorization_url` from `/api/payment/initialize`.

5. Check browser console for errors

### "Payment Initialization Failed"

**Symptoms:**
```
Payment init failed
```

**Solutions:**
1. Check backend logs for detailed error

2. Verify amount is sent correctly:
   ```javascript
   // Amount should be in naira
   body: JSON.stringify({ amount: 28.00 }) // Correct
   // Backend converts to kobo before sending to Paystack
   ```

3. Check Paystack secret key is set:
   ```bash
   echo $PAYSTACK_SECRET_KEY
   # Should print your key
   ```

## 📦 Order & Tracking Issues

### Order Not Creating

**Symptoms:**
- Payment succeeds but no order created
- Tracking ID not generated

**Solutions:**
1. Check backend logs for errors

2. Verify order data is complete:
   ```javascript
   // All required fields must be present
   {
       customer_name: "John Doe",
       customer_email: "john@example.com",
       customer_phone: "1234567890",
       items: JSON.stringify(cart), // Must be JSON string
       total_amount: 28.00,
       payment_status: "paid"
   }
   ```

3. Check database:
   ```sql
   SELECT * FROM orders ORDER BY created_at DESC LIMIT 5;
   ```

### Tracking Page Shows "Not Found"

**Symptoms:**
- Enter tracking ID, get "Order not found"

**Solutions:**
1. Verify tracking ID format:
   ```
   Should be: ORD-1702483821
   Not: ord-1702483821 (case sensitive)
   ```

2. Check if order exists:
   ```sql
   SELECT tracking_id FROM orders WHERE tracking_id = 'ORD-1702483821';
   ```

3. Verify API endpoint:
   ```bash
   curl http://localhost:8080/api/orders/tracking/ORD-1702483821
   ```

4. Check API_URL in tracking.html

### Tracking Updates Not Showing

**Symptoms:**
- Admin updates status but customer doesn't see it

**Solutions:**
1. Refresh tracking page (it's not real-time)

2. Check order_tracking table:
   ```sql
   SELECT * FROM order_tracking WHERE order_id = 1;
   ```

3. Verify tracking endpoint returns updates:
   ```bash
   curl http://localhost:8080/api/orders/tracking/ORD-XXX
   # Should return "tracking" array with updates
   ```

## 🔐 Admin Issues

### Can't Login to Admin

**Symptoms:**
- "Invalid credentials" message

**Solutions:**
1. Use correct demo credentials:
   ```
   Username: admin
   Password: admin123
   ```

2. Check browser console for JavaScript errors

3. Clear browser cache and cookies

4. Check if login.html is being served correctly

### Admin Dashboard Won't Load

**Symptoms:**
- Redirected to login immediately
- Blank admin page

**Solutions:**
1. Check if authenticated:
   ```javascript
   // Open console (F12)
   sessionStorage.getItem('admin_authenticated')
   // Should return 'true'
   ```

2. Login again:
   ```javascript
   sessionStorage.clear()
   // Then login again
   ```

3. Check for JavaScript errors in console

### Image Upload Failing

**Symptoms:**
- "Unable to save image" error
- Menu item creates without image

**Solutions:**
1. Check file size (must be < 10MB)

2. Verify file is an image:
   - Supported: .jpg, .jpeg, .png, .gif, .webp

3. Check uploads folder permissions:
   ```bash
   ls -la backend/uploads/
   # Should have write permissions
   chmod 755 backend/uploads/
   ```

4. Check backend logs for specific error

5. Test with a smaller image (< 1MB)

## 🐛 General Debugging Tips

### Enable Detailed Logging

In `main.go`, add more logs:
```go
log.Printf("Received request: %s %s", r.Method, r.URL.Path)
log.Printf("Request body: %+v", requestData)
log.Printf("Database query result: %+v", result)
```

### Check Browser Console

Always check console (F12) for:
- JavaScript errors
- Failed network requests
- Console.log outputs

### Test API Endpoints

Use curl or Postman:
```bash
# Test menu endpoint
curl http://localhost:8080/api/menu

# Test with data
curl -X POST http://localhost:8080/api/menu \
  -F "title=Test Item" \
  -F "price=10" \
  -F "category=starters" \
  -F "description=Test" \
  -F "image=@/path/to/image.jpg"
```

### Check Database Directly

Connect to Neon PostgreSQL:
```bash
psql "your_neon_connection_string"

# Run queries
\dt  -- List tables
SELECT * FROM menu_items;
SELECT * FROM orders;
```

### Reset Everything

If all else fails:
```bash
# 1. Stop backend
# 2. Clear database (be careful!)
DROP TABLE order_tracking;
DROP TABLE orders;
DROP TABLE menu_items;

# 3. Clear frontend
localStorage.clear()
sessionStorage.clear()

# 4. Restart backend (tables recreate automatically)
go run main.go

# 5. Refresh browser
```

## 📱 Mobile Issues

### Site Not Responsive

**Solutions:**
1. Check viewport meta tag in HTML:
   ```html
   <meta name="viewport" content="width=device-width, initial-scale=1.0">
   ```

2. Test CSS media queries

3. Check mobile console (use device mode in DevTools)

## 🔄 Update Issues

### After Code Changes, Nothing Changes

**Solutions:**
1. **Backend:** Restart server
   ```bash
   # Stop with Ctrl+C
   go run main.go
   ```

2. **Frontend:** Hard refresh browser
   ```
   Ctrl+Shift+R (Windows/Linux)
   Cmd+Shift+R (Mac)
   ```

3. Clear browser cache

4. Check if changes were saved to file

## 💾 Backup & Recovery

### Lost Database Data

**Prevention:**
1. Neon has automatic backups
2. Export data regularly:
   ```bash
   pg_dump "your_neon_url" > backup.sql
   ```

**Recovery:**
```bash
psql "your_neon_url" < backup.sql
```

### Lost Environment Variables

Keep a backup of `.env` file securely (not in Git!)

## 🆘 Still Having Issues?

### Checklist Before Asking for Help

- [ ] Checked all error messages
- [ ] Reviewed browser console
- [ ] Checked backend logs
- [ ] Verified environment variables
- [ ] Tested API endpoints directly
- [ ] Tried with fresh browser session
- [ ] Restarted backend server
- [ ] Cleared browser cache

### Information to Provide

When asking for help, include:
1. **Error message** (full text)
2. **What you tried** (steps to reproduce)
3. **Environment** (OS, Go version, browser)
4. **Backend logs** (relevant parts)
5. **Browser console** (errors)
6. **Code changes** (if any)

### Contact Points

- **GitHub Issues:** Create detailed issue
- **Stack Overflow:** Tag with #golang, #paystack, #postgresql
- **Discord/Slack:** Developer communities

---

## ✅ Quick Fixes Summary

| Issue | Quick Fix |
|-------|-----------|
| Backend won't start | Check DATABASE_URL, restart |
| Menu not loading | Verify API_URL, check backend running |
| CORS error | Update AllowedOrigins in main.go |
| Images not showing | Check uploads folder exists |
| Payment failing | Verify Paystack keys, use test payment |
| Cart not saving | Clear localStorage, check console |
| Admin can't login | Use admin/admin123, clear session |
| Tracking not found | Verify tracking ID case-sensitive |

**Most issues resolve with: Restart backend + Hard refresh browser** 🔄
