# 🚀 Quick Start Guide - Lumière Restaurant

Get your restaurant website up and running in 10 minutes!

## ⚡ Fast Setup (3 Steps)

### Step 1: Setup Neon PostgreSQL (2 minutes)

1. Go to https://neon.tech/ and sign up
2. Create a new project called "restaurant-db"
3. Copy your connection string - it looks like:
   ```
   postgresql://user:pass@ep-xxx.us-east-2.aws.neon.tech/dbname?sslmode=require
   ```
4. Save it - you'll need it in Step 2

### Step 2: Setup paystack (2 minutes)

1. Go to https://dashboard.paystack.com/ and sign up
2. Navigate to **Developers → API Keys**
3. Copy both keys:
   - **Publishable key** (starts with `pk_test_`)
   - **Secret key** (starts with `sk_test_`)
4. Save them - you'll need them next

### Step 3: Run the Application (5 minutes)

#### A. Start the Backend

```bash
# Navigate to your project
cd restaurant-website/backend

# Set environment variables (Mac/Linux)
export DATABASE_URL="your_neon_connection_string_from_step1"
export STRIPE_SECRET_KEY="your_stripe_secret_key_from_step2"

# Or for Windows
set DATABASE_URL=your_neon_connection_string_from_step1
set STRIPE_SECRET_KEY=your_stripe_secret_key_from_step2

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

You should see:
```
Database connected successfully!
Tables created successfully!
Server starting on port 8080
```

#### B. Start the Frontend

**Open a new terminal:**

```bash
cd restaurant-website/frontend
python -m http.server 3000
```

Visit: `http://localhost:3000`

#### C. Update Stripe Publishable Key

Open `frontend/checkout.html` and find line ~32:

```javascript
const STRIPE_PUBLISHABLE_KEY = 'pk_test_YOUR_KEY_HERE';
```

Replace with your Stripe **Publishable key** from Step 2.

## ✅ Test Everything

### 1. Test the Menu
- Visit: `http://localhost:3000/menu.html`
- Menu should be empty (we'll add items next)

### 2. Test Admin Dashboard
- Visit: `http://localhost:3000/admin/login.html`
- Login with:
  - Username: `admin`
  - Password: `admin123`

### 3. Add Your First Menu Item
1. In Admin Dashboard, fill the form:
   - **Title:** "Truffle Risotto"
   - **Description:** "Creamy arborio rice with black truffle"
   - **Price:** 28
   - **Category:** Main Courses
   - **Ingredients:** "Arborio Rice, Truffle, Parmesan"
   - **Image:** Upload any food image
2. Click "Add Menu Item"
3. Visit menu page - your item should appear!

### 4. Test Complete Flow
1. Go to `http://localhost:3000/menu.html`
2. Add items to cart
3. Click cart icon → View Cart
4. Proceed to Checkout
5. Fill details:
   - Card: `4242 4242 4242 4242`
   - Expiry: Any future date
   - CVC: Any 3 digits
6. Complete payment
7. Note your tracking ID
8. Visit tracking page and enter your tracking ID

### 5. Test Order Tracking Update
1. Go to Admin Dashboard
2. Click "Order Tracking" tab
3. Select the order you just made
4. Update status to "Preparing"
5. Add a message like "Your order is being prepared"
6. Check the tracking page - it should update!

## 🎉 You're Done!

Your restaurant website is now fully functional with:
- ✅ Dynamic menu from database
- ✅ Shopping cart
- ✅ Stripe payments
- ✅ Order tracking
- ✅ Admin dashboard

## 🔧 Common Issues

**Backend won't start:**
```bash
# Make sure Go is installed
go version

# If not installed, download from: https://golang.org/dl/

# Try cleaning the module cache
go clean -modcache
go mod tidy
```

**"Database connection failed":**
- Double-check your DATABASE_URL
- Make sure it includes `?sslmode=require` at the end
- Test connection in Neon dashboard

**Menu not loading:**
- Make sure backend is running (check terminal)
- Check browser console (F12) for errors
- Verify API_URL in menu-dynamic.js points to `http://localhost:8080/api`

**Payment not working:**
- Verify you updated STRIPE_PUBLISHABLE_KEY in checkout.html
- Use test card: 4242 4242 4242 4242
- Check Stripe dashboard for payment attempts

**Images not showing:**
- Make sure uploads folder exists in backend directory
- Check file permissions
- Try uploading a smaller image (< 2MB)

## 📁 Project Structure Reminder

```
restaurant-website/
├── frontend/
│   ├── index.html          # Home page
│   ├── menu.html           # Menu (add menu-dynamic.js script)
│   ├── cart.html           # Shopping cart
│   ├── checkout.html       # Stripe payment
│   ├── tracking.html       # Order tracking
│   └── js/
│       └── menu-dynamic.js # NEW - add this file
│
├── admin/
│   ├── login.html          # NEW - Admin login
│   └── admin.html          # NEW - Admin dashboard
│
└── backend/
    ├── main.go             # NEW - Go server
    ├── go.mod              # NEW - Dependencies
    └── uploads/            # Created automatically
```

## 🔐 Update Menu.html

Add this script tag to your `menu.html` before the closing `</body>` tag:

```html
<script src="js/menu-dynamic.js"></script>
```

## 🎨 Customize

- **Colors:** Edit CSS variables in `style.css`
- **Logo:** Replace "Lumière" in navigation
- **Menu Categories:** Add more in admin dropdown
- **Stripe Currency:** Change in `main.go` (line with `Currency:`)

## 📞 Need Help?

1. Check the full README.md for detailed docs
2. Review the browser console (F12) for errors
3. Check backend terminal for error messages
4. Test with Stripe test cards: https://stripe.com/docs/testing

---

**You're all set! Start adding your menu items and taking orders! 🎊**