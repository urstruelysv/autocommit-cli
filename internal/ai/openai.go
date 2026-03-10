package ai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/urstruelysv/autocommit-cli/internal/logger"
)

const (
	openaiAPIURL      = "https://api.openai.com/v1/responses"
	openaiDefaultModel = "gpt-5-mini"
)

// GenerateOpenAICommitMessage uses the OpenAI Responses API (via HTTP POST) to generate a commit message based on the provided diff.
func GenerateOpenAICommitMessage(log logger.Logger, diff string) (string, error) {
	log.Debug("Generating AI commit message via OpenAI...")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = openaiDefaultModel
	}

	prompt := fmt.Sprintf(`Generate a concise conventional commit message (type: subject) for the following Git diff.
The commit message should accurately summarize the changes.
Do not include any explanations or additional text, just the commit message.

Example: feat: add new user authentication endpoint

Diff:
%s`, diff)

	log.Debug("Prompt:\n%s", prompt)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":              model,
		"input":              prompt,
		"max_output_tokens":  maxOutputTokens(),
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	resp, err := postWithRetry(log, openaiAPIURL, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}, requestBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading API response: %w", err)
	}

	log.Debug("API Response: %s", string(responseBody))

	var openaiResponse struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(responseBody, &openaiResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling API response: %w", err)
	}

	if openaiResponse.Error != nil && openaiResponse.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", openaiResponse.Error.Message)
	}

	if strings.TrimSpace(openaiResponse.OutputText) != "" {
		message := strings.TrimSpace(openaiResponse.OutputText)
		log.Debug("Generated commit message: %s", message)
		return message, nil
	}

	for _, item := range openaiResponse.Output {
		for _, part := range item.Content {
			if strings.TrimSpace(part.Text) != "" {
				message := strings.TrimSpace(part.Text)
				log.Debug("Generated commit message: %s", message)
				return message, nil
			}
		}
	}

	return "", fmt.Errorf("no content generated from OpenAI")
}

func maxOutputTokens() int {
	const defaultMax = 120
	if v := strings.TrimSpace(os.Getenv("OPENAI_MAX_OUTPUT_TOKENS")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultMax
}
