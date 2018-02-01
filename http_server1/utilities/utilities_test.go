package utilities

/*
go test
go convey
code coverage

*/

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
		name string
		commandstringsslice []string
		output string
	}{
		{"exporting 8080 grepped", []string{"grep", "-c", "8080", "/go/src/app/Dockerfile"}, "1"},
		{"local log file creation grepped", []string{"grep", "-c", "LOCAL LOG", "/go/src/app/local.log"}, "1"},
	}

	/*Within this loop I will form subtests to test every function in this package.
	That way I don't have to spend time on book keeping of different package names
	that I don't even know now.
	Later on I will make smaller and more meaningful packages and place the corresponding
	tests in those packages
	*/
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T){
			var output string
			output = LocalGrep(tc.commandstringsslice)
			outputexpected := tc.output
			outputexpected += "\n"
			if output != outputexpected {
				t.Fatalf("For test %s, got %q want %q\n", tc.name, output, tc.output)
			}
		})
		
	}
	
}
