# Lumière Restaurant Website - Full Stack Application

A complete restaurant website with Go backend, PostgreSQL database, Paystack payment integration, and order tracking system.

## 🏗️ Project Structure

```
restaurant-website/
├── frontend/              # Your existing HTML/CSS/JS
│   ├── index.html
│   ├── menu.html
│   ├── about.html
│   ├── contact.html
│   ├── packages.html
│   ├── cart.html         # NEW: Shopping cart
│   ├── checkout.html     # NEW: Paystack payment
│   ├── tracking.html     # NEW: Order tracking
│   ├── css/
│   │   └── style.css
│   └── js/
│       ├── script.js
│       └── menu-dynamic.js  # NEW: Dynamic menu loading
│
├── admin/                 # NEW: Admin dashboard
│   └── admin.html
│
└── backend/              # NEW: Go API
    ├── main.go
    ├── go.mod
    └── uploads/          # Created automatically
```

## 🚀 Setup Instructions

### 1. Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **PostgreSQL** (Neon) - [Sign up](https://neon.tech/)
- **Paystack Account** - [Sign up](https://paystack.com/)

### 2. Database Setup (Neon PostgreSQL)

1. Create a Neon PostgreSQL database at https://neon.tech/
2. Copy your connection string (it looks like: `postgresql://user:pass@host/dbname`)
3. The backend will automatically create the necessary tables

### 3. Backend Setup

```bash
# Navigate to backend directory
cd backend

# Initialize Go module (if go.mod doesn't exist)
go mod init restaurant-backend

# Install dependencies
go get github.com/gorilla/mux
go get github.com/lib/pq
go get github.com/rs/cors
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv

# Set environment variables
export DATABASE_URL="your_neon_postgresql_connection_string"
export PAYSTACK_SECRET_KEY="sk_test_your_paystack_secret_key"
export PORT="8080"

# Run the server
go run main.go
```

**For Windows:**
```cmd
set DATABASE_URL=your_neon_postgresql_connection_string
set PAYSTACK_SECRET_KEY=sk_test_your_paystack_secret_key
set PORT=8080
go run main.go
```

### 4. Get Paystack API Key

1. Go to https://dashboard.paystack.com/
2. Navigate to **Developers > API Keys**
3. Copy your **Secret key** (starts with `sk_test_` for testing or `sk_live_` for production)
4. Set it as the backend `PAYSTACK_SECRET_KEY` environment variable

### 5. Frontend Setup

Update API_URL in these files (change `localhost:8080` if needed):
- `frontend/js/menu-dynamic.js`
- `frontend/cart.html`
- `frontend/checkout.html`
- `frontend/tracking.html`
- `admin/admin.html`

**Option 1: Simple Python Server**
```bash
cd frontend
python -m http.server 3000
```

**Option 2: Live Server (VS Code Extension)**
- Install "Live Server" extension
- Right-click on `index.html` > Open with Live Server

Visit: `http://localhost:3000`

### 6. Admin Dashboard

Visit: `http://localhost:3000/admin/admin.html`

**Features:**
- Add menu items with images
- View all orders
- Update order tracking status
- Delete menu items

## 📋 Features

### Customer Features
✅ Dynamic menu loaded from database  
✅ Shopping cart with localStorage  
✅ Paystack payment integration  
✅ Order tracking with real-time updates  
✅ Responsive design  

### Admin Features
✅ Add/delete menu items with image upload  
✅ View all orders  
✅ Update order status and tracking  
✅ Real-time order management  

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the backend directory:

```env
DATABASE_URL=postgresql://user:pass@host/database?sslmode=require
PAYSTACK_SECRET_KEY=sk_test_xxxxxxxxxxxxx
PORT=8080
```

### CORS Configuration

If your frontend is on a different domain, update CORS in `main.go`:

```go
AllowedOrigins: []string{"http://localhost:3000", "https://yourdomain.com"}
```

## 🗄️ Database Schema

### Tables Created Automatically:

**menu_items**
- id (SERIAL PRIMARY KEY)
- title (VARCHAR)
- description (TEXT)
- price (DECIMAL)
- category (VARCHAR)
- image_url (TEXT)
- ingredients (TEXT)
- created_at (TIMESTAMP)

**orders**
- id (SERIAL PRIMARY KEY)
- customer_name (VARCHAR)
- customer_email (VARCHAR)
- customer_phone (VARCHAR)
- items (TEXT - JSON)
- total_amount (DECIMAL)
- status (VARCHAR)
- tracking_id (VARCHAR UNIQUE)
- payment_status (VARCHAR)
- created_at (TIMESTAMP)

**order_tracking**
- id (SERIAL PRIMARY KEY)
- order_id (INTEGER FK)
- status (VARCHAR)
- location (VARCHAR)
- update_message (TEXT)
- updated_at (TIMESTAMP)

## 🎯 API Endpoints

### Menu
- `GET /api/menu` - Get all menu items
- `GET /api/menu?category=starters` - Filter by category
- `POST /api/menu` - Add menu item (Admin)
- `DELETE /api/menu/{id}` - Delete menu item (Admin)

### Orders
- `POST /api/orders` - Create order
- `GET /api/orders` - Get all orders (Admin)
- `GET /api/orders/tracking/{tracking_id}` - Get order by tracking ID
- `POST /api/orders/{id}/tracking` - Update tracking (Admin)

### Payment
- `POST /api/payment/initialize` - Initialize a Paystack transaction
- `GET /api/payment/verify?reference={reference}` - Verify a Paystack transaction

## 🧪 Testing

### Test Paystack Payments

Use Paystack test payment details from your Paystack dashboard or Paystack documentation.

## 📱 Usage Flow

### Customer Journey:
1. Browse menu at `/menu.html`
2. Add items to cart
3. View cart at `/cart.html`
4. Checkout at `/checkout.html`
5. Make payment with Paystack
6. Receive tracking ID
7. Track order at `/tracking.html`

### Admin Workflow:
1. Login to admin at `/admin/admin.html`
2. Add menu items with images
3. View incoming orders
4. Update order status and tracking
5. Customer sees updates in real-time

## 🐛 Troubleshooting

**Backend won't start:**
- Check DATABASE_URL is correct
- Ensure PostgreSQL is accessible
- Verify Go dependencies are installed

**Menu not loading:**
- Check API_URL in frontend files
- Verify backend is running on correct port
- Check browser console for errors

**Payment failing:**
- Verify Paystack keys are correct
- Check you're using test keys for development
- Ensure HTTPS in production

**Images not showing:**
- Check uploads folder permissions
- Verify image paths in database
- Check CORS settings

## 🚀 Deployment

### Backend (Heroku/Railway/Render)
1. Set environment variables
2. Deploy Go application
3. Update frontend API_URL to deployed backend URL

### Frontend (Netlify/Vercel)
1. Build static site
2. Deploy frontend files
3. Update CORS in backend to allow frontend domain

## 📄 License

This project is for educational purposes.

## 👨‍💻 Support

For issues or questions, please check:
- Go documentation: https://golang.org/doc/
- Paystack docs: https://paystack.com/docs
- Neon docs: https://neon.tech/docs

---

**Happy Coding! 🎉**
