package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var (
	docStore     map[string][]MCPDocument
	docStoreLock sync.RWMutex
	mcpEnabled   bool
)

func SetMCPEnabled(enabled bool) {
	mcpEnabled = enabled
}

// InitMCP initializes the MCP subsystem
func InitMCP() {
	log.Println("Initializing Model Context Protocol (MCP) support")

	// Initialize the document store if it doesn't exist
	if docStore == nil {
		docStore = make(map[string][]MCPDocument)
	}

	// Initialize the vector store for semantic search
	InitVectorStore()

	// Register all resource handlers
	RegisterHandlers()
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
		// Transform the document before indexing
		doc = TransformDocument(doc)

		sourceKey := fmt.Sprintf("%s:%s", doc.Source, doc.SourceType)
		docStore[sourceKey] = append(docStore[sourceKey], doc)

		// Also add document to vector store for semantic search
		AddDocumentToVectorStore(doc)
	}
}

// GetDocumentCount returns the total number of documents in the MCP index
func GetDocumentCount() int {
	if docStore == nil {
		return 0
	}

	docStoreLock.RLock()
	defer docStoreLock.RUnlock()

	count := 0
	for _, docs := range docStore {
		count += len(docs)
	}
	return count
}

// GetDocumentTypes returns a list of all document types in the index
func GetDocumentTypes() []string {
	if docStore == nil {
		return []string{}
	}

	docStoreLock.RLock()
	defer docStoreLock.RUnlock()

	types := make([]string, 0, len(docStore))
	for sourceType := range docStore {
		types = append(types, sourceType)
	}

	sort.Strings(types)
	return types
}

// HandleMCPRetrieve handles API requests to retrieve documents from MCP
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

	if query.Query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	limit := 10
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			limit = 10
		}
	}

	// Use vector search if available, otherwise fall back to document retrieval
	var docs []MCPDocument
	var err error

	if vectorStore != nil {
		docs, err = vectorStore.Search(query.Query, limit)
		if err != nil {
			log.Printf("Vector search error: %v, falling back to regular retrieval", err)
			docs, err = retrieveDocuments(query, limit)
		}
	} else {
		docs, err = retrieveDocuments(query, limit)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving documents: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

// retrieveDocuments performs a simple keyword-based search in the document store
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
		firstIndex := -1
		if len(queryTerms) > 0 {
			firstIndex = strings.Index(contentLower, queryTerms[0])
		}
		if firstIndex >= 0 && firstIndex < 50 {
			results[i].Score += 0.1
		}
	}

	// Sort by score (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Truncate to limit
	if len(results) > limit {
		results = results[:limit]
	}

	log.Printf("Returning %d documents after filtering and sorting", len(results))
	return results, nil
}

// matchesMetadata checks if a document's metadata matches the query criteria
func matchesMetadata(doc MCPDocument, metadataQuery map[string]interface{}) bool {
	if metadataQuery == nil || len(metadataQuery) == 0 {
		return true
	}

	for key, queryValue := range metadataQuery {
		docValue, exists := doc.Metadata[key]
		if !exists {
			return false
		}

		// Convert values to strings for comparison
		queryStr := fmt.Sprintf("%v", queryValue)
		docStr := fmt.Sprintf("%v", docValue)

		// Check if values match (case insensitive)
		if !strings.EqualFold(queryStr, docStr) {
			return false
		}
	}

	return true
}

// HandleAIPluginManifest provides the AI Plugin manifest
func HandleAIPluginManifest(w http.ResponseWriter, r *http.Request) {
	baseURL := getBaseURL(r)
	manifest := map[string]interface{}{
		"schema_version":        "v1",
		"name_for_human":        "Kollect Data Explorer",
		"name_for_model":        "KollectDataExplorer",
		"description_for_human": "Access and search data collected by Kollect from various platforms.",
		"description_for_model": "This plugin provides access to data collected from various cloud platforms and services. You can search resources from Kubernetes, AWS, Azure, GCP, Terraform, Vault, and other sources.",
		"auth": map[string]string{
			"type": "none",
		},
		"api": map[string]string{
			"type": "openapi",
			"url":  baseURL + "/openapi.yaml",
		},
		"logo_url":       baseURL + "/logo.png",
		"contact_email":  "support@example.com",
		"legal_info_url": baseURL + "/legal",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manifest)
}

// HandleOpenAPISpec provides the OpenAPI specification
func HandleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	baseURL := getBaseURL(r)

	spec := fmt.Sprintf(`openapi: 3.0.1
info:
  title: Kollect Data Explorer API
  description: Search and explore data collected from various platforms.
  version: "v1"
servers:
  - url: %s
paths:
  /api/mcp/retrieve:
    post:
      operationId: searchResources
      summary: Search collected resources
      description: Search for resources across all collected data.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - query
              properties:
                query:
                  type: string
                  description: Search query
                metadata:
                  type: object
                  description: Additional metadata filters
      responses:
        "200":
          description: Matching resources
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Resource'
components:
  schemas:
    Resource:
      type: object
      properties:
        id:
          type: string
        content:
          type: string
        metadata:
          type: object
        source:
          type: string
        sourceType:
          type: string
        score:
          type: number
`, baseURL)

	w.Header().Set("Content-Type", "text/yaml")
	w.Write([]byte(spec))
}

// Helper function to get the base URL
func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}
