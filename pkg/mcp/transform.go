package mcp

import (
	"strings"
)

// TransformDocument prepares a document for indexing
func TransformDocument(doc MCPDocument) MCPDocument {
	// Normalize content
	doc.Content = strings.TrimSpace(doc.Content)

	// Ensure metadata exists
	if doc.Metadata == nil {
		doc.Metadata = make(map[string]interface{})
	}

	// Add source type to metadata if not present
	if _, exists := doc.Metadata["type"]; !exists {
		doc.Metadata["type"] = doc.SourceType
	}

	// Extract key entities and phrases for better searchability
	extractedTerms := extractKeyTerms(doc.Content)
	if len(extractedTerms) > 0 {
		doc.Metadata["key_terms"] = extractedTerms
	}

	// Normalize source field to lowercase
	doc.Source = strings.ToLower(doc.Source)

	// Normalize source type to lowercase and remove spaces
	doc.SourceType = strings.ToLower(strings.ReplaceAll(doc.SourceType, " ", "_"))

	return doc
}

// extractKeyTerms extracts important terms from content for improved search
func extractKeyTerms(content string) []string {
	// Simple extraction of lines with key-value pairs
	lines := strings.Split(content, "\n")
	var terms []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip title lines (usually end with a colon)
		if strings.HasSuffix(line, ":") && !strings.Contains(line[:len(line)-1], ":") {
			continue
		}

		// Extract key-value pairs
		if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Only add non-empty values that aren't just numbers
			if value != "" && !isNumericOnly(value) {
				terms = append(terms, value)
			}

			// Also add keys for certain important fields
			switch strings.ToLower(key) {
			case "name", "id", "instance", "region", "zone", "type", "version":
				terms = append(terms, key+":"+value)
			}
		}
	}

	return terms
}

// isNumericOnly checks if a string contains only numbers, dots, and common separators
func isNumericOnly(s string) bool {
	// Remove common numeric separators
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")

	// Check if remaining characters are all digits
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}

	return s != ""
}
