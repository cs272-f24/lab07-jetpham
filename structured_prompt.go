package main

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/openai/openai-go"
)

// A struct that will be converted to a Structured Outputs response schema
type IsSimilar struct {
	IsSimilar bool   `json:"isSimilar" jsonschema_description:"If the two texts convey similar information"`
	Reasoning string `json:"reasoning" jsonschema_description:"The reasoning behind the similarity determination"`
}

// Generate the JSON schema at initialization time

func isSimilar(setup Setup, text1, text2 string) (bool, string) {
	var IsSimilarResponseSchema = GenerateSchema[IsSimilar]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("is_similar"),
		Description: openai.F("whether two texts are similar"),
		Schema:      openai.F(IsSimilarResponseSchema),
		Strict:      openai.Bool(true),
	}
	isSimilarPrompt := fmt.Sprintf("Do the following texts convey the same information?\n\nText 1: %s\n\nText 2: %s", text1, text2)
	// Query the Chat Completions API
	chat, err := setup.openAIClient.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(isSimilarPrompt),
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
		panic(err.Error())
	}

	// The model responds with a JSON string, so parse it into a struct
	similarity := IsSimilar{}
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &similarity)
	if err != nil {
		panic(err.Error())
	}

	return similarity.IsSimilar, similarity.Reasoning
}

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
