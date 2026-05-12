package llama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

func (c *Client) Chat(model string, messages []Message) (string, error) {
	req := chatRequest{
		Model:    model,
		Messages: messages,
	}

	return c.doChat(req)
}

func (c *Client) ChatJSON(model string, messages []Message) (string, error) {
	req := chatRequest{
		Model:    model,
		Messages: messages,
		ResponseFormat: &responseFormat{
			Type: "json_object",
		},
	}

	return c.doChat(req)
}

func (c *Client) doChat(reqBody chatRequest) (string, error) {
	body, err := json.Marshal(reqBody)

	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/v1/chat/completions", c.baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"llama.cpp returned status %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var chatResp chatResponse

	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}
