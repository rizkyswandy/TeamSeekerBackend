# TeamSeeker Backend

LFG! Ngoding Go!

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher

## Database Setup

1. Create PostgreSQL database:
```sql
CREATE DATABASE team_seeker;
```

2. Create the student profiles table:
```sql
CREATE TABLE student_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    faculty VARCHAR(100) NOT NULL,
    field_of_study VARCHAR(100) NOT NULL,
    semester INTEGER NOT NULL,
    skills TEXT[] NOT NULL,
    focus TEXT[] NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

## Project Setup

1. Clone the repository:
```bash
git clone https://github.com/rizkyswandy/TeamSeekerBackend.git
cd TeamSeekerBackend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure database connection:
   - Open `cmd/server/main.go`
   - Update the connection string with your PostgreSQL credentials:
     ```go
     connString := "postgres://username:password@localhost:5432/team_seeker?sslmode=disable"
     ```
    - I didn't setup any password so for me it's: 
    ```go
    connString := "postgres://postgres:@localhost:5432/team_seeker?sslmode=disable"
     ```

## Running the Application
```bash
go run cmd/server/main.go
```

## API Endpoints

### Student Profiles
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/profiles` - Create new profile
- `GET /api/profiles` - Get all profiles
- `GET /api/profiles/{id}` - Get profile by ID
- `PUT /api/profiles/{id}` - Update profile
- `DELETE /api/profiles/{id}` - Delete profile
- `GET /api/profiles/search` - Search profiles with filters

## Testing

You can test the API using any HTTP client (e.g., Postman, cURL). Example of creating a profile:

```bash
curl -X POST http://localhost:3001/api/profiles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "faculty": "Computer Science",
    "field_of_study": "Software Engineering",
    "semester": 5,
    "skills": ["Go", "PostgreSQL", "REST API"],
    "focus": ["Backend Development"],
    "is_available": true
  }'
```

## Project Structure
```
TeamSeeker/
├── api/            # API definitions and interfaces
├── cmd/
│   ├── server/    # Main application
│   └── generator/ # Test data generator
├── internal/
│   └── database/  # Database implementations
└── middleware/    # Custom middleware
└── types/         # Common used type
```