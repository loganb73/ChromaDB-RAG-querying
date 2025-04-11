package project05loganb73

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

func SetupAiClient() (aiClient *openai.Client) {
	//openai setup
	openaiKey := os.Getenv("OPENAI_API_KEY")

	aiClient = openai.NewClient(openaiKey)
	return aiClient
}

func GetNamedEntities(aiClient *openai.Client, prompt string) (namedEntityJson string, err error) {

	fullPrompt := "In this question, identify the names of people, locations, and university departments. If you don't find any of the following, write the json field as `[]`." + prompt

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fullPrompt,
		},
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a chatbot 
			which answers question about the USF course catalog and gives responses in JSON format with keys 'people' 'locations' and 'departments'`,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: messages,
	}

	resp, err := aiClient.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
