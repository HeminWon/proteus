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

func buildModelsURLCandidates(baseURL string) []string {
	base := trimTrailingSlash(baseURL)
	if strings.Contains(base, "openrouter.ai/api/v1") {
		return []string{base + "/auth/key"}
	}

	candidates := make([]string, 0, 4)
	seen := map[string]struct{}{}
	appendCandidate := func(url string) {
		if _, ok := seen[url]; ok {
			return
		}
		seen[url] = struct{}{}
		candidates = append(candidates, url)
	}

	if strings.HasSuffix(base, "/v1") {
		appendCandidate(base + "/models")
		appendCandidate(strings.TrimSuffix(base, "/v1") + "/models")
		return candidates
	}

	appendCandidate(base + "/v1/models")
	appendCandidate(base + "/models")
	if strings.HasSuffix(base, "/anthropic") {
		root := strings.TrimSuffix(base, "/anthropic")
		appendCandidate(root + "/v1/models")
		appendCandidate(root + "/models")
	}
	return candidates
}

func readBodySnippet(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	text := strings.TrimSpace(string(body))
	if text == "" {
		return ""
	}
	flat := strings.Join(strings.Fields(text), " ")
	if len(flat) > 200 {
		flat = flat[:200]
	}
	return flat
}

func buildFailDetail(statusCode int, url string, bodySnippet string) string {
	detail := fmt.Sprintf("HTTP %d (%s)", statusCode, url)
	if bodySnippet != "" {
		detail = fmt.Sprintf("HTTP %d (%s) | body=%q", statusCode, url, bodySnippet)
	}
	return detail
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

	urls := buildModelsURLCandidates(baseURL)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client := &http.Client{Timeout: 20 * time.Second}

	var failDetails []string
	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			ms := time.Since(startedAt).Milliseconds()
			return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: "request error: " + err.Error(), LatencyMs: &ms}
		}

		req.Header.Set("x-api-key", token)
		req.Header.Set("authorization", "Bearer "+token)
		req.Header.Set("anthropic-version", "2023-06-01")
		req.Header.Set("accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			ms := time.Since(startedAt).Milliseconds()
			return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: "request error: " + err.Error(), LatencyMs: &ms}
		}

		latency := time.Since(startedAt).Milliseconds()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			return LiveValidationResult{ProviderID: provider.ID, Status: "ok", Detail: fmt.Sprintf("HTTP %d (%s)", resp.StatusCode, url), LatencyMs: &latency}
		}

		bodySnippet := readBodySnippet(resp)
		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
			failDetails = append(failDetails, buildFailDetail(resp.StatusCode, url, bodySnippet))
			continue
		}

		return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: buildFailDetail(resp.StatusCode, url, bodySnippet), LatencyMs: &latency}
	}

	latency := time.Since(startedAt).Milliseconds()
	if len(failDetails) > 0 {
		return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: "all candidates failed: " + strings.Join(failDetails, " ; "), LatencyMs: &latency}
	}
	return LiveValidationResult{ProviderID: provider.ID, Status: "fail", Detail: "no live validation URL candidate", LatencyMs: &latency}
}
