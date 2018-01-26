package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

/*
Launch the cluster
*/
func main() {
	argsWithProg := os.Args[1:]
	if len(argsWithProg) < 1 {
		log.Output(1, "Usage:\ncluster up\ncluster build\ncluster info\ncluster stop\ncluster stopforce\n")
		os.Exit(0)
	}

	switch {
	case argsWithProg[0] == "up":
		c := runLocalCommand(argsWithProg)
		output, ok := <-c
		if !ok {
			fmt.Println("Done!")
			os.Exit(0)
		}
		fmt.Println(output)
	case argsWithProg[0] == "build":
		c := runLocalCommand(argsWithProg)
		output, ok := <-c
		if !ok {
			fmt.Println("Done!")
			os.Exit(0)
		}
		fmt.Println(output)
	default:
		log.Output(1, "Usage:\ncluster run\ncluster build\ncluster stop\ncluster stopforce\n")
		os.Exit(0)
	}

}

func runLocalCommand(cmd []string) <-chan string {
	c := make(chan string)

	go func() {
		var status string

		// Prepare the command string

		switch {
		case cmd[0] == "up":
			shell := exec.Command("C:\\Program Files\\Docker\\Docker\\resources\\bin\\docker-compose.exe",
				"up")
			err := shell.Start()
			if err != nil {
				c <- err.Error()
			}
			status = "Cluster is being brought up!"
		case cmd[0] == "build":
			shell := exec.Command("C:\\Program Files\\Docker\\Docker\\resources\\bin\\docker-compose.exe",
				"build")
			err := shell.Start()
			if err != nil {
				c <- err.Error()
			}
			status = "Building and bring the Cluster up! please sheck status"
		}
		c <- status
		close(c)
	}()
	return c
}
