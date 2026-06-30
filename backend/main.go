package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net" // Added for IP tracking
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync" // Added for thread safety
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/time/rate" // Added for Rate Limiting
)

var db *sql.DB

// --- RATE LIMITING CONFIGURATION START ---
// We create a "visitor" struct to track each user's request rate
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Map to hold visitors (IP addresses) and a mutex to lock it safely
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// Run a background worker to clean up old visitors (saves memory)
func init() {
	go cleanupVisitors()
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		// LIMIT SETTINGS:
		// 2 requests per second allowed on average
		// Burst of 5 requests allowed at once (good for loading images/scripts quickly)
		limiter := rate.NewLimiter(2, 5)

		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update last seen time
	v.lastSeen = time.Now()
	return v.limiter
}

// Remove old IP addresses every minute to free up memory
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// Middleware: The "Bouncer" that checks every request
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user's IP address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests - Please wait a moment", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// --- RATE LIMITING CONFIGURATION END ---

// --- JWT CONFIGURATION ---
// We declare the variable here, but we assign the value in main()
var jwtKey []byte

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// User represents an account in the system (Admin, Customer, or Rider)
type User struct {
	ID         int       `json:"id"`
	GoogleID   string    `json:"google_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      *string   `json:"phone"`
	Role       string    `json:"role"` // 'admin', 'customer', 'rider'
	AvatarURL  string    `json:"avatar_url"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
}

// Rider contains extra operational details linked to a User profile
type Rider struct {
	UserID          int    `json:"user_id"`
	VehicleType     string `json:"vehicle_type"`
	VehiclePlate    string `json:"vehicle_plate"`
	IsAvailable     bool   `json:"is_available"`
	TotalDeliveries int    `json:"total_deliveries"`
}

// GoogleClaims maps the incoming token data from a Google OAuth authentication success
type GoogleClaims struct {
	ID            string `json:"sub"` // FIXED: Google uses "sub" for unique account IDs
	Email         string `json:"email"`
	VerifiedEmail string `json:"email_verified"` // OPTIONAL FIX: Google sends this as a string/bool field "email_verified"
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// Models
type MenuItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	Ingredients string    `json:"ingredients"`
	CreatedAt   time.Time `json:"created_at"`
}

type Order struct {
	ID            int       `json:"id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`
	Items         string    `json:"items"` // JSON string of items
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"`
	TrackingID    string    `json:"tracking_id"`
	PaymentStatus string    `json:"payment_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type OrderTracking struct {
	OrderID       int       `json:"order_id"`
	Status        string    `json:"status"`
	Location      string    `json:"location"`
	UpdateMessage string    `json:"update_message"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Initialize database
func initDB() {
	var err error
	// Neon PostgreSQL connection string
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set. Make sure .env is loaded.")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	log.Println("Database connected successfully!")
	createTables()
}

// FindOrCreateUser checks if a Google user exists; if not, inserts them into the DB.
func FindOrCreateUser(claims GoogleClaims, assignedRole string) (*User, error) {
	var user User

	// 1. Check if user already exists by their unique Google ID
	query := `SELECT id, google_id, name, email, phone, role, avatar_url, is_approved, created_at 
	          FROM users WHERE google_id = $1`
	err := db.QueryRow(query, claims.ID).Scan(
		&user.ID, &user.GoogleID, &user.Name, &user.Email,
		&user.Phone, &user.Role, &user.AvatarURL, &user.IsApproved, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		// 2. User doesn't exist, let's create a new master user account
		insertUserQuery := `
			INSERT INTO users (google_id, name, email, role, avatar_url, is_approved)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, google_id, name, email, phone, role, avatar_url, is_approved, created_at`

		// By default, customers and admins are approved instantly. Riders might await admin confirmation.
		isApproved := true
		if assignedRole == "rider" {
			isApproved = false // Set to false if you want an admin approval step for riders
		}

		err = db.QueryRow(insertUserQuery, claims.ID, claims.Name, claims.Email, assignedRole, claims.Picture, isApproved).Scan(
			&user.ID, &user.GoogleID, &user.Name, &user.Email,
			&user.Phone, &user.Role, &user.AvatarURL, &user.IsApproved, &user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error inserting new user: %v", err)
		}

		// 3. Special rule: If they registered as a rider, provision an accompanying rider metrics row
		if assignedRole == "rider" {
			insertRiderQuery := `INSERT INTO riders (user_id, is_available, total_deliveries) VALUES ($1, false, 0)`
			_, err = db.Exec(insertRiderQuery, user.ID)
			if err != nil {
				return nil, fmt.Errorf("error provisioning rider details row: %v", err)
			}
		}

		return &user, nil
	} else if err != nil {
		return nil, fmt.Errorf("database query error during authentication: %v", err)
	}

	// User exists, return the found record
	return &user, nil
}

// Create tables
func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS menu_items (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			category VARCHAR(100),
			image_url TEXT,
			ingredients TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			customer_name VARCHAR(255) NOT NULL,
			customer_email VARCHAR(255) NOT NULL,
			customer_phone VARCHAR(50),
			items TEXT NOT NULL,
			total_amount DECIMAL(10, 2) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			tracking_id VARCHAR(100) UNIQUE NOT NULL,
			payment_status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS order_tracking (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id),
			status VARCHAR(100),
			location VARCHAR(255),
			update_message TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating table: %v", err)
		}
	}
	log.Println("Tables created successfully!")
}

// --- AUTHENTICATION HANDLERS START ---

// HandleGoogleAuthCallback processes the token sent from frontend Google Sign-In
func HandleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the incoming request payload
	var requestData struct {
		Credential string `json:"credential"`
		Role       string `json:"role"` // 'customer' or 'rider' sent from frontend selection
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	if requestData.Credential == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Google credential token missing"})
		return
	}

	// Default to customer role if nothing valid was submitted
	assignedRole := strings.ToLower(requestData.Role)
	if assignedRole != "customer" && assignedRole != "rider" {
		assignedRole = "customer"
	}

	// 2. Validate the token directly with Google's Token Verification API
	googleTokenURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", requestData.Credential)
	resp, err := http.Get(googleTokenURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to connect to Google validation service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or expired Google token"})
		return
	}

	// 3. Unpack the validated response into our GoogleClaims struct
	var claims GoogleClaims
	err = json.NewDecoder(resp.Body).Decode(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read profile details from Google"})
		return
	}

	// 4. Save/Fetch user information using our database layer function
	user, err := FindOrCreateUser(claims, assignedRole)
	if err != nil {
		log.Printf("DB Error tracking auth user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error managing user account profile"})
		return
	}

	// 5. Generate our own local BelleFood JWT access token for safety
	expirationTime := time.Now().Add(24 * time.Hour)
	tokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not sign application login session token"})
		return
	}

	// 6. Return successful session payload response to the frontend client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"role":        user.Role,
			"avatar_url":  user.AvatarURL,
			"is_approved": user.IsApproved,
		},
	})
}

func adminLogin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check against Environment Variables
	expectedUser := os.Getenv("ADMIN_USERNAME")
	expectedPass := os.Getenv("ADMIN_PASSWORD")

	if expectedUser == "" || expectedPass == "" {
		http.Error(w, "Server configuration error: Admin credentials not set in .env", http.StatusInternalServerError)
		return
	}

	if creds.Username != expectedUser || creds.Password != expectedPass {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create the JWT Token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// Middleware to verify JWT token
func isAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// --- AUTHENTICATION HANDLERS END ---

// API Handlers

// Get all menu items
func getMenuItems(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = db.Query("SELECT id, title, description, price, category, image_url, ingredients, created_at FROM menu_items WHERE category = $1 ORDER BY id DESC", category)
	} else {
		rows, err = db.Query("SELECT id, title, description, price, category, image_url, ingredients, created_at FROM menu_items ORDER BY id DESC")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []MenuItem{}
	for rows.Next() {
		var item MenuItem
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Price, &item.Category, &item.ImageURL, &item.Ingredients, &item.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Add menu item (Admin)
func addMenuItem(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	category := r.FormValue("category")
	ingredients := r.FormValue("ingredients")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	// Handle file upload
	file, handler, err := r.FormFile("image")
	var imageURL string

	if err == nil {
		defer file.Close()

		// Create uploads directory if it doesn't exist
		uploadsDir := "./uploads"
		os.MkdirAll(uploadsDir, os.ModePerm)

		// Create unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
		filepath := filepath.Join(uploadsDir, filename)

		// Save file
		dst, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "Unable to save image", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		io.Copy(dst, file)
		imageURL = "/uploads/" + filename
	}

	// Insert into database
	var id int
	err = db.QueryRow(
		"INSERT INTO menu_items (title, description, price, category, image_url, ingredients) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		title, description, price, category, imageURL, ingredients,
	).Scan(&id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      id,
		"message": "Menu item added successfully",
	})
}

// Delete menu item (Admin)
func deleteMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM menu_items WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Initialize payment with Paystack
func initializePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount float64 `json:"amount"`
		Email  string  `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Paystack expects amount in kobo (smallest currency unit)
	amountInKobo := int(req.Amount * 100)

	// Create Paystack payment
	paystackKey := os.Getenv("PAYSTACK_SECRET_KEY")

	paymentData := map[string]interface{}{
		"email":    req.Email,
		"amount":   amountInKobo,
		"currency": "NGN", // Explicitly using Naira
	}

	jsonData, _ := json.Marshal(paymentData)

	client := &http.Client{}
	paystackReq, err := http.NewRequest("POST", "https://api.paystack.co/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	paystackReq.Header.Set("Authorization", "Bearer "+paystackKey)
	paystackReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(paystackReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Verify Paystack payment
func verifyPayment(w http.ResponseWriter, r *http.Request) {
	reference := r.URL.Query().Get("reference")

	paystackKey := os.Getenv("PAYSTACK_SECRET_KEY")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.paystack.co/transaction/verify/"+reference, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+paystackKey)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Create order
func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate tracking ID
	trackingID := fmt.Sprintf("ORD-%d", time.Now().Unix())
	order.TrackingID = trackingID
	order.Status = "pending"

	err := db.QueryRow(
		`INSERT INTO orders (customer_name, customer_email, customer_phone, items, total_amount, tracking_id, payment_status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		order.CustomerName, order.CustomerEmail, order.CustomerPhone, order.Items, order.TotalAmount, order.TrackingID, order.PaymentStatus,
	).Scan(&order.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create initial tracking entry
	_, trackErr := db.Exec(
		"INSERT INTO order_tracking (order_id, status, location, update_message) VALUES ($1, $2, $3, $4)",
		order.ID, "Order Placed", "Restaurant", "Your order has been received and is being prepared.",
	)
	if trackErr != nil {
		http.Error(w, "Failed to create tracking entry: "+trackErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// Get order by tracking ID
func getOrderByTracking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackingID := vars["tracking_id"]

	var order Order
	err := db.QueryRow(
		"SELECT id, customer_name, customer_email, customer_phone, items, total_amount, status, tracking_id, payment_status, created_at FROM orders WHERE tracking_id = $1",
		trackingID,
	).Scan(&order.ID, &order.CustomerName, &order.CustomerEmail, &order.CustomerPhone, &order.Items, &order.TotalAmount, &order.Status, &order.TrackingID, &order.PaymentStatus, &order.CreatedAt)

	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Get tracking history
	rows, err := db.Query(
		"SELECT status, location, update_message, updated_at FROM order_tracking WHERE order_id = $1 ORDER BY updated_at DESC",
		order.ID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tracking []OrderTracking
	for rows.Next() {
		var t OrderTracking
		rows.Scan(&t.Status, &t.Location, &t.UpdateMessage, &t.UpdatedAt)
		tracking = append(tracking, t)
	}

	response := map[string]interface{}{
		"order":    order,
		"tracking": tracking,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get all orders (Admin)
func getAllOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		"SELECT id, customer_name, customer_email, customer_phone, items, total_amount, status, tracking_id, payment_status, created_at FROM orders ORDER BY created_at DESC",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.CustomerName, &order.CustomerEmail, &order.CustomerPhone, &order.Items, &order.TotalAmount, &order.Status, &order.TrackingID, &order.PaymentStatus, &order.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Update order tracking (Admin)
func updateOrderTracking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var tracking OrderTracking
	if err := json.NewDecoder(r.Body).Decode(&tracking); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update order status
	_, err := db.Exec("UPDATE orders SET status = $1 WHERE id = $2", tracking.Status, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add tracking entry
	_, err = db.Exec(
		"INSERT INTO order_tracking (order_id, status, location, update_message) VALUES ($1, $2, $3, $4)",
		orderID, tracking.Status, tracking.Location, tracking.UpdateMessage,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	// Fix: Load JWT key inside main()
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtKey) == 0 {
		log.Println("Warning: JWT_SECRET_KEY is empty!")
	}

	initDB()

	router := mux.NewRouter()

	// --- APPLY RATE LIMITING ---
	// Wrap all routes with the rate limiter
	router.Use(rateLimitMiddleware)

	// --- PUBLIC API ROUTES ---
	router.HandleFunc("/api/menu", getMenuItems).Methods("GET")
	router.HandleFunc("/api/orders", createOrder).Methods("POST")
	router.HandleFunc("/api/orders/tracking/{tracking_id}", getOrderByTracking).Methods("GET")
	router.HandleFunc("/api/payment/initialize", initializePayment).Methods("POST")
	router.HandleFunc("/api/payment/verify", verifyPayment).Methods("GET")

	// Google Auth Sign-In Endpoint Route
	router.HandleFunc("/api/auth/google", HandleGoogleAuthCallback).Methods("POST")

	// Login Route (Returns the JWT Token)
	router.HandleFunc("/api/login", adminLogin).Methods("POST")

	// --- PROTECTED ADMIN ROUTES (Require JWT) ---
	router.HandleFunc("/api/menu", isAuthorized(addMenuItem)).Methods("POST")
	router.HandleFunc("/api/menu/{id}", isAuthorized(deleteMenuItem)).Methods("DELETE")
	router.HandleFunc("/api/orders", isAuthorized(getAllOrders)).Methods("GET")
	router.HandleFunc("/api/orders/{id}/tracking", isAuthorized(updateOrderTracking)).Methods("POST")

	// Serve uploaded images
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	// Serve Admin Files
	router.PathPrefix("/admin/").Handler(http.StripPrefix("/admin/", http.FileServer(http.Dir("../admin"))))
	// Serve Frontend files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend")))

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
