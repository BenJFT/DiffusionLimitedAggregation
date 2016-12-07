package tools

import "testing"
import "fmt"

func TestSingleSpace(t *testing.T) {
	fmt.Println(SingleSpace("a   b"))
	fmt.Println(SingleSpace("a  b              c"))

	// Output:
	// a b
	// a b c
}
