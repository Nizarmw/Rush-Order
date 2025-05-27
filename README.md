# Rush Order

A modern QR code-based restaurant ordering system built with Go backend and vanilla JavaScript frontend.

## 🚀 Overview

Rush Order is a web-based restaurant ordering system that allows customers to scan QR codes at their tables, browse menus, place orders, and make payments directly from their smartphones. The system includes both customer-facing and admin interfaces for complete restaurant management.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [User Personas](#user-personas)
- [Contributing](#contributing)

## ✨ Features

### Customer Features
- **QR Code Scanning**: Each table has a unique QR code for instant access
- **Digital Menu**: Browse complete menu with photos, descriptions, and prices
- **Shopping Cart**: Add items, modify quantities, and review orders
- **Contactless Payment**: Integrated payment gateway supporting cards and e-wallets
- **Digital Receipt**: Download or print order confirmation

### Admin Features
- **Order Management**: Real-time order monitoring and processing
- **Dashboard**: View pending and completed orders
- **Order Tracking**: Mark orders as completed
- **Session Management**: Secure admin authentication

## 🛠 Tech Stack

### Backend
- **Go** - Main backend language
- **Gin** - Web framework
- **Supabase** - Database and authentication
- **GORM** - ORM for database operations

### Frontend
- **HTML5/CSS3** - Structure and styling
- **Vanilla JavaScript** - Client-side functionality
- **Responsive Design** - Mobile-first approach

### Database
- **PostgreSQL** (via Supabase)

## 📁 Project Structure

```
Rush-Order/
├── back-end/
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── config/
│   │   └── db.go
│   ├── controller/
│   │   ├── admin_controller.go
│   │   ├── cart_controller.go
│   │   ├── CustomerSession_controller.go
│   │   └── produk_controller.go
│   ├── middleware/
│   │   └── auth.go
│   ├── models/
│   │   ├── order_item.go
│   │   ├── order.go
│   │   ├── pegawai.go
│   │   ├── pemesan.go
│   │   └── produk.go
│   ├── routes/
│   │   ├── admin_routes.go
│   │   ├── cart_routes.go
│   │   ├── product_routes.go
│   │   └── session_routes.go
│   ├── service/
│   │   ├── admin_service.go
│   │   ├── cart_service.go
│   │   ├── produk_service.go
│   │   ├── session_service.go
│   │   └── supabase.go
│   └── session/
│       └── types.go
├── front-end/
│   ├── index.html
│   ├── login.html
│   ├── admin/
│   │   ├── index.html
│   │   └── login.html
│   ├── assets/
│   │   └── sounds/
│   │       ├── notification.mp3
│   │       └── success.mp3
│   ├── css/
│   │   ├── admin-style.css
│   │   └── style.css
│   └── js/
│       ├── app-admin.js
│       └── app.js
└── README.md
```

## 🚀 Installation

### Prerequisites
- Go 1.19 or higher
- Node.js (for development tools)
- Supabase account

### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/your-username/Rush-Order.git
cd Rush-Order/back-end
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
# Create .env file with your Supabase credentials
SUPABASE_URL=your_supabase_url
SUPABASE_KEY=your_supabase_key
```

4. Run the server:
```bash
go run main.go
```

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd ../front-end
```

2. Serve the files using a local server:
```bash
# Using Python
python -m http.server 8080

# Using Node.js (with live-server)
npx live-server
```

## 🎯 Usage

### For Customers

1. **Scan QR Code**: Use your smartphone camera to scan the QR code on your table
2. **Enter Name**: Provide your name for order identification
3. **Browse Menu**: Explore the digital menu with detailed descriptions
4. **Add to Cart**: Select items and quantities
5. **Checkout**: Review your order and proceed to payment
6. **Payment**: Complete payment using your preferred method
7. **Confirmation**: Receive digital receipt and wait for your order

### For Restaurant Staff

1. **Login**: Access the admin panel with your credentials
2. **Monitor Orders**: View real-time incoming orders
3. **Process Orders**: Prepare orders based on the queue
4. **Mark Complete**: Update order status when delivered
5. **Track History**: Review completed orders and daily reports

## 📱 User Personas

### Persona 1: Yudi Wahyudi (Office Worker)
- **Background**: Limited lunch break time
- **Goals**: Quick ordering without queuing
- **Frustrations**: Long wait times during peak hours

### Persona 2: Arya Aryanto (Student)
- **Background**: Busy schedule, unfamiliar with local menus
- **Goals**: Clear menu information, accurate orders, cashless payment
- **Frustrations**: Order mistakes due to unclear communication

## 🌐 API Documentation

### Customer Endpoints
- `GET /` - Landing page
- `GET /menu` - Menu display
- `POST /cart` - Add items to cart
- `POST /checkout` - Process payment
- `GET /order/{id}` - Order details

### Admin Endpoints
- `POST /admin/login` - Admin authentication
- `GET /admin/dashboard` - Order management dashboard
- `PUT /admin/orders/{id}/complete` - Mark order as done
- `GET /admin/orders/pending` - Pending orders
- `GET /admin/orders/completed` - Completed orders

## 🔒 Security Features

- Session-based authentication
- Role-based access control
- Secure payment integration
- Data encryption for sensitive information

## 🚀 Future Enhancements

- Real-time notifications with WebSocket
- Mobile app development
- Multi-location support
- Advanced analytics dashboard
- Customer loyalty program
- Inventory management system

## 👥 Team

**Team Name**: Omke Gas


**Rush Order** - Revolutionizing restaurant dining experience through technology 🍽️📱