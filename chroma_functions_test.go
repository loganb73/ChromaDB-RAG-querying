package project05loganb73

import (
	"fmt"
	"testing"
)

func TestSetupDb(t *testing.T) {
	SetupDb()
}

func TestQueryDb(t *testing.T) {
	resp := QueryDb("Where does Bioinformatics meet?", "full-collection")
	fmt.Print(resp)
}
