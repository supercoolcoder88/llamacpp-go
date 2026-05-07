package llama

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type LlamaClient struct {
	Host string
}

func NewLlamaClient(url string) *LlamaClient {
	return &LlamaClient{Host: url}
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (c *LlamaClient) Chat(messages []Message, temperature float64, maxTokens int) (*ChatCompletionResponse, error) {
	if temperature < 0 || temperature > 2 {
		return nil, errors.New("temperature must be between 0 and 2")
	}

	if maxTokens <= 0 {
		return nil, errors.New("max tokens must be greater than 0")
	}

	body, err := json.Marshal(&ChatCompletionRequest{
		Model:       "llama.cpp",
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	})

	if err != nil {
		return nil, err
	}

	r, err := http.Post(c.Host+"/v1/chat/completions", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var response ChatCompletionResponse

	return &response, json.NewDecoder(r.Body).Decode(&response)
}
