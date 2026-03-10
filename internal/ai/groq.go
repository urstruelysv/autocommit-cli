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
	groqAPIURL      = "https://api.groq.com/openai/v1/chat/completions"
	groqDefaultModel = "llama-3.1-8b-instant"
)

// GenerateGroqCommitMessage uses the Groq OpenAI-compatible API to generate a commit message.
func GenerateGroqCommitMessage(log logger.Logger, diff string) (string, error) {
	log.Debug("Generating AI commit message via Groq...")
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	model := strings.TrimSpace(os.Getenv("GROQ_MODEL"))
	if model == "" {
		model = groqDefaultModel
	}

	system := "You are a senior engineer. Return ONLY a concise conventional commit message (type: subject)."
	user := fmt.Sprintf(`Generate a concise conventional commit message (type: subject) for the following Git diff.
The commit message should accurately summarize the changes.
Do not include any explanations or additional text, just the commit message.

Example: feat: add new user authentication endpoint

Diff:
%s`, diff)

	log.Debug("Prompt:\n%s", user)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"temperature": groqTemperature(),
		"max_tokens":  groqMaxTokens(),
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	resp, err := postWithRetry(log, groqAPIURL, map[string]string{
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

	var groqResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(responseBody, &groqResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling API response: %w", err)
	}

	if groqResponse.Error != nil && groqResponse.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", groqResponse.Error.Message)
	}

	if len(groqResponse.Choices) > 0 {
		message := strings.TrimSpace(groqResponse.Choices[0].Message.Content)
		if message != "" {
			log.Debug("Generated commit message: %s", message)
			return message, nil
		}
	}

	return "", fmt.Errorf("no content generated from Groq")
}

func groqMaxTokens() int {
	const defaultMax = 80
	if v := strings.TrimSpace(os.Getenv("GROQ_MAX_TOKENS")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultMax
}

func groqTemperature() float32 {
	const defaultTemp = 0.2
	if v := strings.TrimSpace(os.Getenv("GROQ_TEMPERATURE")); v != "" {
		if parsed, err := strconv.ParseFloat(v, 32); err == nil && parsed >= 0 {
			return float32(parsed)
		}
	}
	return defaultTemp
}
