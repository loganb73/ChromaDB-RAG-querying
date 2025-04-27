package main

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

	fullPrompt := `In this question, identify the names of people, locations, and university departments. Return a JSON object with exactly these three fields: "people", "locations", and "departments". If any field is empty, set it to an empty string. Example format:
	{
		"people": "John Smith",
		"locations": "Science Building",
		"departments": "Computer Science"
	}
	
	Question: ` + prompt

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fullPrompt,
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: `You are a chatbot that extracts named entities from questions about the USF course catalog. You must respond with valid JSON containing exactly these three fields: "people", "locations", and "departments". Use empty strings for fields with no values.`,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4TurboPreview,
		Messages: messages,
	}

	resp, err := aiClient.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		return "", err
	}

	fmt.Printf("AI Response: %s\n", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}

func RagQuery(aiClient *openai.Client, prompt string) (respSlice []string, err error) {
	fmt.Printf("running RagQuery\n")

	namedEntityJson, err := GetNamedEntities(aiClient, prompt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received JSON string: %s\n", namedEntityJson)

	type jsonStruct struct {
		People      string `json:"people"`
		Locations   string `json:"locations"`
		Departments string `json:"departments"`
	}
	var resultStruct jsonStruct
	if err := json.Unmarshal([]byte(namedEntityJson), &resultStruct); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return nil, err
	}

	fmt.Printf("Unmarshaled struct: %+v\n", resultStruct)

	metadata := make(map[string]interface{})

	metadataFound := false
	if resultStruct.People != "" {
		canonicalName := QueryDb(resultStruct.People, "professor-collection")
		if len(canonicalName) > 0 {
			metadata["professor"] = canonicalName[0]
		}
		metadataFound = true
	}
	if resultStruct.Locations != "" {
		canonicalLocation := QueryDb(resultStruct.Locations, "location-collection")
		if len(canonicalLocation) > 0 {
			metadata["location"] = canonicalLocation[0]
		}
		metadataFound = true
	}
	if resultStruct.Departments != "" {
		canonicalDepartment := QueryDb(resultStruct.Departments, "class-collection")
		if len(canonicalDepartment) > 0 {
			metadata["SUBJ"] = canonicalDepartment[0]
		}
		metadataFound = true
	}
	if !metadataFound {
		fmt.Printf("prompt: %s\n", prompt)
		respSlice = QueryDb(prompt, "full-collection")
		return respSlice, nil
	}

	respSlice = QueryWithMetadata(prompt, "full-collection", metadata)

	return respSlice, nil
}
