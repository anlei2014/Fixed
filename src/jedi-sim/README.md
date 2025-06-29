# Jedi-Sim Backend System

## Overview

Jedi-Sim is a Go-based backend system that provides error code reporting and management functionality through HTTP APIs. The system integrates with CAN bus communication for real-time error code processing and maintains error code definitions in JSON format.

## System Architecture

```
jedi-sim/
├── main.go                    # Main application entry point
├── errorCodes.json           # Error code definitions database
├── go.mod                    # Go module dependencies
├── static/                   # Frontend static files
│   ├── index.html           # Main web interface
│   ├── instruction.html     # Documentation page
│   ├── css/                 # Stylesheets
│   └── js/                  # JavaScript files
├── internal/                 # Internal application logic
│   ├── handler/             # HTTP request handlers
│   ├── model/               # Data models and structures
│   └── service/             # Business logic services
├── jediSim/                 # Core JEDI simulation logic
├── msgHandler/              # Message handling utilities
└── README.md               # This documentation
```

## Core Components

### 1. Main Application (`main.go`)

**Purpose**: Application entry point and HTTP server initialization

**Key Functions**:
- `main()`: Application startup and goroutine management
- `startHTTPServer()`: HTTP server configuration and startup

**Server Configuration**:
- **Port**: 9103
- **Static Files**: Served from `./static/` directory
- **API Routes**: `/api/report` and `/api/add-error-code`
- **Root Redirect**: `/` → `/static/index.html`

### 2. HTTP Handlers (`internal/handler/`)

#### Report Handler (`report.go`)

**Endpoint**: `POST /api/report`

**Purpose**: Process error code reports and generate CAN messages

**Request Format**:
```http
POST /api/report
Content-Type: application/x-www-form-urlencoded

errorcode=804
```

**Process Flow**:
1. Validate request method (POST only)
2. Extract error code from form data
3. Validate error code format and existence
4. Convert error code to integer
5. Generate CAN message using `jediSim.GeneratorStatusErrorCode()`
6. Return success/failure response with CAN message data

**Response Format**:
```json
{
  "status": true,
  "code": "804",
  "message": "OK, CAN message has been sent",
  "canMsg": [6, 30, 0, 34, 2, 0, 44, 3, 5, 0]
}
```

#### Add Error Code Handler (`add_error_code.go`)

**Endpoint**: `POST /api/add-error-code`

**Purpose**: Add new error code definitions to the JSON database

**Request Format**:
```http
POST /api/add-error-code
Content-Type: application/json

{
  "Z0": 6,
  "Z1": 30,
  "Z2": 0,
  "Z3_Phase": 2,
  "Z3_Class": 2,
  "Z4Z5_ErrorCode": 804,
  "Z6Z7_ErrorData": 5,
  "Description": "Tube spit (all kV drop/regul errors)"
}
```

**Process Flow**:
1. Validate request method (POST only)
2. Parse JSON request body into `ErrorInfo` struct
3. Read existing `errorCodes.json` file
4. Add new error code using `Z4Z5_ErrorCode` as key
5. Write updated data back to file
6. Return success/failure response

**Response Format**:
```json
{
  "status": true,
  "message": "Error code data saved successfully"
}
```

### 3. Data Models (`internal/model/`)

#### ErrorInfo Structure (`errorCodesInfo.go`)

**Purpose**: Define error code data structure

```go
type ErrorInfo struct {
    Z0             int    // Generator Status (e.g., 0x06 = Error condition)
    Z1             int    // Simplified Error Code (e.g., 30)
    Z2             int    // Display Bitmap (bit 0-7)
    Z3_Phase       int    // Phase (high 4 bits)
    Z3_Class       int    // Error Class (low 4 bits)
    Z4Z5_ErrorCode int    // Generator Error Code (little-endian format)
    Z6Z7_ErrorData int    // Error Data (little-endian format)
    Description    string // Error description
}
```

**Key Functions**:
- `LoadErrorInfoFromJSON(code int, path string)`: Load error info by code from JSON file
- `LoadErrorCodes()`: Initialize error codes (currently placeholder)

### 4. Business Logic (`internal/service/`)

#### Validator Service (`validator.go`)

**Purpose**: Validate error codes and provide business logic

**Key Functions**:
- `ValidateErrorCode(code string)`: Validate error code format and existence
- `LoadErrorInfoFromJSON(code int, path string)`: Load error info with fallback paths

**Validation Process**:
1. Trim whitespace from input
2. Convert string to integer
3. Check existence in `errorCodes.json`
4. Return validation result and message

### 5. JEDI Simulation (`jediSim/`)

#### Core Functions (`jediApp.go`)

**Purpose**: Handle CAN message generation and JEDI simulation logic

**Key Functions**:
- `GeneratorStatusErrorCode(msg [MSG_LENGTH]int, errorCode int)`: Generate CAN messages for error codes

**Process Flow**:
1. Create base message using "NOTIFY_JEDI_STATUS" template
2. Load error information from JSON database
3. Set CAN message fields:
   - Z0: Generator Status
   - Z1: Simplified Error Code
   - Z2: Display Bitmap
   - Z3: Phase and Error Class
   - Z4/Z5: Generator Error Code
   - Z6/Z7: Error Data
4. Send message via CAN bus
5. Return success status and generated message

## File Paths and Configuration

### JSON Database
- **Primary Path**: `./errorCodes.json`
- **Fallback Path**: `src/jedi-sim/errorCodes.json`
- **Format**: JSON object with error codes as keys

### Static Files
- **Directory**: `./static/`
- **Main Page**: `./static/index.html`
- **Documentation**: `./static/instruction.html`
- **Stylesheets**: `./static/css/`
- **JavaScript**: `./static/js/`

### CAN Message Configuration
- **Message Length**: 10 bytes (`MSG_LENGTH = 10`)
- **Template**: "NOTIFY_JEDI_STATUS"
- **Field Mapping**: Z0-Z7 fields mapped to CAN message bytes

## Deployment and Path Configuration

### Build and Deploy Structure

When building the application with `go build .` in the jedi-sim directory, the recommended deployment structure is:

```
your-deployment-directory/
├── jedi-sim.exe          # Compiled executable file
├── static/               # Static files directory
│   ├── index.html
│   ├── instruction.html
│   ├── css/
│   │   ├── style.css
│   │   └── instruction-btn.css
│   └── js/
│       └── main.js
└── errorCodes.json       # Error codes database file
```

### Path Handling Recommendations

#### Current Path Issues

The current code uses development-time paths that need adjustment for production deployment:

**Current Paths in Code**:
```go
// main.go - Static files
staticDir := fmt.Sprintf("%s/static", pwd)

// validator.go - JSON file
path = cwd + "/src/jedi-sim/errorCodes.json"

// jediApp.go - JSON file
errInfo, ok := model.LoadErrorInfoFromJSON(errorCode, "src/jedi-sim/errorCodes.json")
```

#### Recommended Path Solutions

**Option 1: Relative Paths (Recommended for Simple Deployment)**
```go
// Static files
staticDir := "./static"

// JSON file
filePath := "./errorCodes.json"
```

**Option 2: Executable Directory (Recommended for Production)**
```go
import (
    "os"
    "path/filepath"
)

// Get executable directory
execPath, _ := os.Executable()
execDir := filepath.Dir(execPath)

// Static files
staticDir := filepath.Join(execDir, "static")

// JSON file
filePath := filepath.Join(execDir, "errorCodes.json")
```

**Option 3: Unified Path Resolver (Recommended for Flexibility)**
```go
func getResourcePath(filename string) string {
    // 1. Try relative to executable
    execPath, _ := os.Executable()
    execDir := filepath.Dir(execPath)
    execFilePath := filepath.Join(execDir, filename)
    
    if _, err := os.Stat(execFilePath); err == nil {
        return execFilePath
    }
    
    // 2. Try current working directory
    cwd, _ := os.Getwd()
    cwdPath := filepath.Join(cwd, filename)
    
    if _, err := os.Stat(cwdPath); err == nil {
        return cwdPath
    }
    
    // 3. Return default path
    return filename
}
```

### Files Requiring Path Updates

The following files need path configuration updates for production deployment:

1. **main.go** - Static file serving path
2. **internal/service/validator.go** - JSON file loading path
3. **jediSim/jediApp.go** - JSON file loading path
4. **internal/handler/add_error_code.go** - JSON file writing path

### Deployment Scenarios

#### Development Environment
- Use relative paths: `./static/`, `./errorCodes.json`
- Run from project directory

#### Production Environment
- Use executable directory paths
- Deploy as system service
- Ensure proper file permissions

#### Containerized Deployment
- Use absolute paths or environment variables
- Mount volumes for persistent data

## API Endpoints Summary

| Endpoint | Method | Purpose | Input | Output |
|----------|--------|---------|-------|--------|
| `/` | GET | Redirect to main page | - | Redirect to `/static/index.html` |
| `/static/*` | GET | Serve static files | - | HTML/CSS/JS files |
| `/api/report` | POST | Report error code | `errorcode` (form) | JSON response with CAN message |
| `/api/add-error-code` | POST | Add error code | JSON `ErrorInfo` | JSON success/failure response |

## Error Handling

### HTTP Status Codes
- `200 OK`: Successful operation
- `400 Bad Request`: Invalid input data
- `405 Method Not Allowed`: Wrong HTTP method
- `500 Internal Server Error`: Server-side error

### Validation Errors
- Invalid error code format
- Error code not found in database
- Invalid JSON data
- File I/O errors

## Development and Deployment

### Prerequisites
- Go 1.16+
- CAN bus interface (for production)

### Building the Application
```bash
cd src/jedi-sim
go build .
```

### Running the Application
```bash
# Development
go run main.go

# Production (after building)
./jedi-sim.exe
```

### Access Points
- **Web Interface**: http://localhost:9103
- **API Base**: http://localhost:9103/api/

### Logging
- **Format**: Date, time, microseconds, file:line
- **Level**: Debug, Info, Error
- **Output**: Console

## Data Flow

### Error Code Reporting Flow
1. **Frontend**: User submits error code via web form
2. **HTTP Handler**: `/api/report` receives POST request
3. **Validation**: Service validates error code format and existence
4. **CAN Generation**: JEDI simulation generates CAN message
5. **Response**: JSON response with status and CAN message data

### Error Code Addition Flow
1. **Frontend**: User submits error code data via web form
2. **HTTP Handler**: `/api/add-error-code` receives POST request
3. **JSON Parsing**: Parse request body into ErrorInfo struct
4. **File Update**: Read existing JSON, add new entry, write back
5. **Response**: JSON success/failure response

## Security Considerations

- Input validation for all user inputs
- JSON parsing error handling
- File I/O error handling
- HTTP method validation
- Content-Type validation

## Future Enhancements

- Database integration (currently file-based)
- Authentication and authorization
- Rate limiting
- Enhanced error logging
- CAN bus monitoring and diagnostics 