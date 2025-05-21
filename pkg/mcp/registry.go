package mcp

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// ResourceTypeHandler defines how to process a specific resource type
type ResourceTypeHandler struct {
	// Unique identifier for this resource type
	ID string

	// Source platform (aws, kubernetes, etc.)
	Platform string

	// Resource type (e.g., ec2, pod, etc.)
	ResourceType string

	// Extractor function that converts data to documents
	ExtractFunc func(data interface{}) []MCPDocument
}

var (
	// Global registry of resource handlers
	resourceHandlers = make(map[string]ResourceTypeHandler)
)

// RegisterResourceHandler registers a handler for a specific resource type
func RegisterResourceHandler(handler ResourceTypeHandler) {
	id := fmt.Sprintf("%s:%s", handler.Platform, handler.ResourceType)
	resourceHandlers[id] = handler
	log.Printf("Registered MCP handler for %s", id)
}

// GetResourceHandler retrieves a handler for a given platform and type
func GetResourceHandler(platform, resourceType string) (ResourceTypeHandler, bool) {
	id := fmt.Sprintf("%s:%s", platform, resourceType)
	handler, exists := resourceHandlers[id]
	return handler, exists
}

// GetAllHandlersForPlatform retrieves all handlers for a given platform
func GetAllHandlersForPlatform(platform string) []ResourceTypeHandler {
	var handlers []ResourceTypeHandler

	for id, handler := range resourceHandlers {
		if strings.HasPrefix(id, platform+":") {
			handlers = append(handlers, handler)
		}
	}

	return handlers
}

// GenericDocumentExtractor provides a fallback extractor for unknown types
func GenericDocumentExtractor(platform string, resourceType string, data interface{}) []MCPDocument {
	var docs []MCPDocument

	switch value := data.(type) {
	case map[string]interface{}:
		// Process map data
		for key, item := range value {
			doc := MCPDocument{
				ID:      fmt.Sprintf("%s-%s-%s", platform, resourceType, key),
				Content: formatGenericContent(fmt.Sprintf("%s %s", platform, resourceType), item),
				Metadata: map[string]interface{}{
					"type": resourceType,
				},
				Source:     platform,
				SourceType: resourceType,
				CreatedAt:  time.Now(),
			}
			docs = append(docs, doc)
		}

	case []interface{}:
		// Process slice data
		for i, item := range value {
			doc := MCPDocument{
				ID:      fmt.Sprintf("%s-%s-%d", platform, resourceType, i),
				Content: formatGenericContent(fmt.Sprintf("%s %s", platform, resourceType), item),
				Metadata: map[string]interface{}{
					"type":  resourceType,
					"index": i,
				},
				Source:     platform,
				SourceType: resourceType,
				CreatedAt:  time.Now(),
			}
			docs = append(docs, doc)
		}

	default:
		// Handle any other type
		doc := MCPDocument{
			ID:      fmt.Sprintf("%s-%s-data", platform, resourceType),
			Content: fmt.Sprintf("%s %s Data:\n%v", platform, resourceType, data),
			Metadata: map[string]interface{}{
				"type": resourceType,
			},
			Source:     platform,
			SourceType: resourceType,
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	return docs
}
