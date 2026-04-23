package ai

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestAnalyzeBillsReturnsGatewayErrorOnNon2xxResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":{"message":"upstream unavailable"}}`))
	}))
	defer server.Close()

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

func TestAnalyzeBillsReturnsTimeoutErrorWhenGatewayTooSlow(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
		_, _ = w.Write([]byte(`{"content":[{"type":"text","text":"slow response"}]}`))
	}))
	defer server.Close()

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
}

func TestAnalyzeBillsReturnsErrorWhenAnthropicResponseHasNoTextBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"content":[{"type":"tool_use","text":""}]}`))
	}))
	defer server.Close()

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
