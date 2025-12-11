package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func SummarizeWithAI(config *Config, htmlContent, url string) (string, error) {
	text := extractTextFromHTML(htmlContent)
	
	if len(text) > 8000 {
		text = text[:8000]
	}

	userPrompt := config.SummaryPrompt
	if userPrompt == "" {
		userPrompt = "Summarize the key news topics and main stories from this website. Focus on the most important headlines and provide a concise overview in markdown format."
	}

	prompt := fmt.Sprintf(`%s

Website URL: %s

Content:
%s

Format your response in clean markdown with headings, bullet points, and clear structure.`, userPrompt, url, text)

	reqBody := ChatCompletionRequest{
		Model: config.LLMAPIModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", config.LLMAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.LLMAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	summary := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	return summary, nil
}

func ExtractFocusedContent(config *Config, summary string) (string, error) {
	if config.FocusTopics == "" {
		return "", nil
	}

	prompt := fmt.Sprintf(`Extract only the content related to these topics: %s

From this summary, extract and list ONLY the items that are directly related to the specified topics. Return them as a bullet list in markdown format. If there are no relevant items, return "No relevant content found."

Summary:
%s`, config.FocusTopics, summary)

	reqBody := ChatCompletionRequest{
		Model: config.LLMAPIModel,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", config.LLMAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.LLMAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API error: %s - %s", resp.Status, string(body))
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	focused := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	return focused, nil
}
