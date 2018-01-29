package utilities

import (
	"fmt"
	"log"
	"net/url"
)

type node struct {
	hostname string
}

func (n *node) Grep(commandstring []string) <- chan string {
	fmt.Println("Commandstring:", commandstring)
	v := url.Values{}
	v.Set("ask", commandstring[0])
	v.Set("option", commandstring[1])
	v.Add("search", commandstring[2])
	v.Add("file", commandstring[3])
	return RemoteGrep(n.hostname, v)
}

type grepper interface {
	Grep(commandstring ...string) string
}

type Cluster struct {
	local node
	nodes []node
}


func (c *Cluster) NewCluster(configFileName string) {
	l := ReadConfig(configFileName)
	
	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		//fmt.Println(e.Value)
		if str, ok := e.Value.(string); ok {
			node := node{hostname: str}
			c.nodes = append(c.nodes, node)
			fmt.Println("Node's hostname:", node)
		} else {
			log.Fatal("The server names file doesn't have strings")
		}
	}

}

func (c *Cluster) Grep(commandstring []string) string {
	//Get local grep
	lg := LocalGrep(commandstring)
	//Get remote grep
	for _, node := range(c.nodes) {
		fmt.Println("Remote grep", <-node.Grep(commandstring))
	}
	return lg
}

/*
		Launch a goroutine to act as a server for peer's grep requests
	
	
	1. Launch a go routine to get local grep
	2. Launch a go routine to get peer grep results
	3. We have all the grep results, send it to client
	
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
			
			c := utilities.RemoteGrep(str, v)
			fmt.Printf("Response from %s:%s", str, <-c)
		} else {
			
			fmt.Println("The server names file doesn't have strings")
		}

	}
**/