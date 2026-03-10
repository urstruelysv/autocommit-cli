package ai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	var lastRateLimitDetail string

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
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			lastRateLimitDetail = strings.TrimSpace(string(bodyBytes))
			sleepFor := backoff
			if retryAfter := strings.TrimSpace(resp.Header.Get("Retry-After")); retryAfter != "" {
				if secs, err := strconv.Atoi(retryAfter); err == nil && secs > 0 {
					sleepFor = time.Duration(secs) * time.Second
				}
			}

			if lastRateLimitDetail != "" {
				log.Info("Rate limit exceeded. Details: %s", lastRateLimitDetail)
			} else {
				log.Info("Rate limit exceeded.")
			}
			log.Info("Retrying in %v...", sleepFor)

			time.Sleep(sleepFor)
			backoff *= 2
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if lastRateLimitDetail != "" {
		return nil, fmt.Errorf("exceeded max retries: %s", lastRateLimitDetail)
	}
	return nil, fmt.Errorf("exceeded max retries")
}
