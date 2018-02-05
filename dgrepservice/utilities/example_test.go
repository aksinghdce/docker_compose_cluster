package utilities_test

import (
	"app/utilities"
	"fmt"
)

func ExampleLocalGrep() {
	lc := utilities.LocalGrep([]string{"grep", "-c", "8080", "/go/src/app/Dockerfile"})
	fmt.Println(lc)
	// Output:
	// 1
}