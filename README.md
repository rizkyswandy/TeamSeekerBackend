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
    connString := "postgres://ILB:@localhost:5432/team_seeker?sslmode=disable"
     ```

## Running the Application

1. Start the server:
```bash
go run cmd/server/main.go
```

2. (Optional) Generate test data:
```bash
go run cmd/generator/main.go
```
Note: Make sure to update the connection string in `cmd/generator/main.go` as well if you use the data generator. Currently it's generating 10000 data in a blink of an eye. I tried to get that 10.000 data in Postman and it's only 65ms.

## API Endpoints

### Student Profiles
- `POST /api/profiles` - Create new profile
- `GET /api/profiles` - Get all profiles
- `GET /api/profiles/{id}` - Get profile by ID
- `PUT /api/profiles/{id}` - Update profile
- `DELETE /api/profiles/{id}` - Delete profile
- `GET /api/profiles/search` - Search profiles with filters

### Search Filters Example
```json
{
    "faculty": "Computer Science",
    "skills": ["Go", "PostgreSQL"],
    "focus": ["Backend Development"],
    "availability": true
}
```

## Testing

You can test the API using any HTTP client (e.g., Postman, cURL). Example of creating a profile:

```bash
curl -X POST http://localhost:8080/api/profiles \
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

## Troubleshooting

1. Database connection issues:
   - Verify PostgreSQL is running
   - Check connection string in `main.go`
   - Ensure database and table exist

2. Common errors:
   - "connection refused" - Check if PostgreSQL is running
   - "role does not exist" - Verify PostgreSQL username
   - "database does not exist" - Create the database first

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
```