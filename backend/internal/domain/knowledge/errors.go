// Package knowledge defines knowledge domain errors
package knowledge

import "errors"

var (
	// ErrDocumentNotFound indicates document does not exist
	ErrDocumentNotFound = errors.New("document not found")
	// ErrDocumentParseFailed indicates document parsing failed
	ErrDocumentParseFailed = errors.New("document parse failed")
	// ErrUnsupportedDocumentType indicates document type is not supported
	ErrUnsupportedDocumentType = errors.New("unsupported document type")
	// ErrEmptyChunks indicates document has no chunks
	ErrEmptyChunks = errors.New("document chunks is empty")
	// ErrEmbeddingFailed indicates embedding generation failed
	ErrEmbeddingFailed = errors.New("embedding failed")
	// ErrVectorSearchFailed indicates vector search failed
	ErrVectorSearchFailed = errors.New("vector search failed")
	// ErrKnowledgeBaseEmpty indicates knowledge base is empty
	ErrKnowledgeBaseEmpty = errors.New("knowledge base is empty, please upload documents first")
	// ErrDocumentProcessing indicates document is being processed
	ErrDocumentProcessing = errors.New("document is being processed")
)
