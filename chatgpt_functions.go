package project05loganb73

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Chat(prompt string) (response string, err error) {
	//openai setup
	openaiKey := os.Getenv("openai_key_temp")

	aiClient := openai.NewClient(openaiKey)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a bubbly, excited chatbot 
			which answers question about the USF course catalog`,
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
