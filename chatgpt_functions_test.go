package project05loganb73

import (
	"fmt"
	"testing"
)

func TestChat(t *testing.T) {

	resp, err := Chat("who teaches cs727 this semester?")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", resp)
}
