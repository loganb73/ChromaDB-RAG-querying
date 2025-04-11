package project05loganb73

import (
	"encoding/json"
	"fmt"
	"testing"
)

// func TestSetupDb(t *testing.T) {
// 	SetupDb()
// }

func TestQueryDb(t *testing.T) {
	fmt.Printf("running TestQueryDb\n")
	resp := QueryDb("Where does Bioinformatics meet?", "full-collection")
	fmt.Print(resp)
}

func TestPhil(t *testing.T) {
	queryString := "What courses is Phil Peterson teaching in Fall 2024?"
	aiClient := SetupAiClient()

	namedEntities, err := GetNamedEntities(aiClient, queryString)
	if err != nil {
		t.Error(err)
	}

	//remove gpt mess
	//cleanResultString := namedEntities[8 : len(namedEntities)-4]
	fmt.Println(namedEntities)
	type jsonStruct struct {
		People      string `json:"people"`
		Locations   string `json:"locations"`
		Departments string `json:"departments"`
	}
	var resultStruct jsonStruct
	json.Unmarshal([]byte(namedEntities), &resultStruct)

	fmt.Printf("query contains canonical name: %s\n", resultStruct.People)
}
