package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"life-financial-assistant-backend/internal/model"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Config struct {
	BaseURL        string
	APIKey         string
	Model          string
	TimeoutSeconds int
}

type AIAgent struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

type anthropicMessageRequest struct {
	Model     string                    `json:"model"`
	MaxTokens int                       `json:"max_tokens"`
	Messages  []anthropicMessageContent `json:"messages"`
}

type anthropicMessageContent struct {
	Role    string                  `json:"role"`
	Content []anthropicContentBlock `json:"content"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicMessageResponse struct {
	Content []anthropicContentBlock `json:"content"`
}

func NewAIAgent(cfg Config) *AIAgent {
	timeoutSeconds := cfg.TimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = 30
	}

	modelName := strings.TrimSpace(cfg.Model)
	if modelName == "" {
		modelName = "claude-3-5-sonnet-20241022"
	}

	return &AIAgent{
		BaseURL: strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/"),
		APIKey:  strings.TrimSpace(cfg.APIKey),
		Model:   modelName,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (a *AIAgent) AnalyzeBills(bills []model.Bill, context string) (string, error) {
	if len(bills) == 0 {
		return "No bills found to analyze.", nil
	}
	if a.BaseURL == "" {
		return "", errors.New("AI_BASE_URL is not configured")
	}
	if a.APIKey == "" {
		return "", errors.New("AI_API_KEY is not configured")
	}

	requestBody := anthropicMessageRequest{
		Model:     a.Model,
		MaxTokens: 1024,
		Messages: []anthropicMessageContent{{
			Role: "user",
			Content: []anthropicContentBlock{{
				Type: "text",
				Text: a.buildPrompt(bills, context),
			}},
		}},
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, a.BaseURL+"/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", a.APIKey)
	request.Header.Set("anthropic-version", "2023-06-01")

	response, err := a.HTTPClient.Do(request)
	if err != nil {
		if isTimeoutError(err) {
			return "", errors.New("AI 服务请求超时")
		}
		return "", errors.New("AI 服务连接失败")
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		_, _ = io.ReadAll(response.Body)
		return "", errors.New("AI 服务返回异常")
	}

	var anthropicResponse anthropicMessageResponse
	if err := json.NewDecoder(response.Body).Decode(&anthropicResponse); err != nil {
		return "", errors.New("AI 服务返回异常")
	}

	for _, block := range anthropicResponse.Content {
		if block.Type == "text" {
			text := strings.TrimSpace(block.Text)
			if text != "" {
				return text, nil
			}
		}
	}

	return "", errors.New("AI 服务返回异常")
}

func (a *AIAgent) buildPrompt(bills []model.Bill, context string) string {
	type categoryTotal struct {
		Name   string
		Amount float64
	}

	var totalIncome float64
	var totalExpense float64
	categoryMap := make(map[string]float64)
	billLines := make([]string, 0, len(bills))

	for _, bill := range bills {
		if bill.Amount >= 0 {
			totalIncome += bill.Amount
		} else {
			totalExpense += -bill.Amount
		}
		categoryMap[bill.Category] += bill.Amount
		billLines = append(billLines, fmt.Sprintf("- 日期:%s 类型:%s 分类:%s 金额:%.2f 备注:%s", bill.Date.Format("2006-01-02"), bill.Type, bill.Category, bill.Amount, strings.TrimSpace(bill.Remarks)))
	}

	categoryTotals := make([]categoryTotal, 0, len(categoryMap))
	for name, amount := range categoryMap {
		categoryTotals = append(categoryTotals, categoryTotal{Name: name, Amount: amount})
	}
	sort.Slice(categoryTotals, func(i, j int) bool {
		return categoryTotals[i].Name < categoryTotals[j].Name
	})

	categoryLines := make([]string, 0, len(categoryTotals))
	for _, item := range categoryTotals {
		categoryLines = append(categoryLines, fmt.Sprintf("- %s: %.2f", item.Name, item.Amount))
	}

	return strings.Join([]string{
		"请基于以下账单数据生成中文纯文本财务分析报告。",
		"输出中文，直接输出分析结论。",
		"不返回 JSON、不要使用代码块、不要添加多余格式。",
		fmt.Sprintf("报告对象: %s", context),
		fmt.Sprintf("总收入: %.2f", totalIncome),
		fmt.Sprintf("总支出: %.2f", totalExpense),
		"分类统计:",
		strings.Join(categoryLines, "\n"),
		"账单明细:",
		strings.Join(billLines, "\n"),
		"请按以下结构输出：财务概览、主要消费分类、风险或异常提醒、简短建议。",
	}, "\n")
}

func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
