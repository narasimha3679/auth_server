# [Auth Server](https://narasimha3679.github.io/auth_server)

A robust authentication server built with Go, similar to Auth0, providing JWT-based authentication and user management capabilities.

## Features

- User registration and authentication
- JWT token-based authorization with refresh tokens
- Protected routes with middleware
- PostgreSQL database integration
- Password hashing with bcrypt
- RESTful API endpoints

## Project Structure

```
auth-server/
├── config/
│   └── database.go     # Database configuration and connection
├── handlers/
│   └── auth.go         # Authentication handlers (register, login)
├── middleware/
│   └── auth.go         # JWT authentication middleware
├── models/
│   └── user.go         # User model definition
├── utils/
│   └── jwt.go          # JWT utility functions
├── .env                # Environment configuration
├── go.mod              # Go module file
├── main.go             # Application entry point
└── README.md           # Documentation
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Git (optional)

## Dependencies

- gin-gonic/gin: Web framework
- golang-jwt/jwt: JWT token handling
- gorm.io/gorm: ORM library
- gorm.io/driver/postgres: PostgreSQL driver
- golang.org/x/crypto: Password hashing
- joho/godotenv: Environment configuration

## Installation

1. Clone the repository (optional):
```bash
git clone <repository-url>
cd auth-server
```

2. Install dependencies:
```bash
make deps
```

3. Configure environment variables:
Create a `.env` file in the root directory with the following content:
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=auth_db
DB_PORT=5432
JWT_SECRET=your_jwt_secret_key
PORT=8080
```

4. Create PostgreSQL database:
```bash
make create-db
```

## Development Commands

The project includes a Makefile with common development commands:

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies
- `make create-db` - Create PostgreSQL database
- `make drop-db` - Drop PostgreSQL database
- `make help` - Show available commands

## Running the Server

Start the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in your .env file).

## API Endpoints

### Public Routes

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
	"email": "user@example.com",
	"password": "password123",
	"first_name": "John",
	"last_name": "Doe"
}
```

#### Login User
```http
POST /auth/login
Content-Type: application/json

{
	"email": "user@example.com",
	"password": "password123"
}
```

Response:
```json
{
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Response:
```json
{
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Protected Routes

All protected routes require the JWT access token (not refresh token) in the Authorization header:
```http
Authorization: Bearer <your-access-token>
```

Note: Refresh tokens cannot be used to access protected routes. They are only for obtaining new access tokens through the /auth/refresh endpoint.

#### Get User Profile
```http
GET /api/profile
```

Response:
```json
{
	"user": {
		"id": 1,
		"email": "user@example.com",
		"first_name": "John",
		"last_name": "Doe",
		"last_login": "2024-01-01T00:00:00Z",
		"is_active": true,
		"email_verified": false
	}
}
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- 200: Success
- 201: Created (successful registration)
- 400: Bad Request (invalid input)
- 401: Unauthorized (invalid credentials or token)
- 404: Not Found
- 500: Internal Server Error

## Frontend Implementation Guide

### Authentication Flow

1. **Initial Login**
   - User logs in with email/password
   - Server returns both access token (24h) and refresh token (permanent)
   - Store both tokens securely (localStorage/secure storage)
   ```javascript
   const response = await fetch('/auth/login', {
	   method: 'POST',
	   headers: { 'Content-Type': 'application/json' },
	   body: JSON.stringify({ email, password })
   });
   const { access_token, refresh_token } = await response.json();
   localStorage.setItem('access_token', access_token);
   localStorage.setItem('refresh_token', refresh_token);
   ```

2. **Making API Requests**
   - Use access token for all API requests
   - Handle 401 errors by refreshing the token
   ```javascript
   async function apiRequest(url, options = {}) {
	   // Add access token to headers
	   const headers = {
		   ...options.headers,
		   'Authorization': `Bearer ${localStorage.getItem('access_token')}`
	   };

	   let response = await fetch(url, { ...options, headers });
	   
	   // If token expired, refresh it and retry
	   if (response.status === 401) {
		   const refreshed = await refreshAccessToken();
		   if (refreshed) {
			   // Retry with new access token
			   headers.Authorization = `Bearer ${localStorage.getItem('access_token')}`;
			   response = await fetch(url, { ...options, headers });
		   }
	   }
	   
	   return response;
   }
   ```

3. **Token Refresh Mechanism**
   - Use refresh token to get new access token
   - Original refresh token remains valid
   ```javascript
   async function refreshAccessToken() {
	   try {
		   const response = await fetch('/auth/refresh', {
			   method: 'POST',
			   headers: { 'Content-Type': 'application/json' },
			   body: JSON.stringify({
				   refresh_token: localStorage.getItem('refresh_token')
			   })
		   });

		   if (response.ok) {
			   const { access_token } = await response.json();
			   localStorage.setItem('access_token', access_token);
			   return true;
		   }
		   
		   // If refresh fails, redirect to login
		   window.location.href = '/login';
		   return false;
	   } catch (error) {
		   console.error('Token refresh failed:', error);
		   return false;
	   }
   }
   ```

4. **Logout Handling**
   ```javascript
   function logout() {
	   localStorage.removeItem('access_token');
	   localStorage.removeItem('refresh_token');
	   window.location.href = '/login';
   }
   ```

### Key Implementation Notes

1. **Token Storage**
   - Access token: Short-lived (24 hours)
   - Refresh token: Never expires (permanent session)
   - Store tokens securely based on platform:
	 - Web: localStorage/sessionStorage
	 - Mobile: Secure storage/Keychain
	 - Desktop: Encrypted storage

2. **Error Handling**
   - 401 responses indicate expired access token
   - Failed refresh attempts indicate invalid refresh token
   - Handle network errors appropriately

3. **Security Considerations**
   - Store tokens securely
   - Use HTTPS for all API requests
   - Clear tokens on logout
   - Implement rate limiting for refresh attempts

## Security Features

1. Password Hashing
   - Passwords are hashed using bcrypt before storage
   - Secure comparison during login

2. JWT Authentication
   - Access tokens expire after 24 hours
   - Refresh tokens never expire (permanent session)
   - Tokens are signed with a secret key
   - Contains user ID for authentication
   - Automatic token refresh mechanism

3. Protected Routes
   - Middleware verification of JWT tokens
   - Token validation and parsing
   - User context injection

## Database Schema

The User model includes:
- ID (uint, primary key)
- Email (string, unique)
- Password (string, hashed)
- FirstName (string)
- LastName (string)
- LastLogin (timestamp)
- IsActive (boolean)
- EmailVerified (boolean)
- CreatedAt (timestamp)
- UpdatedAt (timestamp)
- DeletedAt (soft delete)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## Testing with Postman

The project includes a Postman collection (`auth_server.postman_collection.json`) for testing the API endpoints. To use it:

1. Import the collection into Postman
2. Set up environment variables:
   - Create a new environment in Postman
   - Add variable `jwt_token` (will be automatically updated after login)
3. Use the collection to test:
   - Register a new user
   - Login and get JWT token
   - Access protected routes

The collection includes all available endpoints with example payloads.

## License

This project is licensed under the MIT License.
