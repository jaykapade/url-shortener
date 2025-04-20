# URL Shortener Service

A modern URL shortener service built with Go, featuring user authentication, custom short URLs, and click tracking.

## Features

- User authentication with JWT
- URL shortening with Base62 encoding
- Click tracking for shortened URLs
- RESTful API endpoints
- PostgreSQL database for data persistence
- Sonyflake ID generation for unique short codes

## Tech Stack

- Go 1.23.2
- Chi router for HTTP routing
- PostgreSQL for database
- JWT for authentication
- Sonyflake for unique ID generation
- SQLC for type-safe SQL
- Docker (optional for development)

## Prerequisites

- Go 1.23.2 or later
- PostgreSQL
- Git

## Setup

1. Clone the repository:

```bash
git clone https://github.com/jaykapade/url-shortener.git
cd url-shortener
```

2. Set up the environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

The `.env` file should contain the following variables:

```env
DB_URL=postgres://username:password@localhost:5432/url_shortener
JWT_SECRET=your_jwt_secret_key
```

Replace the values with your actual database credentials and JWT secret. The `DB_URL` follows the PostgreSQL connection string format:

- `username`: Your PostgreSQL username
- `password`: Your PostgreSQL password
- `localhost`: Database host (change if using a remote database)
- `5432`: Database port (default PostgreSQL port)
- `url_shortener`: Database name

3. Create a PostgreSQL database:

```sql
CREATE DATABASE url_shortener;
```

4. Run database migrations:

```bash
# Using goose (make sure goose is installed)
goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/url_shortener" up
```

5. Install dependencies:

```bash
go mod download
```

6. Run the server:

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

#### Register a new user

```http
POST /register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "your_password"
}
```

#### Login

```http
POST /login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "your_password"
}
```

### URL Shortening

#### Create a short URL

```http
POST /shortener
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
    "full_url": "https://example.com/very/long/url"
}
```

#### Access a shortened URL

```http
GET /{short_code}
```

## Security

- Passwords are hashed using bcrypt
- JWT tokens are used for authentication
- Database connections are secured with PostgreSQL's native security features

## Development

### Database Schema

The project uses two main tables:

1. `users` - Stores user information

   - `id` (UUID)
   - `email` (VARCHAR)
   - `password` (VARCHAR)
   - `created_at` (TIMESTAMP)

2. `links` - Stores shortened URLs
   - `id` (UUID)
   - `short_code` (VARCHAR)
   - `full_url` (TEXT)
   - `user_id` (UUID, foreign key to users)
   - `created_at` (TIMESTAMP)
   - `click_count` (INT)

### Code Generation

The project uses SQLC for generating type-safe database queries. To regenerate the database code:

```bash
sqlc generate
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
