package main

import (
	"log"
	"os/exec"
)

func localgrep(ask, search, file string) string {
	cmd := exec.Command(ask, search, file)
	stdOutStdErr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(stdOutStdErr)
}
