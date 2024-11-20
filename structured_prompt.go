package main

import (
	"encoding/json"

	"github.com/openai/openai-go"
)

func (openAIClient *OpenAIClient) GetCourseFilter(prompt, systemPrompt string) CourseFilter {
	// Based off of https://github.com/openai/openai-go structured output example
	var CourseFilterResponseSchema = GenerateSchema[CourseFilter]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("course_filter"),
		Description: openai.F("Criteria for filtering for a course"),
		Schema:      openai.F(CourseFilterResponseSchema),
		Strict:      openai.Bool(true),
	}

	// Query the Chat Completions API
	chat, err := openAIClient.client.Chat.Completions.New(openAIClient.context, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		// Only certain models can perform structured outputs
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		panic(err)
	}
	if chat == nil || len(chat.Choices) == 0 || chat.Choices[0].Message.Content == "" {
		panic("invalid response from OpenAI API")
	}

	// extract into a well-typed struct
	courseFilter := CourseFilter{}
	_ = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &courseFilter)

	return courseFilter
}
