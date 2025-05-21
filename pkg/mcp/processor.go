package mcp

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// ProcessData converts collected data to MCP documents
func ProcessData(data interface{}, sourceType string) []MCPDocument {
	// Add a check for if MCP is enabled
	if !mcpEnabled {
		// Return empty results if MCP is not enabled
		return []MCPDocument{}
	}

	log.Printf("Processing %s data for MCP, data type: %T", sourceType, data)

	// Ensure MCP is initialized
	if docStore == nil {
		log.Println("MCP not initialized, initializing now...")
		InitMCP()
	}

	// Process data based on type
	var documents []MCPDocument

	// Get all handlers for this platform
	handlers := GetAllHandlersForPlatform(sourceType)

	if len(handlers) > 0 {
		// Use registered handlers to process data
		for _, handler := range handlers {
			docs := handler.ExtractFunc(data)
			if len(docs) > 0 {
				documents = append(documents, docs...)
			}
		}
	} else {
		// Fallback: use generic processing
		log.Printf("No specific handlers found for %s, using generic processing", sourceType)
		documents = processGenericData(data, sourceType)
	}

	log.Printf("MCP: Processed %s data, found %d documents", sourceType, len(documents))

	// Index the documents
	if len(documents) > 0 {
		IndexDocuments(documents)
		log.Printf("MCP: Successfully indexed %d documents", len(documents))
	} else {
		log.Printf("MCP: No documents were generated for source type '%s'", sourceType)
	}

	return documents
}

// Process data of unknown structure
func processGenericData(data interface{}, sourceType string) []MCPDocument {
	var documents []MCPDocument

	val := reflect.ValueOf(data)

	// Handle different types of data
	switch val.Kind() {
	case reflect.Map:
		// Process each key in the map
		iter := val.MapRange()
		for iter.Next() {
			key := fmt.Sprintf("%v", iter.Key().Interface())
			value := iter.Value().Interface()

			// Try to identify resource type from the key
			resourceType := inferResourceTypeFromKey(key)

			// Use generic extractor
			docs := GenericDocumentExtractor(sourceType, resourceType, value)
			documents = append(documents, docs...)
		}

	case reflect.Slice, reflect.Array:
		// Process slice of items
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()

			// Infer resource type or use generic "item"
			resourceType := "item"
			itemVal := reflect.ValueOf(item)
			if itemVal.Kind() == reflect.Map {
				// Try to find type information in the map
				for _, key := range []string{"type", "Type", "kind", "Kind"} {
					typeField := itemVal.MapIndex(reflect.ValueOf(key))
					if typeField.IsValid() {
						resourceType = fmt.Sprintf("%v", typeField.Interface())
						break
					}
				}
			}

			docs := GenericDocumentExtractor(sourceType, resourceType, item)
			documents = append(documents, docs...)
		}

	default:
		// Handle single value
		doc := MCPDocument{
			ID:      fmt.Sprintf("%s-data", sourceType),
			Content: fmt.Sprintf("%s Data:\n%v", sourceType, data),
			Metadata: map[string]interface{}{
				"type": "generic",
			},
			Source:     sourceType,
			SourceType: "unknown",
			CreatedAt:  time.Now(),
		}
		documents = append(documents, doc)
	}

	return documents
}

// Helper to infer resource type from key names
func inferResourceTypeFromKey(key string) string {
	key = strings.ToLower(key)
	// Common naming patterns in Cloud APIs
	switch {
	case strings.Contains(key, "instance"):
		return "instance"
	case strings.Contains(key, "bucket"):
		return "bucket"
	case strings.Contains(key, "pod"):
		return "pod"
	case strings.Contains(key, "node"):
		return "node"
	case strings.Contains(key, "vm"):
		return "vm"
	case strings.Contains(key, "resource"):
		return "resource"
	default:
		return key
	}
}
