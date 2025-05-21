package mcp

import (
	"strings"
)

// TransformDocument prepares a document for indexing
func TransformDocument(doc MCPDocument) MCPDocument {
	// Normalize content
	doc.Content = strings.TrimSpace(doc.Content)

	// Add additional processing as needed

	return doc
}
