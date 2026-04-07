// Package generation provides AI generation services
package generation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// GenerateCasesRequest contains test case generation parameters
type GenerateCasesRequest struct {
	Prompt       string   `json:"prompt" validate:"required,min=10"`
	Context      string   `json:"context"`        // Retrieved document context
	CaseCount    int      `json:"case_count"`     // Number of cases to generate (1-20)
	SceneTypes   []string `json:"scene_types"`    // positive, negative, boundary
	Priority     string   `json:"priority"`       // P0-P3
	CaseType     string   `json:"case_type"`      // functionality, api, etc.
	ModelVersion string   `json:"model_version"`  // LLM model to use
}

// GeneratedCase represents a single generated test case
type GeneratedCase struct {
	Title         string         `json:"title"`
	Preconditions []string       `json:"preconditions"`
	Steps         []string       `json:"steps"`
	Expected      map[string]any `json:"expected"`
	CaseType      string         `json:"case_type"`
	Priority      string         `json:"priority"`
	Reasoning     string         `json:"reasoning"` // AI's reasoning for this case
}

// GenerateCasesResponse contains generated test cases
type GenerateCasesResponse struct {
	Cases        []*GeneratedCase `json:"cases"`
	ModelVersion string           `json:"model_version"`
	TokensUsed   int64            `json:"tokens_used"`
}

// GenerateEmbeddingRequest contains embedding generation parameters
type GenerateEmbeddingRequest struct {
	Text string `json:"text" validate:"required"`
}

// GenerateEmbeddingResponse contains generated embedding
type GenerateEmbeddingResponse struct {
	Embedding []byte `json:"embedding"`
	Model     string `json:"model"`
}

// LLMService provides LLM (Large Language Model) operations
type LLMService interface {
	// GenerateCases generates test cases using LLM
	GenerateCases(ctx context.Context, req *GenerateCasesRequest) (*GenerateCasesResponse, error)

	// GenerateEmbedding generates vector embedding for text
	GenerateEmbedding(ctx context.Context, req *GenerateEmbeddingRequest) (*GenerateEmbeddingResponse, error)

	// GetModelVersion returns the current model version
	GetModelVersion() string
}

// LLMServiceImpl implements LLMService
type LLMServiceImpl struct {
	apiKey       string
	apiEndpoint  string
	modelVersion string
	httpClient   *http.Client
	timeout      time.Duration
}

// NewLLMService creates a new LLMService instance
func NewLLMService(apiKey, apiEndpoint, modelVersion string, timeout time.Duration) LLMService {
	return &LLMServiceImpl{
		apiKey:       apiKey,
		apiEndpoint:  apiEndpoint,
		modelVersion: modelVersion,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// GenerateCases generates test cases using LLM
func (s *LLMServiceImpl) GenerateCases(ctx context.Context, req *GenerateCasesRequest) (*GenerateCasesResponse, error) {
	// Validate prompt length
	if len(req.Prompt) < 10 {
		return nil, errors.New("prompt must be at least 10 characters")
	}

	// Validate case count
	caseCount := req.CaseCount
	if caseCount <= 0 {
		caseCount = 1
	}
	if caseCount > 20 {
		return nil, errors.New("case count must be between 1 and 20")
	}

	// Validate scene types
	validSceneTypes := map[string]bool{
		"positive": true, "negative": true, "boundary": true,
	}
	for _, st := range req.SceneTypes {
		if !validSceneTypes[st] {
			return nil, fmt.Errorf("invalid scene type: %s", st)
		}
	}

	// Validate priority
	validPriorities := map[string]bool{
		"P0": true, "P1": true, "P2": true, "P3": true,
	}
	if req.Priority != "" && !validPriorities[req.Priority] {
		return nil, fmt.Errorf("invalid priority: %s", req.Priority)
	}

	// Build prompt for LLM
	systemPrompt := s.buildSystemPrompt(req)
	userPrompt := s.buildUserPrompt(req)

	// Call LLM API
	// TODO: Implement actual API call
	// For now, return placeholder response
	_ = systemPrompt
	_ = userPrompt

	// Mock response for testing
	cases := make([]*GeneratedCase, 0, caseCount)
	for i := 0; i < caseCount; i++ {
		cases = append(cases, &GeneratedCase{
			Title:         fmt.Sprintf("Generated Test Case %d", i+1),
			Preconditions: []string{"System is running", "User is logged in"},
			Steps:         []string{"Step 1", "Step 2", "Step 3"},
			Expected:      map[string]any{"status": "success"},
			CaseType:       req.CaseType,
			Priority:       req.Priority,
			Reasoning:      "Generated based on context",
		})
	}

	return &GenerateCasesResponse{
		Cases:        cases,
		ModelVersion: s.modelVersion,
		TokensUsed:   1000, // Mock token usage
	}, nil
}

// GenerateEmbedding generates vector embedding for text
func (s *LLMServiceImpl) GenerateEmbedding(ctx context.Context, req *GenerateEmbeddingRequest) (*GenerateEmbeddingResponse, error) {
	// Validate text
	if req.Text == "" {
		return nil, errors.New("text cannot be empty")
	}

	// Check text length (typical limit is 8192 tokens)
	if len(req.Text) > 32000 {
		return nil, errors.New("text exceeds maximum length")
	}

	// TODO: Implement actual embedding API call
	// For now, return placeholder embedding
	embedding := make([]byte, 1536*4) // 1536 dimensions * 4 bytes (float32)
	for i := range embedding {
		embedding[i] = byte(i % 256)
	}

	return &GenerateEmbeddingResponse{
		Embedding: embedding,
		Model:     s.modelVersion,
	}, nil
}

// GetModelVersion returns the current model version
func (s *LLMServiceImpl) GetModelVersion() string {
	return s.modelVersion
}

// buildSystemPrompt builds the system prompt for test case generation
func (s *LLMServiceImpl) buildSystemPrompt(req *GenerateCasesRequest) string {
	return fmt.Sprintf(`You are an expert QA engineer. Generate test cases based on the following specifications:

Case Type: %s
Priority: %s
Scene Types: %v

Generate comprehensive test cases that cover:
1. Happy path scenarios
2. Edge cases
3. Error conditions

Each test case should have:
- Clear title
- Preconditions
- Step-by-step instructions
- Expected results`, req.CaseType, req.Priority, req.SceneTypes)
}

// buildUserPrompt builds the user prompt for test case generation
func (s *LLMServiceImpl) buildUserPrompt(req *GenerateCasesRequest) string {
	prompt := fmt.Sprintf("Generate %d test cases for:\n\n%s\n", req.CaseCount, req.Prompt)
	if req.Context != "" {
		prompt += fmt.Sprintf("\nContext:\n%s\n", req.Context)
	}
	return prompt
}

// parseGeneratedCases parses LLM response into GeneratedCase structs
func (s *LLMServiceImpl) parseGeneratedCases(response string) ([]*GeneratedCase, error) {
	var cases []*GeneratedCase
	if err := json.Unmarshal([]byte(response), &cases); err != nil {
		return nil, fmt.Errorf("parse generated cases: %w", err)
	}
	return cases, nil
}
