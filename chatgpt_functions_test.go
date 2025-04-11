package project05loganb73

import (
	"fmt"
	"testing"
)

type GetNamedEntitiesStruct struct {
	query        string
	expectedJson string
}

func setupNamedEntityTestCases() []GetNamedEntitiesStruct {
	return []GetNamedEntitiesStruct{
		{"What courses is Phil Peterson teaching in Fall 2024?",
			`{
  "people": ["Phil Peterson"],
  "locations": [],
  "departments": []
}`},
		{"Which philosophy courses are offered this semester?",
			`{
  "people": [],
  "locations": [],
  "departments": ["Philosophy"]
}`},
		{"Where does Bioinformatics meet?",
			`{
  "people": [],
  "locations": [],
  "departments": ["Bioinformatics"]
}`},
		{"Can I learn guitar this semester?",
			`{
  "people": [],
  "locations": [],
  "departments": []
}`},
		{"I would like to take a Rhetoric course from Phil Choong. What can I take?",
			`{
  "people": ["Phil Choong"],
  "locations": [],
  "departments": ["Rhetoric"]
}`},
	}
}

func TestGetNamedEntities(t *testing.T) {
	fmt.Printf("running TestGetNamedEntity\n")

	aiClient := SetupAiClient()
	tests := setupNamedEntityTestCases()

	for _, test := range tests {
		resp, err := GetNamedEntities(aiClient, test.query)
		if err != nil {
			t.Error(err)
		}

		//remove gpt mess
		cleanResultString := resp[8 : len(resp)-4]

		if cleanResultString != test.expectedJson {
			t.Errorf("resp:\n%s\ndoesn't match expected:\n%s\n\n", cleanResultString, test.expectedJson)
		}
	}

}
