package ai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/urstruelysv/autocommit-cli/internal/logger"
)

const (
	maxRetries     = 5
	initialBackoff = 2 * time.Second
	requestTimeout = 60 * time.Second
)

// postWithRetry sends a POST request with a retry mechanism for rate limiting.
func postWithRetry(log logger.Logger, url string, headers map[string]string, body []byte) (*http.Response, error) {
	backoff := initialBackoff

	client := &http.Client{Timeout: requestTimeout}

	for i := 0; i < maxRetries; i++ {
		log.Debug("Making API request (attempt %d)", i+1)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("error creating API request: %w", err)
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making API request: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			log.Info("Rate limit exceeded. Retrying in %v...", backoff)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil, fmt.Errorf("exceeded max retries")
}
