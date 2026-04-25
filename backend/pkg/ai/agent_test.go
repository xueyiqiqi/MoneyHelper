package ai

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"life-financial-assistant-backend/internal/logger"
	"life-financial-assistant-backend/internal/model"
)

func TestNewAIAgentAppliesConfigValues(t *testing.T) {
	agent := NewAIAgent(Config{
		BaseURL:        "https://example-proxy.test",
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 45,
	})

	if agent.BaseURL != "https://example-proxy.test" {
		t.Fatalf("expected base URL to be copied into agent, got %q", agent.BaseURL)
	}
	if agent.APIKey != "secret-token" {
		t.Fatalf("expected API key to be copied into agent, got %q", agent.APIKey)
	}
	if agent.Model != "claude-3-5-sonnet-20241022" {
		t.Fatalf("expected model to be copied into agent, got %q", agent.Model)
	}
	if agent.HTTPClient == nil {
		t.Fatal("expected HTTP client to be initialized")
	}
	if agent.HTTPClient.Timeout != 45*time.Second {
		t.Fatalf("expected timeout 45s, got %v", agent.HTTPClient.Timeout)
	}
}

func TestAnalyzeBillsReturnsErrorWhenBaseURLMissing(t *testing.T) {
	agent := NewAIAgent(Config{
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	_, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "personal")
	if err == nil || err.Error() != "AI_BASE_URL is not configured" {
		t.Fatalf("expected missing base URL error, got %v", err)
	}
}

func TestAnalyzeBillsReturnsErrorWhenAPIKeyMissing(t *testing.T) {
	agent := NewAIAgent(Config{
		BaseURL:        "https://example-proxy.test",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	_, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "personal")
	if err == nil || err.Error() != "AI_API_KEY is not configured" {
		t.Fatalf("expected missing API key error, got %v", err)
	}
}

func TestAnalyzeBillsLogsStatusCodeAndBodyPreviewOnNon2xxResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":{"message":"upstream unavailable for maintenance"}}`))
	}))
	defer server.Close()

	var logBuffer bytes.Buffer
	originalErrorLogger := logger.Error
	logger.Error = log.New(&logBuffer, "[ERROR] ", 0)
	t.Cleanup(func() {
		logger.Error = originalErrorLogger
	})

	agent := NewAIAgent(Config{
		BaseURL:        server.URL,
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	_, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "personal")
	if err == nil || err.Error() != "AI 服务返回异常" {
		t.Fatalf("expected gateway error, got %v", err)
	}

	logged := logBuffer.String()
	for _, expected := range []string{
		"AI gateway returned non-2xx response",
		"status=502",
		"model=claude-3-5-sonnet-20241022",
		"/v1/messages",
		"upstream unavailable for maintenance",
	} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("expected log to contain %q, got %s", expected, logged)
		}
	}
}

func TestAnalyzeBillsLogsDecodeFailureWithBodyPreview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not-json-response-from-upstream`))
	}))
	defer server.Close()

	var logBuffer bytes.Buffer
	originalErrorLogger := logger.Error
	logger.Error = log.New(&logBuffer, "[ERROR] ", 0)
	t.Cleanup(func() {
		logger.Error = originalErrorLogger
	})

	agent := NewAIAgent(Config{
		BaseURL:        server.URL,
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	_, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "personal")
	if err == nil || err.Error() != "AI 服务返回异常" {
		t.Fatalf("expected decode failure error, got %v", err)
	}

	logged := logBuffer.String()
	for _, expected := range []string{
		"AI gateway returned invalid JSON",
		"model=claude-3-5-sonnet-20241022",
		"/v1/messages",
		"not-json-response-from-upstream",
	} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("expected log to contain %q, got %s", expected, logged)
		}
	}
}

func TestAnalyzeBillsReturnsParsedTextFromAnthropicResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Fatalf("expected path /v1/messages, got %s", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "secret-token" {
			t.Fatalf("expected x-api-key header, got %q", got)
		}
		if got := r.Header.Get("anthropic-version"); got != "2023-06-01" {
			t.Fatalf("expected anthropic-version header, got %q", got)
		}
		_, _ = w.Write([]byte(`{
			"content": [
				{"type":"text","text":"这是模型生成的财务分析"}
			]
		}`))
	}))
	defer server.Close()

	agent := NewAIAgent(Config{
		BaseURL:        server.URL,
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	result, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "space")
	if err != nil {
		t.Fatalf("expected success, got error %v", err)
	}
	if result != "这是模型生成的财务分析" {
		t.Fatalf("expected parsed text response, got %q", result)
	}
}

func TestAnalyzeBillsLogsContentTypesWhenResponseHasNoTextBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"content":[{"type":"tool_use","text":""},{"type":"thinking","text":"internal"}]}`))
	}))
	defer server.Close()

	var logBuffer bytes.Buffer
	originalErrorLogger := logger.Error
	logger.Error = log.New(&logBuffer, "[ERROR] ", 0)
	t.Cleanup(func() {
		logger.Error = originalErrorLogger
	})

	agent := NewAIAgent(Config{
		BaseURL:        server.URL,
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	_, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "shared")
	if err == nil || err.Error() != "AI 服务返回异常" {
		t.Fatalf("expected malformed response error, got %v", err)
	}

	logged := logBuffer.String()
	for _, expected := range []string{
		"AI gateway returned no text block",
		"content_count=2",
		"content_types=tool_use,thinking",
		"model=claude-3-5-sonnet-20241022",
	} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("expected log to contain %q, got %s", expected, logged)
		}
	}
}

func TestAnalyzeBillsLogsTimeoutDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
		_, _ = w.Write([]byte(`{"content":[{"type":"text","text":"slow response"}]}`))
	}))
	defer server.Close()

	var logBuffer bytes.Buffer
	originalErrorLogger := logger.Error
	logger.Error = log.New(&logBuffer, "[ERROR] ", 0)
	t.Cleanup(func() {
		logger.Error = originalErrorLogger
	})

	agent := NewAIAgent(Config{
		BaseURL:        server.URL,
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 0,
	})
	agent.HTTPClient.Timeout = 50 * time.Millisecond

	_, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "personal")
	if err == nil || err.Error() != "AI 服务请求超时" {
		t.Fatalf("expected timeout error, got %v", err)
	}

	logged := logBuffer.String()
	for _, expected := range []string{
		"AI gateway request timed out",
		"model=claude-3-5-sonnet-20241022",
		"/v1/messages",
	} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("expected log to contain %q, got %s", expected, logged)
		}
	}
}

func TestAnalyzeBillsLogsConnectionFailureDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"content":[{"type":"text","text":"ok"}]}`))
	}))
	serverURL := server.URL
	server.Close()

	var logBuffer bytes.Buffer
	originalErrorLogger := logger.Error
	logger.Error = log.New(&logBuffer, "[ERROR] ", 0)
	t.Cleanup(func() {
		logger.Error = originalErrorLogger
	})

	agent := NewAIAgent(Config{
		BaseURL:        serverURL,
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	_, err := agent.AnalyzeBills([]model.Bill{{Category: "food", Amount: -20}}, "personal")
	if err == nil || err.Error() != "AI 服务连接失败" {
		t.Fatalf("expected connection failure error, got %v", err)
	}

	logged := logBuffer.String()
	for _, expected := range []string{
		"AI gateway request failed",
		"model=claude-3-5-sonnet-20241022",
		"/v1/messages",
	} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("expected log to contain %q, got %s", expected, logged)
		}
	}
}

func TestBuildPromptIncludesContextAndStructuredInstructions(t *testing.T) {
	agent := NewAIAgent(Config{
		BaseURL:        "https://example-proxy.test",
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	prompt := agent.buildPrompt([]model.Bill{{
		Category: "food",
		Type:     "expense",
		Amount:   -20,
		Remarks:  "午饭",
		Date:     time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC),
	}}, "space")

	for _, expected := range []string{
		"报告对象: space",
		"输出中文",
		"不返回 JSON",
		"财务概览",
		"主要消费分类",
		"风险或异常提醒",
		"简短建议",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected prompt to contain %q, got %s", expected, prompt)
		}
	}
}

func TestAnalyzeBillsReturnsNoBillsMessageWithoutGatewayCall(t *testing.T) {
	agent := NewAIAgent(Config{
		BaseURL:        "https://example-proxy.test",
		APIKey:         "secret-token",
		Model:          "claude-3-5-sonnet-20241022",
		TimeoutSeconds: 30,
	})

	result, err := agent.AnalyzeBills(nil, "personal")
	if err != nil {
		t.Fatalf("expected no error for empty bills, got %v", err)
	}
	if result != "No bills found to analyze." {
		t.Fatalf("expected empty-bills message, got %q", result)
	}
}
