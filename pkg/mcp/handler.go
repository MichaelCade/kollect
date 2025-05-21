package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	// Shared data store for MCP documents
	docStore     map[string][]MCPDocument
	docStoreLock sync.RWMutex
)

func InitMCP() {
	log.Println("Initializing Model Context Protocol (MCP) support")
	docStore = make(map[string][]MCPDocument)

	// Initialize the vector store for semantic search
	InitVectorStore()
}

func GetDocumentCount() int {
	docStoreLock.RLock()
	defer docStoreLock.RUnlock()

	count := 0
	for _, docs := range docStore {
		count += len(docs)
	}
	return count
}

// GetDocumentTypes returns the document types in the store
func GetDocumentTypes() []string {
	docStoreLock.RLock()
	defer docStoreLock.RUnlock()

	types := make([]string, 0, len(docStore))
	for key := range docStore {
		types = append(types, key)
	}
	return types
}

// IndexDocuments adds documents to the search index
func IndexDocuments(documents []MCPDocument) {
	// Ensure docStore is initialized
	if docStore == nil {
		log.Println("Warning: docStore was nil, initializing it now")
		docStore = make(map[string][]MCPDocument)
	}

	docStoreLock.Lock()
	defer docStoreLock.Unlock()

	// Group documents by source
	for _, doc := range documents {
		sourceKey := fmt.Sprintf("%s:%s", doc.Source, doc.SourceType)
		docStore[sourceKey] = append(docStore[sourceKey], doc)

		// Also add document to vector store for semantic search
		AddDocumentToVectorStore(doc)
	}
}

func AddDocumentToVectorStore(doc MCPDocument) {
	if vectorStore == nil {
		log.Println("Warning: Vector store not initialized, initializing now")
		InitVectorStore()
	}

	err := vectorStore.IndexDocument(doc)
	if err != nil {
		log.Printf("Error adding document to vector store: %v", err)
	}
}

// HandleMCPRetrieve handles MCP retrieval requests
func HandleMCPRetrieve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var query MCPQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("MCP query received: %s (limit: %d)", query.Query, query.Limit)

	limit := 5
	if query.Limit > 0 {
		limit = query.Limit
	}

	// Retrieve relevant documents
	documents, err := retrieveDocuments(query, limit)
	if err != nil {
		http.Error(w, "Error retrieving documents", http.StatusInternalServerError)
		return
	}

	response := MCPResponse{
		Documents: documents,
		Metadata: map[string]interface{}{
			"total_docs": len(documents),
			"query":      query.Query,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleAIPluginManifest handles the AI plugin manifest endpoint
func HandleAIPluginManifest(w http.ResponseWriter, r *http.Request) {
	manifest := map[string]interface{}{
		"schema_version":        "v1",
		"name_for_human":        "Kollect Infrastructure Explorer",
		"name_for_model":        "kollect",
		"description_for_human": "Plugin for exploring cloud infrastructure and resources across AWS, Azure, GCP, Kubernetes, Terraform, Vault, and Veeam.",
		"description_for_model": "Plugin for accessing information about cloud infrastructure and resources with Kollect. Use this plugin when the user asks about their AWS, Azure, GCP, Kubernetes, Terraform, Vault, or Veeam resources.",
		"auth": map[string]interface{}{
			"type": "none",
		},
		"api": map[string]interface{}{
			"type": "openapi",
			"url":  "/openapi.yaml",
		},
		"logo_url":       "/logo.png",
		"contact_email":  "contact@example.com",
		"legal_info_url": "https://example.com/legal",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manifest)
}

// HandleOpenAPISpec handles the OpenAPI spec endpoint
func HandleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	// OpenAPI spec for your MCP endpoint
	w.Header().Set("Content-Type", "application/yaml")
	w.Write([]byte(`
openapi: 3.0.1
info:
  title: Kollect Infrastructure Explorer
  description: Plugin for exploring cloud infrastructure and resources
  version: 'v1'
servers:
  - url: http://localhost:8080
paths:
  /api/mcp/retrieve:
    post:
      operationId: retrieveDocuments
      summary: Retrieve infrastructure data based on query
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: The search query
                limit:
                  type: integer
                  description: Maximum number of results to return
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  documents:
                    type: array
                    items:
                      type: object
                      properties:
                        id:
                          type: string
                        content:
                          type: string
                        metadata:
                          type: object
                        score:
                          type: number
`))
}

// Update the retrieveDocuments function to make it more useful:
func retrieveDocuments(query MCPQuery, limit int) ([]MCPDocument, error) {
	docStoreLock.RLock()
	defer docStoreLock.RUnlock()

	var results []MCPDocument
	queryLower := strings.ToLower(query.Query)
	queryTerms := strings.Fields(queryLower) // Split query into terms

	log.Printf("Searching for query: '%s' with %d terms across %d document types",
		queryLower, len(queryTerms), len(docStore))

	// More sophisticated matching
	for sourceType, docList := range docStore {
		log.Printf("Searching in %s (%d documents)", sourceType, len(docList))

		// Check if the query directly mentions this resource type
		resourceTypeMatches := false
		resourceType := strings.Split(sourceType, ":")
		if len(resourceType) > 1 {
			singularType := strings.TrimSuffix(resourceType[1], "s") // Convert plural to singular (pods -> pod)
			pluralType := resourceType[1]
			if !strings.HasSuffix(pluralType, "s") {
				pluralType = resourceType[1] + "s" // Convert singular to plural (pod -> pods)
			}

			// Check if query contains the resource type (singular or plural)
			for _, term := range queryTerms {
				if strings.EqualFold(term, singularType) || strings.EqualFold(term, pluralType) {
					resourceTypeMatches = true
					break
				}
			}
		}

		// If query mentions this resource type, include all documents of this type
		if resourceTypeMatches {
			log.Printf("Query directly mentions resource type %s, including all matching documents", sourceType)
			for _, doc := range docList {
				// Apply metadata filtering
				if matchesMetadata(doc, query.Metadata) {
					// Calculate a higher score for type matches
					doc.Score = 0.9 // High priority for type matches
					results = append(results, doc)
					if len(results) >= limit*3 {
						break // Get more than we need for better sorting
					}
				}
			}
		} else {
			// Otherwise do content matching
			for _, doc := range docList {
				contentLower := strings.ToLower(doc.Content)

				// Check if any query term matches
				matchCount := 0
				for _, term := range queryTerms {
					if strings.Contains(contentLower, term) {
						matchCount++
					}
				}

				// Only include if at least one term matches
				if matchCount > 0 {
					// Apply metadata filtering
					if matchesMetadata(doc, query.Metadata) {
						// Calculate how much of the query matches
						matchRatio := float64(matchCount) / float64(len(queryTerms))
						doc.Score = matchRatio * 0.8 // Lower priority for content matches
						results = append(results, doc)
					}
				}
			}
		}
	}

	log.Printf("Found %d matching documents before sorting", len(results))

	// Improve relevance scoring with term frequency and position
	for i := range results {
		contentLower := strings.ToLower(results[i].Content)

		// Add bonus for exact phrase match
		if strings.Contains(contentLower, queryLower) {
			results[i].Score += 0.3
		}

		// Add bonus for matches near the beginning of content
		firstIndex := strings.Index(contentLower, queryTerms[0])
		if firstIndex >= 0 && firstIndex < 50 {
			results[i].Score += 0.1
		}
	}

	// Sort by score (highest first)
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Score < results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	// Truncate to limit
	if len(results) > limit {
		results = results[:limit]
	}

	log.Printf("Returning %d documents after filtering and sorting", len(results))
	return results, nil
}

// Helper function to check if document matches metadata filters
func matchesMetadata(doc MCPDocument, filters map[string]interface{}) bool {
	if filters == nil || len(filters) == 0 {
		return true
	}

	for key, value := range filters {
		docValue, exists := doc.Metadata[key]
		if !exists || docValue != value {
			return false
		}
	}

	return true
}

// Very basic relevance sorting - can be improved
func sortByRelevance(docs []MCPDocument, query string) {
	queryLower := strings.ToLower(query)

	// Count occurrences of the query in each document
	for i := range docs {
		count := strings.Count(strings.ToLower(docs[i].Content), queryLower)
		docs[i].Score = float64(count) / float64(len(docs[i].Content)) * 100
	}

	// Simple bubble sort by score
	for i := 0; i < len(docs)-1; i++ {
		for j := 0; j < len(docs)-i-1; j++ {
			if docs[j].Score < docs[j+1].Score {
				docs[j], docs[j+1] = docs[j+1], docs[j]
			}
		}
	}
}
