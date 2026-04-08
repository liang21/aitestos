// Package knowledge_test tests DocumentType value object
package knowledge_test

import (
	"testing"

	"github.com/liang21/aitestos/internal/domain/knowledge"
)

func TestParseDocumentType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    knowledge.DocumentType
		wantErr bool
	}{
		{
			name:  "PRD type",
			input: "prd",
			want:  knowledge.TypePRD,
		},
		{
			name:  "figma type",
			input: "figma",
			want:  knowledge.TypeFigma,
		},
		{
			name:  "api spec type",
			input: "api_spec",
			want:  knowledge.TypeAPISpec,
		},
		{
			name:  "swagger type",
			input: "swagger",
			want:  knowledge.TypeSwagger,
		},
		{
			name:  "markdown type",
			input: "markdown",
			want:  knowledge.TypeMarkdown,
		},
		{
			name:    "invalid type",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty type",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := knowledge.ParseDocumentType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDocumentType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDocumentType(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDocumentType_String(t *testing.T) {
	tests := []struct {
		name    string
		docType knowledge.DocumentType
		want    string
	}{
		{
			name:    "PRD string",
			docType: knowledge.TypePRD,
			want:    "prd",
		},
		{
			name:    "Figma string",
			docType: knowledge.TypeFigma,
			want:    "figma",
		},
		{
			name:    "API Spec string",
			docType: knowledge.TypeAPISpec,
			want:    "api_spec",
		},
		{
			name:    "Swagger string",
			docType: knowledge.TypeSwagger,
			want:    "swagger",
		},
		{
			name:    "Markdown string",
			docType: knowledge.TypeMarkdown,
			want:    "markdown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.docType.String(); got != tt.want {
				t.Errorf("DocumentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentType_IsSupported(t *testing.T) {
	tests := []struct {
		name    string
		docType knowledge.DocumentType
		want    bool
	}{
		{
			name:    "PRD is supported",
			docType: knowledge.TypePRD,
			want:    true,
		},
		{
			name:    "Figma is supported",
			docType: knowledge.TypeFigma,
			want:    true,
		},
		{
			name:    "API Spec is supported",
			docType: knowledge.TypeAPISpec,
			want:    true,
		},
		{
			name:    "unknown is not supported",
			docType: knowledge.DocumentType("unknown"),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.docType.IsSupported(); got != tt.want {
				t.Errorf("DocumentType.IsSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}
