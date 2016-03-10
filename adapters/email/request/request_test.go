package request

import (
	"fmt"
	"testing"

	"Gateway311/adapters/email/data"
	"Gateway311/adapters/email/logs"
)

var Debug = true

func init() {
	logs.Init(Debug)

	fmt.Println("Reading config...")
	if err := data.Init("config.json"); err != nil {
		fmt.Printf("Init() failed: %s", err)
	}
}

type testResultS struct {
	input string
	isOK  bool
}

func isOK(e error) bool {
	if e == nil {
		return false
	}
	return true
}

func TestXXX(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestXXX] =============================\n\n")

}
