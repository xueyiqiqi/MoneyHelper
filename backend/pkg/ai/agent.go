package ai

import (
	"fmt"
	"life-financial-assistant-backend/internal/model"
)

type AIAgent struct {
	APIKey string
}

func NewAIAgent(apiKey string) *AIAgent {
	return &AIAgent{APIKey: apiKey}
}

func (a *AIAgent) AnalyzeBills(bills []model.Bill, context string) (string, error) {
	if len(bills) == 0 {
		return "No bills found to analyze.", nil
	}

	var total float64
	categories := make(map[string]float64)
	for _, b := range bills {
		total += b.Amount
		categories[b.Category] += b.Amount
	}

	report := fmt.Sprintf("财务分析报告 (%s):\n", context)
	report += fmt.Sprintf("总支出: %.2f\n", total)
	report += "支出分类明细:\n"
	for cat, amount := range categories {
		report += fmt.Sprintf("- %s: %.2f\n", cat, amount)
	}

	if a.APIKey == "" {
		report += "\n(注：请配置 OPENAI_API_KEY 以获取更详细的 AI 分析)"
		return report, nil
	}

	// TODO: Implement actual LLM call here
	report += "\n[AI 洞察]: 您的消费主要集中在 " + context + "。建议设置一个预算。"
	return report, nil
}
