package main

import (
	/*My project packages are kept at /go/src/app
	In order to use my utility packages I need to
	import my packages like this. The project packages path is
	mentioned in Dockerfile*/
	"app/utilities"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

/*
This program returns grep results from local machine. There is only one log file that is grepped
the name of the log file that gets grepped is machine2.log for local and machine1.log for remote
*/
func main() {
	argsWithProg := os.Args[1:]
	fmt.Printf("The command line arguments: %s\n", argsWithProg)
	if len(argsWithProg) < 4 {
		log.Fatal("Your grep command must be of the form '<grep> <options> <pattern> <file>'")
	}

	if strings.Compare(argsWithProg[0], "grep") != 0 {
		log.Fatal("It wasn't a grep command")
	}

	/*
		Launch a goroutine to act as a server for peer's grep requests
	*/
	/**
	1. Launch a go routine to get local grep
	2. Launch a go routine to get peer grep results
	3. We have all the grep results, send it to client
	**/
	//cmd := "grep"
	//search := "tanuki"
	//logFile := "machine1.log"
	cmd := argsWithProg[0]
	option := argsWithProg[1]
	search := argsWithProg[2]
	logFile := argsWithProg[3]
	//var localGrepResult string
	lgo := utilities.LocalGrep(cmd, option, search, logFile)

	fmt.Println("Response from local machine:", lgo)
	//Get grep result from remote machines
	v := url.Values{}
	v.Set("ask", cmd)
	v.Set("option", option)
	v.Add("search", search)
	v.Add("file", logFile)
	l := utilities.ReadConfig("text.txt")
	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		//fmt.Println(e.Value)
		if str, ok := e.Value.(string); ok {
			/* act on str */
			c := utilities.RemoteGrep(str, v)
			fmt.Printf("Response from %s:%s", str, <-c)
		} else {
			/* not string */
			fmt.Println("The server names file doesn't have strings")
		}

	}

}
