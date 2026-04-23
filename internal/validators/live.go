package validators

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HeminWon/proteus/internal/providers"
)

type LiveValidationResult struct {
	ProviderID string
	Status     string
	Detail     string
	LatencyMs  *int64
}

func trimTrailingSlash(value string) string {
	return strings.TrimRight(value, "/")
}

func buildModelsURL(baseURL string) string {
	base := trimTrailingSlash(baseURL)
	if strings.Contains(base, "openrouter.ai/api/v1") {
		return base + "/auth/key"
	}
	if strings.HasSuffix(base, "/v1") {
		return base + "/models"
	}
	return base + "/v1/models"
}

func ValidateProviderLive(provider providers.Provider) LiveValidationResult {
	startedAt := time.Now()
	token := provider.Claude.Env["ANTHROPIC_AUTH_TOKEN"]
	baseURL := provider.Claude.Env["ANTHROPIC_BASE_URL"]

	if baseURL == "" {
		return LiveValidationResult{ProviderID: provider.ID, Status: "skip", Detail: "missing ANTHROPIC_BASE_URL", LatencyMs: nil}
	}
	if token == "" {
		return LiveValidationResult{ProviderID: provider.ID, Status: "skip", Detail: "missing ANTHROPIC_AUTH_TOKEN", LatencyMs: nil}
	}

	url := buildModelsURL(baseURL)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ms := time.Since(startedAt).Milliseconds()
		return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: "request error: " + err.Error(), LatencyMs: &ms}
	}

	req.Header.Set("x-api-key", token)
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("accept", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		ms := time.Since(startedAt).Milliseconds()
		return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: "request error: " + err.Error(), LatencyMs: &ms}
	}
	defer resp.Body.Close()

	latency := time.Since(startedAt).Milliseconds()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return LiveValidationResult{ProviderID: provider.ID, Status: "ok", Detail: fmt.Sprintf("HTTP %d (%s)", resp.StatusCode, url), LatencyMs: &latency}
	}

	bodySnippet := ""
	body, readErr := io.ReadAll(resp.Body)
	if readErr == nil {
		text := strings.TrimSpace(string(body))
		if text != "" {
			flat := strings.Join(strings.Fields(text), " ")
			if len(flat) > 200 {
				flat = flat[:200]
			}
			bodySnippet = flat
		}
	}

	detail := fmt.Sprintf("HTTP %d (%s)", resp.StatusCode, url)
	if bodySnippet != "" {
		detail = fmt.Sprintf("HTTP %d (%s) | body=%q", resp.StatusCode, url, bodySnippet)
	}
	return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: detail, LatencyMs: &latency}
}
