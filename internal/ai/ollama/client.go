package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	model   string
	client  *http.Client
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewClient(baseURL, model string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "mistral"
	}

	return &Client{
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) IsSpam(emailAddress string, subject string, user_prompt string) (bool, error) {
	prompt := c.buildPrompt(emailAddress, subject, user_prompt)

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	response, err := c.generate(ctx, prompt)
	if err != nil {
		return false, fmt.Errorf("failed to generate response: %w", err)
	}

	return c.parseResponse(response), nil
}

func (c *Client) buildPrompt(emailAddress string, subject string, prompt string) string {
	return fmt.Sprintf(`You are a spam email classifier. Analyze the following email and determine if it's spam.
	From: %s
	Subject: %s
	Answer with ONLY "SPAM" if it's spam or "HAM" if it's not spam. No explanations. %s`, emailAddress, subject, prompt)
}

func (c *Client) generate(ctx context.Context, prompt string) (string, error) {
	reqBody := generateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, nil
}

func (c *Client) parseResponse(response string) bool {
	cleaned := strings.ToLower(strings.TrimSpace(response))

	return strings.Contains(cleaned, "spam")
}
