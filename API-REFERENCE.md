# 📚 API Reference

Complete documentation for the Lumière Restaurant API.

**Base URL:** `http://localhost:8080/api`  
**Production URL:** `https://your-domain.com/api`

## Authentication

Currently, admin endpoints use session-based authentication (to be upgraded to JWT).

## Response Format

All API responses return JSON:

### Success Response
```json
{
  "id": 1,
  "title": "Truffle Risotto",
  "price": 28.00
}
```

### Error Response
```json
{
  "error": "Menu item not found",
  "code": 404
}
```

---

## 📋 Menu Endpoints

### Get All Menu Items

```http
GET /api/menu
```

**Description:** Retrieve all menu items from database

**Query Parameters:**
- `category` (optional) - Filter by category (starters, main, specials, desserts, drinks)

**Request:**
```bash
curl http://localhost:8080/api/menu
curl http://localhost:8080/api/menu?category=starters
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "title": "Truffle Risotto",
    "description": "Creamy arborio rice with black truffle shavings",
    "price": 28.00,
    "category": "main",
    "image_url": "/uploads/1702483821_risotto.jpg",
    "ingredients": "Arborio Rice, Black Truffle, Parmesan",
    "created_at": "2024-12-10T10:30:21Z"
  },
  {
    "id": 2,
    "title": "Beef Carpaccio",
    "description": "Thinly sliced raw beef with arugula",
    "price": 18.00,
    "category": "starters",
    "image_url": "/uploads/1702483822_carpaccio.jpg",
    "ingredients": "Beef, Arugula, Parmesan, Truffle Oil",
    "created_at": "2024-12-10T11:15:33Z"
  }
]
```

**Possible Errors:**
- `500` - Database connection error

---

### Add Menu Item (Admin)

```http
POST /api/menu
```

**Description:** Create a new menu item with image upload

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `title` (required) - Item name
- `description` (required) - Item description
- `price` (required) - Price in NGN (decimal)
- `category` (required) - Category name
- `ingredients` (optional) - Comma-separated list
- `image` (required) - Image file (JPG, PNG, GIF, WebP)

**Request:**
```bash
curl -X POST http://localhost:8080/api/menu \
  -F "title=Wagyu Steak" \
  -F "description=Premium cut grilled to perfection" \
  -F "price=45.00" \
  -F "category=main" \
  -F "ingredients=Wagyu Beef, Butter, Rosemary" \
  -F "image=@/path/to/steak.jpg"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "id": 3,
  "message": "Menu item added successfully"
}
```

**Possible Errors:**
- `400` - Invalid price format
- `400` - Missing required fields
- `500` - Unable to save image
- `500` - Database error

**Image Requirements:**
- Max size: 10MB
- Formats: .jpg, .jpeg, .png, .gif, .webp
- Stored in: `backend/uploads/`

---

### Delete Menu Item (Admin)

```http
DELETE /api/menu/{id}
```

**Description:** Remove a menu item from database

**URL Parameters:**
- `id` (required) - Menu item ID

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/menu/3
```

**Response:** `200 OK`
```json
{
  "success": true
}
```

**Possible Errors:**
- `404` - Menu item not found
- `500` - Database error

---

## 🛒 Order Endpoints

### Create Order

```http
POST /api/orders
```

**Description:** Create a new customer order

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "customer_phone": "+1234567890",
  "items": "[{\"id\":1,\"title\":\"Truffle Risotto\",\"price\":28,\"quantity\":2}]",
  "total_amount": 56.00,
  "payment_status": "paid"
}
```

**Request:**
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_phone": "+1234567890",
    "items": "[{\"id\":1,\"title\":\"Truffle Risotto\",\"price\":28,\"quantity\":2}]",
    "total_amount": 56.00,
    "payment_status": "paid"
  }'
```

**Response:** `200 OK`
```json
{
  "id": 15,
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "customer_phone": "+1234567890",
  "items": "[{\"id\":1,\"title\":\"Truffle Risotto\",\"price\":28,\"quantity\":2}]",
  "total_amount": 56.00,
  "status": "pending",
  "tracking_id": "ORD-1702483821",
  "payment_status": "paid",
  "created_at": "2024-12-10T12:30:21Z"
}
```

**Notes:**
- `tracking_id` is auto-generated
- Initial `status` is "pending"
- `items` must be a JSON string

**Possible Errors:**
- `400` - Invalid request data
- `500` - Database error

---

### Get All Orders (Admin)

```http
GET /api/orders
```

**Description:** Retrieve all customer orders

**Request:**
```bash
curl http://localhost:8080/api/orders
```

**Response:** `200 OK`
```json
[
  {
    "id": 15,
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_phone": "+1234567890",
    "items": "[{\"id\":1,\"title\":\"Truffle Risotto\",\"price\":28,\"quantity\":2}]",
    "total_amount": 56.00,
    "status": "preparing",
    "tracking_id": "ORD-1702483821",
    "payment_status": "paid",
    "created_at": "2024-12-10T12:30:21Z"
  },
  {
    "id": 14,
    "customer_name": "Jane Smith",
    "customer_email": "jane@example.com",
    "customer_phone": "+1234567891",
    "items": "[{\"id\":2,\"title\":\"Beef Carpaccio\",\"price\":18,\"quantity\":1}]",
    "total_amount": 18.00,
    "status": "delivered",
    "tracking_id": "ORD-1702483720",
    "payment_status": "paid",
    "created_at": "2024-12-10T11:20:00Z"
  }
]
```

**Possible Errors:**
- `500` - Database error

---

### Get Order by Tracking ID

```http
GET /api/orders/tracking/{tracking_id}
```

**Description:** Retrieve order details and tracking history

**URL Parameters:**
- `tracking_id` (required) - Order tracking ID (e.g., ORD-1702483821)

**Request:**
```bash
curl http://localhost:8080/api/orders/tracking/ORD-1702483821
```

**Response:** `200 OK`
```json
{
  "order": {
    "id": 15,
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_phone": "+1234567890",
    "items": "[{\"id\":1,\"title\":\"Truffle Risotto\",\"price\":28,\"quantity\":2}]",
    "total_amount": 56.00,
    "status": "preparing",
    "tracking_id": "ORD-1702483821",
    "payment_status": "paid",
    "created_at": "2024-12-10T12:30:21Z"
  },
  "tracking": [
    {
      "order_id": 15,
      "status": "preparing",
      "location": "Kitchen",
      "update_message": "Your order is being prepared by our chef",
      "updated_at": "2024-12-10T12:35:00Z"
    },
    {
      "order_id": 15,
      "status": "pending",
      "location": "Restaurant",
      "update_message": "Your order has been received",
      "updated_at": "2024-12-10T12:30:21Z"
    }
  ]
}
```

**Notes:**
- Tracking array is ordered by most recent first
- Each order starts with "pending" status

**Possible Errors:**
- `404` - Order not found
- `500` - Database error

---

### Update Order Tracking (Admin)

```http
POST /api/orders/{id}/tracking
```

**Description:** Add a tracking update for an order

**URL Parameters:**
- `id` (required) - Order ID (not tracking ID)

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "status": "preparing",
  "location": "Kitchen",
  "update_message": "Your order is being prepared by our chef"
}
```

**Status Values:**
- `pending` - Order received
- `confirmed` - Order confirmed
- `preparing` - Being prepared
- `ready` - Ready for pickup/delivery
- `out_for_delivery` - Out for delivery
- `delivered` - Delivered

**Request:**
```bash
curl -X POST http://localhost:8080/api/orders/15/tracking \
  -H "Content-Type: application/json" \
  -d '{
    "status": "out_for_delivery",
    "location": "On Route",
    "update_message": "Your order is on its way to you!"
  }'
```

**Response:** `200 OK`
```json
{
  "success": true
}
```

**Notes:**
- Also updates the `status` field in orders table
- Creates new record in order_tracking table
- Timestamp is added automatically

**Possible Errors:**
- `400` - Invalid request data
- `404` - Order not found
- `500` - Database error

---

## 💳 Payment Endpoints

### Initialize Paystack Payment

```http
POST /api/payment/initialize
```

**Description:** Initialize a Paystack transaction for checkout

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "amount": 5600,
  "email": "john@example.com"
}
```

**Request:**
```bash
curl -X POST http://localhost:8080/api/payment/initialize \
  -H "Content-Type: application/json" \
  -d '{"amount": 5600, "email": "john@example.com"}'
```

**Response:** `200 OK`
```json
{
  "status": true,
  "message": "Authorization URL created",
  "data": {
    "authorization_url": "https://checkout.paystack.com/...",
    "access_code": "access_code",
    "reference": "reference"
  }
}
```

**Notes:**
- Amount is in naira on the frontend; the backend converts it to kobo before sending it to Paystack.
- The frontend redirects the customer to `data.authorization_url`.
- Currency: NGN.

**Possible Errors:**
- `400` - Invalid amount
- `500` - Paystack API error

**Frontend Usage:**
```javascript
// 1. Initialize Paystack transaction
const response = await fetch('/api/payment/initialize', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ amount: 5600, email: 'john@example.com' })
});
const data = await response.json();

// 2. Redirect customer to Paystack checkout
window.location.href = data.data.authorization_url;
```

### Verify Paystack Payment

```http
GET /api/payment/verify?reference={reference}
```

**Description:** Verify a Paystack transaction after the customer returns from checkout.

**Request:**
```bash
curl "http://localhost:8080/api/payment/verify?reference=PAYSTACK_REFERENCE"
```

**Response:** Paystack verification response.

---

## 📁 Static File Endpoints

### Serve Uploaded Images

```http
GET /uploads/{filename}
```

**Description:** Serve uploaded menu item images

**URL Parameters:**
- `filename` - Image filename

**Example:**
```
http://localhost:8080/uploads/1702483821_risotto.jpg
```

**Response:** Image file (JPG, PNG, GIF, WebP)

**Possible Errors:**
- `404` - File not found

**Notes:**
- Files are stored in `backend/uploads/`
- Filenames include timestamp to prevent conflicts
- Format: `{timestamp}_{originalname}`

---

## 📊 Data Models

### MenuItem

```typescript
{
  id: number;              // Auto-generated
  title: string;           // Max 255 chars
  description: string;     // Text
  price: number;           // Decimal(10,2)
  category: string;        // Max 100 chars
  image_url: string;       // Relative path
  ingredients: string;     // Comma-separated
  created_at: string;      // ISO 8601 timestamp
}
```

### Order

```typescript
{
  id: number;              // Auto-generated
  customer_name: string;   // Max 255 chars
  customer_email: string;  // Max 255 chars
  customer_phone: string;  // Max 50 chars
  items: string;           // JSON array as string
  total_amount: number;    // Decimal(10,2)
  status: string;          // Order status
  tracking_id: string;     // Unique, format: ORD-{timestamp}
  payment_status: string;  // pending/paid/failed
  created_at: string;      // ISO 8601 timestamp
}
```

### OrderTracking

```typescript
{
  order_id: number;        // Foreign key to orders
  status: string;          // Status value
  location: string;        // Current location
  update_message: string;  // Human-readable message
  updated_at: string;      // ISO 8601 timestamp
}
```

---

## 🔒 Security Considerations

### CORS

Currently allows all origins (`*`) for development:

```go
AllowedOrigins: []string{"*"}
```

**Production:** Restrict to your domains:

```go
AllowedOrigins: []string{
    "https://yourdomain.com",
    "https://www.yourdomain.com"
}
```

### Input Validation

All endpoints validate:
- Required fields presence
- Data types (string, number, decimal)
- Email format
- File types (images only)

### SQL Injection Protection

Uses parameterized queries:

```go
db.QueryRow("SELECT * FROM orders WHERE tracking_id = $1", trackingID)
```

---

## 📝 Example Workflows

### Complete Order Flow

```bash
# 1. Customer browses menu
curl http://localhost:8080/api/menu

# 2. Initialize Paystack payment
curl -X POST http://localhost:8080/api/payment/initialize \
  -H "Content-Type: application/json" \
  -d '{"amount": 5600, "email": "john@example.com"}'

# 3. After successful payment, create order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_phone": "+1234567890",
    "items": "[{\"id\":1,\"title\":\"Truffle Risotto\",\"price\":28,\"quantity\":2}]",
    "total_amount": 56.00,
    "payment_status": "paid"
  }'
# Returns: tracking_id = "ORD-1702483821"

# 4. Customer tracks order
curl http://localhost:8080/api/orders/tracking/ORD-1702483821
```

### Admin Workflow

```bash
# 1. View all orders
curl http://localhost:8080/api/orders

# 2. Update order status
curl -X POST http://localhost:8080/api/orders/15/tracking \
  -H "Content-Type: application/json" \
  -d '{
    "status": "preparing",
    "location": "Kitchen",
    "update_message": "Your order is being prepared"
  }'

# 3. Add new menu item
curl -X POST http://localhost:8080/api/menu \
  -F "title=New Dish" \
  -F "description=Amazing food" \
  -F "price=25.00" \
  -F "category=main" \
  -F "image=@dish.jpg"

# 4. Delete menu item
curl -X DELETE http://localhost:8080/api/menu/3
```

---

## 🧪 Testing

### Paystack Test Payments

Use Paystack test payment details from your Paystack dashboard or Paystack documentation.

### Postman Collection

Import this collection for easy testing:

```json
{
  "info": {
    "name": "Lumière Restaurant API"
  },
  "item": [
    {
      "name": "Get Menu",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/api/menu"
      }
    }
  ]
}
```

---

## 📞 Support

- **Issues:** Check TROUBLESHOOTING.md
- **Documentation:** See README.md
- **API Questions:** Review this document

---

**API Version:** 1.0  
**Last Updated:** December 2024
