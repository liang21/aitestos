// Package llm provides LLM (Large Language Model) service implementations
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/liang21/aitestos/internal/config"
	generationservice "github.com/liang21/aitestos/internal/service/generation"
)

const (
	defaultTimeout = 120 * time.Second
)

// Client implements LLM service operations
type Client struct {
	config     *config.LLMConfig
	httpClient *http.Client
}

// NewClient creates a new LLM client
func NewClient(cfg *config.LLMConfig) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

// GenerateCases generates test cases based on requirements using LLM
func (c *Client) GenerateCases(ctx context.Context, req *generationservice.GenerateCasesRequest) (*generationservice.GenerateCasesResponse, error) {
	// Build the prompt for test case generation
	prompt := c.buildTestCasePrompt(req)

	// Call LLM API
	response, err := c.callChatAPI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("call LLM API: %w", err)
	}

	// Parse response into test cases
	cases, err := c.parseTestCasesResponse(response)
	if err != nil {
		return nil, fmt.Errorf("parse test cases: %w", err)
	}

	return &generationservice.GenerateCasesResponse{
		Cases:        c.convertToServiceCases(cases),
		ModelVersion: c.getModelIdentifier(),
		TokensUsed:   int64(c.estimateTokens(prompt + response)),
	}, nil
}

// GenerateEmbedding generates vector embedding for text
func (c *Client) GenerateEmbedding(ctx context.Context, req *generationservice.GenerateEmbeddingRequest) (*generationservice.GenerateEmbeddingResponse, error) {
	// Call embedding API
	embedding, model, err := c.callEmbeddingAPI(ctx, req.Text)
	if err != nil {
		return nil, fmt.Errorf("call embedding API: %w", err)
	}

	// Convert []float32 to []byte
	embeddingBytes := float32SliceToBytes(embedding)

	return &generationservice.GenerateEmbeddingResponse{
		Embedding: embeddingBytes,
		Model:     model,
	}, nil
}

// GetModelVersion returns the current model version
func (c *Client) GetModelVersion() string {
	return c.getModelIdentifier()
}

// Close closes the LLM client
func (c *Client) Close() error {
	return nil
}

// Internal types for API calls

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage,omitempty"`
}

type choice struct {
	Index        int     `json:"index"`
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type embeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingResponse struct {
	Object string      `json:"object"`
	Data   []embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  usage       `json:"usage"`
}

type embedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// Internal test case type for parsing
type internalGeneratedCase struct {
	Title         string         `json:"title"`
	Preconditions []string       `json:"preconditions"`
	Steps         []string       `json:"steps"`
	Expected      map[string]any `json:"expected"`
	CaseType      string         `json:"case_type"`
	Priority      string         `json:"priority"`
	Reasoning     string         `json:"reasoning"`
}

// API call methods

func (c *Client) callChatAPI(ctx context.Context, prompt string) (string, error) {
	reqBody := chatRequest{
		Model: c.config.Model,
		Messages: []message{
			{
				Role:    "system",
				Content: "You are a test case generation expert. Generate comprehensive test cases based on the given requirements. Output only valid JSON without markdown code blocks.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint := c.getChatEndpoint()
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqJSON))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (c *Client) callEmbeddingAPI(ctx context.Context, text string) ([]float32, string, error) {
	reqBody := embeddingRequest{
		Model: c.config.EmbeddingModel,
		Input: text,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint := c.getEmbeddingEndpoint()
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqJSON))
	if err != nil {
		return nil, "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API error: status %d: %s", resp.StatusCode, string(body))
	}

	var embResp embeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(embResp.Data) == 0 {
		return nil, "", fmt.Errorf("no embeddings in response")
	}

	return embResp.Data[0].Embedding, embResp.Model, nil
}

// Helper methods

func (c *Client) getChatEndpoint() string {
	switch c.config.Provider {
	case "deepseek":
		return "https://api.deepseek.com/v1/chat/completions"
	case "azure":
		return fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=2024-02-15-preview",
			c.config.BaseURL, c.config.Model)
	case "openai", "":
		return "https://api.openai.com/v1/chat/completions"
	default:
		return "https://api.openai.com/v1/chat/completions"
	}
}

func (c *Client) getEmbeddingEndpoint() string {
	switch c.config.Provider {
	case "deepseek":
		return "https://api.deepseek.com/v1/embeddings"
	case "azure":
		return fmt.Sprintf("%s/openai/deployments/%s/embeddings?api-version=2024-02-15-preview",
			c.config.BaseURL, c.config.EmbeddingModel)
	case "openai", "":
		return "https://api.openai.com/v1/embeddings"
	default:
		return "https://api.openai.com/v1/embeddings"
	}
}

func (c *Client) getModelIdentifier() string {
	return fmt.Sprintf("%s/%s", c.config.Provider, c.config.Model)
}

func (c *Client) estimateTokens(text string) int {
	return len(text) / 4
}

func (c *Client) buildTestCasePrompt(req *generationservice.GenerateCasesRequest) string {
	contextText := ""
	if req.Context != "" {
		contextText = fmt.Sprintf("\n\nContext from existing documents:\n%s\n", req.Context)
	}

	priorityText := ""
	if req.Priority != "" {
		priorityText = fmt.Sprintf("\nPriority level: %s", req.Priority)
	}

	caseTypeText := ""
	if req.CaseType != "" {
		caseTypeText = fmt.Sprintf("\nTest case type: %s", req.CaseType)
	}

	sceneTypesText := ""
	if len(req.SceneTypes) > 0 {
		sceneTypesText = fmt.Sprintf("\nScene types to cover: %v", req.SceneTypes)
	}

	return fmt.Sprintf(`Generate %d test cases based on the following requirements:%s%s%s%s

Requirements:
%s

Please generate test cases in the following JSON format (output ONLY the JSON array, no markdown):
[
  {
    "title": "Test case title",
    "preconditions": ["Precondition 1", "Precondition 2"],
    "steps": ["Step 1", "Step 2", "Step 3"],
    "expected": {"result": "success", "response_code": 200},
    "case_type": "functional",
    "priority": "high",
    "reasoning": "This case verifies..."
  }
]

Make the test cases:
- Specific and actionable
- Cover both positive and negative scenarios
- Include edge cases where appropriate`,
		req.CaseCount, contextText, priorityText, caseTypeText, sceneTypesText, req.Prompt)
}

func (c *Client) parseTestCasesResponse(response string) ([]*internalGeneratedCase, error) {
	response = extractJSON(response)

	var cases []*internalGeneratedCase
	if err := json.Unmarshal([]byte(response), &cases); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	return cases, nil
}

func (c *Client) convertToServiceCases(internalCases []*internalGeneratedCase) []*generationservice.GeneratedCase {
	serviceCases := make([]*generationservice.GeneratedCase, len(internalCases))
	for i, c := range internalCases {
		serviceCases[i] = &generationservice.GeneratedCase{
			Title:         c.Title,
			Preconditions: c.Preconditions,
			Steps:         c.Steps,
			Expected:      c.Expected,
			CaseType:      c.CaseType,
			Priority:      c.Priority,
			Reasoning:     c.Reasoning,
		}
	}
	return serviceCases
}

func extractJSON(response string) string {
	start := 0
	end := len(response)

	if idx := findIndex(response, "["); idx != -1 {
		start = idx
	}

	// Find matching closing bracket
	depth := 0
	for i := start; i < len(response); i++ {
		if response[i] == '[' {
			depth++
		} else if response[i] == ']' {
			depth--
			if depth == 0 {
				end = i + 1
				break
			}
		}
	}

	return response[start:end]
}

func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func float32SliceToBytes(data []float32) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	result := make([]byte, len(data)*4)
	for i, v := range data {
		bits := math.Float32bits(v)
		result[i*4] = byte(bits)
		result[i*4+1] = byte(bits >> 8)
		result[i*4+2] = byte(bits >> 16)
		result[i*4+3] = byte(bits >> 24)
	}
	return result
}
