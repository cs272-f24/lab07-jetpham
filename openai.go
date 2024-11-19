package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct {
	client  *openai.Client
	context context.Context
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
		context: context.Background(),
	}
}

func (c *OpenAIClient) CreateCompletion(prompt string) (string, error) {
	resp, err := c.client.Chat.Completions.New(c.context, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
