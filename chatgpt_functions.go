package project05loganb73

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	fullPrompt := "In this question, identify the names of people, locations, and university departments. If you don't find any of the following, write the json field as ``. Do not put square brackets around filled fields." + prompt

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

func RagQuery(aiClient *openai.Client, prompt string) (resp string, err error) {
	fmt.Printf("running RagQuery\n")

	namedEntityJson, err := GetNamedEntities(aiClient, prompt)
	if err != nil {
		log.Fatal(err)
	}

	type jsonStruct struct {
		People      string `json:"people"`
		Locations   string `json:"locations"`
		Departments string `json:"departments"`
	}
	var resultStruct jsonStruct
	json.Unmarshal([]byte(namedEntityJson), &resultStruct)

	metadata := make(map[string]interface{})

	if resultStruct.People != "" {
		canonicalName := QueryDb(resultStruct.People, "professor-collection")
		metadata["professor"] = canonicalName
	}
	if resultStruct.Locations != "" {
		canonicalLocation := QueryDb(resultStruct.People, "location-collection")
		metadata["location"] = canonicalLocation
	}
	if resultStruct.Departments != "" {
		canonicalDepartment := QueryDb(resultStruct.People, "class-collection")
		metadata["department"] = canonicalDepartment
	}

	return "end of function", nil
}

func GetCanonicalName(fuzzyName string) (canonicalName string) {
	switch fuzzyName {

	}
	return canonicalName
}
