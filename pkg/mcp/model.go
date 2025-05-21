package mcp

import "time"

// MCPDocument represents a document in the MCP system
type MCPDocument struct {
	// Unique identifier for the document
	ID string `json:"id"`

	// Content of the document
	Content string `json:"content"`

	// Metadata for the document (type, attributes, etc.)
	Metadata map[string]interface{} `json:"metadata"`

	// Source of the document (aws, azure, gcp, k8s, etc.)
	Source string `json:"source"`

	// Source type more specifically (ec2, s3, pod, etc.)
	SourceType string `json:"source_type"`

	// When the document was created
	CreatedAt time.Time `json:"created_at"`

	// Relevance score for search results (not stored)
	Score float64 `json:"score,omitempty"`
}

// MCPQuery represents a query to the MCP system
type MCPQuery struct {
	// The query string
	Query string `json:"query"`

	// Optional metadata filters
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Maximum number of results to return
	Limit int `json:"limit,omitempty"`

	// Whether to use vector search
	VectorSearch bool `json:"vector_search,omitempty"`
}

// MCPResponse is the response to an MCP query
type MCPResponse struct {
	Documents []MCPDocument          `json:"documents"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
