package mcp

import (
	"time"
)

// MCPDocument represents a document that can be retrieved via MCP
type MCPDocument struct {
	ID         string                 `json:"id"`
	Content    string                 `json:"content"`
	Metadata   map[string]interface{} `json:"metadata"`
	Score      float64                `json:"score,omitempty"`
	Source     string                 `json:"source"`
	SourceType string                 `json:"source_type"`
	CreatedAt  time.Time              `json:"created_at"`
	VectorID   string                 `json:"vector_id,omitempty"`
}

// MCPQuery represents a query for retrieving documents
type MCPQuery struct {
	Query        string                 `json:"query"`
	Limit        int                    `json:"limit,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	VectorSearch bool                   `json:"vector_search,omitempty"`
}

// MCPResponse is the response to an MCP query
type MCPResponse struct {
	Documents []MCPDocument          `json:"documents"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
