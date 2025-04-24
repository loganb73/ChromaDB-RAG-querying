package main

import (
	"fmt"
	"log"
)

func main() {
	// Initialize OpenAI client
	aiClient := SetupAiClient()
	if aiClient == nil {
		log.Fatal("Failed to setup OpenAI client. Make sure OPENAI_API_KEY environment variable is set.")
	}

	// Initialize ChromaDB
	_ = SetupDb()

	// Example query
	prompt := "What courses is Phil Peterson teaching in Fall 2024?"
	fmt.Printf("Query: %s\n", prompt)

	resp, err := RagQuery(aiClient, prompt)
	if err != nil {
		log.Fatalf("Error during RAG query: %v", err)
	}

	fmt.Printf("Response: %s\n", resp)
}
