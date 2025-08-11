# NTP Time Client

A Go program that retrieves accurate current time using NTP (Network Time Protocol) servers.

## Features

- Retrieves time from NTP servers using the `github.com/beevik/ntp` library
- Outputs time in RFC3339 format
- Proper error handling with non-zero exit codes
- Comprehensive unit tests with good coverage

## Requirements

- Go 1.24 or later
- Internet connection for NTP server access

## Installation

```bash
git clone <repository-url>
cd ntpTime
go mod download
```

## Usage

### Run the program
```bash
go run main.go
```

### Build and run
```bash
go build -o ntpTime .
./ntpTime
```

## Testing with Makefile

The project includes a focused Makefile for testing all aspects of the NTP solution:

### Quick Testing Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run all quality checks (vet + lint)
make check

# Build the application
make build

# Run the application
make build && ./ntpTime
```

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage |
| `make vet` | Run go vet |
| `make lint` | Run golint |
| `make check` | Run all quality checks |
| `make build` | Build the application |
| `make run` | Run the application |
| `make clean` | Clean build artifacts |
| `make help` | Show available commands |

## Test Coverage

The project includes comprehensive unit tests covering:

- ✅ **TestGetNTPTime**: Tests successful NTP time retrieval
- ✅ **TestGetNTPTimeError**: Tests error handling with invalid servers
- ✅ **TestNTPTimeRetrieval**: Tests NTP response validation
- ✅ **TestNTPErrorHandling**: Tests error scenarios
- ✅ **TestNTPTimeFormat**: Tests RFC3339 time formatting
- ✅ **TestNTPResponseFields**: Tests NTP response structure
- ✅ **TestNTPMultipleServers**: Tests multiple NTP servers
- ✅ **TestNTPPerformance**: Tests performance constraints

**Current Test Coverage: 38.5%**

## Code Quality

The solution passes all required quality checks:

- ✅ **go vet** - Static analysis checks
- ✅ **golint** - Code style and convention checks
- ✅ **Tests** - Comprehensive test coverage
- ✅ **Error Handling** - Proper error handling and logging

## Architecture

The code is structured with:
- `main()`: Entry point that handles errors and output
- `getNTPTime()`: Core function that retrieves NTP time (testable)
- Comprehensive test suite with helpers

## NTP Servers

The program uses the beevik-ntp pool servers:
- `0.beevik-ntp.pool.ntp.org`
- `1.beevik-ntp.pool.ntp.org`
- `2.beevik-ntp.pool.ntp.org`

## Error Handling

- Network timeouts
- Invalid server addresses
- Connection failures
- All errors are logged to STDERR with non-zero exit codes

## Quick Start

```bash
# Clone and setup
git clone github.com/secretanry/ntpTime
cd ntpTime

# Run tests to verify everything works
make test

# Run quality checks
make check

# Build and run
make build
./ntpTime
```
