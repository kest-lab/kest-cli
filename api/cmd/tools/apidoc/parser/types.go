package parser

// Route represents a parsed route definition
type Route struct {
	Method      string   // HTTP method: GET, POST, PUT, DELETE
	Path        string   // URL path: /api/v1/users
	Name        string   // Route name: v1.users.index
	Handler     string   // Handler function: userHandler.List
	Module      string   // Module name: user
	Middlewares []string // Applied middlewares: auth, rate-limit
	IsPublic    bool     // Whether route is public (no auth)
}

// DTO represents a parsed Data Transfer Object
type DTO struct {
	Name    string     // DTO name: UserRegisterRequest
	Module  string     // Module name: user
	Type    string     // request, response, or model
	Fields  []DTOField // Fields in the DTO
	Comment string     // Doc comment
}

// DTOField represents a field in a DTO
type DTOField struct {
	Name       string // Field name: Username
	Type       string // Field type: string
	JSONName   string // JSON tag name: username
	Binding    string // Binding validation: required,min=3,max=50
	Required   bool   // Whether field is required
	Comment    string // Field comment
	Validation string // Validation rules description
}

// Endpoint represents a complete API endpoint with all information
type Endpoint struct {
	Route       Route  // Route information
	Request     *DTO   // Request DTO (if any)
	Response    *DTO   // Response DTO (if any)
	Summary     string // Brief description
	Description string // Detailed description
	Tags        string // API tags/category
}

// ModuleDoc represents documentation for a module
type ModuleDoc struct {
	Name      string     // Module name
	Endpoints []Endpoint // Endpoints in this module
}
