package utilities

import (
	"testing"
	"os/exec"
)

func TestLocalGrep(t *testing.T) {
	//Test whether grep command is found
	_, err := exec.LookPath("grep")
  	if err != nil {
  		t.Fatalf("grep command not found")
  	}
	
	/*Test basic grep functionality
	*/
	tt := []struct{
		commandstringsslice []string
		output string
	}{
		{[]string{"grep", "-c", "8080", "/go/src/app/Dockerfile"}, "1"},
		{[]string{"grep", "-c", "LOCAL LOG", "/go/src/app/local.log"}, "1"},
	}

	for _, tc := range tt {
		var output string
		output = LocalGrep(tc.commandstringsslice)
		outputexpected := tc.output
		outputexpected += "\n"
		if output != outputexpected {
			t.Fatalf(`Sprintf("%%s", empty("7")) = %q want %q`, output, "7")
		}
	}
	
}
