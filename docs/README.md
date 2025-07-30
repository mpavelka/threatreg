# Threatreg Documentation

This directory contains tools and generated documentation for the Threatreg service, including API documentation and service layer documentation.

## Files

- `docs.go` - Generated OpenAPI/Swagger documentation (auto-generated)
- `swagger.json` - OpenAPI specification in JSON format (auto-generated)
- `swagger.yaml` - OpenAPI specification in YAML format (auto-generated)
- `generate_service_docs.py` - Python script that scans Go service files and generates HTML documentation
- `service_documentation.html` - Generated HTML documentation with expandable sections (auto-generated)

## Usage

### Generating OpenAPI/Swagger Documentation

The REST API documentation is generated using Swagger annotations in the handler code. To regenerate the API documentation after making changes to API endpoints or annotations:

```bash
# Install the swag tool (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation from the project root
swag init --dir ./ --generalInfo docs.go --output ./docs
```

The generation process will:
1. Scan all Go files for Swagger annotations (starting with `// @`)
2. Parse the main API documentation from `docs.go`
3. Extract endpoint documentation from handler functions
4. Generate comprehensive OpenAPI specification files

#### Viewing API Documentation

The API documentation is served automatically when running the REST API server:

```bash
# Start the REST API server
./bin/threatreg restapi

# Access interactive Swagger UI at any of these endpoints:
# http://localhost:8080/swagger/index.html
# http://localhost:8080/docs/index.html  
# http://localhost:8080/api-docs/index.html
```

#### Adding API Documentation

When adding new endpoints or modifying existing ones, add Swagger annotations above your handler functions:

```go
// CreateProduct handles POST /api/v1/products
// @Summary Create a new product
// @Description Create a new product with the provided name and description
// @Tags Products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product creation request"
// @Success 201 {object} handlers.SuccessResponse{data=models.Product}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /products [post]
func CreateProduct(c *gin.Context) {
    // implementation
}
```

### Generating Service Documentation

To regenerate the service documentation after making changes to service functions or their docstrings:

```bash
# From the project root directory
python3 docs/generate_service_docs.py
```

The script will:
1. Scan all Go files in `internal/service/` (excluding test files)
2. Extract public functions (those starting with capital letters) and their documentation
3. Organize functions by service category
4. Generate a comprehensive HTML document with expandable sections

### Viewing Documentation

Open the generated HTML file in any web browser:

```bash
# Open in default browser (macOS)
open docs/service_documentation.html

# Or navigate directly
file:///path/to/threatreg/docs/service_documentation.html
```

## Documentation Features

- **Organized by Category**: Functions are grouped by service type (Control Management, Domain Management, etc.)
- **Expandable Sections**: Click category headers to expand/collapse sections
- **Function Details**: Each function shows its signature and documentation
- **Search-Friendly**: Standard HTML that works with browser search (Ctrl+F / Cmd+F)
- **Responsive Design**: Works on desktop and mobile devices
- **Statistics**: Shows total function count and generation date

## Writing Good Docstrings

When adding or updating service functions, follow these docstring conventions:

```go
// FunctionName performs a specific operation with clear description.
// Explains what the function does, parameters, return values, and any important behavior.
// Multiple lines are supported for detailed explanations.
func FunctionName(param1 string, param2 uuid.UUID) (*Model, error) {
    // implementation
}
```

### Docstring Guidelines

1. **Start with function name** and a clear, concise description
2. **Explain parameters** if they're not self-evident
3. **Describe return values**, especially error conditions
4. **Mention side effects** or important behavior
5. **Use complete sentences** with proper punctuation
6. **Keep it concise** but informative

## Categories

The documentation organizes functions into these categories:

- **Control Management** - Security controls and countermeasures
- **Domain Management** - Domain organization and instance grouping
- **Instance Management** - Application instances and deployments
- **Product Management** - Software products and systems
- **Threat Management** - Security threats and vulnerabilities
- **Relationship Management** - Entity relationships and connections
- **Tag Management** - Tagging and categorization system
- **Threat Resolution Management** - Threat mitigation and resolution tracking
- **Threat Pattern Management** - Threat pattern definitions and rules
- **Threat Pattern Conditions** - Pattern condition logic and validation
- **Instance Threat Pattern Evaluation** - Pattern matching and evaluation

## Maintenance

The documentation generator should be run whenever:
- New public functions are added to service files
- Existing function signatures change
- Documentation comments are added or updated
- Service files are reorganized

Consider adding this to your development workflow or CI/CD pipeline to keep documentation current.