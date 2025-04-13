package project05loganb73

import (
	"fmt"
	"testing"
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
	fmt.Printf("running TestQueryDb\n")
	resp := QueryDb("Where does Bioinformatics meet?", "full-collection")
	fmt.Print(resp)
}
