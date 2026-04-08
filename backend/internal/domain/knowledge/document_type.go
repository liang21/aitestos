// Package knowledge defines DocumentType value object
package knowledge

import "errors"

// DocumentType is a value object representing the type of document
type DocumentType string

const (
	// TypePRD represents Product Requirements Document
	TypePRD DocumentType = "prd"
	// TypeFigma represents Figma design file
	TypeFigma DocumentType = "figma"
	// TypeAPISpec represents API specification
	TypeAPISpec DocumentType = "api_spec"
	// TypeSwagger represents Swagger/OpenAPI spec
	TypeSwagger DocumentType = "swagger"
	// TypeMarkdown represents Markdown document
	TypeMarkdown DocumentType = "markdown"
)

// ParseDocumentType validates and creates a DocumentType
func ParseDocumentType(s string) (DocumentType, error) {
	switch DocumentType(s) {
	case TypePRD, TypeFigma, TypeAPISpec, TypeSwagger, TypeMarkdown:
		return DocumentType(s), nil
	default:
		return "", errors.New("invalid document type")
	}
}

// String returns the string representation
func (t DocumentType) String() string {
	return string(t)
}

// IsSupported returns true if the document type is supported
func (t DocumentType) IsSupported() bool {
	switch t {
	case TypePRD, TypeFigma, TypeAPISpec, TypeSwagger, TypeMarkdown:
		return true
	default:
		return false
	}
}
