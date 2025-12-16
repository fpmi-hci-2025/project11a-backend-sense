# Manual API Testing

This directory contains Go-based test client for comprehensive API endpoint testing.

## Prerequisites

1. Docker and Docker Compose installed
2. Go 1.25 or later
3. Services running via `docker-compose up -d`
4. API accessible at `http://localhost:8080`

## Structure

- `client/` - HTTP client for API requests
- `testdata/` - Test data management (tokens, IDs)
- `testers/` - Test functions for each endpoint group
  - `auth.go` - Authentication endpoints
  - `publication.go` - Publication endpoints
  - `comment.go` - Comment endpoints
  - `feed.go` - Feed endpoints
  - `profile.go` - Profile endpoints
  - `media.go` - Media endpoints
- `main.go` - Main test runner
- `test_data.json` - Test data storage (auto-generated)
- `test_results.md` - Test results and status
- `errors.md` - Detailed error documentation

## Usage

### Build

```bash
cd tests/manual
go build -o test_runner .
```

### Run All Tests

```bash
./test_runner
```

### Run Specific Test Groups

```bash
# Skip auth tests
./test_runner -skip-auth

# Run only auth and publication tests
./test_runner -skip-comments -skip-feed -skip-profile -skip-media

# Custom API URL
./test_runner -url http://localhost:8080
```

### Options

- `-url` - API base URL (default: http://localhost:8080)
- `-data` - Test data file path (default: test_data.json)
- `-skip-auth` - Skip authentication tests
- `-skip-publications` - Skip publication tests
- `-skip-comments` - Skip comment tests
- `-skip-feed` - Skip feed tests
- `-skip-profile` - Skip profile tests
- `-skip-media` - Skip media tests

## Test Data

Test data is automatically stored in `test_data.json`. Tokens and IDs are populated during test execution and reused across test runs.

## Output

- Console output shows progress and results for each test
- `test_data.json` is updated with tokens and IDs after each test run
- Review `test_results.md` and `errors.md` for detailed documentation

## Example Output

```
========================================
Sense API Manual Testing Suite
========================================
Base URL: http://localhost:8080

=== Testing Auth Endpoints ===
1. Testing POST /auth/register (user1)
   ✓ User1 registered: ID=..., Token=...
...
```
