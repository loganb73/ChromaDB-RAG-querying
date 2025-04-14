package project05loganb73

import (
	"encoding/json"
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

func TestPhil(t *testing.T) {
	queryString := "I would like to take a Rhetoric course from Phil Choong. What can I take?"
	aiClient := SetupAiClient()

	namedEntities, err := GetNamedEntities(aiClient, queryString)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(namedEntities)
	type jsonStruct struct {
		People      string `json:"people"`
		Locations   string `json:"locations"`
		Departments string `json:"departments"`
	}
	var resultStruct jsonStruct
	json.Unmarshal([]byte(namedEntities), &resultStruct)

	if resultStruct.People != "" {
		fmt.Printf("query contains canonical name: %s\n", resultStruct.People)
	}
	if resultStruct.Locations != "" {
		fmt.Printf("query contains canonical location: %s\n", resultStruct.Locations)
	}
	if resultStruct.Departments != "" {
		fmt.Printf("query contains canonical department: %s\n", resultStruct.Departments)
	}

}

func TestRagQuery(t *testing.T) {
	aiClient := SetupAiClient()
	RagQuery(aiClient, "Which philosophy courses are offered this semester?")
}
