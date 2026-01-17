package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	geminiAPIURL    = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-pro:generateContent?key=%s"
	maxRetries      = 5
	initialBackoff  = 2 * time.Second
)

// GenerateAICommitMessage uses the Gemini API (via HTTP POST) to generate a commit message based on the provided diff.
func GenerateAICommitMessage(diff string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	url := fmt.Sprintf(geminiAPIURL, apiKey)

	prompt := fmt.Sprintf(`Generate a concise conventional commit message (type: subject) for the following Git diff.
The commit message should accurately summarize the changes.
Do not include any explanations or additional text, just the commit message.

Example: feat: add new user authentication endpoint

Diff:
%s`, diff)

	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	resp, err := postWithRetry(url, requestBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading API response: %w", err)
	}

	var geminiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(responseBody, &geminiResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling API response: %w", err)
	}

	if len(geminiResponse.Candidates) > 0 && len(geminiResponse.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(geminiResponse.Candidates[0].Content.Parts[0].Text), nil
	}

	return "", fmt.Errorf("no content generated from Gemini")
}

// postWithRetry sends a POST request with a retry mechanism for rate limiting.
func postWithRetry(url string, body []byte) (*http.Response, error) {
	var resp *http.Response
	var err error
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("error making API request: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			fmt.Printf("Rate limit exceeded. Retrying in %v...\n", backoff)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil, fmt.Errorf("exceeded max retries")
}