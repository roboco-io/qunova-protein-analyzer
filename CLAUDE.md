# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Protein Structure Analyzer** system designed to parse and analyze protein structure files (PDB, CIF, mmCIF). The system consists of:
- **Backend**: Go-based REST API using Gin framework
- **Frontend**: React component with Tailwind CSS (to be migrated to Next.js for MVP)
- **Parser**: Python BioPython for actual protein file parsing

## Development Methodology

This project follows the **Tidy First** methodology for backend development:

### Core Principles
1. **Separation of Changes**: Never mix structural changes and behavioral changes in a single commit
2. **Structure Before Behavior**: Always perform structural cleanup (refactoring) before adding new functionality
3. **TDD Cycle**: Strictly follow Test-Driven Development with clear test → code → refactor cycles
4. **Clean Commits**: Only commit when all tests pass, no warnings exist, and work represents a logical unit

### Workflow Guidelines
- **Before adding features**: First identify and clean up any structural issues
- **Commit Strategy**: Separate commits for:
  - Structural changes (renaming, extracting methods, removing duplication)
  - Behavioral changes (new features, bug fixes)
- **Testing First**: Write failing tests before implementation
- **Minimal Implementation**: Use simplest code that passes tests
- **Refactor After**: Clean up only after tests are green

### Code Quality Standards
- Remove code duplication immediately when spotted
- Use clear, intentional naming for variables, functions, and types
- Keep functions small with single responsibilities
- Maintain code clarity over cleverness

## Architecture

### Core Components Interaction
```
Frontend (React/Next.js) → Backend API (Go/Gin) → Python Parser (BioPython)
                                ↓
                         Storage Service
                                ↓
                         Analysis Service
```

### Backend Service Architecture
- **Handlers** (`internal/handlers/`): HTTP request handling, response formatting
- **Services** (`internal/services/`): Business logic layer
  - `StorageService`: File management and persistence
  - `ParserService`: Interfaces with Python parser or mock data
  - `AnalysisService`: Orchestrates analysis workflow, maintains state
- **Models** (`internal/models/`): Data structures for protein analysis
- **Middleware** (`internal/middleware/`): Error handling, logging

### Key Design Decisions
1. **Mock Parser Mode**: Currently uses mock data by default (`USE_MOCK_PARSER=true`) for rapid development without BioPython dependency
2. **Async Processing**: Analysis runs asynchronously with progress tracking
3. **Microservice Ready**: Python parser as separate process, easily extractable to standalone service
4. **In-Memory Storage**: Analysis results stored in memory (suitable for MVP, replace with DB for production)

## Development Commands

### Backend Development
```bash
cd backend

# Initial setup (installs dependencies, creates .env)
make setup

# Run locally (port 8080)
make run

# Development with hot reload
make dev

# Run tests (always run before commits per Tidy First)
make test

# Format and lint code (structural changes)
make fmt
make lint

# Build binary
make build
```

### TDD + Tidy First Development Workflow

#### Complete TDD Cycle for New Features
```bash
# RED: Write a failing test first
# Example: Adding new file validation feature
cat > internal/handlers/upload_test.go << 'EOF'
func TestValidateFileFormat_InvalidExtension(t *testing.T) {
    handler := NewUploadHandler(nil)
    
    // Test should fail initially
    valid := handler.validateFileFormat("invalid.txt")
    assert.False(t, valid, "should reject invalid file extension")
}
EOF

# Run the test - it should FAIL (RED)
go test -v ./internal/handlers -run TestValidateFileFormat_InvalidExtension

# GREEN: Write minimal code to pass the test
# Add the simplest implementation in upload.go:
func (h *UploadHandler) validateFileFormat(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    validExtensions := []string{".pdb", ".cif", ".mmcif"}
    for _, valid := range validExtensions {
        if ext == valid {
            return true
        }
    }
    return false
}

# Run test again - it should PASS (GREEN)
go test -v ./internal/handlers -run TestValidateFileFormat_InvalidExtension

# REFACTOR: Clean up code while keeping tests green
# Commit the behavioral change
git add . && git commit -m "feat: add file format validation"

# Now perform structural improvements (separate commit)
# Extract constant, improve naming, etc.
git add . && git commit -m "refactor: extract file extension constants"
```

#### TDD Testing Patterns for Each Layer

**Handler Tests (HTTP Layer)**
```bash
# Test file upload endpoint
go test -v ./internal/handlers -run TestUploadHandler

# Test analysis endpoints  
go test -v ./internal/handlers -run TestAnalysisHandler

# Run specific test with verbose output
go test -v ./internal/handlers -run TestUploadFile_ValidPDB
```

**Service Tests (Business Logic)**
```bash
# Test storage service
go test -v ./internal/services -run TestStorageService

# Test analysis service
go test -v ./internal/services -run TestAnalysisService

# Test parser service
go test -v ./internal/services -run TestParserService
```

**Model Tests (Data Structures)**
```bash
# Test protein structure models
go test -v ./internal/models -run TestProteinStructure
```

#### Test Categories and Commands

**Unit Tests (Fast, Isolated)**
```bash
# Run all unit tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run tests for specific package
go test -v ./internal/handlers/

# Run single test function
go test -v ./internal/handlers/ -run TestUploadFile
```

**Integration Tests (API Level)**
```bash
# Run integration test script
./test_api.sh

# Test specific API endpoint manually
curl -X POST -F "file=@test.pdb" http://localhost:8080/api/v1/upload
```

#### TDD Guidelines for This Project

1. **Test Structure**: Follow AAA pattern (Arrange, Act, Assert)
```go
func TestAnalysisService_ProcessStructure(t *testing.T) {
    // Arrange
    service := NewAnalysisService(mockParser)
    testFile := "test_protein.pdb"
    
    // Act
    result := service.ProcessStructure("test-id", testFile, defaultOptions)
    
    // Assert
    assert.Equal(t, "completed", result.Status)
    assert.NotNil(t, result.Structure)
}
```

2. **Mock Dependencies**: Use interfaces for testability
```go
type ParserInterface interface {
    ParseStructureFile(filePath string) (*ProteinStructure, error)
}

// Test with mock parser
mockParser := &MockParser{}
service := NewAnalysisService(mockParser)
```

3. **Test Naming Convention**: `TestFunctionName_Scenario_ExpectedResult`
```go
TestUploadHandler_ValidPDBFile_ReturnsFileID
TestAnalysisService_InvalidFile_ReturnsError
TestStorageService_DuplicateFile_OverwritesExisting
```

4. **Red-Green-Refactor Cycle**:
   - **RED**: Write failing test that describes desired behavior
   - **GREEN**: Write minimal code to make test pass
   - **REFACTOR**: Improve code structure while keeping tests green
   - **COMMIT**: Separate commits for behavioral vs structural changes

### Docker Operations
```bash
# Run entire stack (frontend + backend)
docker-compose up

# Backend only with Docker
cd backend
make docker-build
make docker-run
make docker-logs
make docker-stop
```

### API Testing
```bash
cd backend
./test_api.sh  # Automated API endpoint testing
```

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

- `POST /upload` - Upload protein file
- `POST /analyze/{fileId}` - Start analysis (returns analysisId)
- `GET /analysis/{analysisId}` - Get analysis result
- `GET /analysis/{analysisId}/chains` - Get chain data
- `GET /analysis/{analysisId}/domains` - Get domain data
- `GET /analysis/{analysisId}/ligands` - Get ligand data
- `GET /export/{analysisId}` - Export full data as JSON
- `GET /export/{analysisId}/{dataType}` - Export specific data type

## Environment Configuration

Key environment variables in `backend/.env`:
- `USE_MOCK_PARSER`: Toggle between mock data and real BioPython parser
- `PORT`: Server port (default: 8080)
- `UPLOAD_MAX_SIZE`: Maximum file size in bytes (default: 100MB)
- `GIN_MODE`: Set to "debug" for development, "release" for production

## Current State and Next Steps

### What's Implemented
- Complete Go backend with all API endpoints
- Mock data for rapid testing without BioPython
- Docker containerization with multi-stage builds
- Frontend React component (standalone)
- Python parser script ready for integration

### MVP Priorities
1. **Frontend Migration**: Convert `protein-structure-analyzer.tsx` to Next.js app structure
2. **Integration**: Connect frontend with backend API endpoints
3. **Real Parser**: Enable BioPython parser (`USE_MOCK_PARSER=false`) once tested
4. **3D Visualization**: Add Mol* for protein structure visualization

### Production Considerations
- Replace in-memory storage with PostgreSQL or MongoDB
- Implement proper authentication and user sessions
- Add Redis for caching analysis results
- Deploy Python parser as separate microservice for scaling
- Implement S3 or similar for file storage instead of local filesystem

## Testing Strategy

### TDD Testing Hierarchy

**1. Unit Tests (Fastest, Most Isolated)**
- Test individual functions and methods
- Mock all external dependencies
- Focus on business logic correctness
- Run frequently during development

**2. Integration Tests (Medium Speed)**
- Test service interactions
- Use real file system for storage tests
- Test API handlers with mock services
- Verify component integration

**3. End-to-End Tests (Slowest, Most Realistic)**
- Test complete user workflows
- Use real HTTP requests to running server
- Include file upload and analysis pipeline
- Run before releases

### Test-First Development Examples

**Example 1: Adding Chain Type Detection**
```bash
# 1. RED: Write failing test
echo 'func TestDetermineChainType_ProteinChain_ReturnsProtein(t *testing.T) {
    parser := NewProteinParser()
    chainType := parser.determineChainType(mockProteinChain)
    assert.Equal(t, "Protein", chainType)
}' >> internal/services/parser_test.go

# 2. Run test - should FAIL
go test -v ./internal/services -run TestDetermineChainType_ProteinChain

# 3. GREEN: Implement minimal functionality
# Add method to parser.go

# 4. REFACTOR: Clean up after test passes
```

**Example 2: Adding File Size Validation**
```bash
# 1. RED: Write test that expects size validation
echo 'func TestUploadHandler_FileTooLarge_ReturnsError(t *testing.T) {
    handler := NewUploadHandler(storageService)
    largeFile := createMockFile(200 * 1024 * 1024) // 200MB
    
    response := handler.UploadFile(context, largeFile)
    
    assert.Equal(t, 400, response.StatusCode)
    assert.Contains(t, response.Body, "file size exceeds")
}' >> internal/handlers/upload_test.go

# 2. Run test - should FAIL
go test -v ./internal/handlers -run TestUploadHandler_FileTooLarge

# 3. GREEN: Add size validation to upload handler
# 4. REFACTOR: Extract size constants
```

### Testing Commands by Development Phase

**During Feature Development (TDD Cycle)**
```bash
# Run specific test repeatedly during development
go test -v ./internal/handlers -run TestNewFeature

# Watch mode for continuous testing
watch -n 1 "go test -v ./internal/handlers -run TestNewFeature"

# Test with coverage for specific package
go test -v -cover ./internal/handlers
```

**Before Committing (Tidy First Requirements)**
```bash
# All tests must pass before any commit
make test

# Check test coverage
go test -v -cover ./...

# Lint and format (structural tidying)
make fmt && make lint
```

**Manual Testing Flow**
1. Start backend: `cd backend && make run`
2. Upload a PDB file via POST to `/api/v1/upload`
3. Use returned `fileId` to start analysis
4. Poll `/api/v1/analysis/{analysisId}` until status is "completed"
5. Retrieve specific data via chains/domains/ligands endpoints

**Automated Testing**
- Run `./test_api.sh` for complete API workflow testing
- Uses mock data by default for consistent results
- Integration test covering full upload → analysis → export flow

## Common Issues and Solutions

1. **Port 8080 in use**: Change PORT in `.env` file
2. **Python parser fails**: Ensure BioPython installed: `pip3 install biopython`
3. **Docker build fails**: Check Docker daemon is running
4. **Analysis never completes**: Check if mock parser is enabled in `.env`