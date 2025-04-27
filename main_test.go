package main

import (
	"fmt"
	"testing"
)

func TestPHIL(t *testing.T) {
	fmt.Printf("\n=== Starting TestPHIL ===\n")

	// Then create AI client and run query
	aiClient := SetupAiClient()
	fmt.Printf("AI Client created\n")

	// Now run the actual test
	respSlice, err := RagQuery(aiClient, "Which philosophy courses are offered this semester?")
	if err != nil {
		t.Errorf("RagQuery error: %v", err)
	}

	fmt.Printf("Response slice length: %d\n", len(respSlice))
	fmt.Printf("Response contents: %#v\n", respSlice)

	if len(respSlice) == 0 {
		t.Error("Expected non-empty response slice")
	}

	fmt.Printf("=== End TestPhil ===\n\n")
}

func TestPhil(t *testing.T) {
	fmt.Printf("\n=== Starting TestPhil ===\n")

	// Then create AI client and run query
	aiClient := SetupAiClient()
	fmt.Printf("AI Client created\n")

	// Now run the actual test
	respSlice, err := RagQuery(aiClient, "What courses is Phil Peterson teaching in Fall 2024?")
	if err != nil {
		t.Errorf("RagQuery error: %v", err)
	}

	fmt.Printf("Response slice length: %d\n", len(respSlice))
	fmt.Printf("Response contents: %#v\n", respSlice)

	if len(respSlice) == 0 {
		t.Error("Expected non-empty response slice")
	}

	fmt.Printf("=== End TestPhil ===\n\n")
}

func TestBio(t *testing.T) {
	fmt.Printf("\n=== Starting TestBio ===\n")

	// Then create AI client and run query
	aiClient := SetupAiClient()
	fmt.Printf("AI Client created\n")

	// Now run the actual test
	respSlice, err := RagQuery(aiClient, "Where does Bioinformatics meet?")
	if err != nil {
		t.Errorf("RagQuery error: %v", err)
	}

	fmt.Printf("Response slice length: %d\n", len(respSlice))
	fmt.Printf("Response contents: %#v\n", respSlice)

	if len(respSlice) == 0 {
		t.Error("Expected non-empty response slice")
	}

	fmt.Printf("=== End TestBio ===\n\n")
}

func TestGuitar(t *testing.T) {
	fmt.Printf("\n=== Starting TestGuitar ===\n")

	// Then create AI client and run query
	aiClient := SetupAiClient()
	fmt.Printf("AI Client created\n")

	// Now run the actual test
	respSlice, err := RagQuery(aiClient, "Can I learn guitar this semester?")
	if err != nil {
		t.Errorf("RagQuery error: %v", err)
	}

	fmt.Printf("Response slice length: %d\n", len(respSlice))
	fmt.Printf("Response contents: %#v\n", respSlice)

	if len(respSlice) == 0 {
		t.Error("Expected non-empty response slice")
	}

	fmt.Printf("=== End TestGuitar ===\n\n")
}

func TestMultiple(t *testing.T) {
	fmt.Printf("\n=== Starting TestMultiple ===\n")

	// Then create AI client and run query
	aiClient := SetupAiClient()
	fmt.Printf("AI Client created\n")

	// Now run the actual test
	respSlice, err := RagQuery(aiClient, "I would like to take a Rhetoric course from Phil Choong. What can I take?")
	if err != nil {
		t.Errorf("RagQuery error: %v", err)
	}

	fmt.Printf("Response slice length: %d\n", len(respSlice))
	fmt.Printf("Response contents: %#v\n", respSlice)

	if len(respSlice) == 0 {
		t.Error("Expected non-empty response slice")
	}

	fmt.Printf("=== End TestMultiple ===\n\n")
}
