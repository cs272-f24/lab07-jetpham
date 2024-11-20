package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/openai/openai-go"
)

func toolCallingAgent(setup Setup, prompt string) string {
	// using the official example:
	//https://github.com/openai/openai-go/blob/main/examples/chat-completion-tool-calling/main.go

	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Tools: openai.F([]openai.ChatCompletionToolParam{
			{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name:        openai.String("get_courses"),
					Description: openai.String("get course information"),
					Parameters: openai.F(openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"prompt": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"prompt"},
					}),
				}),
			},
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	}

	completion, err := setup.openAIClient.client.Chat.Completions.New(context.TODO(), params)
	if err != nil {
		log.Printf("Error creating chat completion: %v", err)
		return ""
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	// If there was not tool calls, crashout
	if len(toolCalls) == 0 {
		log.Printf("No function call")
		return completion.Choices[0].Message.Content
	}

	// If there was tool calls, continue
	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_courses" {
			// Extract the prompt from the function call arguments
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				log.Printf("Error unmarshalling arguments: %v", err)
				continue
			}
			prompt := args["prompt"].(string)

			log.Printf("%v(\"%s\")", toolCalls[0].Function.Name, prompt)

			// Call the getCourses function with the arguments requested by the model
			courses, err := getCourses(&setup, prompt)
			if err != nil {
				log.Printf("Error getting courses: %v", err)
				continue
			}
			coursesJSON, _ := json.Marshal(courses)
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, string(coursesJSON)))
		}
	}

	completion, err = setup.openAIClient.client.Chat.Completions.New(context.TODO(), params)
	if err != nil {
		log.Printf("Error creating chat completion: %v", err)
		return ""
	}

	return completion.Choices[0].Message.Content
}
