package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Initialize OpenAI client
	aiClient := SetupAiClient()
	if aiClient == nil {
		log.Fatal("Failed to setup OpenAI client. Make sure OPENAI_API_KEY environment variable is set.")
	}

	// Initialize ChromaDB
	_ = SetupDb()

	// Create a scanner for reading user input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the RAG Query Interface!")
	fmt.Println("Enter your questions (type 'exit' to quit):")

	// Continuous prompt loop
	for {
		fmt.Print("\nEnter your question: ")
		if !scanner.Scan() {
			break
		}

		prompt := strings.TrimSpace(scanner.Text())
		if prompt == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if prompt == "" {
			continue
		}

		fmt.Printf("\nProcessing query: %s\n", prompt)
		resp, err := RagQuery(aiClient, prompt)
		if err != nil {
			fmt.Printf("Error during RAG query: %v\n", err)
			continue
		}

		fmt.Println("\nResponse:")
		for _, r := range resp {
			fmt.Println(r)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
