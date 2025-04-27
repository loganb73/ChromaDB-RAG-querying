package main_test

import (
	"fmt"
	"testing"

	. "github.com/loganb73/cs272/project05-loganb73"
)

func TestSetupDb(t *testing.T) {
	SetupDb()
}

func TestQueryDb(t *testing.T) {
	fmt.Printf("running TestQueryDb\n")
	resp := QueryDb("Where does Bioinformatics meet?", "full-collection")
	fmt.Print(resp)
}

func TestMetadataQuery(t *testing.T) {
	fmt.Printf("running TestMetadataQuery\n")
	metadata := make(map[string]interface{})
	metadata["professor"] = "Philip Peterson"
	metadata["location"] = "MH"
	metadata["subject"] = "CS"

	resp := QueryWithMetadata("What does Philip Peterson teach in MH?", "full-collection", metadata)
	fmt.Print(resp)
}
