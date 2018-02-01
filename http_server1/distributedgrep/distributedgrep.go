package main

import (
	/*My project packages are kept at /go/src/app
	when they run inside of a docker container environment.
	In order to use my utility packages I need to
	import my packages like this. There is a problem that I noticed, 
	because app/utilities is not in the GOPATH of my windows go environment
	I can't use the visual studio code's tools for go project management

	The project packages path for docker containers is mentioned in Dockerfile*/
	"app/utilities"
	"fmt"
	"log"
	"os"
)

/*
This program returns grep results from local machine. There is only one log file that is grepped
the name of the log file that gets grepped is machine2.log for local and machine1.log for remote
*/
func main() {
	argsWithProg := os.Args[1:]
	if len(argsWithProg) < 2 {
		log.Fatal("<grep> <options> <pattern> <file>")
	}
	
	if argsWithProg[0] != "grep" {
		log.Fatal("Panic: Wasn't a grep command")
	}

	var cluster utilities.Cluster
	//text.txt is configuration file that contains the names of the nodes in the cluster
	cluster.NewCluster("nodenames.txt")
	fmt.Printf(cluster.Grep(argsWithProg))
}
