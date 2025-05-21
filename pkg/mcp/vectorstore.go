package mcp

import (
	"log"
)

// Embedder interface for text embedding
type Embedder interface {
	Embed(text string) ([]float32, error)
}

// SimpleEmbedder is a placeholder implementation
type SimpleEmbedder struct{}

func (e *SimpleEmbedder) Embed(text string) ([]float32, error) {
	// This is a placeholder - just returns a simple vector based on text length
	return []float32{float32(len(text))}, nil
}

// NewEmbedder creates a new embedder
func NewEmbedder() Embedder {
	return &SimpleEmbedder{}
}

// VectorStore manages vector embeddings of documents
type VectorStore struct {
	embedder Embedder
	vectors  map[string][]float32 // Document ID -> vector
}

// Global instance
var vectorStore *VectorStore

// InitVectorStore initializes the vector store
func InitVectorStore() {
	vectorStore = &VectorStore{
		embedder: NewEmbedder(),
		vectors:  make(map[string][]float32),
	}
	log.Println("Vector store initialized")
}

// IndexDocument adds a document to the vector store
func (vs *VectorStore) IndexDocument(doc MCPDocument) error {
	vector, err := vs.embedder.Embed(doc.Content)
	if err != nil {
		return err
	}

	// Store the vector
	vs.storeVector(doc.ID, vector)
	return nil
}

// storeVector saves a vector to the store
func (vs *VectorStore) storeVector(id string, vector []float32) {
	vs.vectors[id] = vector
}

// SearchVectors finds similar vectors
func (vs *VectorStore) SearchVectors(query string, limit int) ([]string, error) {
	// Get the query vector
	queryVector, err := vs.embedder.Embed(query)
	if err != nil {
		return nil, err
	}

	// Find similar documents
	return vs.searchVectors(queryVector, limit), nil
}

// searchVectors finds similar vectors using cosine similarity
func (vs *VectorStore) searchVectors(queryVector []float32, limit int) []string {
	// In a real implementation, this would do similarity search
	// For now, just return the first 'limit' document IDs
	var ids []string
	for id := range vs.vectors {
		ids = append(ids, id)
		if len(ids) >= limit {
			break
		}
	}
	return ids
}

// GetDocumentByID retrieves a document by ID
func getDocumentByID(id string) (MCPDocument, bool) {
	docStoreLock.RLock()
	defer docStoreLock.RUnlock()

	// Search through all document collections
	for _, docList := range docStore {
		for _, doc := range docList {
			if doc.ID == id {
				return doc, true
			}
		}
	}

	return MCPDocument{}, false
}
