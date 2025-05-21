package mcp

import (
	"log"
	"strings"
)

var (
	vectorStore *VectorStore
)

// VectorStore provides vector-based search capabilities
type VectorStore struct {
	// This is a simple in-memory implementation
	// In a real application, you might use a dedicated vector database
	documents []MCPDocument
}

// InitVectorStore initializes the vector store
func InitVectorStore() {
	if vectorStore == nil {
		log.Println("Initializing MCP vector store")
		vectorStore = &VectorStore{
			documents: make([]MCPDocument, 0),
		}
	}
}

// AddDocumentToVectorStore adds a document to the vector store
func AddDocumentToVectorStore(doc MCPDocument) {
	if vectorStore == nil {
		log.Println("Warning: Vector store not initialized, initializing now")
		InitVectorStore()
	}

	// Add the document to the vector store
	vectorStore.IndexDocument(doc)
}

// IndexDocument adds a document to the vector store
func (vs *VectorStore) IndexDocument(doc MCPDocument) error {
	// In a real vector store, you would compute embeddings here
	// For this simple implementation, we just store the document
	vs.documents = append(vs.documents, doc)
	return nil
}

// Search performs a vector search
func (vs *VectorStore) Search(query string, limit int) ([]MCPDocument, error) {
	if len(vs.documents) == 0 {
		return []MCPDocument{}, nil
	}

	// This is a very basic search implementation
	// In a real vector store, you would:
	// 1. Convert the query to an embedding
	// 2. Find the nearest neighbors in the embedding space

	var results []MCPDocument
	queryLower := strings.ToLower(query)

	for _, doc := range vs.documents {
		contentLower := strings.ToLower(doc.Content)
		if strings.Contains(contentLower, queryLower) {
			// Set a simple relevance score based on number of occurrences
			occurrences := strings.Count(contentLower, queryLower)
			doc.Score = float64(occurrences) / float64(len(doc.Content))

			results = append(results, doc)

			if len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

// Clear removes all documents from the vector store
func (vs *VectorStore) Clear() {
	vs.documents = make([]MCPDocument, 0)
}

// GetDocumentCount returns the number of documents in the vector store
func (vs *VectorStore) GetDocumentCount() int {
	return len(vs.documents)
}

// GetDocument retrieves a document by ID
func (vs *VectorStore) GetDocument(id string) (MCPDocument, bool) {
	for _, doc := range vs.documents {
		if doc.ID == id {
			return doc, true
		}
	}
	return MCPDocument{}, false
}
