# Auth Server

A robust authentication server built with Go, similar to Auth0, providing JWT-based authentication and user management capabilities.

## Features

- User registration and authentication
- JWT token-based authorization
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
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Protected Routes

All protected routes require the JWT token in the Authorization header:
```http
Authorization: Bearer <your-jwt-token>
```

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

## Security Features

1. Password Hashing
   - Passwords are hashed using bcrypt before storage
   - Secure comparison during login

2. JWT Authentication
   - Tokens expire after 24 hours
   - Signed with a secret key
   - Contains user ID for authentication

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