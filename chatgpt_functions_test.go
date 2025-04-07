package project05loganb73

import (
	"fmt"
	"testing"
)

func TestGetNamedEntity(t *testing.T) {
	fmt.Printf("running TestGetNamedEntity\ns")

	aiClient := SetupAiClient()

	resp, err := GetNamedEntity(aiClient, "What Courses is Phil Peterson teaching in Fall 2024?")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", resp)
}
