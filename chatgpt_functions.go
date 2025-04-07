package project05loganb73

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

func SetupAiClient() (aiClient *openai.Client) {
	//openai setup
	openaiKey := os.Getenv("openai_key_temp")

	aiClient = openai.NewClient(openaiKey)
	return aiClient
}

func GetNamedEntity(aiClient *openai.Client, prompt string) (namedEntityJson string, err error) {

	fullPrompt := "respond with the full name of the person named in this question. " + prompt

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fullPrompt,
		},
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a chatbot 
			which answers question about the USF course catalog and gives responses in JSON format`,
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
